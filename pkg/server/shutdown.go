package textserver

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Shutdown holds initiate and done channels and graceful timeout variable
type Shutdown struct {
	initiateChannel chan struct{}
	doneChannel     chan struct{}
	exitChannel     chan struct{}
	timeout         time.Duration
}

// NewShutdown is a constructor for Shutdown object
func NewShutdown(timeout time.Duration) Shutdown {
	return Shutdown{
		initiateChannel: make(chan struct{}),
		exitChannel:     make(chan struct{}, 1),
		doneChannel:     make(chan struct{}),
		timeout:         timeout,
	}
}

// InitiateChannel blocks until a shutdown sequence is initiated
func (s *Shutdown) InitiateChannel() <-chan struct{} {
	return s.initiateChannel
}

// DoneChannel signals that the graceful shutdown has been completed and we can exit
func (s *Shutdown) DoneChannel() chan<- struct{} {
	return s.doneChannel
}

// GracefullyExit initiates a manual exit
func (s *Shutdown) GracefullyExit() {
	s.exitChannel <- struct{}{}
}

// WaitForSignal listen for a system termination
func (s *Shutdown) WaitForSignal() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-signalChannel:
	case <-s.exitChannel:
	}

	log.Println("received interrupt signal")
	s.initiateChannel <- struct{}{}
	select {
	case <-signalChannel:
		log.Println("forcing shutdown")
		os.Exit(1)
	case <-s.doneChannel:
		log.Println("cleanup done, exiting")
	case <-time.After(s.timeout):
		log.Println("cleanup timed out, exiting")
	}
}
