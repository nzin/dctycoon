package dctycoon

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/img"
	"github.com/nzin/dctycoon/supplier"
	"strconv"
	"sort"
)

const (
	TILE_WIDTH_STEP  = 96
	TILE_HEIGHT_STEP = 48
	TILE_HEIGHT      = 257
)

// base "class" for all tiles
type DcElement interface {
	// should be passive, rack, ...
	ElementType() string               // to know which type to save
	Save() string                      // json
	Draw(rotation uint32) *sdl.Surface // face can be 0,1,2,3 (i.e. 0, 90, 180, 270)
	Power() float64                    // ampere
}

type ItemInventoryArray []*supplier.InventoryItem
func (a ItemInventoryArray) Len() int { return len(a) }
func (a ItemInventoryArray) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ItemInventoryArray) Less(i, j int) bool { return a[i].Zplaced > a[j].Zplaced }

type RackElement struct {
	surface          *sdl.Surface
	items            []*supplier.InventoryItem
	previousrotation uint32
}

func (self *RackElement) AddItem(item *supplier.InventoryItem) {
	self.items=append(self.items,item)
	sort.Sort(ItemInventoryArray(self.items))
	self.surface=nil
}

func (self *RackElement) RemoveItem(item *supplier.InventoryItem) {
	for p,i := range(self.items) {
		if i==item {
			self.items=append(self.items[:p],self.items[p+1:]...)
			self.surface=nil
		}
	}
}

func (self *RackElement) ElementType() string {
	return "rack"
}

func (self *RackElement) Save() string {
	return ""
}

func (self *RackElement) Draw(rotation uint32) *sdl.Surface {
	if (self.surface != nil) && (self.previousrotation != rotation) {
		self.surface.Free()
		self.surface = nil
	}
	if self.surface == nil {
		self.surface = getSprite("resources/rack.bottom" + strconv.Itoa(int(rotation)) + ".png")
		for _,item := range (self.items) {
			img := getSprite("resources/" + item.Serverconf.ConfType.ServerSprite + strconv.Itoa(int(rotation)) + ".png")
			rectSrc := sdl.Rect{0, 0, img.W, img.H}
			rectDst := sdl.Rect{0, TILE_HEIGHT - img.H - ((42-item.Zplaced)+item.Serverconf.ConfType.NbU+1)*4, img.W, img.H}
			img.Blit(&rectSrc, self.surface, &rectDst)
		}

		top := getSprite("resources/rack.top" + strconv.Itoa(int(rotation)) + ".png")
		rectSrc := sdl.Rect{0, 0, top.W, top.H}
		rectDst := sdl.Rect{0, TILE_HEIGHT - top.H, top.W, top.H}
		top.Blit(&rectSrc, self.surface, &rectDst)
		self.previousrotation = rotation
	}
	return self.surface
}

func (self *RackElement) Power() float64 {
	power := float64(0)
/*
	for _, e := range self.rackmount {
		power += e.power
	}
*/
	return power
}

func CreateRackElement() *RackElement {
	r := &RackElement{
		surface:   nil,
		items: make([]*supplier.InventoryItem,0),
	}
	return r
}

type ElectricalElement struct {
	flavor           string  // ac, battery, generatorA,generatorB
	power            float64 // negative if it is a generator
	capacity         int32   // kWh if it is a battery
	surface          *sdl.Surface
	previousrotation uint32
}

func (self *ElectricalElement) ElementType() string {
	return self.flavor
}

func (self *ElectricalElement) Save() string {
	return ""
}

func (self *ElectricalElement) Draw(rotation uint32) *sdl.Surface {
	if rotation != self.previousrotation && self.surface != nil {
		self.surface.Free()
		self.surface = nil
	}
	if self.surface == nil {
		self.surface = getSprite("resources/" + self.flavor + strconv.Itoa(int(rotation)) + ".png")
		self.previousrotation = rotation
	}
	return self.surface
}

func (self *ElectricalElement) Power() float64 {
	return self.power
}

func CreateElectricalElement(dcelementtype string) *ElectricalElement {
	//power := payload["power"].(float64)
	//capacity := int32(payload["capacity"].(float64))
	ee := &ElectricalElement{
		flavor:           dcelementtype,
		power:            0,
		capacity:         0,
		surface:          nil,
		previousrotation: 0,
	}
	return ee
}

type Tile struct {
	wall     [2]string // "" when nothing
	floor    string
	element  DcElement
	surface  *sdl.Surface
	rotation uint32 // rotation of the inner element: floor+element (not the walls)
}

func (self *Tile) ItemInstalled(item *supplier.InventoryItem) {
	if item.Typeitem==supplier.PRODUCT_SERVER && item.Serverconf.ConfType.NbU>0 {
		if self.element!=nil && self.element.ElementType()=="rack" {
			rack:=self.element.(*RackElement)
			rack.AddItem(item)
			self.surface=nil
		}
	}
}

func (self *Tile) ItemUninstalled(item *supplier.InventoryItem) {
	if item.Typeitem==supplier.PRODUCT_SERVER && item.Serverconf.ConfType.NbU>0 {
		if self.element!=nil && self.element.ElementType()=="rack" {
			rack:=self.element.(*RackElement)
			rack.RemoveItem(item)
			self.surface=nil
		}
	}
}

func (self *Tile) DcElement() DcElement {
	return self.element
}

func (self *Tile) IsElementAt(x, y int32) bool {
	if self.element == nil {
		return false
	}
	elt := self.element.Draw(self.rotation)
	y -= TILE_HEIGHT - elt.H
	if (x < 0) || (y < 0) || (x >= elt.W) || (y > elt.H) {
		return false
	}
	_, _, _, alpha := GetSurfacePixel(elt, x, y)
	if alpha > 0 {
		return true
	}
	return false
}

func (self *Tile) Draw() *sdl.Surface {
	if self.surface == nil {
		self.surface, _ = sdl.CreateRGBSurface(0, 105, TILE_HEIGHT, 32, 0x00ff0000, 0x0000ff00, 0x000000ff, 0xff000000)
		floor := getSprite("resources/" + self.floor + strconv.Itoa(int(self.rotation)) + ".png")
		rectSrc := sdl.Rect{0, 0, floor.W, floor.H}
		rectDst := sdl.Rect{0, TILE_HEIGHT - floor.H, floor.W, floor.H}
		floor.Blit(&rectSrc, self.surface, &rectDst)
		if (self.wall[0] != "") || (self.wall[1] != "") {
			wall := getSprite("resources/wallX.png")
			rectSrc := sdl.Rect{0, 0, wall.W, wall.H}
			rectDst := sdl.Rect{0, TILE_HEIGHT - wall.H, wall.W, wall.H}
			wall.Blit(&rectSrc, self.surface, &rectDst)
		}
		if self.wall[0] != "" {
			wall := getSprite("resources/" + self.wall[0] + "L.png")
			rectSrc := sdl.Rect{0, 0, wall.W, wall.H}
			rectDst := sdl.Rect{0, TILE_HEIGHT - wall.H, wall.W, wall.H}
			wall.Blit(&rectSrc, self.surface, &rectDst)
		}
		if self.wall[1] != "" {
			wall := getSprite("resources/" + self.wall[1] + "R.png")
			rectSrc := sdl.Rect{0, 0, wall.W, wall.H}
			rectDst := sdl.Rect{0, TILE_HEIGHT - wall.H, wall.W, wall.H}
			wall.Blit(&rectSrc, self.surface, &rectDst)
		}

		if self.element != nil {
			elt := self.element.Draw(self.rotation)
			rectSrc := sdl.Rect{0, 0, elt.W, elt.H}
			rectDst := sdl.Rect{0, TILE_HEIGHT - elt.H, elt.W, elt.H}
			elt.Blit(&rectSrc, self.surface, &rectDst)
		}
	}
	return self.surface
}

func (self *Tile) Rotate(rotation uint32) {
	if self.surface != nil {
		self.surface.Free()
	}
	self.surface = nil
	self.rotation = rotation
}

func (self *Tile) Power() float64 {
	return self.element.Power()
}

func CreateGrassTile() *Tile {
	tile := &Tile{
		wall:     [2]string{"", ""},
		rotation: 0,
		floor:    "green",
		element:  nil,
	}
	return tile
}

func CreateElectricalTile(wall0, wall1, floor string, rotation uint32, dcelementtype string) *Tile {
	if rotation > 3 {
		rotation = 0
	}
	var element DcElement
	if dcelementtype == "rack" {
		element = CreateRackElement()
	} else if dcelementtype != "" {
		element = CreateElectricalElement(dcelementtype)
	}
	tile := &Tile{
		wall:     [2]string{wall0, wall1},
		rotation: rotation,
		floor:    floor,
		element:  element,
	}
	return tile
}

var spritecache map[string]*sdl.Surface

func getSprite(image string) *sdl.Surface {
	if spritecache == nil {
		spritecache = make(map[string]*sdl.Surface)
	}
	sprite := spritecache[image]
	if sprite == nil {
		var err error
		sprite, err = img.Load(image)
		if sprite == nil || err != nil {
			fmt.Println("Error loading ", image, err)
			panic(err)
		}
		spritecache[image] = sprite
	}
	return sprite
}
