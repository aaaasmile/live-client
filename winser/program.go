package winser

import (
	"github.com/kardianos/service"
	"github.com/aaaasmile/live-client/web/idl"
)

type RunServiceFn func(<-chan struct{}, service.Logger, *idl.ServiceHandlerSettings)

type program struct {
	exit     chan struct{}
	Settings *idl.ServiceHandlerSettings
	FnRun    RunServiceFn
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() error {
	// Non blocking routine could be used for all process that are blocking the service start routine
	p.FnRun(p.exit, logger, p.Settings)
	return nil
}

func (p *program) Stop(s service.Service) error {
	logger.Info("Stop the service")
	close(p.exit)
	return nil
}
