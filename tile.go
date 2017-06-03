package dctycoon

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/sdl_image"
    "fmt"
    "strconv"
)

const (
    TILE_WIDTH_STEP = 96
    TILE_HEIGHT_STEP = 48
    TILE_HEIGHT = 257
)

// base "class" for all tiles
type DcElement interface { // should be passive, rack, ...
    Save() // json
    Load() // json
    Draw() *sdl.Surface
    Rotate(face uint32) // face can be 0,1,2,3 (i.e. 0, 90, 180, 270)
    Power() int32 // ampere
}



type RackPart struct {
    size          int32 // 1u,2u,4u, 8u...
    name          string // space, rack2u, rack4u, blade, switch, KVM ...
    //sprite        string // name of the png 
    power         int32  // in Ampere
    // what about elements: disk, cpu, RAM...
    disksize      int32 // in Go
    disknum       int32
    cpunum        int32 
    ramsize       int32 // in Go
    virtualizable bool // does it has VT
}

func CreateRackPart(payload map[string]interface{}) *RackPart {
    name:=payload["name"].(string)
    size:=int32(payload["size"].(float64))
    if (name=="space") {
        rp := &RackPart {
            size: size,
            name: name,
            power: 0,
            disksize: 0,
            disknum: 0,
            cpunum: 0,
            ramsize: 0,
            virtualizable: false,
        }
        return rp
    }
    power:=int32(payload["power"].(float64))
    disksize:=int32(payload["disksize"].(float64))
    disknum:=int32(payload["disknum"].(float64))
    cpunum:=int32(payload["cpunum"].(float64))
    ramsize:=int32(payload["ramsize"].(float64))
    virtualizable:=payload["vt"].(bool)
    
    rp :=&RackPart {
        size: size,
        name: name,
        power: power,
        disksize: disksize,
        disknum: disknum,
        cpunum: cpunum,
        ramsize: ramsize,
        virtualizable: virtualizable,
    }
    return rp
}



type RackElement struct {
    rackmount [] *RackPart // must fill 42u from bottom to top
    surface   *sdl.Surface
    rotation  uint32
}

func (self *RackElement) Save() {
}

func (self *RackElement) Load() {
}

func (self *RackElement) Draw() *sdl.Surface {
    if (self.surface==nil) {
        self.surface= getSprite("resources/rack.bottom"+strconv.Itoa(int(self.rotation))+".png")
        var offset int32=0
        for _,rp := range self.rackmount {
            if rp.name!="space" {
                img:=getSprite("resources/"+rp.name+strconv.Itoa(int(self.rotation))+".png")
                rectSrc := sdl.Rect{0,0,img.W,img.H}
                rectDst := sdl.Rect{0,TILE_HEIGHT-img.H-(offset+rp.size+1)*4,img.W,img.H}
                img.Blit(&rectSrc,self.surface,&rectDst)
            }
            offset+=rp.size
        }
        top:= getSprite("resources/rack.top"+strconv.Itoa(int(self.rotation))+".png")
        rectSrc := sdl.Rect{0,0,top.W,top.H}
        rectDst := sdl.Rect{0,TILE_HEIGHT-top.H,top.W,top.H}
        top.Blit(&rectSrc,self.surface,&rectDst)
    }
    return self.surface
}

func (self *RackElement) Rotate(rotation uint32) {
    self.rotation=rotation
}

func (self *RackElement) Power() int32 {
    power := int32(0)
    for _,e := range self.rackmount {
        power+=e.power
    }
    return power
}



func CreateRackElement(payload map[string]interface{}) *RackElement {
    rotation := uint32(payload["rotation"].(float64))
    rackmount := make([]*RackPart,0)
    rackParts := payload["rackmount"].([]interface{})
    for _,rp := range rackParts {
        rackmount=append(rackmount,CreateRackPart(rp.(map[string]interface{})))
    }
    r := &RackElement { 
        rackmount: rackmount,
        surface: nil,
        rotation: rotation,
    }
    return r
}



type ElectricalElement struct {
    flavor    string // ac, battery, generatorA,generatorB
    power     int32  // negative if it is a generator
    capacity  int32  // kWh if it is a battery
    surface   *sdl.Surface
    rotation  uint32
}

func (self *ElectricalElement) Save() {
}

func (self *ElectricalElement) Load() {
}

func (self *ElectricalElement) Draw() *sdl.Surface {
    return self.surface
}

func (self *ElectricalElement) Rotate(face uint32) {
}

func (self *ElectricalElement) Power() int32 {
    return self.power
}



func CreateElectricalElement(flavor string, payload map[string]interface{}) *ElectricalElement {
    rotation := uint32(payload["rotation"].(float64))
    power := int32(payload["power"].(float64))
    capacity := int32(payload["capacity"].(float64))
    surface := getSprite("resources/"+flavor+strconv.Itoa(int(rotation))+".png")
    ee := &ElectricalElement { 
        flavor: flavor,
        power: power,
        capacity: capacity,
        surface: surface,
        rotation: rotation,
    }
    return ee
}



type Tile struct {
    wall           [2]string // "" when nothing
    floor          string
    element        DcElement
    surface        *sdl.Surface
}

func (self *Tile) Save() {
}

func (self *Tile) Load() {
}

func (self *Tile) Draw() *sdl.Surface {
    if self.surface == nil {
        self.surface,_ = sdl.CreateRGBSurface(0,105,TILE_HEIGHT,32,0x00ff0000,0x0000ff00,0x000000ff,0xff000000)
        floor := getSprite("resources/"+self.floor+".png")
        rectSrc := sdl.Rect{0,0,floor.W,floor.H}
        rectDst := sdl.Rect{0,TILE_HEIGHT-floor.H,floor.W,floor.H}
        floor.Blit(&rectSrc,self.surface,&rectDst)
        if (self.wall[0] != "") || (self.wall[1] != "") {
            wall := getSprite("resources/wallX.png")
            rectSrc := sdl.Rect{0,0,wall.W,wall.H}
            rectDst := sdl.Rect{0,TILE_HEIGHT-wall.H,wall.W,wall.H}
            wall.Blit(&rectSrc,self.surface,&rectDst)
        }
        if self.wall[0] != "" {
            wall := getSprite("resources/"+self.wall[0]+"L.png")
            rectSrc := sdl.Rect{0,0,wall.W,wall.H}
            rectDst := sdl.Rect{0,TILE_HEIGHT-wall.H,wall.W,wall.H}
            wall.Blit(&rectSrc,self.surface,&rectDst)
        }
        if self.wall[1] != "" {
            wall := getSprite("resources/"+self.wall[1]+"R.png")
            rectSrc := sdl.Rect{0,0,wall.W,wall.H}
            rectDst := sdl.Rect{0,TILE_HEIGHT-wall.H,wall.W,wall.H}
            wall.Blit(&rectSrc,self.surface,&rectDst)
        }

        if (self.element!=nil) {
            elt:=self.element.Draw()
            rectSrc := sdl.Rect{0,0,elt.W,elt.H}
            rectDst := sdl.Rect{0,TILE_HEIGHT-elt.H,elt.W,elt.H}
            elt.Blit(&rectSrc,self.surface,&rectDst)
        }
    }
    return self.surface
}

// individual rotate
func (self *Tile) Rotate(face uint32) {
    self.surface=nil
    if self.element!=nil {
        self.element.Rotate(face)
    }
}

func (self *Tile) Power() int32 {
    return self.element.Power()
}




func CreateGrassTile() *Tile {
    tile := &Tile {
        wall: [2]string{"",""},
        floor: "green",
        element: nil,
    }
    return tile
}



func CreateElectricalTile(wall0,wall1,floor,dcelementtype string, dcelement map[string]interface{}) *Tile {
    var element DcElement
    if dcelementtype=="rack" {
        element = CreateRackElement(dcelement)
    } else if dcelementtype!="" {
        element = CreateElectricalElement(dcelementtype,dcelement)
    }
    tile := &Tile {
        wall: [2]string{wall0,wall1},
        floor: floor,
        element: element,
    }
    return tile
}



var spritecache map[string]*sdl.Surface

func getSprite(image string) *sdl.Surface{
    if spritecache == nil {
        spritecache = make(map[string]*sdl.Surface)
    }
    sprite:=spritecache[image]
    if (sprite==nil) {
       var err error
       sprite,err=img.Load(image)
       if sprite==nil || err!=nil {
           fmt.Println("Error loading ",image,err)
           panic(err)
       }
       spritecache[image]=sprite
    }
    return sprite
}



