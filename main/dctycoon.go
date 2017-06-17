package main

import (
	"encoding/json"
	"fmt"
	"github.com/nzin/dctycoon"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/sws"
	"os"
	"time"
)

func main() {

	root := sws.Init(800, 600)
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

	gamemap := v["map"].(map[string]interface{})
	dc.LoadMap(gamemap)
	supplier.Trends = supplier.TrendLoad(v["trends"].(map[string]interface{}))
	root.AddChild(dc)
	root.SetFocus(dc)
	supplierwidget.Show()

	fmt.Println(supplier.Trends.Cpuprice.CurrentValue(time.Now()))

	for sws.PoolEvent() == false {
	}
	data := dc.SaveMap()
	gamefile, err = os.Create("backup.map")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	gamefile.WriteString("{")
	gamefile.WriteString(fmt.Sprintf(`"map": %s,`, data) + "\n")
	gamefile.WriteString(fmt.Sprintf(`"trends": %s`, supplier.TrendSave(supplier.Trends)) + "\n")
	gamefile.WriteString("}\n")

	gamefile.Close()
}
