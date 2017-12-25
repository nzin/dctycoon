package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nzin/dctycoon"
	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/dctycoon/timer"
	"github.com/nzin/sws"
)

func main() {
	quit := false

	root := sws.Init(800, 600)
	dctycoon.GlobalLocation = "siliconvalley"

	timer.GlobalEventPublisher = timer.NewEventPublisher(root)
	timer.GlobalGameTimer = timer.NewGameTimer()
	accounting.GlobalLedger = accounting.NewLedger(timer.GlobalGameTimer, dctycoon.AvailableLocation[dctycoon.GlobalLocation].Taxrate, dctycoon.AvailableLocation[dctycoon.GlobalLocation].Bankinterestrate)
	supplier.GlobalInventory = supplier.NewInventory()

	dc := dctycoon.NewDcWidget(root.Width(), root.Height(), root, supplier.GlobalInventory)
	supplierwidget := dctycoon.NewMainSupplierWidget(root)
	inventorywidget := dctycoon.NewMainInventoryWidget(root)
	accountingui := accounting.NewMainAccountingWidget(root)
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
	dctycoon.GlobalLocation = v["location"].(string)

	// initiate the game timer
	timer.GlobalGameTimer.Load(v["clock"].(map[string]interface{}))

	// initiate the ledger
	accounting.GlobalLedger.Load(v["ledger"].(map[string]interface{}), dctycoon.AvailableLocation[dctycoon.GlobalLocation].Taxrate, dctycoon.AvailableLocation[dctycoon.GlobalLocation].Bankinterestrate)
	//accountingui.SetBankinterestrate(dctycoon.AvailableLocation[dctycoon.GlobalLocation].Bankinterestrate)

	gamemap := v["map"].(map[string]interface{})
	dc.LoadMap(gamemap)

	supplier.GlobalInventory.Load(v["inventory"].(map[string]interface{}))

	root.AddChild(dc)
	root.SetFocus(dc)

	// dock
	dock := dctycoon.NewDockWidget(timer.GlobalGameTimer)
	dock.Move(root.Width()-dock.Width(), 0)
	root.AddChild(dock)

	supplier.Trends = supplier.TrendLoad(v["trends"].(map[string]interface{}))
	dock.SetShopCallback(func() {
		supplierwidget.Show()
	})

	dock.SetQuitCallback(func() {
		quit = true
	})

	dock.SetLedgerCallback(func() {
		accountingui.Show()
	})

	dock.SetInventoryCallback(func() {
		inventorywidget.Show()
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
	gamefile.WriteString(fmt.Sprintf(`"clock": %s,`, timer.GlobalGameTimer.Save()+"\n"))
	gamefile.WriteString(fmt.Sprintf(`"inventory": %s,`, supplier.GlobalInventory.Save()+"\n"))
	gamefile.WriteString(fmt.Sprintf(`"ledger": %s`, accounting.GlobalLedger.Save()+"\n"))
	gamefile.WriteString("}\n")

	gamefile.Close()
}
