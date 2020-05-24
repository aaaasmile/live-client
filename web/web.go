package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"

	"time"

	"github.com/aaaasmile/live-client/conf"
	"github.com/aaaasmile/live-client/util"
	"github.com/aaaasmile/live-client/web/idl"
	"github.com/aaaasmile/live-client/web/srv"

	"github.com/kardianos/service"
)

func RunService(chShutdown <-chan struct{}, logger service.Logger, settings *idl.ServiceHandlerSettings) {
	configfile := settings.ConfigFile
	if logger == nil {
		logger = service.ConsoleLogger
	}

	conf.ReadConfig(util.GetFullPath(configfile))
	log.Println("Configuration is read - Prepare service init")

	var wait time.Duration
	serverurl := conf.Current.ServiceURL
	protoHtt := "http"

	idl.LocalServiceURL = fmt.Sprintf("%s://%s%s", protoHtt, strings.Replace(serverurl, "0.0.0.0", "localhost", 1), conf.Current.RootURLPattern)
	idl.LocalServiceURL = strings.Replace(idl.LocalServiceURL, "127.0.0.1", "localhost", 1)
	logger.Infof("Server started with URL %s", serverurl)
	log.Println("Try this url: ", idl.LocalServiceURL)

	http.Handle(conf.Current.RootURLPattern+"static/", http.StripPrefix(conf.Current.RootURLPattern+"static", http.FileServer(http.Dir(util.GetFullPath("static")))))
	http.HandleFunc(conf.Current.RootURLPattern, srv.HandleIndex)

	srv := &http.Server{
		Addr: serverurl,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 60,
		IdleTimeout:  time.Second * 60,
		Handler:      nil,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println("Server is not listening anymore: ", err)
		}
	}()

	if settings.AutoStartPage {
		go openBrowser(idl.LocalServiceURL)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt) //We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	log.Println("Enter in server loop")
loop:
	for {
		select {
		case <-sig:
			log.Println("stop because interrupt")
			break loop
		case <-chShutdown:
			log.Println("stop because service shutdown")
			break loop
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Bye, service")
}

func openBrowser(url string) error {
	var cmd string
	var args []string
	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "start"}
	} else {
		log.Fatal("OS not recognized")
		return fmt.Errorf("OS not supported %s", runtime.GOOS)
	}
	args = append(args, url)
	log.Println("open a browser url ", url)
	return exec.Command(cmd, args...).Start()

}
