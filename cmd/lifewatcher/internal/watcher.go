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

func (fp itemFootprint) Match(area area.Area, i data.Item) bool {
	return fp.area == area && fp.position == i.Position && fp.name == i.Name && fp.quality == i.Quality
}

func NewWatcher(gr *memory.GameReader, rules []nip.Rule) *Watcher {
	return &Watcher{gr: gr, rules: rules}
}

func (w *Watcher) Start(ctx context.Context) error {
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

			log.Printf("%s: Life: %d Mana: %d", time.Now().Format(time.RFC3339), d.PlayerUnit.HPPercent(), d.PlayerUnit.MPPercent())

			speaker.Play(audioBuffer.Streamer(0, audioBuffer.Len()))
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
