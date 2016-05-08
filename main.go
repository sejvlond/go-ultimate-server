package main

import (
	"os"
	"sync"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/docopt/docopt-go"
)

type LOGGER logrus.FieldLogger

type App struct {
	cfg          *Config
	lgr          LOGGER
	workerWG     sync.WaitGroup
	shutdownChan chan struct{}
	signalWG     sync.WaitGroup
	signalChan   chan os.Signal
}

func (a *App) StartSignalHandler() error {
	dispatch := map[os.Signal]SignalHandlerFunc{
		os.Interrupt:    a.ShutDown,
		os.Kill:         a.ShutDown,
		syscall.SIGTERM: a.ShutDown,
	}
	signal, err := NewSignalHandler(a.lgr.WithField("name", "SIGNAL HANDLER"),
		a.signalChan, &a.signalWG, dispatch)
	if err != nil {
		return err
	}
	a.signalWG.Add(1)
	go signal.Run()
	return nil
}

func (a *App) StartServer() error {
	server, err := NewServer(a.lgr.WithField("name", "SERVER"),
		&a.cfg.Server, a.ShutDown, a.shutdownChan, &a.workerWG)
	if err != nil {
		return err
	}
	a.workerWG.Add(1)
	go server.Run()
	return nil
}

func (a *App) Start() {
	a.shutdownChan = make(chan struct{})
	a.signalChan = make(chan os.Signal)
	if err := a.StartSignalHandler(); err != nil {
		a.lgr.Errorf("Signal handler initialization error: %q", err)
		goto shutdown
	}
	if err := a.StartServer(); err != nil {
		a.lgr.Errorf("Server initialization error: %q", err)
		goto shutdown
	}
	goto stopping // validly here - skip shutdown
shutdown:
	a.ShutDown()
stopping:
	a.workerWG.Wait()
	return
}

func (a *App) ShutDown() {
	a.lgr.Infof("Shutdown initialized")
	close(a.shutdownChan)
}

func (a *App) Stop() {
	close(a.signalChan)
	a.signalWG.Wait()
	a.lgr.Infof("Bye")
}

func main() {
	lgr := logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}

	usage := `server

Usage:
    server -c <config_file>
    server -h | --help

Options:
    -c --config         configuration file
    -h --help           Show this screen.`

	args, err := docopt.Parse(usage, nil, true, "", false)
	if err != nil {
		lgr.Fatalf("Error parsing arguments %q", err)
	}

	var cfg *Config
	if useCfg := args["--config"].(bool); useCfg {
		configFile, ok := args["<config_file>"].(string)
		if !ok {
			lgr.Fatalf("Error config argument %q", err)
		}
		cfg, err = NewConfig(configFile)
		if err != nil {
			lgr.Fatalf("Error loading config file %q", err)
		}
	} else {
		lgr.Fatalf("Config was not loaded")
	}

	app := App{
		lgr: lgr.WithField("name", "MAIN"),
		cfg: cfg,
	}
	app.Start()
	app.Stop()
	return
}
