package dataprepare

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/aaaasmile/live-client/db/sqlite"
	"github.com/aaaasmile/live-client/web/idl"
	"github.com/aaaasmile/live-client/web/srv/checker"
)

type ObjectInfoPre struct {
	Force           bool
	Debug           bool
	Store           *checker.Store
	LiteDB          sqlite.LiteDB
	countTouch      int
	TempDir         string
	RemoteServerURL string
	SyncRepo        string
}

func (p *ObjectInfoPre) PrepareData() (chan idl.ResErr, error) {
	log.Println("Prepare data for the store ", p.Store.ObjTypeInProv.String())

	if err := p.Store.PopulateObjs(&p.LiteDB); err != nil {
		return nil, err
	}

	var prvch, nextch chan idl.ResErr
	if p.Force {
		log.Println("Force store rebuild")
		// flash the delete
		prvch = p.Store.StartSyncWithProv(&p.LiteDB, "Delete All", nil)
		p.Store.ResetInfoWithTrack()
		p.Store.EndTrackingChanges()
	}

	nextch = p.Store.StartSyncWithProv(&p.LiteDB, "Update", prvch)

	switch p.Store.ObjTypeInProv {
	case idl.OTPServerFile:
		if err := p.fetchRemoteObjects(); err != nil {
			p.Store.Abort()
			return nil, err
		}
		p.Store.EndTrackingChanges()
	case idl.OTPSourceFile:
		if err := p.scanSourceFiles(); err != nil {
			p.Store.Abort()
			return nil, err
		}
		p.Store.EndTrackingChanges()
	default:
		return nil, fmt.Errorf("Store type not supported %s", p.Store.ObjTypeInProv.String())
	}

	log.Println("Provider changes submitted to persist ", p.Store.ObjTypeInProv.String())

	return nextch, nil
}

func (p *ObjectInfoPre) fetchRemoteObjects() error {

	log.Println("WARNING remote fetch is not implemented")

	return nil
}

func (p *ObjectInfoPre) scanSourceFiles() error {
	log.Printf("Scan source file repository (store len is %d)", len(p.Store.InfoObjects))
	//fmt.Println("*** settings", p.CheckContent, p)
	p.countTouch = 0
	chanSources := make(idl.ChanSourceFiles, 20) // 20 is the buffer for the result

	err := p.getSourceFiles(p.SyncRepo, chanSources, p.sourceFilter)
	if err != nil {
		log.Println("GetSourceFiles returns with error", err)
		return err
	}
	count := 0
	//p.Store.Debug = true
	for itemSource := range chanSources {
		if itemSource.Err != nil {
			log.Println("WARN Ignore processing on file because error: ", itemSource.Err)
		} else {
			objNew := idl.NewObjectInfoFromSF(itemSource.SourceFile)
			//fmt.Println("*** objNew in source ", objNew.SourceFile)
			p.Store.InsertOrUpdateSingleObj(objNew)
			count++
		}
	}

	log.Printf("%d source files full scanned. Touched and ignored %d. Tot sources in store %d",
		count, p.countTouch, len(p.Store.InfoObjects))

	return nil
}

func (p *ObjectInfoPre) sourceFilter(fileInfo os.FileInfo) bool {
	// function returns true if the file needs to be fully rescanned
	fname := fileInfo.Name()

	srcItem, err := p.getSourceFileOnFileName(fname)
	if err != nil {
		log.Println("WARN Ignore item because error ", fname, err)
		return false
	}
	if srcItem == nil {
		//fmt.Println("**File not in store", fname)
		return true
	}
	// if strings.Contains(fname, "C11.txt") {
	// 	fmt.Println("**File ", fname, srcItem, fileInfo.ModTime())
	// }
	if srcItem != nil {
		objNew := idl.NewObjectInfoFromSF(*srcItem)
		if srcItem.FileSize == int(fileInfo.Size()) {
			tmf := fileInfo.ModTime()
			if srcItem.FileModTime.Unix() == tmf.Unix() {
				p.countTouch++
				p.Store.InsertOrUpdateSingleObj(objNew)
				return false
			}
		}
	}
	return true
}

func (p *ObjectInfoPre) getSourceFileOnFileName(fname string) (*idl.SourceFile, error) {
	tmp := idl.SourceFile{}
	if err := tmp.FillFromFname(fname); err != nil {
		return nil, err
	}

	if p.Store.InfoObjects[tmp.ObjectID] != nil {
		return &p.Store.InfoObjects[tmp.ObjectID].SourceFile, nil
	}
	return nil, nil
}

type FunAskForScan func(fileInfo os.FileInfo) bool

func (p *ObjectInfoPre) getSourceFiles(dirToScan string, chsources idl.ChanSourceFiles, askItemFn FunAskForScan) error {
	files, err := ioutil.ReadDir(dirToScan)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	tokens := make(chan struct{}, 5) // limit the parallel file processing. 5 seems the best.

	// 5 go routines in parallel working on 80 files at the same time seems to be the best performance.
	// Note that number of routines  build up are len(files) / limitFiles and if limitFiles = 1, then 10024 routines are build up with only 5 processing at the same time.
	// With 80 only 125 routines are build up and -race is happy.
	limitFiles := 80 // using 1 file spawn one go rutine in wait. Only 5 are processing at the same time. Race exeed the limit of 8192 simultaneously alive goroutines.
	toprocess := make([]os.FileInfo, 0, limitFiles)

	for _, ffInfo := range files {
		if !ffInfo.IsDir() {
			//fmt.Println("*** ", ffInfo.Name())
			if askItemFn(ffInfo) {
				// full scan of the source file
				toprocess = append(toprocess, ffInfo)
				if len(toprocess) >= limitFiles {
					p.startScanFiles(&wg, toprocess, tokens, dirToScan, chsources)
					toprocess = make([]os.FileInfo, 0, limitFiles)
				}
			}
		}
	}
	if len(toprocess) > 0 {
		log.Println("Remaing files to scan ", len(toprocess))
		p.startScanFiles(&wg, toprocess, tokens, dirToScan, chsources)
	}
	// closer
	go func() {
		wg.Wait()
		log.Println("FileSource processing is terminated")
		close(chsources)
	}()

	return nil
}

func (p *ObjectInfoPre) startScanFiles(wg *sync.WaitGroup, toprocess []os.FileInfo, tokens chan struct{}, dirToScan string, chsources idl.ChanSourceFiles) {
	wg.Add(1)
	go func(filesToProc []os.FileInfo) {
		defer wg.Done()

		tokens <- struct{}{}
		for _, fileInfo := range filesToProc {
			log.Println("Scan file ", fileInfo.Name())
			var res idl.SourceFileWithErr
			sf := idl.SourceFile{
				Filename:    fileInfo.Name(),
				FileSize:    int(fileInfo.Size()),
				FileModTime: fileInfo.ModTime(),
			}
			res.Err = sf.FillFromFname(fileInfo.Name())
			res.SourceFile = sf
			chsources <- res
		}
		<-tokens
	}(toprocess)
}
