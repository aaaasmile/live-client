package srv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aaaasmile/live-client/web/idl"

	"github.com/aaaasmile/live-client/conf"
	"github.com/aaaasmile/live-client/db/sqlite"
	"github.com/aaaasmile/live-client/web/srv/checker"
)

type HandlerPrjReq struct {
	storeSourceFile     *checker.Store
	storeSeverFile      *checker.Store
	liteDB              sqlite.LiteDB
	DebugSQL            bool
	Debug               bool
	chSourceFilePersist chan idl.ResErr
	chServerFilePersist chan idl.ResErr
}

type DoWithSelectionReq struct {
	Selected []string `json:"selected"`
	Debug    bool     `json:"debug"`
}

type DoWithSingleSelectionReq struct {
	Selected string `json:"selected"`
	Debug    bool   `json:"debug"`
}

func NewHandlerSynchReq(debugSQL bool) (*HandlerPrjReq, error) {
	st2 := checker.NewStore(idl.OTPSourceFile)
	st1 := checker.NewStore(idl.OTPServerFile)
	res := HandlerPrjReq{
		storeSourceFile: &st2,
		storeSeverFile:  &st1,
		liteDB: sqlite.LiteDB{
			DebugSQL:     debugSQL,
			SqliteDBPath: conf.Current.SqliteDBPath,
		},
		DebugSQL: debugSQL,
	}

	err := res.liteDB.OpenSqliteDatabase()
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (h *HandlerPrjReq) Close() {
	log.Println("Handler close")
}

var handlePrjReq *HandlerPrjReq

func getURLForRoute(uri string) string {
	arr := strings.Split(uri, "/")
	//fmt.Println("split: ", arr, len(arr))
	for i := len(arr) - 1; i >= 0; i-- {
		ss := arr[i]
		if ss != "" {
			if !strings.HasPrefix(ss, "?") {
				//fmt.Printf("Url for route is %s\n", ss)
				return ss
			}
		}
	}
	return uri
}

func handlePostOperations(w http.ResponseWriter, req *http.Request) error {
	start := time.Now()
	if conf.Current.DebugVerbose {
		log.Println("POST index", req.RequestURI)
	}
	var err error
	if handlePrjReq == nil {
		handlePrjReq, err = NewHandlerSynchReq(conf.Current.DebugSQL)
		if err != nil {
			return err
		}
	}

	lastPath := getURLForRoute(req.RequestURI)
	log.Println("Check the last path ", lastPath)
	switch lastPath {
	case "CallSync":
		err = handlePrjReq.HandleCallSync(w, req)
	case "ViewDiff":
		err = handlePrjReq.HandleViewCurrentDiff(w, req)
	case "OpenVsCode":
		err = handlePrjReq.handleOpenVsCode(w, req)
	case "ExportToFile":
		err = handlePrjReq.handleExportToFile(w, req)
	case "ImportSelectionToNav":
		err = handlePrjReq.handleImportSelectionToServer(w, req)
	case "CompareDiff":
		err = handlePrjReq.handleCompareDiff(w, req)
	case "NewFile":
		err = handlePrjReq.handleNewFile(w, req)
	default:
		return fmt.Errorf("%s Not supported", lastPath)
	}
	if err != nil {
		return err
	}

	t := time.Now()
	elapsed := t.Sub(start)
	log.Printf("Post handler total call duration: %v\n", elapsed)
	return nil
}

func writeJsonResp(w http.ResponseWriter, resp interface{}) error {
	blobresp, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(blobresp)

	return nil
}

func parseBodyReq(req *http.Request, paraDef interface{}) error {
	rawbody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	debug := conf.Current.DebugVerbose
	if v, ok := paraDef.(DoWithSingleSelectionReq); ok {
		debug = debug || v.Debug
	}
	if v, ok := paraDef.(DoWithSelectionReq); ok {
		debug = debug || v.Debug
	}
	if debug {
		log.Println("Body is: ", string(rawbody))
	}

	return json.Unmarshal(rawbody, &paraDef)
}

func errorResHandler(itemRes idl.ResErr, ok bool) {
	if ok {
		if itemRes.Err != nil {
			log.Println("--+- Error --+- result ", itemRes.Err)
		}
	}
}
