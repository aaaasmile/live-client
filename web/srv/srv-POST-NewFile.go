package srv

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aaaasmile/live-client/web/idl"

	"github.com/aaaasmile/live-client/conf"
)

func (hc *HandlerPrjReq) handleNewFile(w http.ResponseWriter, req *http.Request) error {
	paraDef := struct {
		FileName string `json:"filename"`
	}{}

	if err := parseBodyReq(req, &paraDef); err != nil {
		return err
	}

	log.Println("NewFile with param: ", paraDef)
	if len(paraDef.FileName) == 0 {
		return fmt.Errorf("The new filename is empty")
	}
	if strings.Contains(paraDef.FileName, "-") {
		return fmt.Errorf("the - is not allowed here as file name")
	}

	newSrc := idl.SourceFile{}
	if err := newSrc.CreateNewFile(conf.Current.SyncRepo, paraDef.FileName); err != nil {
		return err
	}

	return hc.synchAndCompare(w, &ParamCallSync{})
}
