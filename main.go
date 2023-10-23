package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Hefero/D2R-AutoPotion-Go/cmd/config"
	"github.com/Hefero/D2R-AutoPotion-Go/cmd/lifewatcher"
	"github.com/Hefero/D2R-AutoPotion-Go/pkg/memory"
)

func main() {

	a := app.New()
	w := a.NewWindow("Hello")

	var manager = lifewatcher.Manager{}
	manager.Timer = time.Now()

	var XP = lifewatcher.ExperienceCalc{}

	errL := config.Load()
	if errL != nil {
		log.Fatalf("Error loading configuration file config.yaml: %s", errL.Error())
	}

	process, err := memory.NewProcess()
	for err != nil {
		if err != nil {
			fmt.Printf("error starting process: player needs to be inside a running game %s, retrying in 5 seconds\n", err.Error())
			fmt.Print("\033[A")
		}
		time.Sleep(5 * time.Second)
		process, err = memory.NewProcess()
	}

	gr := memory.NewGameReader(process)

	watcher := lifewatcher.NewWatcher(gr)

	ctx := contextWithSigterm(context.Background())

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
		}),
	))
	go func() {

		for {
			watcher, err = StartWatcher(*watcher, ctx, &manager, &XP)
			//if err != nil {
			//	hello.SetText(err.Error())
			//}
			if err == nil {
				if !XP.FirstStart {
					//duration.Round(time.Second).String()
					var textLabel = strconv.FormatFloat(XP.Hours, 'E', -1, 32) + ":" + strconv.FormatFloat(XP.Minutes, 'E', -1, 32) + "h"
					updateLabel(hello, textLabel)
				}
				//updateLabel(hello, "test")
			}
		}
	}()

	w.ShowAndRun()

}

func updateLabel(label *widget.Label, text string) {
	formatted := time.Now().Format(text)
	label.SetText(formatted)
}

func StartWatcher(watcher lifewatcher.Watcher, ctx context.Context, manager *lifewatcher.Manager, XP *lifewatcher.ExperienceCalc) (*lifewatcher.Watcher, error) {
	err := watcher.Start(ctx, manager, XP)
	return &watcher, err
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
