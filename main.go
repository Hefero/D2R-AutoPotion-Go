package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Hefero/D2R-AutoPotion-Go/cmd/config"
	"github.com/Hefero/D2R-AutoPotion-Go/cmd/lifewatcher"
	"github.com/Hefero/D2R-AutoPotion-Go/pkg/memory"
	"github.com/faiface/beep"
)

func main() {

	a := app.New()
	w := a.NewWindow("Hefero Diablo 2 Ressurrected AutoPotion")

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

	audioBufferL, err := lifewatcher.InitAudio("cmd/lifewatcher/assets/life.wav")
	audioBufferM, err := lifewatcher.InitAudio("cmd/lifewatcher/assets/mana.wav")
	audioBufferR, err := lifewatcher.InitAudio("cmd/lifewatcher/assets/rejuv.wav")

	gr := memory.NewGameReader(process)

	watcher := lifewatcher.NewWatcher(gr)

	ctx := contextWithSigterm(context.Background())

	hello := widget.NewLabel("Diablo 2 Ressurrected AutoPotion")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Start", func() {
			path, err := os.Getwd()
			if err != nil {
				log.Println(err)
			}
			cmd := exec.Command(path + "\\gui.exe")
			if err := cmd.Run(); err != nil {
				fmt.Println(err)
			}
		}),
	))
	go func() {
		for {
			watcher, err = StartWatcher(*watcher, ctx, &manager, &XP, audioBufferL, audioBufferM, audioBufferR)
			//if err != nil {
			//	hello.SetText(err.Error())
			//}
			if err == nil {
				if !XP.FirstStart {
					//duration.Round(time.Second).String()
					//var textLabel = strconv.FormatFloat(XP.Hours, 'f', 2, 64)
					//updateLabel(hello, "Running")
				}
				//updateLabel(hello, "test")
			}
		}
	}()

	w.ShowAndRun()

}

func updateLabel(label *widget.Label, text string) {
	label.SetText(text)
}

func StartWatcher(watcher lifewatcher.Watcher, ctx context.Context, manager *lifewatcher.Manager, XP *lifewatcher.ExperienceCalc, audioBufferL *beep.Buffer, audioBufferM *beep.Buffer, audioBufferR *beep.Buffer) (*lifewatcher.Watcher, error) {
	err := watcher.Start(ctx, manager, XP, audioBufferL, audioBufferM, audioBufferR)
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
