<h1 align="center">Diablo II Ressurrected Auto Potion</h1>

---

Tool for Diablo 2: Resurrected written in GoLang. 
<br />
Automatically uses HP and Mana Potion
<br />
Parameter configuration on config/config.yaml
<br />

[Windows .zip](https://github.com/Hefero/D2R-AutoPotion-Go/releases/download/v1/D2R-AutoPotion-Go.zip) - Binary compiled for Windows (just run main.exe)

Or with source using Golang
- Download GoLang Install [Download Golang](https://go.dev/doc/install)
- To execute in Terminal from source:
```ruby
$ go run main.go
```
- To compile main.exe:
```ruby
$ go build main.go
```

### Libraries

- [data](https://github.com/Hefero/D2R-AutoPotion-Go/tree/main/pkg/data) - D2R Game data structures
- [memory](https://github.com/Hefero/D2R-AutoPotion-Go/tree/main/pkg/memory) - D2R memory reader (it provides the data
  structures)

### Tools

- [https://github.com/hectorgimenez/d2go/cmd/itemwatcher](https://github.com/hectorgimenez/d2go/tree/main/cmd/itemwatcher) - Small tool that plays a sound
  when an item passing the filtering process is dropped
