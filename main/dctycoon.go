package main

import (
	"encoding/json"
	"fmt"
	"github.com/nzin/dctycoon"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/dctycoon/timer"
	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/sws"
	"os"
)

func main() {
	quit:=false

	root := sws.Init(800, 600)
	dctycoon.GlobalLocation="siliconvalley"

	timer.GlobalEventPublisher=timer.CreateEventPublisher(root)
	accounting.GlobalLedger=accounting.CreateLedger(dctycoon.AvailableLocation[dctycoon.GlobalLocation].Taxrate)
	timer.GlobalGameTimer=timer.CreateGameTimer()

	dc := dctycoon.CreateDcWidget(root.Width(), root.Height())
	supplierwidget := dctycoon.CreateSupplier(root)
	accountingui := accounting.CreateAccounting(root)
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

	// initiate the location
	dctycoon.GlobalLocation=v["location"].(string)
	
	// initiate the ledger
	accounting.GlobalLedger.Load(v["ledger"].(map[string]interface{}),dctycoon.AvailableLocation[dctycoon.GlobalLocation].Taxrate)
	accountingui.SetBankinterestrate(dctycoon.AvailableLocation[dctycoon.GlobalLocation].Bankinterestrate)
	
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
	
	dock.SetLedgerCallback(func() {
		accountingui.Show()
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
	gamefile.WriteString(fmt.Sprintf(`"location": "%s",`, dctycoon.GlobalLocation) + "\n")
	gamefile.WriteString(fmt.Sprintf(`"map": %s,`, data) + "\n")
	gamefile.WriteString(fmt.Sprintf(`"trends": %s,`, supplier.TrendSave(supplier.Trends)) + "\n")
	gamefile.WriteString(fmt.Sprintf(`"clock": %s`, timer.GlobalGameTimer.Save() + "\n"))
	gamefile.WriteString(fmt.Sprintf(`"ledger": %s`, accounting.GlobalLedger.Save() + "\n"))
	gamefile.WriteString("}\n")

	gamefile.Close()
}
