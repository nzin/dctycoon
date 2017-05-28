package dctycoon

import (
    "github.com/nzin/sws"
    "github.com/veandco/go-sdl2/sdl"
)



type DcWidget struct {
    sws.SWS_CoreWidget
    tiles       [][]*Tile
    xRoot,yRoot int32
}



func (self *DcWidget) Repaint() {
    self.FillRect(0,0,self.Width(),self.Height(),0xff000000)
    for y,_ := range self.tiles {
        for x,tile := range self.tiles[y] {
            if tile!=nil {
                surface := (*tile).Draw()
                rectSrc := sdl.Rect{0,0,surface.W,surface.H}
                rectDst := sdl.Rect{self.xRoot+(self.Surface().W/2)+70*int32(x)-70*int32(y),self.yRoot+40*int32(x)+40*int32(y),surface.W,surface.H}
                surface.Blit(&rectSrc,self.Surface(),&rectDst)
            }
        }
    }
    sws.PostUpdate()
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
