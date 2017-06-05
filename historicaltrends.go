package dctycoon

import (
    "time"
    "sort"
    "fmt"
    "math/rand"
)



/*
 * We will have different trends to follow
 * - number of cores/CPU
 * - vt 
 * - disk size / plateau
 * - ram size / slot
 * 
 * noisy trends (float64)
 * - cpu price / core + noise
 * - disk price / Go + noise
 * - ram price / Go + noise
 *
 * - and the same for customer needs?
 */
type TrendItem struct {
    pit    time.Time
    value  interface{}
}

type TrendList []TrendItem



func (self TrendList) Len() int {
  return len(self)
}

func (self TrendList) Swap(i,j int) {
    self[i],self[j] = self[j],self[i]
}

func (self TrendList) Less(i, j int) bool {
    return self[i].pit.Before(self[j].pit)
}

func (self TrendList) Sort() {
    sort.Sort(self)
}

func (self TrendList) CurrentValue(now time.Time) interface{} {
    if (len(self)==0) { panic("no elements in the array") }
    
    index:=0
    for (index<len(self) && self[index].pit.Before(now)) {
        index++
    }
    if (index==0) { return self[0].value }
    return self[index-1].value
}

func TrendListLoad(json []interface{}) TrendList {
    tl := make(TrendList,len(json))
    for i,t := range json {
        te:=t.(map[string]interface{})
        var date time.Time
        var year,month,day int
        fmt.Sscanf(te["pit"].(string),"%d-%d-%d",&year,&month,&day)
        tl[i]=TrendItem{
            pit:   date,
            value: te["value"],
        }
    }
    tl.Sort()
    return tl
}



/*****************/



type PriceTrendItem struct {
    pit    time.Time
    value  float64
}

type PriceTrendList []PriceTrendItem

func (self PriceTrendList) Len() int {
  return len(self)
}

func (self PriceTrendList) Swap(i,j int) {
    self[i],self[j] = self[j],self[i]
}

func (self PriceTrendList) Less(i, j int) bool {
    return self[i].pit.Before(self[j].pit)
}

func (self PriceTrendList) Sort() {
    sort.Sort(self)
}

func (self PriceTrendList) CurrentValue(now time.Time) float64 {
    if (len(self)==0) { panic("no elements in the array") }
   
    index:=0
    for (index<len(self) && self[index].pit.Before(now)) {
        index++
    }
    if (index==0) { return self[0].value }
    if (index==len(self)) { return self[index-1].value }
    interval:=(self[index].pit.Sub(self[index-1].pit)).Hours()
    since:=now.Sub(self[index-1].pit).Hours()
    return self[index-1].value*(since/interval)+self[index].value*((interval-since)/interval)
}


func PriceTrendListLoad(json []interface{}) PriceTrendList {
    tl := make(PriceTrendList,len(json))
    for i,t := range json {
        te:=t.(map[string]interface{})
        var date time.Time
        var year,month,day int
        fmt.Sscanf(te["pit"].(string),"%d-%d-%d",&year,&month,&day)
        tl[i]=PriceTrendItem{
            pit:   date,
            value: te["value"].(float64),
        }
    }
    tl.Sort()
    return tl
}



/*********************/



type NoiseTrendItem struct {
    pit    time.Time
    value  float64
}

type NoiseTrendList []NoiseTrendItem

func (self NoiseTrendList) Len() int {
  return len(self)
}

func (self NoiseTrendList) Swap(i,j int) {
    self[i],self[j] = self[j],self[i]
}

func (self NoiseTrendList) Less(i, j int) bool {
    return self[i].pit.Before(self[j].pit)
}

func (self NoiseTrendList) Sort() {
    sort.Sort(self)
}

func (self *NoiseTrendList) CurrentValue(now time.Time) float64 {
    if self==nil || len(*self)==0 {
        temp:=make(NoiseTrendList,1)
        temp[0].pit=now
        temp[0].value=0.0
        self=&temp
    }
    endarray:=len(*self)-1
    if (! now.Before((*self)[endarray].pit) ) {
       random:=rand.Float64()
       if random<0.1 { random=0.1 }
       elt:=NoiseTrendItem{
           pit:   now.AddDate(0,0,int(100*random)),
           value: 1.0-random,
       }
       *self=append(*self,elt)
    }

    index:=0
    for (index<len(*self) && (*self)[index].pit.Before(now)) {
        index++
    }
    if (index==0) { return (*self)[0].value }
    if (index==len(*self)) { return (*self)[index-1].value }

    interval:=((*self)[index].pit.Sub((*self)[index-1].pit)).Hours()
    since:=now.Sub((*self)[index-1].pit).Hours()
    return (*self)[index-1].value*(since/interval)+(*self)[index].value*((interval-   since)/interval)
}

func NoiseTrendListLoad(json []interface{}) *NoiseTrendList {
    tl := make(NoiseTrendList,len(json))
    for i,t := range json {
        te:=t.(map[string]interface{})
        var date time.Time
        var year,month,day int
        fmt.Sscanf(te["pit"].(string),"%d-%d-%d",&year,&month,&day)
        tl[i]=NoiseTrendItem{
            pit:   date,
            value: te["value"].(float64),
        }
    }
    tl.Sort()
    return &tl
}



