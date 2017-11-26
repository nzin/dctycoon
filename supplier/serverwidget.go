package supplier

import (
	"fmt"
	"strings"

	"github.com/nzin/sws"
	//	"github.com/veandco/go-sdl2/sdl"
)

type InventoryLineWidget struct {
	sws.CoreWidget
	Checkbox  *sws.CheckboxWidget
	desc      *sws.LabelWidget
	placement *sws.LabelWidget
	item      *InventoryItem
}

func NewInventoryLineWidget(item *InventoryItem) *InventoryLineWidget {
	ramSizeText := fmt.Sprintf("%d Mo", item.Serverconf.NbSlotRam*item.Serverconf.RamSize)
	if item.Serverconf.NbSlotRam*item.Serverconf.RamSize >= 2048 {
		ramSizeText = fmt.Sprintf("%d Go", item.Serverconf.NbSlotRam*item.Serverconf.RamSize/1024)
	}
	text := item.Serverconf.ConfType.ServerName
	placement := " - "
	if item.Xplaced != -1 {
		placement = fmt.Sprintf("%d/%d", item.Xplaced, item.Yplaced)
	}
	line := &InventoryLineWidget{
		CoreWidget: *sws.NewCoreWidget(625, 25),
		Checkbox:   sws.NewCheckboxWidget(),
		desc:       sws.NewLabelWidget(200, 25, text),
		placement:  sws.NewLabelWidget(100, 25, placement),
		item:       item,
	}
	line.Checkbox.SetColor(0xffffffff)
	line.AddChild(line.Checkbox)

	line.desc.SetColor(0xffffffff)
	line.desc.Move(25, 0)
	line.AddChild(line.desc)

	line.placement.SetColor(0xffffffff)
	line.placement.Move(225, 0)
	line.AddChild(line.placement)

	cores := sws.NewLabelWidget(100, 25, fmt.Sprintf("%d", item.Serverconf.NbProcessors*item.Serverconf.NbCore))
	cores.SetColor(0xffffffff)
	cores.Move(325, 0)
	line.AddChild(cores)

	ram := sws.NewLabelWidget(100, 25, ramSizeText)
	ram.SetColor(0xffffffff)
	ram.Move(425, 0)
	line.AddChild(ram)

	diskText := fmt.Sprintf("%d Mo", item.Serverconf.NbDisks*item.Serverconf.DiskSize)
	if item.Serverconf.NbDisks*item.Serverconf.DiskSize > 4096 {
		diskText = fmt.Sprintf("%d Go", item.Serverconf.NbDisks*item.Serverconf.DiskSize/1024)
	}
	if item.Serverconf.NbDisks*item.Serverconf.DiskSize > 4*1024*1024 {
		diskText = fmt.Sprintf("%d To", item.Serverconf.NbDisks*item.Serverconf.DiskSize/(1024*1024))
	}
	disk := sws.NewLabelWidget(100, 25, diskText)
	disk.SetColor(0xffffffff)
	disk.Move(525, 0)
	line.AddChild(disk)

	return line
}

const (
	ASSIGNED_UNASSIGNED = 0
	ASSIGNED_PHYSICAL   = 1
	ASSIGNED_VPS        = 2
)

type ServerFilter struct {
	assigned *int32
	racked   *bool
	inuse    *bool
}

type ServerWidget struct {
	sws.CoreWidget
	inventory              *Inventory
	instock                []*InventoryItem
	searchUnassignedButton *sws.ButtonWidget
	searchPhysicalButton   *sws.ButtonWidget
	searchVpsButton        *sws.ButtonWidget
	searchbar              *sws.InputWidget
	selected               map[*InventoryLineWidget]bool
	selectallButton        *sws.CheckboxWidget
	listing                *sws.VBoxWidget
	scrolllisting          *sws.ScrollWidget
	currentFilter          ServerFilter
}

func (self *ServerWidget) SelectLine(line *InventoryLineWidget, selected bool) {
	if selected {
		if len(self.selected) == 0 {
			// show items related action button
		}
		self.selected[line] = true
	} else {
		delete(self.selected, line)
		if len(self.selected) == 0 {
			// remove items related action button
		}
	}
}

//
// InventorySubscriber interface
//
func (self *ServerWidget) ItemInTransit(*InventoryItem) {
}

func (self *ServerWidget) ItemInStock(item *InventoryItem) {
	if item.Typeitem != PRODUCT_SERVER {
		return
	}
	// we don't want to add it twice
	for i, c := range self.instock {
		if c == item {
			if i == 0 {
				self.instock = self.instock[1:]
			} else {
				self.instock = append(self.instock[:i], self.instock[i+1:]...)
			}
		}
	}
	self.instock = append(self.instock, item)
	// add to the listing
	if self.searchFilter(item) {
		line := NewInventoryLineWidget(item)
		line.Checkbox.SetClicked(func() {
			self.SelectLine(line, line.Checkbox.Selected)
			self.selectallButton.SetSelected(self.selectallButton.Selected)
		})
		self.listing.AddChild(line)
	}
}

func (self *ServerWidget) ItemRemoveFromStock(item *InventoryItem) {
	if item.Typeitem != PRODUCT_SERVER {
		return
	}
	// remove from the listing
	for _, l := range self.listing.GetChildren() {
		line := l.(*InventoryLineWidget)
		if line.item == item {
			self.listing.RemoveChild(line)
		}
	}

	for i, c := range self.instock {
		if c == item {
			if i == 0 {
				self.instock = self.instock[1:]
			} else {
				self.instock = append(self.instock[:i], self.instock[i+1:]...)
			}
			return
		}
	}
}

func (self *ServerWidget) ItemInstalled(*InventoryItem) {
}

func (self *ServerWidget) ItemUninstalled(*InventoryItem) {
}

func (self *ServerWidget) searchFilter(item *InventoryItem) bool {
	// assigned = [unassigned|physical|vps]
	if self.currentFilter.assigned != nil {
		switch *self.currentFilter.assigned {
		case ASSIGNED_UNASSIGNED:
			if item.pool != nil {
				return false
			}
		case ASSIGNED_PHYSICAL:
			if item.pool == nil || item.pool.IsVps() == true {
				return false
			}
		case ASSIGNED_VPS:
			if item.pool == nil || item.pool.IsVps() == false {
				return false
			}
		}
	}

	// racked = [true|false]
	if self.currentFilter.racked != nil {
		switch *self.currentFilter.racked {
		case true:
			if item.Xplaced != -1 {
				return false
			}
		case false:
			if item.Xplaced == -1 {
				return false
			}
		}
	}

	return true
}

func (self *ServerWidget) Search(search string) {
	self.searchbar.SetText(search)

	tokens := strings.Fields(search)
	var filter ServerFilter
	error := false
	for _, token := range tokens {
		if strings.Contains(token, ":") {
			s := strings.SplitN(token, ":", 2)
			key := s[0]
			value := s[1]

			switch key {
			case "assigned":
				switch value {
				case "unassigned":
					var assigned int32 = ASSIGNED_UNASSIGNED
					filter.assigned = &assigned
				case "physical":
					var assigned int32 = ASSIGNED_PHYSICAL
					filter.assigned = &assigned
				case "vps":
					var assigned int32 = ASSIGNED_VPS
					filter.assigned = &assigned
				default:
					error = true
				}
			case "racked":
				switch value {
				case "true":
					var racked bool = true
					filter.racked = &racked
				case "false":
					var racked bool = false
					filter.racked = &racked
				default:
					error = true
				}
			case "inuse":
				switch value {
				case "true":
					var inuse bool = true
					filter.inuse = &inuse
				case "false":
					var inuse bool = false
					filter.inuse = &inuse
				default:
					error = true
				}
			default:
				error = true
			}
		} else {
			error = true
		}
	}

	if error == true {
		self.searchbar.SetColor(0xffff5555)
		return
	}
	self.currentFilter = filter
	self.searchbar.SetColor(0xffffffff)

	self.listing.RemoveAllChildren()
	self.selected = make(map[*InventoryLineWidget]bool)

	for _, c := range self.instock {
		if self.searchFilter(c) == true {
			line := NewInventoryLineWidget(c)
			line.Checkbox.SetClicked(func() {
				self.SelectLine(line, line.Checkbox.Selected)
				self.selectallButton.SetSelected(false)
			})
			self.listing.AddChild(line)
		}
	}
	self.selectallButton.SetSelected(false)
	self.PostUpdate()
}

func (self *ServerWidget) Resize(width, height int32) {
	self.CoreWidget.Resize(width, height)
	if height > 150 {
		self.scrolllisting.Resize(width, height-125)
	}
}

func NewServerWidget(inventory *Inventory) *ServerWidget {
	corewidget := sws.NewCoreWidget(600, 400)
	widget := &ServerWidget{
		CoreWidget:             *corewidget,
		inventory:              inventory,
		instock:                make([]*InventoryItem, 0, 0),
		searchUnassignedButton: sws.NewButtonWidget(150, 50, "Arrival"),
		searchPhysicalButton:   sws.NewButtonWidget(150, 50, "Physical pool"),
		searchVpsButton:        sws.NewButtonWidget(150, 50, "Vps pool"),
		searchbar:              sws.NewInputWidget(500, 25, "assigned:unassigned"),
		selected:               make(map[*InventoryLineWidget]bool),
		selectallButton:        sws.NewCheckboxWidget(),
		listing:                sws.NewVBoxWidget(600, 10),
		scrolllisting:          sws.NewScrollWidget(600, 400),
	}
	var assigned int32 = ASSIGNED_UNASSIGNED
	widget.currentFilter.assigned = &assigned

	inventory.AddInventorySubscriber(widget)

	widget.searchUnassignedButton.SetClicked(func() {
		widget.Search("assigned:unassigned")
	})
	widget.searchUnassignedButton.Move(10, 5)
	widget.AddChild(widget.searchUnassignedButton)

	widget.searchPhysicalButton.SetClicked(func() {
		widget.Search("assigned:physical")
	})
	widget.searchPhysicalButton.Move(170, 5)
	widget.AddChild(widget.searchPhysicalButton)

	widget.searchVpsButton.SetClicked(func() {
		widget.Search("assigned:vps")
	})
	widget.searchVpsButton.Move(330, 5)
	widget.AddChild(widget.searchVpsButton)

	widget.searchbar.SetEnterCallback(func() {
		widget.Search(widget.searchbar.GetText())
	})

	widget.searchbar.Move(10, 60)
	widget.AddChild(widget.searchbar)

	// description line
	widget.selectallButton.Move(0, 100)
	widget.AddChild(widget.selectallButton)

	globaldesc := sws.NewLabelWidget(200, 25, "Description")
	globaldesc.Move(25, 100)
	widget.AddChild(globaldesc)

	globalplacement := sws.NewLabelWidget(100, 25, "Placement")
	globalplacement.Move(225, 100)
	widget.AddChild(globalplacement)

	globalnbcores := sws.NewLabelWidget(100, 25, "Nb cores")
	globalnbcores.Move(325, 100)
	widget.AddChild(globalnbcores)

	globalram := sws.NewLabelWidget(100, 25, "RAM")
	globalram.Move(425, 100)
	widget.AddChild(globalram)

	globaldisk := sws.NewLabelWidget(100, 25, "Disk")
	globaldisk.Move(525, 100)
	widget.AddChild(globaldisk)

	widget.scrolllisting.Move(0, 125)
	widget.scrolllisting.ShowHorizontalScrollbar(false)
	widget.scrolllisting.SetInnerWidget(widget.listing)
	widget.AddChild(widget.scrolllisting)

	return widget
}
