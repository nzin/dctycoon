package dctycoon

import (
	"sort"
	"strconv"
	"strings"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/supplier"
	log "github.com/sirupsen/logrus"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	TILE_WIDTH_STEP  = 96
	TILE_HEIGHT_STEP = 48
	TILE_HEIGHT      = 257
)

// base "class" for all tiles:
// - RackElement
// - SimpleElement
// - DecorationElement
type TileElement interface {
	// should be passive, rack, ...
	ElementType() int32                             // which type sit on: PRODUCT_RACK, PRODUCT_AC, ...
	InventoryItem() *supplier.InventoryItem         // if there is one
	Draw(rotation, flasheffect uint32) *sdl.Surface // face can be 0,1,2,3 (i.e. 0, 90, 180, 270)
	Power() float64                                 // ampere
}

type ItemInventoryArray []*supplier.InventoryItem

func (a ItemInventoryArray) Len() int           { return len(a) }
func (a ItemInventoryArray) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ItemInventoryArray) Less(i, j int) bool { return a[i].Zplaced > a[j].Zplaced }

type RackElement struct {
	surface          *sdl.Surface
	item             *supplier.InventoryItem
	items            []*supplier.InventoryItem
	previousrotation uint32
	previousflash    uint32
}

// InventoryItem return the Rack item itself
func (self *RackElement) InventoryItem() *supplier.InventoryItem {
	return self.item
}

// GetRackServers return the racked servers inside the rack
func (self *RackElement) GetRackServers() []*supplier.InventoryItem {
	return self.items
}

func (self *RackElement) AddItem(item *supplier.InventoryItem) {
	self.items = append(self.items, item)
	sort.Sort(ItemInventoryArray(self.items))
	self.surface = nil
}

func (self *RackElement) RemoveItem(item *supplier.InventoryItem) {
	for p, i := range self.items {
		if i == item {
			self.items = append(self.items[:p], self.items[p+1:]...)
			self.surface = nil
		}
	}
}

func (self *RackElement) ElementType() int32 {
	return supplier.PRODUCT_RACK
}

func (self *RackElement) Draw(rotation, flasheffect uint32) *sdl.Surface {
	if (self.surface != nil) && (self.previousrotation != rotation || self.previousflash != flasheffect) {
		self.surface.Free()
		self.surface = nil
	}
	if self.surface == nil {
		var err error
		bottom := getSprite("assets/ui/rack.bottom" + strconv.Itoa(int(rotation)) + ".png")
		self.surface, err = sdl.CreateRGBSurface(0, bottom.W, bottom.H, 32, 0x00ff0000, 0x0000ff00, 0x000000ff, 0xff000000)
		if err != nil {
			panic(err)
		}
		rectSrc := sdl.Rect{0, 0, bottom.W, bottom.H}
		rectDst := sdl.Rect{0, 0, bottom.W, bottom.H}
		bottom.Blit(&rectSrc, self.surface, &rectDst)

		inside, _ := sdl.CreateRGBSurface(0, bottom.W, bottom.H, 32, 0x00ff0000, 0x0000ff00, 0x000000ff, 0xff000000)
		for _, item := range self.items {
			img := getSprite("assets/ui/" + item.Serverconf.ConfType.ServerSprite + strconv.Itoa(int(rotation)) + ".png")
			rectSrc := sdl.Rect{0, 0, img.W, img.H}
			rectDst := sdl.Rect{0, TILE_HEIGHT - img.H - ((42-item.Zplaced)-item.Serverconf.ConfType.NbU+2)*4, img.W, img.H}
			img.Blit(&rectSrc, inside, &rectDst)
		}

		// do we have to flash the image
		if flasheffect != 0 {
			global.FlashImage(inside, flasheffect)
		}

		rectSrc = sdl.Rect{0, 0, inside.W, inside.H}
		rectDst = sdl.Rect{0, 0, inside.W, inside.H}
		inside.Blit(&rectSrc, self.surface, &rectDst)

		top := getSprite("assets/ui/rack.top" + strconv.Itoa(int(rotation)) + ".png")
		rectSrc = sdl.Rect{0, 0, top.W, top.H}
		rectDst = sdl.Rect{0, TILE_HEIGHT - top.H, top.W, top.H}
		top.Blit(&rectSrc, self.surface, &rectDst)
		self.previousrotation = rotation
	}
	return self.surface
}

func (self *RackElement) Power() float64 {
	power := float64(0)
	for _, i := range self.items {
		power += i.Serverconf.PowerConsumption()
	}
	return power
}

func NewRackElement(item *supplier.InventoryItem) *RackElement {
	r := &RackElement{
		surface: nil,
		item:    item,
		items:   make([]*supplier.InventoryItem, 0),
	}
	return r
}

type SimpleElement struct {
	inventoryitem    *supplier.InventoryItem // ac, battery, generator, tower
	capacity         int32                   // kWh if it is a battery
	surface          *sdl.Surface
	previousrotation uint32
}

func (self *SimpleElement) InventoryItem() *supplier.InventoryItem {
	return self.inventoryitem
}

func (self *SimpleElement) ElementType() int32 {
	return self.inventoryitem.Typeitem
}

func (self *SimpleElement) Draw(rotation, flasheffect uint32) *sdl.Surface {
	if rotation != self.previousrotation && self.surface != nil {
		self.surface.Free()
		self.surface = nil
	}
	if self.surface == nil {
		var err error
		itemsprite := getSprite("assets/ui/" + self.inventoryitem.GetSprite() + strconv.Itoa(int(rotation)) + ".png")
		self.surface, err = sdl.CreateRGBSurface(0, itemsprite.W, itemsprite.H, 32, 0x00ff0000, 0x0000ff00, 0x000000ff, 0xff000000)
		if err != nil {
			panic(err)
		}
		rectSrc := sdl.Rect{0, 0, itemsprite.W, itemsprite.H}
		rectDst := sdl.Rect{0, 0, itemsprite.W, itemsprite.H}
		itemsprite.Blit(&rectSrc, self.surface, &rectDst)
		self.previousrotation = rotation
	}
	return self.surface
}

func (self *SimpleElement) Power() float64 {
	if self.inventoryitem.Typeitem == supplier.PRODUCT_GENERATOR {
		return -50000
	}
	if self.inventoryitem.Typeitem == supplier.PRODUCT_SERVER {
		return self.inventoryitem.Serverconf.PowerConsumption()
	}
	return 0
}

func NewSimpleElement(item *supplier.InventoryItem) *SimpleElement {
	//power := payload["power"].(float64) // will depend on "flavor"
	//capacity := int32(payload["capacity"].(float64)) // will depend on "flavor"
	ee := &SimpleElement{
		inventoryitem:    item,
		capacity:         0,
		surface:          nil,
		previousrotation: 0,
	}
	return ee
}

type DecorationElement struct {
	name    string
	surface *sdl.Surface
}

func (self *DecorationElement) GetName() string {
	return self.name
}

func (self *DecorationElement) InventoryItem() *supplier.InventoryItem {
	return nil
}

func (self *DecorationElement) ElementType() int32 {
	return supplier.PRODUCT_DECORATION
}

func (self *DecorationElement) Draw(rotation, flasheffect uint32) *sdl.Surface {
	if self.surface == nil {
		self.surface = getSprite("assets/ui/" + self.name + ".png")
	}
	return self.surface
}

func (self *DecorationElement) Power() float64 {
	return 0
}

func NewDecorationElement(name string) *DecorationElement {
	decoration := &DecorationElement{
		name:    name,
		surface: nil,
	}
	return decoration
}

type Tile struct {
	wall               [2]string // "" when nothing
	floor              string
	element            TileElement // either RackElement or SimpleElement
	surface            *sdl.Surface
	surfaceWithoutWall *sdl.Surface
	rotation           uint32 // rotation of the inner element: floor+element (not the walls)
	flasheffect        uint32
}

func (self *Tile) ItemInstalled(item *supplier.InventoryItem) {
	if item.Typeitem == supplier.PRODUCT_SERVER && item.Serverconf.ConfType.NbU > 0 {
		if self.element != nil && self.element.ElementType() == supplier.PRODUCT_RACK {
			rack := self.element.(*RackElement)
			rack.AddItem(item)
			self.freeSurface()
		}
	}
	// tower
	if item.Typeitem == supplier.PRODUCT_SERVER && item.Serverconf.ConfType.NbU <= 0 {
		self.element = NewSimpleElement(item)
		self.freeSurface()
	}

	if item.Typeitem == supplier.PRODUCT_RACK && self.element == nil {
		self.element = NewRackElement(item)
		self.freeSurface()
	}
	if item.Typeitem != supplier.PRODUCT_RACK && item.Typeitem != supplier.PRODUCT_SERVER {
		self.element = NewSimpleElement(item)
		self.freeSurface()
	}
}

func (self *Tile) ItemUninstalled(item *supplier.InventoryItem) {
	log.Debug("Tile::ItemUninstalled(", item, ")")
	// if we are removing a racked server
	if item.Typeitem == supplier.PRODUCT_SERVER && item.Serverconf.ConfType.NbU > 0 {
		if self.element != nil && self.element.ElementType() == supplier.PRODUCT_RACK {
			rack := self.element.(*RackElement)
			rack.RemoveItem(item)
			self.freeSurface()
		}
	}
	// if we are removing a tower server
	if item.Typeitem == supplier.PRODUCT_SERVER && item.Serverconf.ConfType.NbU == -1 {
		self.element = nil
		self.freeSurface()
	}
	// the rack here is supposed to be empty
	if item.Typeitem == supplier.PRODUCT_RACK && self.element != nil {
		self.element = nil
		self.freeSurface()
	}
	if item.Typeitem != supplier.PRODUCT_RACK && item.Typeitem != supplier.PRODUCT_SERVER {
		self.element = nil
		self.freeSurface()
	}
}

func (self *Tile) TileElement() TileElement {
	return self.element
}

// IsFloorOutside used to know if we can place a AC on it
func (self *Tile) IsFloorOutside() bool {
	return strings.HasPrefix(self.floor, "green")
}

// IsFloorInsideNotAir used to know if we are on a server tile
// but not on a air flow to install anything
func (self *Tile) IsFloorInsideNotAirFlow() bool {
	return self.floor == "inside"
}

// IsFloorInsideAir used to know if we are on a air flow tile
func (self *Tile) IsFloorInsideAirFlow() bool {
	return self.floor == "inside.air"
}

func (self *Tile) SwitchToAirFlow() {
	if self.floor == "inside" {
		self.floor = "inside.air"
		self.freeSurface()
	}
}

func (self *Tile) SwitchToNotAirFlow() {
	if self.floor == "inside.air" {
		self.floor = "inside"
		self.freeSurface()
	}
}

func (self *Tile) freeSurface() {
	if self.surface != nil {
		self.surface.Free()
	}
	if self.surfaceWithoutWall != nil {
		self.surfaceWithoutWall.Free()
	}
	self.surface = nil
	self.surfaceWithoutWall = nil
}

func (self *Tile) Draw() *sdl.Surface {
	if self.surface == nil {
		self.surface = self.draw(true)
	}
	return self.surface
}

func (self *Tile) DrawWithoutWall() *sdl.Surface {
	if self.surfaceWithoutWall == nil {
		self.surfaceWithoutWall = self.draw(false)
	}
	return self.surfaceWithoutWall
}

func (self *Tile) draw(withWall bool) *sdl.Surface {
	surface, _ := sdl.CreateRGBSurface(0, 105, TILE_HEIGHT, 32, 0x00ff0000, 0x0000ff00, 0x000000ff, 0xff000000)
	floor := getSprite("assets/ui/" + self.floor + strconv.Itoa(int(self.rotation)) + ".png")
	rectSrc := sdl.Rect{0, 0, floor.W, floor.H}
	rectDst := sdl.Rect{0, TILE_HEIGHT - floor.H, floor.W, floor.H}
	floor.Blit(&rectSrc, surface, &rectDst)
	if withWall {
		if (self.wall[0] != "") || (self.wall[1] != "") {
			wall := getSprite("assets/ui/wallX.png")
			rectSrc := sdl.Rect{0, 0, wall.W, wall.H}
			rectDst := sdl.Rect{0, TILE_HEIGHT - wall.H, wall.W, wall.H}
			wall.Blit(&rectSrc, surface, &rectDst)
		}
		if self.wall[0] != "" {
			wall := getSprite("assets/ui/" + self.wall[0] + "L.png")
			rectSrc := sdl.Rect{0, 0, wall.W, wall.H}
			rectDst := sdl.Rect{0, TILE_HEIGHT - wall.H, wall.W, wall.H}
			wall.Blit(&rectSrc, surface, &rectDst)
		}
		if self.wall[1] != "" {
			wall := getSprite("assets/ui/" + self.wall[1] + "R.png")
			rectSrc := sdl.Rect{0, 0, wall.W, wall.H}
			rectDst := sdl.Rect{0, TILE_HEIGHT - wall.H, wall.W, wall.H}
			wall.Blit(&rectSrc, surface, &rectDst)
		}
	}

	if self.element != nil {
		elt := self.element.Draw(self.rotation, self.flasheffect)
		rectSrc := sdl.Rect{0, 0, elt.W, elt.H}
		rectDst := sdl.Rect{0, TILE_HEIGHT - elt.H, elt.W, elt.H}
		elt.Blit(&rectSrc, surface, &rectDst)
	}
	return surface
}

//
// Rotate is used to rotate the tile by 90 degrees
// you can call Draw() after calling this function
func (self *Tile) Rotate(rotation uint32) {
	self.freeSurface()
	self.rotation = rotation
}

//
// SetFlashEffect is used to "flash/brighten" the rack content (the racked servers)
// flash is a value between 8 (completely white) and 0 (normal image).
// see global.FlashImage()
// you can call Draw() after calling this function
func (self *Tile) SetFlashEffect(flash uint32) {
	self.freeSurface()
	self.flasheffect = flash
}

func (self *Tile) Power() float64 {
	if self.element == nil {
		return 0
	}
	return self.element.Power()
}

func NewGrassTile() *Tile {
	tile := &Tile{
		wall:        [2]string{"", ""},
		rotation:    0,
		flasheffect: 0,
		floor:       "green",
		element:     nil,
	}
	return tile
}

func NewTile(wall0, wall1, floor string, rotation uint32, decorationname string) *Tile {
	if rotation > 3 {
		rotation = 0
	}
	var element TileElement
	if decorationname != "" {
		element = NewDecorationElement(decorationname)
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

//
// getSprite is used to cache loaded images assets
func getSprite(image string) *sdl.Surface {
	if spritecache == nil {
		spritecache = make(map[string]*sdl.Surface)
	}
	sprite := spritecache[image]
	if sprite == nil {
		var err error
		//		sprite, err = img.Load(image)
		sprite, err = global.LoadImageAsset(image)
		if sprite == nil || err != nil {
			log.Error("Error loading ", image, err)
			panic(err)
		}
		spritecache[image] = sprite
	}
	return sprite
}
