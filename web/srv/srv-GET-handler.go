package srv

import (
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/aaaasmile/live-client/conf"
	"github.com/aaaasmile/live-client/web/idl"
)

type PageCtx struct {
	Buildnr        string
	RootUrl        string
	ProjectName    string
	Database       string
	Repo           string
	SplitFile      string
	SplittedDir    string
	DbLite         string
	ImportFilename string
}

func handleIndexGet(w http.ResponseWriter, req *http.Request) error {
	u, _ := url.Parse(req.RequestURI)

	pagectx := PageCtx{
		RootUrl: conf.Current.RootURLPattern,
		Buildnr: idl.Buildnr,
	}
	templName := "template/vue/index.html"

	tmplIndex := template.Must(template.New("AppIndex").ParseFiles(templName))

	err := tmplIndex.ExecuteTemplate(w, "base", pagectx)
	if err != nil {
		return err
	}
	log.Println("GET requested ", u)
	return nil
}
