package supplier

/*
 * Pool management widget (gmail style) to be able to attribute servers to Physical(s)/VPS(s) pools.
 * see MainInventoryWiget
 */
import (
	"fmt"
	"strings"
	"time"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/ui"
	"github.com/nzin/sws"
)

type PoolManagementLineWidget struct {
	sws.CoreWidget
	Checkbox  *sws.CheckboxWidget
	desc      *sws.LabelWidget
	placement *sws.LabelWidget
	cores     *sws.LabelWidget
	ram       *sws.LabelWidget
	disk      *sws.LabelWidget
	allocated *sws.LabelWidget
	item      *InventoryItem
}

func NewPoolManagementLineWidget(item *InventoryItem) *PoolManagementLineWidget {
	text := item.Serverconf.ConfType.ServerName
	placement := " - "
	if item.Xplaced != -1 {
		placement = fmt.Sprintf("%d/%d", item.Xplaced, item.Yplaced)
	}

	allocated := "-"
	if item.Coresallocated != 0 {
		allocated = "yes"
	}

	line := &PoolManagementLineWidget{
		CoreWidget: *sws.NewCoreWidget(725, 25),
		Checkbox:   sws.NewCheckboxWidget(),
		desc:       sws.NewLabelWidget(200, 25, text),
		placement:  sws.NewLabelWidget(100, 25, placement),
		cores:      sws.NewLabelWidget(100, 25, fmt.Sprintf("%d", item.Serverconf.NbProcessors*item.Serverconf.NbCore)),
		ram:        sws.NewLabelWidget(100, 25, global.AdjustMega(item.Serverconf.NbSlotRam*item.Serverconf.RamSize)),
		disk:       sws.NewLabelWidget(100, 25, global.AdjustMega(item.Serverconf.NbDisks*item.Serverconf.DiskSize)),
		allocated:  sws.NewLabelWidget(100, 25, allocated),
		item:       item,
	}
	line.AddChild(line.Checkbox)

	line.desc.Move(25, 0)
	line.AddChild(line.desc)

	line.placement.Move(225, 0)
	line.AddChild(line.placement)

	line.cores.Move(325, 0)
	line.AddChild(line.cores)

	line.ram.Move(425, 0)
	line.AddChild(line.ram)

	line.disk.Move(525, 0)
	line.AddChild(line.disk)

	line.allocated.Move(625, 0)
	line.AddChild(line.allocated)

	line.UpdateBgColor()

	return line
}

//
// Update the bg color depending on the pool the item belongs to
//
func (self *PoolManagementLineWidget) UpdateBgColor() {
	bgcolor := uint32(0xffffffff)
	if self.item.Pool != nil {
		if self.item.Pool.IsVps() {
			bgcolor = global.VPS_COLOR
		} else {
			bgcolor = global.PHYSICAL_COLOR
		}
	}
	self.Checkbox.SetColor(bgcolor)
	self.desc.SetColor(bgcolor)
	self.placement.SetColor(bgcolor)
	self.cores.SetColor(bgcolor)
	self.ram.SetColor(bgcolor)
	self.disk.SetColor(bgcolor)
	self.allocated.SetColor(bgcolor)
}

func (self *PoolManagementLineWidget) UpdateStatus() {
	placement := " - "
	if self.item.Xplaced != -1 {
		placement = fmt.Sprintf("%d/%d", self.item.Xplaced, self.item.Yplaced)
	}
	self.placement.SetText(placement)

	allocated := "-"
	if self.item.Coresallocated != 0 {
		allocated = "yes"
	}

	self.allocated.SetText(allocated)
}

func (self *PoolManagementLineWidget) AddChild(child sws.Widget) {
	self.CoreWidget.AddChild(child)
	child.SetParent(self)
}

func (self *PoolManagementLineWidget) MousePressDown(x, y int32, button uint8) {
	self.Checkbox.MousePressDown(1, 1, button)
}

func (self *PoolManagementLineWidget) MousePressUp(x, y int32, button uint8) {
	self.Checkbox.MousePressUp(1, 1, button)
}

type PoolManagementFilter struct {
	assigned  *int32
	installed *bool
	inuse     *bool
}

type PoolManagementWidget struct {
	sws.CoreWidget
	inventory              *Inventory
	root                   *sws.RootWidget
	instock                []*InventoryItem
	poolassignation        *sws.DropdownWidget
	searchUnassignedButton *sws.ButtonWidget
	searchPhysicalButton   *sws.ButtonWidget
	searchVpsButton        *sws.ButtonWidget
	searchbar              *sws.InputWidget
	selected               map[*PoolManagementLineWidget]bool
	selectallButton        *sws.CheckboxWidget
	listing                *sws.VBoxWidget
	scrolllisting          *sws.ScrollWidget
	currentFilter          PoolManagementFilter
	addToPhysical          *sws.ButtonWidget
	addToVps               *sws.ButtonWidget
	addToUnassigned        *sws.ButtonWidget
}

//
// select the line, update action buttons to show
//
func (self *PoolManagementWidget) SelectLine(line *PoolManagementLineWidget, selected bool) {
	if selected {
		line.Checkbox.SetSelected(true)
		if len(self.selected) == 0 {
			// show items related action button
		}
		self.selected[line] = true
	} else {
		line.Checkbox.SetSelected(false)
		delete(self.selected, line)
		if len(self.selected) == 0 {
			// remove items related action button
		}
	}

	// check which action is possible
	var showPhysical, showVps, showUnassigned bool
	for l, lSelected := range self.selected {
		if lSelected {
			if l.item.Pool == nil {
				showPhysical = true
				showVps = true
			} else {
				if l.item.Coresallocated == 0 {
					showUnassigned = true
					if l.item.Pool.IsVps() {
						showPhysical = true
					} else {
						showVps = true
					}
				}
			}
		}
	}
	if showPhysical {
		self.AddChild(self.addToPhysical)
	} else {
		self.RemoveChild(self.addToPhysical)
	}
	if showVps {
		self.AddChild(self.addToVps)
	} else {
		self.RemoveChild(self.addToVps)
	}
	if showUnassigned {
		self.AddChild(self.addToUnassigned)
	} else {
		self.RemoveChild(self.addToUnassigned)
	}
}

func (self *PoolManagementWidget) callbackToPool(pool ServerPool) {
	for l, lSelected := range self.selected {
		if lSelected {
			self.inventory.AssignPool(l.item, pool)

			l.UpdateBgColor()
			self.updateLineInSearch(l.item)
			self.SelectLine(l, false)
			self.selectallButton.SetSelected(false)
		}
	}
}

//
// InventorySubscriber interface
//
func (self *PoolManagementWidget) ItemInTransit(*InventoryItem) {
}

func (self *PoolManagementWidget) ItemInStock(item *InventoryItem) {
	if item.Typeitem != PRODUCT_SERVER || item.Serverconf.ConfType.Scrap == true {
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
	self.updateLineInSearch(item)
}

func (self *PoolManagementWidget) ItemRemoveFromStock(item *InventoryItem) {
	if item.Typeitem != PRODUCT_SERVER {
		return
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
	self.updateLineInSearch(item)
}

func (self *PoolManagementWidget) ItemInstalled(item *InventoryItem) {
	self.ItemInStock(item)
}

func (self *PoolManagementWidget) ItemUninstalled(item *InventoryItem) {
	self.ItemRemoveFromStock(item)
}

func (self *PoolManagementWidget) ItemChangedPool(*InventoryItem) {
}

// is part of PoolSubscriber interface
func (self *PoolManagementWidget) InventoryItemAdd(*InventoryItem) {}

// InventoryItemRemove is part of PoolSubscriber interface
func (self *PoolManagementWidget) InventoryItemRemove(*InventoryItem) {}

// InventoryItemAllocate is part of PoolSubscriber interface
func (self *PoolManagementWidget) InventoryItemAllocate(item *InventoryItem) {
	for _, l := range self.listing.GetChildren() {
		line := l.(*PoolManagementLineWidget)
		if line.item == item {
			line.UpdateStatus()
		}
	}
}

// InventoryItemRelease is part of PoolSubscriber interface
func (self *PoolManagementWidget) InventoryItemRelease(item *InventoryItem) {
	for _, l := range self.listing.GetChildren() {
		line := l.(*PoolManagementLineWidget)
		if line.item == item {
			line.UpdateStatus()
		}
	}
}

//
// this function will add/remove the item from the listing
// if the item still correspond (or not) to the search filter
func (self *PoolManagementWidget) updateLineInSearch(item *InventoryItem) {
	var foundLine *PoolManagementLineWidget
	for _, l := range self.listing.GetChildren() {
		line := l.(*PoolManagementLineWidget)
		if line.item == item {
			foundLine = line
		}
	}
	// include in the search filter?
	if self.searchFilter(item) {
		if foundLine == nil {
			line := NewPoolManagementLineWidget(item)
			line.Checkbox.SetClicked(func() {
				self.SelectLine(line, line.Checkbox.Selected)
				self.selectallButton.SetSelected(false) //self.selectallButton.Selected)
			})
			self.listing.AddChild(line)
			self.AddChild(self.scrolllisting)
		} else {
			// in case updateLineInSearch() was called because of an item installed/uninstalled
			foundLine.UpdateStatus()
		}
	} else {
		if foundLine != nil {
			self.SelectLine(foundLine, false)
			self.listing.RemoveChild(foundLine)
			if len(self.listing.GetChildren()) == 0 {
				self.RemoveChild(self.scrolllisting)
			}
		}
	}
}

func (self *PoolManagementWidget) searchFilter(item *InventoryItem) bool {
	// assigned = [unassigned|physical|vps]
	if self.currentFilter.assigned != nil {
		switch *self.currentFilter.assigned {
		case ASSIGNED_UNASSIGNED:
			if item.Pool != nil {
				return false
			}
		case ASSIGNED_PHYSICAL:
			if item.Pool == nil || item.Pool.IsVps() == true {
				return false
			}
		case ASSIGNED_VPS:
			if item.Pool == nil || item.Pool.IsVps() == false {
				return false
			}
		}
	}

	// installed = [true|false]
	if self.currentFilter.installed != nil {
		switch *self.currentFilter.installed {
		case true:
			if item.Xplaced == -1 {
				return false
			}
		case false:
			if item.Xplaced != -1 {
				return false
			}
		}
	}

	return true
}

func (self *PoolManagementWidget) Search(search string) {
	self.searchbar.SetText(search)

	tokens := strings.Fields(search)
	var filter PoolManagementFilter
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
			case "installed":
				switch value {
				case "true":
					var installed bool = true
					filter.installed = &installed
				case "false":
					var installed bool = false
					filter.installed = &installed
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
		self.searchbar.SetInnerColor(0xffff5555)
		return
	}
	self.currentFilter = filter
	self.searchbar.SetInnerColor(0xffffffff)

	self.listing.RemoveAllChildren()
	self.RemoveChild(self.scrolllisting)
	self.selected = make(map[*PoolManagementLineWidget]bool)
	self.RemoveChild(self.addToPhysical)
	self.RemoveChild(self.addToVps)
	self.RemoveChild(self.addToUnassigned)

	for _, c := range self.instock {
		if self.searchFilter(c) == true {
			line := NewPoolManagementLineWidget(c)
			line.Checkbox.SetClicked(func() {
				self.SelectLine(line, line.Checkbox.Selected)
				self.selectallButton.SetSelected(false)
			})
			self.listing.AddChild(line)
			self.AddChild(self.scrolllisting)
		}
	}
	self.selectallButton.SetSelected(false)
	self.PostUpdate()
}

func (self *PoolManagementWidget) Resize(width, height int32) {
	self.CoreWidget.Resize(width, height)
	if height > 150 {
		if width > 725 {
			width = 725
		}
		if width < 20 {
			width = 20
		}
		self.scrolllisting.Resize(width, height-150)
	}
}

func NewPoolManagementWidget(root *sws.RootWidget) *PoolManagementWidget {
	corewidget := sws.NewCoreWidget(800, 400)
	widget := &PoolManagementWidget{
		CoreWidget:             *corewidget,
		inventory:              nil,
		root:                   root,
		instock:                make([]*InventoryItem, 0, 0),
		poolassignation:        sws.NewDropdownWidget(150, 25, []string{"Unassigned", "Physical pool", "VPS pool"}),
		searchUnassignedButton: sws.NewButtonWidget(170, 40, "Unassigned"),
		searchPhysicalButton:   sws.NewButtonWidget(170, 40, "Physical pool"),
		searchVpsButton:        sws.NewButtonWidget(170, 40, "Vps pool"),
		searchbar:              sws.NewInputWidget(605, 25, "assigned:unassigned"),
		selected:               make(map[*PoolManagementLineWidget]bool),
		selectallButton:        sws.NewCheckboxWidget(),
		listing:                sws.NewVBoxWidget(625, 10),
		scrolllisting:          sws.NewScrollWidget(625, 400),
		addToPhysical:          sws.NewButtonWidget(170, 25, "> To physical pool"),
		addToVps:               sws.NewButtonWidget(170, 25, "> To vps pool"),
		addToUnassigned:        sws.NewButtonWidget(170, 25, "> Back to unassigned"),
	}
	labelAssignement := sws.NewLabelWidget(200, 25, "New servers are assigned to pool:")
	labelAssignement.Move(5, 0)
	widget.AddChild(labelAssignement)

	widget.poolassignation.Move(205, 0)
	widget.AddChild(widget.poolassignation)
	widget.poolassignation.SetCallbackValueChanged(func() {
		widget.inventory.SetDefaultPoolAllocation(widget.poolassignation.ActiveChoice)
	})

	var assigned int32 = ASSIGNED_PHYSICAL
	widget.currentFilter.assigned = &assigned

	widget.searchPhysicalButton.SetClicked(func() {
		widget.Search("assigned:physical")
	})
	widget.searchPhysicalButton.Move(10, 25)
	widget.searchPhysicalButton.SetButtonColor(global.PHYSICAL_COLOR)
	widget.AddChild(widget.searchPhysicalButton)

	widget.searchVpsButton.SetClicked(func() {
		widget.Search("assigned:vps")
	})
	widget.searchVpsButton.Move(190, 25)
	widget.searchVpsButton.SetButtonColor(global.VPS_COLOR)
	widget.AddChild(widget.searchVpsButton)

	widget.searchUnassignedButton.SetClicked(func() {
		widget.Search("assigned:unassigned")
	})
	widget.searchUnassignedButton.Move(370, 25)
	widget.AddChild(widget.searchUnassignedButton)

	widget.searchbar.SetEnterCallback(func() {
		widget.Search(widget.searchbar.GetText())
	})

	widget.searchbar.Move(10, 70)
	widget.AddChild(widget.searchbar)

	widget.addToPhysical.Move(10, 100)
	widget.addToPhysical.SetButtonColor(global.PHYSICAL_COLOR)
	widget.addToPhysical.SetCentered(false)
	widget.addToPhysical.SetClicked(func() {
		widget.callbackToPool(widget.inventory.GetDefaultPhysicalPool())
	})

	widget.addToVps.Move(190, 100)
	widget.addToVps.SetButtonColor(global.VPS_COLOR)
	widget.addToVps.SetCentered(false)
	widget.addToVps.SetClicked(func() {
		widget.callbackToPool(widget.inventory.GetDefaultVpsPool())
	})

	widget.addToUnassigned.Move(370, 100)
	widget.addToUnassigned.SetCentered(false)
	widget.addToUnassigned.SetClicked(func() {
		widget.callbackToPool(nil)
	})

	// description line
	widget.selectallButton.Move(0, 125)
	widget.AddChild(widget.selectallButton)
	widget.selectallButton.SetClicked(func() {
		state := widget.selectallButton.Selected
		for _, l := range widget.listing.GetChildren() {
			line := l.(*PoolManagementLineWidget)
			line.Checkbox.SetSelected(state)
			widget.SelectLine(line, state)
		}
	})

	globaldesc := sws.NewLabelWidget(200, 25, "Description")
	globaldesc.Move(25, 125)
	widget.AddChild(globaldesc)

	globalplacement := sws.NewLabelWidget(100, 25, "Placement")
	globalplacement.Move(225, 125)
	widget.AddChild(globalplacement)

	globalnbcores := sws.NewLabelWidget(100, 25, "Nb cores")
	globalnbcores.Move(325, 125)
	widget.AddChild(globalnbcores)

	globalram := sws.NewLabelWidget(100, 25, "RAM")
	globalram.Move(425, 125)
	widget.AddChild(globalram)

	globaldisk := sws.NewLabelWidget(100, 25, "Disk")
	globaldisk.Move(525, 125)
	widget.AddChild(globaldisk)

	globalallocated := sws.NewLabelWidget(100, 25, "Allocated")
	globalallocated.Move(625, 125)
	widget.AddChild(globalallocated)

	na := ui.NewNothingWidget(725, 25)
	na.Move(0, 150)
	widget.AddChild(na)

	widget.scrolllisting.Move(0, 150)
	widget.scrolllisting.ShowHorizontalScrollbar(false)
	widget.scrolllisting.SetInnerWidget(widget.listing)
	//	widget.AddChild(widget.scrolllisting)

	return widget
}

func (self *PoolManagementWidget) SetGame(inventory *Inventory, currenttime time.Time) {
	if self.inventory != nil {
		inventory.RemoveInventorySubscriber(self)
		self.instock = make([]*InventoryItem, 0, 0)
	}
	self.inventory = inventory
	inventory.AddInventorySubscriber(self)
	inventory.GetDefaultPhysicalPool().AddPoolSubscriber(self)
	inventory.GetDefaultVpsPool().AddPoolSubscriber(self)
	self.Search("assigned:physical")
	// for material not placed but in stock
	for _, item := range self.inventory.Items {
		if item.HasArrived(currenttime) {
			self.ItemInStock(item)
		}
	}
	self.poolassignation.SetActiveChoice(inventory.GetDefaultPoolAllocation())
	height := self.Height()
	width := self.Width()
	if height > 150 {
		if width > 725 {
			width = 725
		}
		if width < 20 {
			width = 20
		}
		self.scrolllisting.Resize(width, height-150)
	}
}
