package main


type rackelement struct {
    size   int 32 // 2u,4u, 8u...
    name   string // space, rack2u, rack4u, blade, switch, KVM ...
    sprite string // name of the png 
    power  int32  // in Ampere
    // what about elements: disk, cpu, RAM...
}

interface dcElement { // should be passive, rack, ...
    func export() // json
    func import() // json
    func draw() 
    func power() int32 // ampere
}

type rack struct {
    rackmount [] rackelement // must fill 42u from top to bottom
}

func (*rack) export() {
}

func (*rack) import() {
}

func (*rack) draw() {
}

func (*rack) power() {
}



type electricalelement struct {
    flavor string // ac, battery, generatorA,generatorB
    int power // negative if it is a generator
    int capacity // kWh if it is a battery
}

func (*elecricalelement) export() {
}

func (*elecricalelement) import() {
}

func (*elecricalelement) draw() {
}

func (*elecricalelement) power() {
}



type Tile struct {
    wall string[4] // "" when nothing
    floor string
    element dcElement
    // sdl.surface
}

