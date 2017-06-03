package dctycoon

import (
    "github.com/nzin/sws"
    "github.com/veandco/go-sdl2/sdl"
    "time"
)



type DcWidget struct {
    sws.SWS_CoreWidget
    tiles       [][]*Tile
    xRoot,yRoot int32
}



func (self *DcWidget) Repaint() {
    mapheight:=len(self.tiles)
    mapwidth:=len(self.tiles[0])
    self.FillRect(0,0,self.Width(),self.Height(),0xff000000)
    for y:=0; y<mapheight; y++ {
        for x:=0; x<mapwidth; x++ {
            tile := self.tiles[y][x]
            if tile!=nil {
                surface := (*tile).Draw()
                rectSrc := sdl.Rect{0,0,surface.W,surface.H}
                rectDst := sdl.Rect{self.xRoot+(self.Surface().W/2)+(TILE_WIDTH_STEP/2)*int32(x)-(TILE_WIDTH_STEP/2)*int32(y),self.yRoot+(TILE_HEIGHT_STEP/2)*int32(x)+(TILE_HEIGHT_STEP/2)*int32(y),surface.W,surface.H}
                surface.Blit(&rectSrc,self.Surface(),&rectDst)
            }
        }
    }
    sws.PostUpdate()
}



var te *sws.TimerEvent

func (self *DcWidget) HasFocus(focus bool) {
    if (focus==false) {
        sws.StopRepeat(te)
        te=nil
    }
}

func (self *DcWidget) MouseMove(x,y,xrel,yrel int32) {
    if te!=nil {
        sws.StopRepeat(te)
        te=nil
    }
    if (x<10) {
        te=sws.TimerAddEvent(time.Now(),50*time.Millisecond,func() {
            self.MoveLeft()
        })
        return
    }
    if (y<10) {
        te=sws.TimerAddEvent(time.Now(),50*time.Millisecond,func() {
            self.MoveUp()
        })
        return
    }
    if (x>self.Width()-10) {
        te=sws.TimerAddEvent(time.Now(),50*time.Millisecond,func() {
            self.MoveRight()
        })
        return
    }
    if (y>self.Height()-10) {
        te=sws.TimerAddEvent(time.Now(),50*time.Millisecond,func() {
            self.MoveDown()
        })
        return
    }
}


func (self *DcWidget) MoveLeft() {
    width:=int32(len(self.tiles[0])+len(self.tiles)+1)*TILE_WIDTH_STEP/2
    if self.xRoot-width/2+self.Width()/2<0 {
        self.xRoot+=20
        sws.PostUpdate()
    }
}

func (self *DcWidget) MoveUp() {
    height:=int32(len(self.tiles[0])+len(self.tiles))*TILE_HEIGHT_STEP/2+TILE_HEIGHT
    if self.yRoot-height/2+self.Height()/2<0 {
        self.yRoot+=20
        sws.PostUpdate()
    }
}

func (self *DcWidget) MoveRight() {
    width:=int32(len(self.tiles[0])+len(self.tiles)+1)*TILE_WIDTH_STEP/2
    if self.xRoot+width>self.Width() {
        self.xRoot-=20
        sws.PostUpdate()
    }
}

func (self *DcWidget) MoveDown() {
    height:=int32(len(self.tiles[0])+len(self.tiles))*TILE_HEIGHT_STEP/2+TILE_HEIGHT
    if self.yRoot+height>self.Height() {
        self.yRoot-=20
        sws.PostUpdate()
    }
}

func (self *DcWidget) KeyDown(key sdl.Keycode, mod uint16) {
    if (key == sdl.K_LEFT) {
        self.MoveLeft()
    }
    if (key == sdl.K_UP) {
        self.MoveUp()
    }
    if (key == sdl.K_RIGHT) {
        self.MoveRight()
    }
    if (key == sdl.K_DOWN) {
        self.MoveDown()
    }


}


/*
 * {
 *   "width": 10,
 *   "height": 10,
 *   "tiles": [
 *     {"x":0, "y":0, "wall0":"","wall1":"","floor":"inside","dcelementtype":"rack","dcelement":{...}},
 *     {"x":1, "y":0, "wall0":"","wall1":"","floor":"inside"},
 *   ]
 * }
 */
func (self *DcWidget) LoadMap(dc map[string]interface{}) {
    width := int32(dc["width"].(float64))
    height := int32(dc["height"].(float64))
    self.tiles = make([][]*Tile,height)
    for y:= range self.tiles {
        self.tiles[y] = make([]*Tile,width)
        for x:= range self.tiles[y] {
            self.tiles[y][x] = CreateGrassTile()
        }
    }
    tiles := dc["tiles"].([]interface{})
    for _,t := range tiles {
        tile := t.(map[string]interface{})
        x := int32(tile["x"].(float64))
        y := int32(tile["y"].(float64))
        wall0 := tile["wall0"].(string)
        wall1 := tile["wall1"].(string)
        floor := tile["floor"].(string)
        var dcelementtype string
        var dcelement map[string]interface{}
        if (tile["dcelementtype"] != nil) {
            dcelementtype = tile["dcelementtype"].(string)
        }
        if (tile["dcelement"]!=nil) {
            dcelement = tile["dcelement"].(map[string]interface{})
        }
        if (dcelementtype=="" || dcelementtype=="rack") { // basic floor
            self.tiles[y][x]=CreateElectricalTile(wall0,wall1,floor, dcelementtype, dcelement)
        } 
    }
}



func (self *DcWidget) SaveMap() map[string]interface{} {
    return nil
}



func CreateDcWidget(w,h int32) *DcWidget {
    corewidget := sws.CreateCoreWidget(w,h)
    widget := &DcWidget { SWS_CoreWidget: *corewidget,
        tiles: [][]*Tile{{}},
        xRoot: 0,
        yRoot: 0,
    }
    return widget
}
