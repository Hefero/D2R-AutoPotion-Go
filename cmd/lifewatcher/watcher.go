package lifewatcher

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Hefero/D2R-AutoPotion-Go/cmd/config"
	"github.com/Hefero/D2R-AutoPotion-Go/pkg/data"
	"github.com/Hefero/D2R-AutoPotion-Go/pkg/data/stat"
	"github.com/Hefero/D2R-AutoPotion-Go/pkg/data/state"
	"github.com/Hefero/D2R-AutoPotion-Go/pkg/memory"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/micmonay/keybd_event"
)

type Watcher struct {
	gr *memory.GameReader
	pr *data.Params
}

type Manager struct {
	lastRejuv     time.Time
	lastRejuvMerc time.Time
	lastHeal      time.Time
	lastMana      time.Time
	lastMercHeal  time.Time
	lastDebugMsg  time.Time
}

func NewWatcher(gr *memory.GameReader) *Watcher {
	return &Watcher{gr: gr}
}

func (w *Watcher) Start(ctx context.Context) error {
	var manager = Manager{}
	audioBufferL, err := initAudio("cmd/lifewatcher/assets/life.wav")
	audioBufferM, err := initAudio("cmd/lifewatcher/assets/mana.wav")
	audioBufferR, err := initAudio("cmd/lifewatcher/assets/rejuv.wav")

	XP := [10]int{}
	XP_aux := [10]int{}
	XParray := [10]int{}
	XPbefore := 0

	timer := time.Now()
	var first30s = true

	d, err := w.gr.GetData()
	if err != nil {
		fmt.Printf("\r                                              ") //clean line
		fmt.Printf("\rnot In Game\n")
		fmt.Print("\033[A")
		time.Sleep(1 * time.Second)
	}
	if err == nil {
		XPbefore = d.PlayerUnit.Stats[stat.Experience]
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			d, err = w.gr.GetData()
			if err != nil {
				fmt.Printf("\r                                  ") //clean line
				fmt.Printf("\rnot In Game\n")
				fmt.Print("\033[A")
				time.Sleep(1 * time.Second)
			}

			if err == nil {
				if time.Since(manager.lastDebugMsg) > (time.Second * 2) {
					if XPbefore == 0 {
						XPbefore = d.PlayerUnit.Stats[stat.Experience]
					}
					fmt.Printf("\r%2.0f PercentLife:%*d PercentMana:%*d", time.Since(timer).Seconds(), 3, d.PlayerUnit.HPPercent(), 3, d.PlayerUnit.MPPercent())
					manager.lastDebugMsg = time.Now()

					if time.Since(timer) > (time.Second * 15) {
						timer = time.Now()
						if first30s {
							for i := 0; i < len(XP); i++ {
								XP[i] = d.PlayerUnit.Stats[stat.Experience] - XPbefore
							}
							first30s = false
						}
						diff := d.PlayerUnit.Stats[stat.Experience] - XPbefore

						XP_aux = [10]int{diff, XP[0], XP[1], XP[2], XP[3], XP[4], XP[5], XP[6], XP[7], XP[8]}
						XP = XP_aux

						for i := 0; i < len(XParray); i++ {
							XParray[i] = 0
							for j := 0; j < i; j++ {
								XParray[i] += XP[j] / i
							}
							//if (i % 1) > 0 {
							fmt.Printf(" xp_%d:%*d", i, 8, XParray[i]*4)
							//}
						}
						XPbefore = d.PlayerUnit.Stats[stat.Experience]
					}
					fmt.Print("\n\033[A")
				}

				if !d.PlayerUnit.Area.IsTown() {

					var healingInterval float32 = config.Config.Timings.HealingInterval

					if d.PlayerUnit.States.HasState(state.Poison) {
						healingInterval += 2
					}

					usedRejuv := false
					if time.Since(manager.lastRejuv) > (time.Duration(config.Config.Timings.RejuvInterval)*time.Second) && (d.PlayerUnit.HPPercent() <= config.Config.Health.RejuvPotionAtLife || d.PlayerUnit.MPPercent() < config.Config.Health.RejuvPotionAtMana) {
						UseRejuv()
						usedRejuv := true
						if usedRejuv {
							manager.lastRejuv = time.Now()
						}
						speaker.Play(audioBufferR.Streamer(0, audioBufferR.Len()))
					}

					if !usedRejuv {

						if d.PlayerUnit.HPPercent() <= config.Config.Health.HealingPotionAt && time.Since(manager.lastHeal) > (time.Duration(healingInterval)*time.Second) {
							UseHP()
							manager.lastHeal = time.Now()
							speaker.Play(audioBufferL.Streamer(0, audioBufferL.Len()))
						}

						if d.PlayerUnit.MPPercent() <= config.Config.Health.ManaPotionAt && time.Since(manager.lastMana) > (time.Duration(config.Config.Timings.ManaInterval)*time.Second) {
							UseMana()
							manager.lastMana = time.Now()
							speaker.Play(audioBufferM.Streamer(0, audioBufferM.Len()))
						}
					}

					// Mercenary
					if d.MercHPPercent() > 0 {
						usedMercRejuv := false
						if time.Since(manager.lastRejuvMerc) > (time.Duration(config.Config.Timings.RejuvInterval)*time.Second) && d.MercHPPercent() <= config.Config.Health.MercRejuvPotionAt {
							UseMercRejuv()
							usedMercRejuv := true
							if usedMercRejuv {
								manager.lastRejuvMerc = time.Now()
							}
							speaker.Play(audioBufferR.Streamer(0, audioBufferR.Len()))
						}

						if !usedMercRejuv {

							if d.MercHPPercent() <= config.Config.Health.MercHealingPotionAt && time.Since(manager.lastMercHeal) > (time.Duration(config.Config.Timings.HealingMercInterval)*time.Second) {
								UseHPMerc()
								manager.lastMercHeal = time.Now()
								speaker.Play(audioBufferL.Streamer(0, audioBufferL.Len()))
							}
						}
					}
				}
			}

		}
	}
}

func initAudio(path string) (*beep.Buffer, error) {
	f, err := os.Open(path)
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
	kb.SetKeys(getKey(config.Config.Bindings.PotionREJUV))
	err = kb.Launching()
	kb.HasSHIFT(false)
}

func UseRejuv() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
	}
	kb.HasSHIFT(false)
	kb.SetKeys(getKey(config.Config.Bindings.PotionREJUV))
	err = kb.Launching()
}
