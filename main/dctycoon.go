package main

import (
	"flag"
	"os"

	"github.com/nzin/dctycoon"
	"github.com/nzin/sws"
	log "github.com/sirupsen/logrus"
)

func initLog(loglevel, filename string) {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	if filename != "" {
		f, err := os.Open(filename)
		if err != nil {
			log.SetOutput(os.Stdout)
		} else {
			log.SetOutput(f)
		}
	} else {
		log.SetOutput(os.Stdout)
	}

	log.SetLevel(log.ErrorLevel)
	switch loglevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	}

}

func main() {
	quit := false

	loglevel := flag.String("loglevel", "", "[debug,info,warning,error] Default to error")
	logfile := flag.String("logfile", "", "optional if we want the log to not be on stdout")
	debug := flag.Bool("debug", false, "Development only")
	flag.Parse()
	initLog(*loglevel, *logfile)

	root := sws.Init(800, 600)

	game := dctycoon.NewGame(&quit, root, *debug)
	game.ShowOpening()
	/*
		//	game.InitGame(12000, "siliconvalley")
		game.LoadGame("example.map")

		for sws.PoolEvent() == false && quit == false {
		}

		game.SaveGame("backup.map")

		game.LoadGame("backup.map")
		quit = false */
	for sws.PoolEvent() == false && quit == false {
	}
}
