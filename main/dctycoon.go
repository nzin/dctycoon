package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/nzin/dctycoon"
	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/dctycoon/timer"
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
	flag.Parse()
	initLog(*loglevel, *logfile)

	gamefile, err := os.Open("example.map")
	if err != nil {
		log.Error("main(): ", err.Error())
		os.Exit(1)
	}
	var v map[string]interface{}
	jsonParser := json.NewDecoder(gamefile)
	if err = jsonParser.Decode(&v); err != nil {
		log.Error("main(): parsing game file ", err.Error())
		os.Exit(1)
	}
	gamefile.Close()

	root := sws.Init(800, 600)

	timer.GlobalEventPublisher = timer.NewEventPublisher(root)
	timer.GlobalGameTimer = timer.NewGameTimer()
	trends := supplier.TrendLoad(v["trends"].(map[string]interface{}), timer.GlobalEventPublisher, timer.GlobalGameTimer)

	// specific to player
	location := v["location"].(string)
	accounting.GlobalLedger = accounting.NewLedger(timer.GlobalGameTimer, supplier.AvailableLocation[location].Taxrate, supplier.AvailableLocation[location].Bankinterestrate)
	supplier.GlobalInventory = supplier.NewInventory(timer.GlobalGameTimer)

	dc := dctycoon.NewDcWidget(root.Width(), root.Height(), root, supplier.GlobalInventory)
	supplierwidget := dctycoon.NewMainSupplierWidget(trends, root)
	inventorywidget := dctycoon.NewMainInventoryWidget(root)
	accountingui := accounting.NewMainAccountingWidget(root)

	// initiate the game timer
	timer.GlobalGameTimer.Load(v["clock"].(map[string]interface{}))

	// initiate the ledger
	accounting.GlobalLedger.Load(v["ledger"].(map[string]interface{}), supplier.AvailableLocation[location].Taxrate, supplier.AvailableLocation[location].Bankinterestrate)
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
	gamefile.WriteString(fmt.Sprintf(`"location": "%s",`, location) + "\n")
	gamefile.WriteString(fmt.Sprintf(`"map": %s,`, data) + "\n")
	gamefile.WriteString(fmt.Sprintf(`"trends": %s,`, supplier.TrendSave(supplier.Trends)) + "\n")
	gamefile.WriteString(fmt.Sprintf(`"clock": %s,`, timer.GlobalGameTimer.Save()+"\n"))
	gamefile.WriteString(fmt.Sprintf(`"inventory": %s,`, supplier.GlobalInventory.Save()+"\n"))
	gamefile.WriteString(fmt.Sprintf(`"ledger": %s`, accounting.GlobalLedger.Save()+"\n"))
	gamefile.WriteString("}\n")

	gamefile.Close()
}
