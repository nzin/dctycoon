package dctycoon

import (
    "github.com/nzin/sws"
    "github.com/veandco/go-sdl2/sdl"
    "time"
    "fmt"
)



type DcWidget struct {
    sws.SWS_CoreWidget
    tiles         [][]*Tile
    xRoot,yRoot   int32
    activeTile    *Tile
    activeElement DcElement
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



func GetSurfacePixel(surface *sdl.Surface, x,y int32) (red,green,blue,alpha uint8) {
    if (x<0 || x>=surface.W || y<0 || y>= surface.H) {
        return 0,0,0,0
    }
    err := surface.Lock()
    if err!=nil {
        panic(err)
    }
    bpp := surface.Format.BytesPerPixel
    bytes:=surface.Pixels() 
    red=bytes[int(y) * int(surface.Pitch) + int(x) * int(bpp)]
    green=bytes[int(y) * int(surface.Pitch) + int(x) * int(bpp) +1]
    blue=bytes[int(y) * int(surface.Pitch) + int(x) * int(bpp)  +2]
    alpha=bytes[int(y) * int(surface.Pitch) + int(x) * int(bpp) +3]
    
    surface.Unlock()
    return
}



func (self *DcWidget) MousePressDown(x,y int32, button uint8) {
    self.activeTile=nil
    self.activeElement=nil
    mapheight:=len(self.tiles)
    mapwidth:=len(self.tiles[0])
    for ty:=mapheight-1; ty>=0; ty-- {
        for tx:=mapwidth-1; tx>=0; tx-- {
            tile:=self.tiles[ty][tx]
            surface := (*tile).Draw()
            xShift:=self.xRoot+(self.Surface().W/2)+(TILE_WIDTH_STEP/2)*int32(tx)-(TILE_WIDTH_STEP/2)*int32(ty)
            yShift:=self.yRoot+(TILE_HEIGHT_STEP/2)*int32(tx)+(TILE_HEIGHT_STEP/2)*int32(ty)

            if (x>=xShift) &&
               (y>=yShift) &&
               (x<xShift+surface.W) &&
               (y<yShift+surface.H) {
               _,_,_,alpha := GetSurfacePixel(surface,x-xShift,y-yShift)
               if alpha>0 {
                   //fmt.Println("activeTile: "+tile.floor,tx,ty)
                   self.activeTile=tile
                   // now, do we have an active item inside the tile
                   if tile.IsElementAt(x-xShift,y-yShift) {
                       //fmt.Println("active element!!!")
                       self.activeElement=tile.DcElement()
                   }
                   break
               }
           }
        }
        if self.activeTile!=nil {
            break
        }
    }
}



func (self *DcWidget) MousePressUp(x,y int32, button uint8) {
    if (self.activeElement!=nil) {
    }
}



var te *sws.TimerEvent

func (self *DcWidget) HasFocus(focus bool) {
    if (focus==false) {
        te.StopRepeat()
        te=nil
    }
}

func (self *DcWidget) MouseMove(x,y,xrel,yrel int32) {
    if te!=nil {
        te.StopRepeat()
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
        rotation := uint32(tile["rotation"].(float64))
        var dcelementtype string
        var dcelement map[string]interface{}
        if (tile["dcelementtype"] != nil) {
            dcelementtype = tile["dcelementtype"].(string)
        }
        if (tile["dcelement"]!=nil) {
            dcelement = tile["dcelement"].(map[string]interface{})
        }
        if (dcelementtype=="" || dcelementtype=="rack") { // basic floor
            self.tiles[y][x]=CreateElectricalTile(wall0,wall1,floor,rotation, dcelementtype, dcelement)
        } 
    }
}



func (self *DcWidget) SaveMap() string {
    s:=fmt.Sprintf(`{"width":%d, "height":%d, "tiles": [`,len(self.tiles[0]), len(self.tiles))
    previous:=false
    for y,_ := range self.tiles {
        for x,_ := range self.tiles[y] {
            t:=self.tiles[y][x]
            value:=""
            if t.element==nil {
                if (t.wall[0]!="" || t.wall[1]!="" || t.floor!="green") {
                    value=fmt.Sprintf(`{"x":%d, "y":%d, "wall0":"%s", "wall1":"%s", "floor":"%s","rotation":%d}`,
                      x,
                      y,
                      t.wall[0],
                      t.wall[1],
                      t.floor,
                      t.rotation,
                    )
                }
            } else {
                value=fmt.Sprintf(`{"x":%d, "y":%d, "wall0":"%s", "wall1":"%s", "floor":"%s", "rotation":%d, "dcelementtype":"%s", "dcelement":%s}`,
                    x,
                    y,
                    t.wall[0],
                    t.wall[1],
                    t.floor,
                    t.rotation,
                    t.element.ElementType(),
                    t.element.Save(),
                )
            }
            if value!="" {
                if previous==true { s+=",\n"}
                previous=true
                s+=value
            }
        }
    }
    s+="]}"
    return s
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
