package ui

import (
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"
)

//
// NothingWidget is used to print a "N/A" banner widget
type NothingWidget struct {
	sws.CoreWidget
}

func (self *NothingWidget) Repaint() {
	self.SetDrawColor(128, 128, 128, 255)

	for y := int32(0); y < self.Height(); y++ {
		for x := int32(0); x < self.Width()+20; x += 40 {
			offset := (y % 40) - 20
			self.DrawLine(x+offset, y, x+offset+20, y)
		}
	}

	self.FillRect(self.Width()/2-15, self.Height()/2-10, 30, 20, 0)
	self.WriteText(self.Width()/2-15, self.Height()/2-10, "N/A", sdl.Color{0, 0, 0, 255})
	self.SetDirtyFalse()
}

func NewNothingWidget(width, height int32) *NothingWidget {
	corewidget := sws.NewCoreWidget(width, height)
	widget := &NothingWidget{
		CoreWidget: *corewidget,
	}
	return widget
}
