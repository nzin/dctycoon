package dctycoon

import (
    "github.com/nzin/sws"
    "github.com/veandco/go-sdl2/sdl"
)



type DcWidget struct {
    sws.SWS_CoreWidget
    tiles       [][]*Tile
    xRoot,yRoot int32
    rotate      uint32
}



func (self *DcWidget) Repaint() {
    mapheight:=len(self.tiles)
    mapwidth:=len(self.tiles[0])
    self.FillRect(0,0,self.Width(),self.Height(),0xff000000)
    if (self.rotate==0) {
        for y:=0; y<mapheight; y++ {
            for x:=0; x<mapwidth; x++ {
                tile := self.tiles[y][x]
                if tile!=nil {
                    surface := (*tile).Draw(self.rotate)
                    rectSrc := sdl.Rect{0,0,surface.W,surface.H}
                    rectDst := sdl.Rect{self.xRoot+(self.Surface().W/2)+70*int32(x)-70*int32(y),self.yRoot+40*int32(x)+40*int32(y),surface.W,surface.H}
                    surface.Blit(&rectSrc,self.Surface(),&rectDst)
                }
            }
        }
    }
    if (self.rotate==1) {
        for x:=0; x<mapwidth; x++ {
            for y:=mapheight-1; y>=0; y-- {
                tile := self.tiles[y][x]
                if tile!=nil {
                    surface := (*tile).Draw(self.rotate)
                    rectSrc := sdl.Rect{0,0,surface.W,surface.H}
                    rectDst := sdl.Rect{self.xRoot+(self.Surface().W/2)+70*int32(mapheight-1-y)-70*   int32(x),self.yRoot+40*int32(mapheight-1-y)+40*int32(x),surface.W,surface.H}
                    surface.Blit(&rectSrc,self.Surface(),&rectDst)
                }
            }
        }
    }
    if (self.rotate==2) {
        for y:=mapheight-1; y>=0; y-- {
            for x:=mapwidth-1; x>=0; x-- {
                tile := self.tiles[y][x]
                if tile!=nil {
                    surface := (*tile).Draw(self.rotate)
                    rectSrc := sdl.Rect{0,0,surface.W,surface.H}
                    rectDst := sdl.Rect{self.xRoot+(self.Surface().W/2)+70*int32(mapwidth-1-x)-70*int32(mapheight-1-y),self.yRoot+40*int32(mapwidth-1-x)+40*int32(mapheight-1-y),surface.W,surface.H}
                    surface.Blit(&rectSrc,self.Surface(),&rectDst)
                }
            }
        }
    }
    if (self.rotate==3) {
        for y:=0; y<mapheight; y++ {
            for x:=mapwidth-1; x>=0; x-- {
                tile := self.tiles[y][x]
                if tile!=nil {
                    surface := (*tile).Draw(self.rotate)
                    rectSrc := sdl.Rect{0,0,surface.W,surface.H}
                    rectDst := sdl.Rect{self.xRoot+(self.Surface().W/2)+70*int32(y)-70*   int32(mapwidth-1-x),self.yRoot+40*int32(y)+40*int32(mapwidth-1-x),surface.W,surface.H}
                    surface.Blit(&rectSrc,self.Surface(),&rectDst)
                }
            }
        }
    }
    sws.PostUpdate()
}



func (self *DcWidget) KeyDown(key sdl.Keycode, mod uint16) {
/*    if key == sdl.K_LEFT {
        if (mod == sdl.KMOD_LSHIFT || mod == sdl.KMOD_RSHIFT) {
            if self.initialCursorPosition>0 {
                self.initialCursorPosition--
            }
*/
    if (key=='r') {
        self.rotate=(self.rotate+1)%4
        sws.PostUpdate()
    }

}


/*
 * {
 *   "width": 10,
 *   "height": 10,
 *   "tiles": [
 *     {"x":0, "y":0, "wall0":"","wall1":"","wall2":"","wall3":"","floor":"inside","dcelementtype":"rack","dcelement":{...}},
 *     {"x":1, "y":0, "wall0":"","wall1":"","wall2":"","wall3":"","floor":"inside"},
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
        wall2 := tile["wall2"].(string)
        wall3 := tile["wall3"].(string)
        floor := tile["floor"].(string)
        var dcelementtype string
        var dcelement map[string]interface{}
        if (tile["dcelementtype"] != nil) {
            dcelementtype = tile["dcelementtype"].(string)
        }
        if (tile["dcelement"]!=nil) {
            dcelement = tile["dcelement"].(map[string]interface{})
        }
        if (dcelementtype=="") { // basic floor
            self.tiles[y][x]=CreateElectricalTile(wall0,wall1,wall2,wall3,floor, dcelementtype, dcelement)
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
