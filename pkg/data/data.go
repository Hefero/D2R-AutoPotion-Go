package data

import (
	"strings"

	"github.com/Hefero/D2R-AutoPotion-Go/pkg/data/area"
	"github.com/Hefero/D2R-AutoPotion-Go/pkg/data/skill"
	"github.com/Hefero/D2R-AutoPotion-Go/pkg/data/stat"
	"github.com/Hefero/D2R-AutoPotion-Go/pkg/data/state"
)

// since stat.MaxLife is returning max life without stats, we are setting the max life value that we read from the
// game memory, overwriting this value each time it increases. It's not a good solution but it will provide
// more accurate values for the life %. This value is checked for each memory iteration.

const (
	goldPerLevel = 10000

	// Monster Types
	MonsterTypeNone        MonsterType = "None"
	MonsterTypeChampion    MonsterType = "Champion"
	MonsterTypeMinion      MonsterType = "Minion"
	MonsterTypeUnique      MonsterType = "Unique"
	MonsterTypeSuperUnique MonsterType = "SuperUnique"
)

type Data struct {
	AreaOrigin Position
	Corpse     Corpse
	Monsters   Monsters
	// First slice represents X and second Y
	CollisionGrid  [][]bool
	PlayerUnit     PlayerUnit
	NPCs           NPCs
	Items          Items
	Objects        Objects
	AdjacentLevels []Level
	Rooms          []Room
	OpenMenus      OpenMenus
	Roster         Roster
	HoverData      HoverData
	TerrorZones    []area.Area
	Params_        Params
}

type Params struct {
	MaxLife   int `default:"0"`
	MaxLifeBO int `default:"0"`
	MaxMana   int `default:"0"`
	MaxManaBO int `default:"0"`
}

var Params_ Params

type Room struct {
	Position
	Width  int
	Height int
}

type HoverData struct {
	IsHovered bool
	UnitID
	UnitType int
}

func (r Room) GetCenter() Position {
	return Position{
		X: r.Position.X + r.Width/2,
		Y: r.Position.Y + r.Height/2,
	}
}

func (r Room) IsInside(p Position) bool {
	if p.X >= r.X && p.X <= r.X+r.Width {
		return p.Y >= r.Y && p.Y <= r.Y+r.Height
	}

	return false
}

func (d Data) MercHPPercent() int {
	for _, m := range d.Monsters {
		if m.IsMerc() {
			// Hacky thing to read merc life properly
			maxLife := m.Stats[stat.MaxLife] >> 8
			life := float64(m.Stats[stat.Life] >> 8)
			if m.Stats[stat.Life] <= 32768 {
				life = float64(m.Stats[stat.Life]) / 32768.0 * float64(maxLife)
			}

			return int(life / float64(maxLife) * 100)
		}
	}

	return 0
}

type RosterMember struct {
	Name     string
	Area     area.Area
	Position Position
}
type Roster []RosterMember

func (r Roster) FindByName(name string) (RosterMember, bool) {
	for _, rm := range r {
		if strings.EqualFold(rm.Name, name) {
			return rm, true
		}
	}

	return RosterMember{}, false
}

type Level struct {
	Area       area.Area
	Position   Position
	IsEntrance bool // This means the area can not be accessed just walking through it, needs to be clicked
}

type Class uint

const (
	Amazon Class = iota
	Sorceress
	Necromancer
	Paladin
	Barbarian
	Druid
	Assassin
)

type Corpse struct {
	Found     bool
	IsHovered bool
	Position  Position
}

type Position struct {
	X int
	Y int
}

type PlayerUnit struct {
	Name       string
	ID         UnitID
	Area       area.Area
	Position   Position
	Stats      map[stat.ID]int
	Skills     map[skill.Skill]int
	States     state.States
	Class      Class
	LeftSkill  skill.Skill
	RightSkill skill.Skill
}

func (pu PlayerUnit) MaxGold() int {
	return goldPerLevel * pu.Stats[stat.Level]
}

// TotalGold returns the amount of gold, including inventory and stash
func (pu PlayerUnit) TotalGold() int {
	return pu.Stats[stat.Gold] + pu.Stats[stat.StashGold]
}

func (pu PlayerUnit) HPPercent() int {
	_, found := pu.Stats[stat.MaxLife]
	if !found {
		return 100
	}

	if Params_.MaxLifeBO == 0 && Params_.MaxLife == 0 {
		Params_.MaxLife = pu.Stats[stat.Life]
		Params_.MaxLifeBO = pu.Stats[stat.Life]
	}

	if pu.States.HasState(state.Battleorders) {
		if Params_.MaxLifeBO < pu.Stats[stat.Life] {
			Params_.MaxLifeBO = pu.Stats[stat.Life]
		}
		return int((float64(pu.Stats[stat.Life]) / float64(Params_.MaxLifeBO)) * 100)
	}
	if !pu.States.HasState(state.Battleorders) {
		if Params_.MaxLife < pu.Stats[stat.Life] {
			Params_.MaxLife = pu.Stats[stat.Life]
		}
		return int((float64(pu.Stats[stat.Life]) / float64(Params_.MaxLife)) * 100)
	}

	return int((float64(pu.Stats[stat.Life]) / float64(pu.Stats[stat.Life])) * 100)
}

func (pu PlayerUnit) MPPercent() int {
	_, found := pu.Stats[stat.MaxMana]
	if !found {
		return 100
	}

	if Params_.MaxManaBO == 0 && Params_.MaxMana == 0 {
		Params_.MaxMana = pu.Stats[stat.Mana]
		Params_.MaxManaBO = pu.Stats[stat.Mana]
	}

	if pu.States.HasState(state.Battleorders) {
		if Params_.MaxManaBO < pu.Stats[stat.Mana] {
			Params_.MaxManaBO = pu.Stats[stat.Mana]
		}
		return int((float64(pu.Stats[stat.Mana]) / float64(Params_.MaxManaBO)) * 100)
	}
	if !pu.States.HasState(state.Battleorders) {
		if Params_.MaxMana < pu.Stats[stat.Mana] {
			Params_.MaxMana = pu.Stats[stat.Mana]
		}
		return int((float64(pu.Stats[stat.Mana]) / float64(Params_.MaxMana)) * 100)
	}

	return int((float64(pu.Stats[stat.Mana]) / float64(pu.Stats[stat.MaxMana])) * 100)
}

func (pu PlayerUnit) HasDebuff() bool {
	debuffs := []state.State{
		state.Amplifydamage,
		state.Attract,
		state.Confuse,
		state.Conversion,
		state.Decrepify,
		state.Dimvision,
		state.Ironmaiden,
		state.Lifetap,
		state.Lowerresist,
		state.Terror,
		state.Weaken,
		state.Convicted,
		state.Conviction,
		state.Poison,
		state.Cold,
		state.Slowed,
		state.BloodMana,
		state.DefenseCurse,
	}

	for _, s := range pu.States {
		for _, d := range debuffs {
			if s == d {
				return true
			}
		}
	}

	return false
}

type PointOfInterest struct {
	Name     string
	Position Position
}

type OpenMenus struct {
	Inventory     bool
	LoadingScreen bool
	NPCInteract   bool
	NPCShop       bool
	Stash         bool
	Waypoint      bool
	MapShown      bool
	SkillTree     bool
	Character     bool
	QuitMenu      bool
	Cube          bool
	SkillSelect   bool
	Anvil         bool
}
