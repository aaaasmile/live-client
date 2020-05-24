package srv

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aaaasmile/live-client/conf"
	"github.com/aaaasmile/live-client/web/idl"
	"github.com/aaaasmile/live-client/web/srv/remote"
)

func (hc *HandlerPrjReq) handleExportToFile(w http.ResponseWriter, req *http.Request) error {
	paraDef := DoWithSelectionReq{}

	if err := parseBodyReq(req, &paraDef); err != nil {
		return err
	}

	log.Println("Export to file with param: ", len(paraDef.Selected))

	return hc.execExport(w, paraDef.Selected)
}

func (hc *HandlerPrjReq) execExport(w http.ResponseWriter, selected []string) error {
	if len(selected) == 0 {
		return fmt.Errorf("Nothing to export")
	}

	selServerObj := idl.ObjectInfoColl{}
	for _, key := range selected {
		obj := hc.storeSeverFile.InfoObjects[key]
		if obj != nil {
			selServerObj = append(selServerObj, obj)
		}
	}
	if len(selServerObj) == 0 {
		return fmt.Errorf("Nothing valid to export")
	}

	tempFiles, err := remote.ExportObjectsInTmp(selServerObj, conf.Current.TempDir, conf.Current.RemoteServerURL)
	if err != nil {
		return err
	}

	log.Println("Created export files", tempFiles)

	return hc.postProcessTmpFiles(w, tempFiles)
}

func (hc *HandlerPrjReq) postProcessTmpFiles(w http.ResponseWriter, tempFiles []string) error {
	allFiles := []string{}

	for _, tempFile := range tempFiles {
		return fmt.Errorf("TODO copy file from temp to repo, %s", tempFile)
	}

	log.Println("Exported files", allFiles)

	paraDef := ParamCallSync{}
	return hc.synchAndCompare(w, &paraDef)
}

func (hc *HandlerPrjReq) copyFile(src, dst string) error {
	if hc.Debug {
		log.Println("Copy file ", src, dst)
	}

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}
