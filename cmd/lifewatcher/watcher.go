package lifewatcher

import (
	"context"
	"os"
	"time"

	"github.com/Hefero/D2R-AutoPotion-Go/cmd/config"
	"github.com/Hefero/D2R-AutoPotion-Go/pkg/memory"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/micmonay/keybd_event"
)

type Watcher struct {
	gr *memory.GameReader
}

type Manager struct {
	lastRejuv     time.Time
	lastRejuvMerc time.Time
	lastHeal      time.Time
	lastMana      time.Time
	lastMercHeal  time.Time
}

const (
	healingInterval     = time.Second * 4
	healingMercInterval = time.Second * 6
	manaInterval        = time.Second * 5
	rejuvInterval       = time.Second * 1
)

func NewWatcher(gr *memory.GameReader) *Watcher {
	return &Watcher{gr: gr}
}

func (w *Watcher) Start(ctx context.Context) error {
	var manager = Manager{}
	audioBuffer, err := initAudio()
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			d := w.gr.GetData()
			usedRejuv := false
			if time.Since(manager.lastRejuv) > rejuvInterval && (d.PlayerUnit.HPPercent() <= config.Config.Health.RejuvPotionAtLife || d.PlayerUnit.MPPercent() < config.Config.Health.RejuvPotionAtMana) {
				UseRejuv()
				usedRejuv := true
				if usedRejuv {
					manager.lastRejuv = time.Now()
				}
				speaker.Play(audioBuffer.Streamer(0, audioBuffer.Len()))
			}

			if !usedRejuv {

				if d.PlayerUnit.HPPercent() <= config.Config.Health.HealingPotionAt && time.Since(manager.lastHeal) > healingInterval {
					UseHP()
					manager.lastHeal = time.Now()
					speaker.Play(audioBuffer.Streamer(0, audioBuffer.Len()))
				}

				if d.PlayerUnit.MPPercent() <= config.Config.Health.ManaPotionAt && time.Since(manager.lastMana) > manaInterval {
					UseMana()
					manager.lastMana = time.Now()
					speaker.Play(audioBuffer.Streamer(0, audioBuffer.Len()))
				}
			}

			// Mercenary
			if d.MercHPPercent() > 0 {
				usedMercRejuv := false
				if time.Since(manager.lastRejuvMerc) > rejuvInterval && d.MercHPPercent() <= config.Config.Health.MercRejuvPotionAt {
					UseMercRejuv()
					usedMercRejuv := true
					if usedMercRejuv {
						manager.lastRejuvMerc = time.Now()
					}
					speaker.Play(audioBuffer.Streamer(0, audioBuffer.Len()))
				}

				if !usedMercRejuv {

					if d.MercHPPercent() <= config.Config.Health.MercHealingPotionAt && time.Since(manager.lastMercHeal) > healingMercInterval {
						UseHPMerc()
						manager.lastMercHeal = time.Now()
						speaker.Play(audioBuffer.Streamer(0, audioBuffer.Len()))
					}
				}
			}

		}
	}
}

func initAudio() (*beep.Buffer, error) {
	f, err := os.Open("cmd/lifewatcher/assets/ching.wav")
	if err != nil {
		return nil, err
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		return nil, err
	}
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()

	return buffer, nil
}

func getKey(key int) int {
	switch key {
	case 1:
		return keybd_event.VK_1
	case 2:
		return keybd_event.VK_2
	case 3:
		return keybd_event.VK_3
	case 4:
		return keybd_event.VK_4
	default:
		return 0
	}
}

func UseHP() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
	}
	kb.HasSHIFT(false)
	kb.SetKeys(getKey(config.Config.Bindings.PotionHP))
	err = kb.Launching()
}

func UseMana() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
	}
	kb.HasSHIFT(false)
	kb.SetKeys(getKey(config.Config.Bindings.PotionMANA))
	err = kb.Launching()
}

func UseHPMerc() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
	}
	kb.HasSHIFT(true)
	kb.SetKeys(getKey(config.Config.Bindings.PotionHP))
	err = kb.Launching()
	kb.HasSHIFT(false)
}

func UseMercRejuv() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
	}
	kb.HasSHIFT(true)
	kb.SetKeys(keybd_event.VK_2)
	err = kb.Launching()
	kb.HasSHIFT(false)
}

func UseRejuv() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
	}
	kb.HasSHIFT(false)
	kb.SetKeys(keybd_event.VK_2)
	err = kb.Launching()
}
