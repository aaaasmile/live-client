package srv

import (
	"fmt"
	"log"
	"net/http"
)

var ()

func HandleIndex(w http.ResponseWriter, req *http.Request) {
	var err error
	switch req.Method {
	case "GET":
		handleIndexGet(w, req)
	case "POST":
		err = handlePostOperations(w, req)
	}
	if err != nil {
		log.Println("Handler error: ", err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
	}
}
