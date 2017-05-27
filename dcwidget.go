package dctycoon

import (
    "github.com/nzin/sws"
    "github.com/veandco/go-sdl2/sdl"
)



type DcWidget struct {
    sws.SWS_CoreWidget
    tiles       [][]*DcElement
    xRoot,yRoot int32
}



func (self *DcWidget) Repaint() {
    self.FillRect(0,0,self.Width(),self.Height(),0xff000000)
    for y,_ := range self.tiles {
        for x,tile := range self.tiles[y] {
            if tile!=nil {
                surface := (*tile).Draw()
                rectSrc := sdl.Rect{0,0,surface.W,surface.H}
                rectDst := sdl.Rect{self.xRoot+35*int32(x),self.yRoot+20*int32(y),surface.W,surface.H}
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
    self.tiles = make([][]*DcElement,height)
    for y:= range self.tiles {
        self.tiles[y] = make([]*DcElement,width)
    }
}



func (self *DcWidget) SaveMap() map[string]interface{} {
    return nil
}



func CreateDcWidget(w,h int32) *DcWidget {
    corewidget := sws.CreateCoreWidget(w,h)
    widget := &DcWidget { SWS_CoreWidget: *corewidget,
        tiles: [][]*DcElement{{}},
        xRoot: 0,
        yRoot: 0,
    }
    return widget
}
