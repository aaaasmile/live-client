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
	Buildnr  string
	RootUrl  string
	Database string
	Repo     string
	DbLite   string
}

func handleIndexGet(w http.ResponseWriter, req *http.Request) error {
	u, _ := url.Parse(req.RequestURI)

	pagectx := PageCtx{
		RootUrl: conf.Current.RootURLPattern,
		Buildnr: idl.Buildnr,
		Repo:    conf.Current.RepoSync,
		DbLite:  conf.Current.SqliteDBPath,
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
