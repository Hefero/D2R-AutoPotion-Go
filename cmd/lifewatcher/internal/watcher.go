package lifewatcher

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/hectorgimenez/d2go/pkg/data"
	"github.com/hectorgimenez/d2go/pkg/data/area"
	"github.com/hectorgimenez/d2go/pkg/data/item"
	"github.com/hectorgimenez/d2go/pkg/memory"
	"github.com/hectorgimenez/d2go/pkg/nip"
	"github.com/micmonay/keybd_event"
)

type Watcher struct {
	gr                     *memory.GameReader
	rules                  []nip.Rule
	alreadyNotifiedItemIDs []itemFootprint
}

type itemFootprint struct {
	detectedAt time.Time
	area       area.Area
	position   data.Position
	name       item.Name
	quality    item.Quality
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
	rejuvInterval       = time.Second * 2
)

func (fp itemFootprint) Match(area area.Area, i data.Item) bool {
	return fp.area == area && fp.position == i.Position && fp.name == i.Name && fp.quality == i.Quality
}

func NewWatcher(gr *memory.GameReader, rules []nip.Rule) *Watcher {
	return &Watcher{gr: gr, rules: rules}
}

func (w *Watcher) Start(ctx context.Context) error {
	var manager = Manager{}
	w.alreadyNotifiedItemIDs = make([]itemFootprint, 0)
	audioBuffer, err := initAudio()
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			time.Sleep(100 * time.Millisecond)

			d := w.gr.GetData()

			HPPercent := d.PlayerUnit.HPPercent()
			MPPercent := d.PlayerUnit.MPPercent()

			log.Printf("%s: Life: %d Mana: %d", time.Now().Format(time.RFC3339), HPPercent, MPPercent)

			usedRejuv := false
			if time.Since(manager.lastRejuv) > rejuvInterval && (d.PlayerUnit.HPPercent() <= 40 || d.PlayerUnit.MPPercent() < 30) {
				UseRejuv()
				if usedRejuv {
					manager.lastRejuv = time.Now()
				}
				speaker.Play(audioBuffer.Streamer(0, audioBuffer.Len()))
			}

			if !usedRejuv {

				if d.PlayerUnit.HPPercent() <= 85 && time.Since(manager.lastHeal) > healingInterval {
					UseHP()
					manager.lastHeal = time.Now()
					speaker.Play(audioBuffer.Streamer(0, audioBuffer.Len()))
				}

				if d.PlayerUnit.MPPercent() <= 65 && time.Since(manager.lastMana) > manaInterval {
					UseMana()
					manager.lastMana = time.Now()
					speaker.Play(audioBuffer.Streamer(0, audioBuffer.Len()))
				}
			}

			// Mercenary
			if d.MercHPPercent() > 0 {
				usedMercRejuv := false
				if time.Since(manager.lastRejuvMerc) > rejuvInterval && d.MercHPPercent() <= 30 {
					UseMercRejuv()
					if usedMercRejuv {
						manager.lastRejuvMerc = time.Now()
					}
					speaker.Play(audioBuffer.Streamer(0, audioBuffer.Len()))
				}

				if !usedMercRejuv {

					if d.MercHPPercent() <= 60 && time.Since(manager.lastMercHeal) > healingMercInterval {
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
	f, err := os.Open("assets/ching.wav")
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

func UseHP() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
	}
	kb.HasSHIFT(false)
	kb.SetKeys(keybd_event.VK_1)
	err = kb.Launching()
}

func UseMana() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
	}
	kb.HasSHIFT(false)
	kb.SetKeys(keybd_event.VK_4)
	err = kb.Launching()
}

func UseHPMerc() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
	}
	kb.HasSHIFT(true)
	kb.SetKeys(keybd_event.VK_1)
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
