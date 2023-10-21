package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/hectorgimenez/d2go/cmd/config"
	"github.com/hectorgimenez/d2go/cmd/lifewatcher"
	"github.com/hectorgimenez/d2go/pkg/memory"
)

func main() {
	process, err := memory.NewProcess()
	if err != nil {
		log.Fatalf("error starting process: %s", err.Error())
	}

	errL := config.Load()
	if errL != nil {
		log.Fatalf("Error loading configuration: %s", errL.Error())
	}

	gr := memory.NewGameReader(process)

	watcher := lifewatcher.NewWatcher(gr)

	ctx := contextWithSigterm(context.Background())
	err = watcher.Start(ctx)
	if err != nil {
		log.Fatalf("error during process: %s", err.Error())
	}
}

func contextWithSigterm(ctx context.Context) context.Context {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()

		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt)

		select {
		case <-signalCh:
		case <-ctx.Done():
		}
	}()

	return ctxWithCancel
}
