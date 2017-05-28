package main

import (
    "github.com/nzin/sws"
    "github.com/nzin/dctycoon"
    "os"
    "fmt"
    "encoding/json"
)



func main() {
    root := sws.Init(800,600)
    dc := dctycoon.CreateDcWidget(root.Width(),root.Height())
    gamefile,err:=os.Open("example.map")
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

    dc.LoadMap(v)
    root.AddChild(dc)
    root.SetFocus(dc)
    
    for sws.PoolEvent() == false {
    }
}
