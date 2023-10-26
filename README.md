<h1 align="center">Diablo II Ressurrected Auto Potion</h1>

---

Tool for Diablo 2: Resurrected written in GoLang. Based on [d2rgo library](https://github.com/hectorgimenez/d2go/) by [hectorgimenez](https://github.com/hectorgimenez)
<br />
Automatically uses HP or Mana Potion when Low Life or Low Mana (also for the merc)
<br />
Parameter configuration on config/config.yaml
<br />

### Instructions
### - Windows Ready Binary
- [Windows .zip](https://github.com/Hefero/D2R-AutoPotion-Go/releases/download/v1/D2R-AutoPotion-Go.zip) - Binary compiled for Windows (just run main.exe)
### Or
### - Build/Execute from Source
- Download and Install [GoLang](https://go.dev/doc/install)
- Clone repository:
```ruby
$ git clone https://github.com/Hefero/D2R-AutoPotion-Go
```
- Install [Autohotkey](https://www.autohotkey.com/) and compile gui.ahk using Ahk v2
- That can be done via below command line on windows prompt or right click on gui.ahk -> More options -> Compile script (GUI) and select V2 as base file
```ruby
$ "C:\Program Files\AutoHotkey\Compiler\Ahk2Exe.exe" /in gui.ahk /base "C:\Program Files\AutoHotkey\v2\AutoHotkey32.exe"
```
- Compile and execute main.exe:
```ruby
$ cd D2R-AutoPotion-Go
$ go build -ldflags -H=windowsgui main.go
$ main.exe
```

### Libraries

- [data](https://github.com/Hefero/D2R-AutoPotion-Go/tree/main/pkg/data) - D2R Game data structures
- [memory](https://github.com/Hefero/D2R-AutoPotion-Go/tree/main/pkg/memory) - D2R memory reader (it provides the data
  structures)

### Other Tools

- [https://github.com/hectorgimenez/d2go/cmd/itemwatcher](https://github.com/hectorgimenez/d2go/tree/main/cmd/itemwatcher) - Small tool that plays a sound
  when an item passing the filtering process is dropped

  <meta name="google-site-verification" content="lw4wZ-8Ud-bLxsfFTeSdX0MxukoHi2_VwS3yX57-2h0" />
