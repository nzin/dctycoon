# Introduction

This game is heavily influence by the unfinished Server Tycoon (http://www.servertycoon.com/).

i.e. I wanted to see the guys behind Server Tycoon create and release their game, but apparently they stopped developping it. So I decided to tentatively create my own. And because I began to learn Go in march 2017, I tried to use Go to develop it.

# Compilation Requirements

- you need go version 1.8.x minimum
- Needs sws 0.3.10 library (https://github.com/nzin/sws) and especially grab sdl2 (see https://github.com/veandco/go-sdl2)
- you should download dependancies with `go get -u ./...`
- and finaly you should compile with `./compile.sh`

## For Linux

Get Ubuntu 17.10, change the resolution for at least 1024x768 and execute:
```
sudo apt-get install git golang-go libsdl2{,-mixer,-image,-ttf,-gfx}-dev
export GOPATH=~/go
go get github.com/nzin/dctycoon # there will be an error message "undefined Asset". Don't pay attention to it.
cd go/src/github.com/nzin/dctycoon
./compile.sh
```

## For Mac

```
brew install sdl2{,_image,_ttf,_mixer} pkg-config go@1.8
export GOPATH=~/go
go get github.com/nzin/dctycoon # there will be an error message "undefined Asset". Don't pay attention to it.
cd go/src/github.com/nzin/dctycoon
./compile.sh
```

# Licences

This software is covered by the GNU GPLv3 (see LICENCE file)

It use the Twitter Bootstrap v4 framework (https://github.com/twbs/bootstrap) covered by the MIT licence

It used the "Font Awesome" (https://www.flaticon.com/packs/font-awesome) covered by Creative Common BY 3.0
