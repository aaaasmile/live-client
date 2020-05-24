package srv

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aaaasmile/live-client/web/srv/checker"
)

type ParamCallSync struct {
	Debug bool `json:"debug"`
}

func (hc *HandlerPrjReq) HandleCallSync(w http.ResponseWriter, req *http.Request) error {
	paraDef := ParamCallSync{}
	if err := parseBodyReq(req, &paraDef); err != nil {
		return err
	}
	log.Println("Synch with options: ", paraDef)

	hc.Debug = paraDef.Debug

	return fmt.Errorf("TODO...")
}

// 	if hc.chSourceFilePersist != nil {
// 		<-hc.chSourceFilePersist
// 	}

// 	go func() {
// 		// Source File checker
// 		var err error
// 		start := time.Now()
// 		preData := dataprepare.ObjectInfoPre{
// 			Force:        paraDef.ForceSource,
// 			Debug:        conf.Current.DebugVerbose || hc.Debug,
// 			ProjectInfo:  hc.ProjectInfo,
// 			Store:        hc.storeSourceFile,
// 			IgnoreStore:  hc.storeIgnoreList,
// 			LiteDB:       hc.liteDB,
// 			CheckContent: paraDef.CheckContent,
// 		}

// 		hc.chSourceFilePersist, err = preData.PrepareData(chNavPostProc)
// 		if err != nil {
// 			cherr <- err
// 		}
// 		// Post processing for source files is not needed
// 		log.Printf("prepare SourceFileData call duration: %v\n", time.Now().Sub(start))
// 		chfinished <- struct{}{}
// 	}()

// 	// wait until both processes have been finished
// 	finCount := 2
// loop:
// 	for {
// 		select {
// 		case err := <-cherr:
// 			hc.chSourceFilePersist = nil
// 			hc.chNavPersist = nil
// 			log.Println("Error in prepare data")
// 			return err
// 		case <-chfinished:
// 			finCount--
// 			if finCount <= 0 {
// 				log.Println("Prepare data finished")
// 				break loop
// 			}
// 		}
// 	}

// 	return hc.doCompare(w)
// }

func (hc *HandlerPrjReq) doCompare(w http.ResponseWriter) error {
	ck := checker.Checker{}
	ck.CreateResultView(hc.storeSeverFile, hc.storeSourceFile)
	return writeJsonResp(w, ck)
}
