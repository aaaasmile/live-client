package srv

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aaaasmile/live-client/conf"
	"github.com/aaaasmile/live-client/web/idl"
	"github.com/aaaasmile/live-client/web/srv/remote"
)

func (hc *HandlerPrjReq) handleImportSelectionToServer(w http.ResponseWriter, req *http.Request) error {
	paraDef := DoWithSelectionReq{}

	if err := parseBodyReq(req, &paraDef); err != nil {
		return err
	}

	log.Println("Import selection to Nav params len is ", len(paraDef.Selected))

	if len(paraDef.Selected) == 0 {
		return fmt.Errorf("Nothing valid to import received")
	}

	selSrcs := idl.ObjectInfoColl{}
	for _, key := range paraDef.Selected {
		obj := hc.storeSourceFile.InfoObjects[key]
		if obj != nil {
			selSrcs = append(selSrcs, obj)
		}
	}

	if len(paraDef.Selected) == 0 {
		return fmt.Errorf("Nothing valid to import in store. Is the project synch with the backend?")
	}

	err := remote.ImportObjectsIntoServer(selSrcs, conf.Current.RemoteServerURL)
	if err != nil {
		return err
	}
	paraSync := ParamCallSync{}
	return hc.synchAndCompare(w, &paraSync)

}
