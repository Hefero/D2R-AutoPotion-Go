<h1 align="center">Diablo II Ressurrected Auto Potion</h1>

---

Tool for Diablo 2: Resurrected written in Go. 
<br />
Automatically uses HP and Mana Potion
<br />
Parameter configuration on config/config.yaml
<br />

[Windows .zip](https://github.com/Hefero/D2R-AutoPotion-Go/releases/download/v1/D2R-AutoPotion-Go.zip) - Binary compiled for Windows

### Libraries

- [data](https://github.com/Hefero/D2R-AutoPotion-Go/tree/main/pkg/data) - D2R Game data structures
- [memory](https://github.com/Hefero/D2R-AutoPotion-Go/tree/main/pkg/memory) - D2R memory reader (it provides the data
  structures)
- [nip](https://github.com/Hefero/D2R-AutoPotion-Go/tree/main/pkg/nip) - A very basic NIP file parser
- [itemfilter](https://github.com/Hefero/D2R-AutoPotion-Go/tree/main/pkg/itemfilter) - Based on game data, it provides an item
  pickup filtering

### Tools

- [cmd/itemwatcher](https://github.com/hectorgimenez/d2go/cmd/itemwatcher) - Small tool that plays a sound
  when an item passing the filtering process is dropped
