package main

import (
	"encoding/json"
	"fmt"
	"github.com/nzin/dctycoon"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/dctycoon/timer"
	"github.com/nzin/sws"
	"os"
)

func main() {
	quit:=false

	root := sws.Init(800, 600)

	timer.GlobalEventPublisher=timer.CreateEventPublisher(root)
	dctycoon.GlobalLedger=dctycoon.CreateLedger()
	timer.GlobalGameTimer=timer.CreateGameTimer()

	dc := dctycoon.CreateDcWidget(root.Width(), root.Height())
	supplierwidget := dctycoon.CreateSupplier(root)
	gamefile, err := os.Open("example.map")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var v map[string]interface{}
	jsonParser := json.NewDecoder(gamefile)
	if err = jsonParser.Decode(&v); err != nil {
		fmt.Println("parsing game file", err.Error())
		os.Exit(1)
	}
	gamefile.Close()

	// initiate the ledger
	dctycoon.GlobalLedger.Load(v["ledger"].(map[string]interface{}))
	
	// initiate the game timer
	timer.GlobalGameTimer.Load(v["clock"].(map[string]interface{}))

	gamemap := v["map"].(map[string]interface{})
	dc.LoadMap(gamemap)
	root.AddChild(dc)
	root.SetFocus(dc)
	
	// dock 
	dock:=dctycoon.CreateDockWidget(timer.GlobalGameTimer)
	dock.Move(root.Width()-dock.Width(),0)
	root.AddChild(dock)

	supplier.Trends = supplier.TrendLoad(v["trends"].(map[string]interface{}))
	dock.SetShopCallback(func() {
		supplierwidget.Show()
	})

	dock.SetQuitCallback(func() {
		quit=true
	})

	//fmt.Println(supplier.Trends.Cpuprice.CurrentValue(time.Now()))

	for sws.PoolEvent() == false && quit == false {
	}
	data := dc.SaveMap()
	gamefile, err = os.Create("backup.map")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	gamefile.WriteString("{")
	gamefile.WriteString(fmt.Sprintf(`"map": %s,`, data) + "\n")
	gamefile.WriteString(fmt.Sprintf(`"trends": %s,`, supplier.TrendSave(supplier.Trends)) + "\n")
	gamefile.WriteString(fmt.Sprintf(`"clock": %s`, timer.GlobalGameTimer.Save() + "\n"))
	gamefile.WriteString(fmt.Sprintf(`"ledger": %s`, dctycoon.GlobalLedger.Save() + "\n"))
	gamefile.WriteString("}\n")

	gamefile.Close()
}
