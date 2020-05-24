package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aaaasmile/live-client/web/idl"
	"github.com/aaaasmile/live-client/winser"
)

func main() {
	var configfile = flag.String("config", "config.toml", "Configuration file path")
	var ver = flag.Bool("version", false, "Prints current version")
	var auto = flag.Bool("auto", false, "Opent a new browser window with the dashboard")
	var cmd = flag.String("cmd", "", `
	where the string is one of this service control commands: 
	start, stop, restart, install, uninstall or like.
	Commands are used to manage the windows service. After the install, usually, you have to configure the "Log On" user before starting the service.
	Only the command "like" starts the application in console but it is different as an empty command.
	An empty command is used to start the application without the windows service stuff.`)
	var serviceName = flag.String("servicename", idl.WebServiceName, fmt.Sprintf("Set the Windows service install name (default %s)", idl.WebServiceName))
	flag.Parse()

	if *ver {
		fmt.Printf("%s, version: %s", idl.Appname, idl.Buildnr)
		os.Exit(0)
	}

	log.Printf("** Start the program: %s vers: %s **\n", os.Args[0], idl.Buildnr)
	servHand := winser.ServiceHandler{
		Settings: idl.ServiceHandlerSettings{
			Command:         *cmd,
			ConfigFile:      *configfile,
			ServiceName:     *serviceName,
			AutoStartPage:   *auto,
		},
	}
	if err := servHand.HandleService(); err != nil {
		log.Println("Error in HandleService", err)
	}

	return
}
