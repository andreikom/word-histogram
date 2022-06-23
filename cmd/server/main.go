package main

import (
	"context"
	"time"

	"word-histrogram/pkg/server"
	"word-histrogram/pkg/storage"
	"word-histrogram/pkg/tasks"
	"word-histrogram/pkg/worker"
)

const (
	nWorkers        = 32
	queueSize       = 10000
	addServerPort   = 1337
	getServerPort   = 1338
	shutdownTimeout = 15 * time.Second
)

func main() {
	pool := worker.New(nWorkers, queueSize)
	defer pool.Wait()
	defer pool.Close()

	server, err := textserver.New(
		addServerPort, getServerPort, pool,
		&tasks.Handler{
			WordStore: storage.NewInMemory(),
		})
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go server.ServeAdd(ctx)
	go server.ServeGet(ctx)

	exit := textserver.NewShutdown(shutdownTimeout)

	go func(cancel context.CancelFunc) {
		<-exit.InitiateChannel()
		cancel()
		exit.DoneChannel() <- struct{}{}
	}(cancel)
	exit.WaitForSignal()
}
