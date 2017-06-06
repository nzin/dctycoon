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
    Pit    time.Time
    Value  interface{}
}

type TrendList []TrendItem



func (self TrendList) Len() int {
  return len(self)
}

func (self TrendList) Swap(i,j int) {
    self[i],self[j] = self[j],self[i]
}

func (self TrendList) Less(i, j int) bool {
    return self[i].Pit.Before(self[j].Pit)
}

func (self TrendList) Sort() {
    sort.Sort(self)
}

func (self TrendList) CurrentValue(now time.Time) interface{} {
    if (len(self)==0) { panic("no elements in the array") }
    
    index:=0
    for (index<len(self) && self[index].Pit.Before(now)) {
        index++
    }
    if (index==0) { return self[0].Value }
    return self[index-1].Value
}

func TrendListLoad(json []interface{}) TrendList {
    tl := make(TrendList,len(json))
    for i,t := range json {
        te:=t.(map[string]interface{})
        var year,month,day int
        fmt.Sscanf(te["pit"].(string),"%d-%d-%d",&year,&month,&day)
        tl[i]=TrendItem{
            Pit:   time.Date(year,time.Month(month),day,0,0,0,0,time.UTC),
            Value: te["value"],
        }
    }
    tl.Sort()
    return tl
}



func TrendListSave(t TrendList) string {
    str:=`[`
    for i,te := range t {
        if i>0 {str+=","}
        str+=fmt.Sprintf(`{"pit":"%d-%d-%d", "value":%v}`,te.Pit.Year(),te.Pit.Month(),te.Pit.Day(),te.Value)
    }
    str+=`]`
    return str
}





/*****************/



type PriceTrendItem struct {
    Pit    time.Time
    Value  float64
}

type PriceTrendList []PriceTrendItem

func (self PriceTrendList) Len() int {
  return len(self)
}

func (self PriceTrendList) Swap(i,j int) {
    self[i],self[j] = self[j],self[i]
}

func (self PriceTrendList) Less(i, j int) bool {
    return self[i].Pit.Before(self[j].Pit)
}

func (self PriceTrendList) Sort() {
    sort.Sort(self)
}


type PriceTrend struct {
    Trend PriceTrendList
    Noise PriceTrendList
}


func (self *PriceTrend) CurrentValue(now time.Time) float64 {
    if (self.Trend==nil) || (len(self.Trend)==0) {
        temp:=make(PriceTrendList,1)
        temp[0].Pit=now
        temp[0].Value=0.0
        self.Trend=temp
    }
   
    index:=0
    for (index<len(self.Trend) && self.Trend[index].Pit.Before(now)) {
        index++
    }
    
    var Value float64
    if (index==0) {
        Value=self.Trend[0].Value
    } else if (index==len(self.Trend)) {
        Value=self.Trend[index-1].Value
    } else {
        interval:=(self.Trend[index].Pit.Sub(self.Trend[index-1].Pit)).Hours()
        since:=now.Sub(self.Trend[index-1].Pit).Hours()
        Value= self.Trend[index-1].Value*((interval-since)/interval)+self.Trend[index].Value*(since/interval)
    }
    
    // now compute the noise
    
    if (self.Noise==nil) || len(self.Noise)==0 {
        temp:=make(PriceTrendList,1)
        temp[0].Pit=now
        temp[0].Value=0.0
        self.Noise=temp
    }
    endarray:=len(self.Noise)-1
    for (now.Before((self.Noise)[endarray].Pit) == false) {
       random:=rand.Float64()
       if random<0.1 { random=0.1 }
       elt:=PriceTrendItem{
           Pit:   (self.Noise)[endarray].Pit.AddDate(0,0,int(100*random)),
           Value: 1.0-random,
       }
       self.Noise=append(self.Noise,elt)
       endarray=len(self.Noise)-1
    }

    for (index<len(self.Noise) && (self.Noise)[index].Pit.Before(now)) {
        index++
    }
    var noise float64
    if (index==0) {
        noise= (self.Noise)[0].Value
    } else if (index==len(self.Noise)) {
        return (self.Noise)[index-1].Value
    } else {
        interval:=((self.Noise)[index].Pit.Sub((self.Noise)[index-1].Pit)).Hours()
        since:=now.Sub((self.Noise)[index-1].Pit).Hours()
        noise= (self.Noise)[index-1].Value*((interval-since)/interval)+(self.Noise)[index].Value*(since/interval)
    }
    return Value*(noise+1.0)
}



func PriceTrendListLoad(json map[string]interface{}) PriceTrend {
    pt:=PriceTrend{
    }
    trend:=json["trend"].([]interface{})
    noise:=json["noise"].([]interface{})
    
    tl := make(PriceTrendList,len(trend))
    for i,t := range trend {
        te:=t.(map[string]interface{})
        var year,month,day int
        fmt.Sscanf(te["pit"].(string),"%d-%d-%d",&year,&month,&day)
        tl[i]=PriceTrendItem{
            Pit:   time.Date(year,time.Month(month),day,0,0,0,0,time.UTC),
            Value: te["value"].(float64),
        }
    }
    tl.Sort()
    pt.Trend=tl
    
    nl := make(PriceTrendList,len(noise))
    for i,n := range noise {
        ne:=n.(map[string]interface{})
        var year,month,day int
        fmt.Sscanf(ne["pit"].(string),"%d-%d-%d",&year,&month,&day)
        nl[i]=PriceTrendItem{
            Pit:   time.Date(year,time.Month(month),day,0,0,0,0,time.UTC),
            Value: ne["value"].(float64),
        }
    }
    nl.Sort()
    pt.Noise=nl

    return pt
}



func PriceTrendListSave(pt PriceTrend) string {
    str:=`{ "trend":[`
    for i,te := range pt.Trend {
        if i>0 {str+=","}
        str+=fmt.Sprintf(`{"pit":"%d-%d-%d", "value":%g}`,te.Pit.Year(),te.Pit.   Month(),te.Pit.Day(),te.Value)
    }
    str+=`],"noise":[`
    for i,ne := range pt.Noise {
        if i>0 {str+=","}
        str+=fmt.Sprintf(`{"pit":"%d-%d-%d", "value":%g}`,ne.Pit.Year(),ne.Pit.   Month(),ne.Pit.Day(),ne.Value)
    }
    str+=`]}`
    return str
}



type Trend struct {
    Corepercpu TrendList
    Vt         TrendList
    Disksize   TrendList
    Ramsize    TrendList
    
    Cpuprice   PriceTrend
    Diskprice  PriceTrend
    Ramprice   PriceTrend
}



func TrendLoad(json map[string]interface{}) *Trend {
    t:=&Trend{
      Corepercpu:TrendListLoad(json["corepercpu"].([]interface{})),
      Vt:        TrendListLoad(json["vt"].([]interface{})),
      Disksize:  TrendListLoad(json["disksize"].([]interface{})),
      Ramsize:   TrendListLoad(json["ramsize"].([]interface{})),
    
      Cpuprice:  PriceTrendListLoad(json["cpuprice"].(map[string]interface{})),
      Diskprice: PriceTrendListLoad(json["diskprice"].(map[string]interface{})),
      Ramprice:  PriceTrendListLoad(json["ramprice"].(map[string]interface{})),
    }
    
    return t
}



func TrendSave(t *Trend) string {
    str:="{\n"
    str+=fmt.Sprintf(`"corepercpu": %s,`,TrendListSave(t.Corepercpu))+"\n"
    str+=fmt.Sprintf(`"vt": %s,`,TrendListSave(t.Vt))+"\n"
    str+=fmt.Sprintf(`"disksize": %s,`,TrendListSave(t.Disksize))+"\n"
    str+=fmt.Sprintf(`"ramsize": %s,`,TrendListSave(t.Ramsize))+"\n"
    
    str+=fmt.Sprintf(`"cpuprice": %s,`,PriceTrendListSave(t.Cpuprice))+"\n"
    str+=fmt.Sprintf(`"diskprice": %s,`,PriceTrendListSave(t.Diskprice))+"\n"
    str+=fmt.Sprintf(`"ramprice": %s`,PriceTrendListSave(t.Ramprice))+"\n"
    str+="}"
    return str
}

