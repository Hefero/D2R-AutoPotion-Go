package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Hefero/D2R-AutoPotion-Go/cmd/config"
	"github.com/Hefero/D2R-AutoPotion-Go/cmd/lifewatcher"
	"github.com/Hefero/D2R-AutoPotion-Go/pkg/memory"
)

func main() {

	errL := config.Load()
	if errL != nil {
		log.Fatalf("Error loading configuration file config.yaml: %s", errL.Error())
	}

	process, err := memory.NewProcess()
	for err != nil {
		if err != nil {
			log.Printf("error starting process: player needs to be inside a running game %s, retrying in 5 seconds", err.Error())
		}
		time.Sleep(5 * time.Second)
		process, err = memory.NewProcess()
	}

	gr := memory.NewGameReader(process)

	watcher := lifewatcher.NewWatcher(gr)

	ctx := contextWithSigterm(context.Background())
	err = watcher.Start(ctx)

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
