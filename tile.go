package dctycoon

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/sdl_image"
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

func (self *Rack) Rotate(face int32) {
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

func (self *ElectricalElement) Rotate(face int32) {
}

func (self *ElectricalElement) Power() int32 {
    return self.power
}



type Tile struct {
    wall [4]string // "" when nothing
    floor string
    element DcElement
}



var spritecache map[string]*sdl.Surface

func getSprite(image string) *sdl.Surface{
    sprite:=spritecache[image]
    if (sprite==nil) {
       var err error
       sprite,err=img.Load(image)
       if sprite==nil {
           panic(err)
       }
       spritecache[image]=sprite
    }
    return sprite
}



