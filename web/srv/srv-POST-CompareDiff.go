package srv

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/aaaasmile/live-client/conf"
	"github.com/aaaasmile/live-client/web/idl"
	"github.com/aaaasmile/live-client/web/srv/remote"
)

func (hc *HandlerPrjReq) handleOpenVsCode(w http.ResponseWriter, req *http.Request) error {
	paraDef := struct {
		Repo string `json:"repo"`
	}{}
	if err := parseBodyReq(req, &paraDef); err != nil {
		return err
	}

	cmd := conf.Current.VsCodePath
	if cmd == "" {
		return fmt.Errorf("VsCode path not found")
	}
	if paraDef.Repo == "" {
		return fmt.Errorf("VsCode Target dir is empty")
	}

	args := []string{paraDef.Repo}
	_, err := exec.Command(cmd, args...).Output()
	if err != nil {
		log.Printf("Error on executing VsCode: %v", err)
		return err
	}

	okResp := struct {
		Status string
	}{Status: "OK"}
	return writeJsonResp(w, okResp)
}

func (hc *HandlerPrjReq) handleCompareDiff(w http.ResponseWriter, req *http.Request) error {
	paraDef := DoWithSingleSelectionReq{}

	if err := parseBodyReq(req, &paraDef); err != nil {
		return err
	}

	log.Println("compare diff with param: ", paraDef)

	resp := struct {
		Status string
	}{Status: "OK"}

	key := paraDef.Selected
	fsrc, fsrv := "", ""
	if srcItem, ok := hc.storeSourceFile.InfoObjects[key]; ok {
		fsrc = path.Join(conf.Current.SyncRepo, srcItem.SourceFile.Name)
	}

	selServerObj := idl.ObjectInfoColl{}
	if serverItem, ok := hc.storeSeverFile.InfoObjects[key]; ok {
		selServerObj = append(selServerObj, serverItem)
		fsrv = serverItem.Key
	}
	if fsrc == "" && fsrv == "" {
		return fmt.Errorf("Key without any items in store")
	} else {
		log.Println("Comapring ", fsrc, fsrv)
	}

	tempFiles, err := remote.ExportObjectsInTmp(selServerObj, conf.Current.TempDir, conf.Current.RemoteServerURL)
	if err != nil {
		return err
	}

	fsrv = tempFiles[0]
	if fsrv == "" {
		fsrv, err = hc.touchEmptyFile()
		if err != nil {
			return err
		}
	}
	if fsrc == "" {
		fsrc, err = hc.touchEmptyFile()
		if err != nil {
			return err
		}
	}
	err = startBeyondCompare(fsrc, fsrv, conf.Current.BeyondComparePath)
	if err != nil {
		log.Println("Error on starting beyond compare")
		return err
	}

	return writeJsonResp(w, resp)
}

func (hc *HandlerPrjReq) touchEmptyFile() (string, error) {
	fname := path.Join(conf.Current.TempDir, "empty")
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

func startBeyondCompare(f1, f2, cmdPat string) error {
	var cmd string
	var args []string
	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c"} // do not use /start
	} else {
		log.Fatal("OS not recognized")
		return fmt.Errorf("OS not supported %s", runtime.GOOS)
	}
	args = append(args, cmdPat)
	args = append(args, f1)
	args = append(args, f2)

	log.Println("open beyond compare ", args)
	return exec.Command(cmd, args...).Start()
}
