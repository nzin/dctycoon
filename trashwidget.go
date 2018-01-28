package dctycoon

import (
	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"
)

type TrashWidget struct {
	sws.CoreWidget
	inventory *supplier.Inventory
	trashIcon *sdl.Surface
}

func (self *TrashWidget) DragDrop(x, y int32, payload sws.DragPayload) bool {
	if payload.GetType() == global.DRAG_RACK_SERVER_FROM_TOWER {
		servermove := payload.(*ServerMovePayload)
		return self.inventory.DiscardItem(servermove.item)
	}
	return false
}

func (self *TrashWidget) Repaint() {
	wGap := self.Width() - self.trashIcon.W
	hGap := self.Height() - self.trashIcon.H
	rectSrc := sdl.Rect{0, 0, self.trashIcon.W, self.trashIcon.H}
	rectDst := sdl.Rect{(wGap / 2), (hGap / 2), self.Width() - (wGap / 2), self.Height() - (hGap / 2)}
	self.trashIcon.Blit(&rectSrc, self.Surface(), &rectDst)

	self.SetDrawColor(0, 0, 0, 255)
	self.DrawLine(0, 0, self.Width()-1, 0)
	self.DrawLine(self.Width()-1, 0, self.Width()-1, self.Height()-1)
	self.DrawLine(self.Width()-1, self.Height()-1, 0, self.Height()-1)
	self.DrawLine(0, self.Height()-1, 0, 0)
}

func (self *TrashWidget) SetInventory(inventory *supplier.Inventory) {
	self.inventory = inventory
}

func NewTrashWidget(w, h int32) *TrashWidget {
	core := sws.NewCoreWidget(w, h)
	widget := &TrashWidget{
		CoreWidget: *core,
		inventory:  nil,
		trashIcon:  nil,
	}

	widget.trashIcon, _ = global.LoadImageAsset("assets/ui/trash.png")

	return widget
}
