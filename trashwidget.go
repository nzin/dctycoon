package dctycoon

import (
	"fmt"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"
)

type TrashWidget struct {
	sws.CoreWidget
	root      *sws.RootWidget
	inventory *supplier.Inventory
	trashIcon *sdl.Surface
}

func (self *TrashWidget) DragDrop(x, y int32, payload sws.DragPayload) bool {
	fmt.Println("TrashWidget::DragDrop(", x, ",", y, ",", payload, ")", payload.GetType())
	if payload.GetType() == global.DRAG_RACK_SERVER_FROM_TOWER {
		servermove := payload.(*ServerMovePayload)
		discarded := self.inventory.DiscardItem(servermove.item)
		if discarded == false {
			iconsurface, _ := global.LoadImageAsset("assets/ui/icon-triangular-big.png")
			sws.ShowModalErrorSurfaceicon(self.root, "Cannot trash it", iconsurface, "This item cannot be trashed because it is used by a customer", nil)
		}
		return discarded
	}
	if payload.GetType() == global.DRAG_ELEMENT_PAYLOAD {
		element := payload.(*ElementDragPayload)
		return self.inventory.DiscardItem(element.item)
	}
	if payload.GetType() == global.DRAG_RACK_SERVER {
		server := payload.(*ServerDragPayload)
		return self.inventory.DiscardItem(server.item)
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

func NewTrashWidget(w, h int32, root *sws.RootWidget) *TrashWidget {
	core := sws.NewCoreWidget(w, h)
	widget := &TrashWidget{
		CoreWidget: *core,
		root:       root,
		inventory:  nil,
		trashIcon:  nil,
	}

	widget.trashIcon, _ = global.LoadImageAsset("assets/ui/trash.png")

	return widget
}
