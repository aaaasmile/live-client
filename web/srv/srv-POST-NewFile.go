package srv

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/aaaasmile/live-client/conf"
	"github.com/aaaasmile/live-client/util"
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

	intName, id := util.GetInternalFilename(paraDef.FileName, "1")
	fsrc, err := hc.touchNewFile(conf.Current.SyncRepo, intName)
	if err != nil {
		return err
	}
	log.Println("Created an empty file ", intName, id, fsrc)
	return hc.synchAndCompare(w, &ParamCallSync{})
}

func (hc *HandlerPrjReq) touchNewFile(destDir string, baseName string) (string, error) {
	fname := path.Join(destDir, baseName)
	_, err := os.Stat(fname)
	if os.IsNotExist(err) {
		file, err := os.Create(fname)
		if err != nil {
			return "", err
		}
		defer file.Close()
	}
	return fname, nil
}
