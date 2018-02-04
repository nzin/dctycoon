package dctycoon

import (
	"fmt"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/sws"
)

type DebugOpponentWidget struct {
	sws.CoreWidget
	refreshbutton      *sws.ButtonWidget
	table              *sws.TableWidget
	data               [][]string
	datachangecallback func()
	bigmapbutton       *sws.ButtonWidget
	smallmapbutton     *sws.ButtonWidget
	game               *Game
}

func NewDebugOpponentWidget(w, h int32, game *Game) *DebugOpponentWidget {
	core := sws.NewCoreWidget(w, h)
	widget := &DebugOpponentWidget{
		CoreWidget:     *core,
		refreshbutton:  sws.NewButtonWidget(200, 25, "Refresh"),
		data:           make([][]string, 0, 0),
		bigmapbutton:   sws.NewButtonWidget(200, 25, "Switch to 24x24 map"),
		smallmapbutton: sws.NewButtonWidget(200, 25, "Switch to 3x4 map"),
		game:           game,
	}
	widget.refreshbutton.SetClicked(widget.refresh)
	widget.AddChild(widget.refreshbutton)

	widget.table = sws.NewTableWidget(200, 200, widget)
	widget.table.Move(0, 25)
	widget.AddChild(widget.table)

	widget.bigmapbutton.Move(0, 250)
	widget.AddChild(widget.bigmapbutton)
	widget.bigmapbutton.SetClicked(func() {
		game.MigrateMap("24_24_standard.json")
	})

	widget.smallmapbutton.Move(0, 275)
	widget.AddChild(widget.smallmapbutton)
	widget.smallmapbutton.SetClicked(func() {
		game.MigrateMap("3_4_room.json")
	})

	return widget
}

func (self *DebugOpponentWidget) refresh() {
	self.data = make([][]string, 0, 0)
	for _, o := range self.game.GetNPActors() {
		for _, i := range o.GetInventory().Items {
			if i.Typeitem == supplier.PRODUCT_SERVER {
				array := make([]string, 5, 5)
				array[0] = o.GetName()
				array[1] = global.AdjustMega(i.Serverconf.NbDisks * i.Serverconf.DiskSize)
				array[2] = global.AdjustMega(i.Serverconf.NbSlotRam * i.Serverconf.RamSize)
				array[3] = fmt.Sprintf("%d", i.Serverconf.NbCore)
				array[4] = "false"
				if i.Diskallocated != 0 {
					array[4] = "true"
				}
				self.data = append(self.data, array)
			}
		}
	}
	self.datachangecallback()
}

func (self *DebugOpponentWidget) Resize(w, h int32) {
	self.CoreWidget.Resize(w, h)
	if h > 50 {
		self.table.Resize(w, h-25)
	}
}

func (self *DebugOpponentWidget) GetNbColumns() int32 {
	return 5
}

func (self *DebugOpponentWidget) GetNbRows() int32 {
	return int32(len(self.data))
}

func (self *DebugOpponentWidget) GetHeader(column int32) (string, int32) {
	switch column {
	case 0:
		return "owner", 200
	case 1:
		return "disk", 100
	case 2:
		return "ram", 100
	case 3:
		return "cpu", 100
	case 4:
		return "allocated", 100
	}
	return "undefined", 100
}

func (self *DebugOpponentWidget) GetCell(column, row int32) string {
	return self.data[row][column]
}

// when the table grows/shrink
func (self *DebugOpponentWidget) SetRowUpdateCallback(callback func()) {
	self.datachangecallback = callback
}

// when the table need to be refreshed
func (self *DebugOpponentWidget) SetDataChangeCallback(callback func()) {
}
