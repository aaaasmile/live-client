package srv

import (
	"log"
	"net/http"
	"time"

	"github.com/aaaasmile/live-client/conf"
	"github.com/aaaasmile/live-client/web/idl"
	"github.com/aaaasmile/live-client/web/srv/checker"
	"github.com/aaaasmile/live-client/web/srv/dataprepare"
)

type ParamCallSync struct {
	Debug       bool `json:"debug"`
	ForceSource bool `json:"forcesource"`
	ForceServer bool `json:"forceserver"`
}

func (hc *HandlerPrjReq) HandleCallSync(w http.ResponseWriter, req *http.Request) error {
	paraDef := ParamCallSync{}
	if err := parseBodyReq(req, &paraDef); err != nil {
		return err
	}
	log.Println("Synch with options: ", paraDef)

	hc.Debug = paraDef.Debug

	return hc.synchAndCompare(w, &paraDef)
}

func (hc *HandlerPrjReq) synchAndCompare(w http.ResponseWriter, paraDef *ParamCallSync) error {
	if hc.chServerFilePersist != nil {
		<-hc.chServerFilePersist
	}
	cherr := make(chan error, 2)
	chfinished := make(chan struct{}, 2)

	go func() {
		// Remote checker
		var err error
		start := time.Now()
		preData := dataprepare.ObjectInfoPre{
			Force:           paraDef.ForceServer,
			Debug:           conf.Current.DebugVerbose || hc.Debug,
			Store:           hc.storeSeverFile,
			LiteDB:          hc.liteDB,
			SyncRepo:        conf.Current.SyncRepo,
			RemoteServerURL: conf.Current.RemoteServerURL,
			TempDir:         conf.Current.TempDir,
		}
		hc.chServerFilePersist, err = preData.PrepareData()
		if err != nil {
			cherr <- err
		}

		log.Printf("Server items prepare call duration: %v\n", time.Now().Sub(start))
		chfinished <- struct{}{}
	}()

	if hc.chSourceFilePersist != nil {
		<-hc.chSourceFilePersist
	}

	go func() {
		// Local Source File checker
		var err error
		start := time.Now()
		preData := dataprepare.ObjectInfoPre{
			Force:           paraDef.ForceSource,
			Debug:           conf.Current.DebugVerbose || hc.Debug,
			Store:           hc.storeSourceFile,
			LiteDB:          hc.liteDB,
			SyncRepo:        conf.Current.SyncRepo,
			RemoteServerURL: conf.Current.RemoteServerURL,
			TempDir:         conf.Current.TempDir,
		}

		hc.chSourceFilePersist, err = preData.PrepareData()
		if err != nil {
			cherr <- err
		}
		// Post processing for source files is not needed
		log.Printf("prepare SourceFileData call duration: %v\n", time.Now().Sub(start))
		chfinished <- struct{}{}
	}()

	// wait until both processes have been finished
	finCount := 2
loop:
	for {
		select {
		case err := <-cherr:
			hc.chSourceFilePersist = nil
			hc.chServerFilePersist = nil
			log.Println("Error in prepare data")
			return err
		case <-chfinished:
			finCount--
			if finCount <= 0 {
				log.Println("Prepare data finished")
				break loop
			}
		}
	}

	return hc.doCompare(w)
}

func (hc *HandlerPrjReq) HandleViewCurrentDiff(w http.ResponseWriter, req *http.Request) error {
	storeServer := checker.NewStore(idl.OTPServerFile)
	storeSourceFile := checker.NewStore(idl.OTPSourceFile)
	lite := hc.liteDB

	if hc.chServerFilePersist != nil {
		<-hc.chServerFilePersist
	}

	if err := storeServer.PopulateObjs(&lite); err != nil {
		return err
	}

	if hc.chSourceFilePersist != nil {
		<-hc.chSourceFilePersist
	}

	if err := storeSourceFile.PopulateObjs(&lite); err != nil {
		return err
	}

	hc.storeSeverFile = &storeServer
	hc.storeSourceFile = &storeSourceFile

	log.Printf("Store len: Remote %d, Source file %d", len(storeServer.InfoObjects), len(storeSourceFile.InfoObjects))

	return hc.doCompare(w)
}

func (hc *HandlerPrjReq) doCompare(w http.ResponseWriter) error {
	ck := checker.Checker{}
	ck.CreateResultView(hc.storeSeverFile, hc.storeSourceFile)
	return writeJsonResp(w, ck)
}
