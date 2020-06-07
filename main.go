package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"webchat-golang-profile/internal/config"
	"webchat-golang-profile/internal/logger"
	"webchat-golang-profile/internal/rediscache"
	"webchat-golang-profile/internal/server"
)

var (
	serviceList []*server.BaseService
	conf        *config.AppConfig
	wg          sync.WaitGroup
	log         *logger.Logger
)

func init() {
	log = logger.NewLogger("main")
}

//ExitSuccess is exit code 0 and ExitFailure is exit code 1
const (
	ExitSuccess = iota
	ExitFailure
)

func main() {
	//	log = logger.NewLogger("Main")
	log.I("Starting service ...")
	err := Initialize()
	if err != nil {
		log.E("Failed to initialize service: %v", err)
		os.Exit(1)
	}

	err = StartService()
	if err != nil {
		log.E("Failed to start service: %v", err)
		os.Exit(1)
	}
	log.E("Exiting service ...")
}

// Initialize is to setup the basic
func Initialize() error {
	log.I("initiate the service...")

	// Configuration loading
	var configFileName string = "configs/config.json"

	conf = config.GetInstance()
	if !config.Load(configFileName) {
		err := fmt.Errorf("Failed to load config file: %s", configFileName)
		return err
	}
	log.D("Configuration has been loaded.")

	// Setup log level
	// log.SetupLogger(conf.Logging.Enable, conf.Logging.Level)
	logger.SetupLogger(conf.Logging.Enable, conf.Logging.Level)

	// Setup signal handlers for interruption and termination
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range sigCh {
			if sig == syscall.SIGINT || sig == syscall.SIGTERM {
				log.D("Graceful Termination Time = %d", conf.GracefulTermTimeMillis)
				Finalize()
				time.Sleep(time.Duration(conf.GracefulTermTimeMillis) * time.Millisecond)
				os.Exit(ExitFailure)
			}
		}
	}()

	return nil
}

// StartService starts all the component of this service.
func StartService() error {
	log.I("start the service...")

	conf = config.GetInstance()

	// if there are more services, those can be appened here
	serviceList = append(serviceList, server.NewBaseService(&server.ProfileService{}, &wg, conf))

	for _, service := range serviceList {
		go service.Run()
	}

	wg.Wait()

	return nil
}

// Finalize cleans up this service including wrapping up current
// DB transaction and closing open DB connection before shutting
// down this service
func Finalize() {
	//	db.Close()
	rediscache.Close()
	for _, service := range serviceList {
		service.Stop()
	}
	log.E("Shutdown service...")
}
