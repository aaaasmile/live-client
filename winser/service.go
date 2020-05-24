package winser

import (
	"fmt"
	"log"
	"os"

	"github.com/kardianos/service"
	"github.com/aaaasmile/live-client/util"
	"github.com/aaaasmile/live-client/web"
	"github.com/aaaasmile/live-client/web/idl"
)

var logger service.Logger

type ServiceHandler struct {
	Settings idl.ServiceHandlerSettings
}

func (sh *ServiceHandler) HandleService() error {
	set := sh.Settings
	cmd := set.Command
	var err error
	if !service.Interactive() {
		// started from windows service, this log file need to be writable from user logon service
		var f *os.File

		fmt.Println("Some output will be redirected to a file log")
		f, err = os.OpenFile(util.GetUserLogFile(set.ServiceName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()

		log.SetOutput(f)
		util.UseRelativeRoot = false
		log.Println("->Service is managed, start file logging. App ver nr: ", idl.Buildnr)
	}
	if cmd == "" && service.Interactive() {
		log.Println("Start the service directly as console process")
		web.RunService(nil, nil, &set)
	} else {
		if err := sh.handleAsManagedService(web.RunService); err == nil {
			log.Printf("Command '%s' executed with success.", cmd)
		}
	}
	return err
}

func (sh *ServiceHandler) handleAsManagedService(runner RunServiceFn) error {
	cmd := sh.Settings.Command
	serviceName := sh.Settings.ServiceName
	svcConfig := &service.Config{
		Name:        fmt.Sprintf("%sService", serviceName),
		DisplayName: fmt.Sprintf("%s Web Service", serviceName),
		Description: fmt.Sprintf("This is the %s Web Service", serviceName),
	}

	prg := &program{
		Settings: &sh.Settings,
		FnRun:    runner,
	}
	kardService, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	chErrs := make(chan error, 5)
	logger, err = kardService.Logger(chErrs)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Program called with %v num %d, cmd is %s\n", os.Args[0], len(os.Args), cmd)

	go func() { // Error writer in background
		for {
			err := <-chErrs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(cmd) != 0 && cmd != "like" {
		err := service.Control(kardService, cmd)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return err
	}
	log.Println("Run the service using Run()")
	err = kardService.Run()
	if err != nil {
		logger.Error(err)
		return err
	}

	log.Println("Finally HandleAsManagedService is terminated")
	return nil
}
