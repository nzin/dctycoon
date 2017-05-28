package dctycoon

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/sdl_image"
    "fmt"
    "strconv"
)

// base "class" for all tiles
type DcElement interface { // should be passive, rack, ...
    Save() // json
    Load() // json
    Draw() *sdl.Surface
    Rotate(face uint32) // face can be 0,1,2,3 (i.e. 0, 90, 180, 270)
    Power() int32 // ampere
}



type Rackelement struct {
    size          int32 // 2u,4u, 8u...
    name          string // space, rack2u, rack4u, blade, switch, KVM ...
    sprite        string // name of the png 
    power         int32  // in Ampere
    // what about elements: disk, cpu, RAM...
    disksize      int32 // in Go
    disknum       int32
    cpunum        int32 
    ramsize       int32 // in Go
    virtualizable bool // does it has VT
}



type Rack struct {
    rackmount [] Rackelement // must fill 42u from top to bottom
    surface   *sdl.Surface
    rotation  uint32
}

func (self *Rack) Save() {
}

func (self *Rack) Load() {
}

func (self *Rack) Draw() *sdl.Surface {
    return self.surface
}

func (self *Rack) Rotate(rotation uint32) {
    self.rotation=rotation
}

func (self *Rack) Power() int32 {
    power := int32(0)
    for _,e := range self.rackmount {
        power+=e.power
    }
    return power
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
    wall           [4]string // "" when nothing
    floor          string
    element        DcElement
    surface        *sdl.Surface
    previousRotate uint32
}

func (self *Tile) Save() {
}

func (self *Tile) Load() {
}

func (self *Tile) Draw(rotate uint32) *sdl.Surface {
    if self.surface == nil || self.previousRotate!=rotate{
        self.surface,_ = sdl.CreateRGBSurface(0,140,182,32,0x00ff0000,0x0000ff00,0x000000ff,0xff000000)
        floor := getSprite("resources/"+self.floor+strconv.Itoa(int(rotate))+".png")
        rectSrc := sdl.Rect{0,0,floor.W,floor.H}
        rectDst := sdl.Rect{0,182-floor.H,floor.W,floor.H}
        floor.Blit(&rectSrc,self.surface,&rectDst)

//    return self.element.Draw()
    }
    self.previousRotate=rotate
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
        wall: [4]string{"","","",""},
        floor: "green",
        element: nil,
    }
    return tile
}



func CreateElectricalTile(wall0,wall1,wall2,wall3,floor,dcelementtype string, dcelement map[string]interface{}) *Tile {
    var element *ElectricalElement
    if dcelementtype!="" {
        element = CreateElectricalElement(dcelementtype,dcelement)
    }
    tile := &Tile {
        wall: [4]string{wall0,wall1,wall2,wall3},
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



