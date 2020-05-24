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
	ForceSource  bool `json:"forcesource"`
	ForceNavObj  bool `json:"forcenavobj"`
	Debug        bool `json:"debug"`
	CheckContent bool `json:"checkcontent"`
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

	if err := hc.storeIgnoreList.PopulateObjs(&hc.liteDB); err != nil {
		return err
	}

	if hc.chNavPersist != nil {
		<-hc.chNavPersist
	}
	chNavPostProc := make(chan []*idl.ObjectInfo)
	cherr := make(chan error, 2)
	chfinished := make(chan struct{}, 2)

	go func() {
		// NAV checker
		var err error
		start := time.Now()
		preData := dataprepare.ObjectInfoPre{
			Force:        paraDef.ForceNavObj,
			Debug:        conf.Current.DebugVerbose || hc.Debug,
			ProjectInfo:  hc.ProjectInfo,
			Store:        hc.storeNav,
			IgnoreStore:  hc.storeIgnoreList,
			LiteDB:       hc.liteDB,
			CheckContent: paraDef.CheckContent,
		}
		hc.chNavPersist, err = preData.PrepareData(nil)
		if err != nil {
			cherr <- err
		}
		if paraDef.CheckContent {
			log.Println("Waiting for post processing")
			if oiPosts, ok := <-chNavPostProc; ok {
				if hc.chNavPersist != nil {
					<-hc.chNavPersist
				}
				hc.chNavPersist, err = preData.PostProcNav(oiPosts)
				if err != nil {
					cherr <- err
				}
			}
		}
		log.Printf("NavData prepare call duration: %v\n", time.Now().Sub(start))
		chfinished <- struct{}{}
	}()

	if hc.chSourceFilePersist != nil {
		<-hc.chSourceFilePersist
	}

	go func() {
		// Source File checker
		var err error
		start := time.Now()
		preData := dataprepare.ObjectInfoPre{
			Force:        paraDef.ForceSource,
			Debug:        conf.Current.DebugVerbose || hc.Debug,
			ProjectInfo:  hc.ProjectInfo,
			Store:        hc.storeSourceFile,
			IgnoreStore:  hc.storeIgnoreList,
			LiteDB:       hc.liteDB,
			CheckContent: paraDef.CheckContent,
		}

		hc.chSourceFilePersist, err = preData.PrepareData(chNavPostProc)
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
			hc.chNavPersist = nil
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
	storeNav := checker.NewStore(idl.OTPNavObj)
	storeSourceFile := checker.NewStore(idl.OTPSourceFile)
	lite := hc.liteDB

	if hc.chNavPersist != nil {
		<-hc.chNavPersist
	}

	if err := storeNav.PopulateObjs(&lite); err != nil {
		return err
	}

	if hc.chSourceFilePersist != nil {
		<-hc.chSourceFilePersist
	}

	if err := storeSourceFile.PopulateObjs(&lite); err != nil {
		return err
	}

	hc.storeNav = &storeNav
	hc.storeSourceFile = &storeSourceFile

	log.Printf("Store len nav %d, source file %d", len(storeNav.InfoObjects), len(storeSourceFile.InfoObjects))

	return hc.doCompare(w)
}

func (hc *HandlerPrjReq) doCompare(w http.ResponseWriter) error {
	ck := checker.Checker{}
	ck.CreateResultView(hc.storeNav, hc.storeSourceFile)
	return writeJsonResp(w, ck)
}
