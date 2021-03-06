package dctycoon

import (
	"math"
	"time"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"

	log "github.com/sirupsen/logrus"
)

// WorldmapWidget is a worldmap image widget
// where you can select a Location (see supplier.LocationType and supplier.AvailableLocation)
type WorldmapWidget struct {
	sws.CoreWidget
	selected             string
	hotspot              string
	background           *sdl.Surface
	xshift               int32
	yshift               int32
	scale                float64
	hotspotcallback      func(selected, hotspot string)
	spotzoomeffect       int32
	evtspotzoomeffect    *sws.TimerEvent
	hotspotzoomeffect    int32
	evthotspotzoomeffect *sws.TimerEvent
}

func (self *WorldmapWidget) MouseMove(x, y, xrel, yrel int32) {
	var currenthotspot string
	for locationid, l := range supplier.AvailableLocation {
		xlocation := float64(l.Xmap)*self.scale + float64(self.xshift)
		ylocation := float64(l.Ymap)*self.scale + float64(self.yshift)

		if math.Abs((float64(x)-xlocation)*(float64(x)-xlocation)+(float64(y)-ylocation)*(float64(y)-ylocation)) < 200.0 {
			currenthotspot = locationid
		}
	}
	if currenthotspot != self.hotspot {
		self.hotspot = currenthotspot
		if self.hotspotcallback != nil {
			self.hotspotcallback(self.selected, self.hotspot)
		}

		self.hotspotzoomeffect = 0
		if self.evthotspotzoomeffect != nil {
			self.evthotspotzoomeffect.StopRepeat()
			self.evthotspotzoomeffect = nil
		}
		self.evthotspotzoomeffect = sws.TimerAddEvent(time.Now(), 30*time.Millisecond, func(evt *sws.TimerEvent) {
			self.hotspotzoomeffect++
			if self.hotspotzoomeffect > 10 {
				evt.StopRepeat()
			}
			self.PostUpdate()
		})
	}
	self.PostUpdate()
}

func (self *WorldmapWidget) MousePressDown(x, y int32, button uint8) {
	self.selected = self.hotspot
	self.PostUpdate()
}

func (self *WorldmapWidget) Repaint() {
	// background image
	self.FillRect(0, 0, self.Width(), self.Height(), 0xff000000)
	self.background.Blit(&sdl.Rect{X: 0, Y: 0, W: self.background.W, H: self.background.H}, self.Surface(), &sdl.Rect{X: 0, Y: 0, W: self.background.W, H: self.background.H})

	// different spots
	for locationid, l := range supplier.AvailableLocation {
		x := int32(float64(l.Xmap)*self.scale) + self.xshift
		y := int32(float64(l.Ymap)*self.scale) + self.yshift

		alpha := self.spotzoomeffect*10 - (l.Xmap / 8)
		if alpha < 0 {
			alpha = 0
		}
		if alpha > 255 {
			alpha = 255
		}
		self.SetDrawColor(0x46, 0xc8, 0xe8, uint8(alpha))
		if locationid == self.hotspot {
			self.SetDrawColor(0xff, 0xa0, 0xa0, 255)
		}
		if locationid == self.selected {
			self.SetDrawColor(0xff, 0x20, 0x20, 255)
		}
		//		for dy := y - 3; dy < y+3; dy++ {
		//			self.DrawLine(x-3, dy, x+3, dy)
		//		}
		stretch := (255 - alpha) / 20
		for dy := int32(0); dy < 4; dy++ {
			self.DrawLine(x-dy-stretch, y-dy-stretch, x+dy+stretch, y-dy-stretch)
			self.DrawLine(x-dy-stretch, y+dy+stretch, x+dy+stretch, y+dy+stretch)

			self.DrawLine(x-dy-stretch, y-dy-stretch, x-dy-stretch, y+dy+stretch)
			self.DrawLine(x+dy+stretch, y-dy-stretch, x+dy+stretch, y+dy+stretch)
		}
		if self.hotspotzoomeffect > 0 && locationid == self.hotspot && self.hotspotzoomeffect < 10 {
			self.SetDrawColor(0xff, 0xa0, 0xa0, 255-uint8(20*self.hotspotzoomeffect))
			for dy := int32(0); dy < 4; dy++ {
				self.DrawLine(x-dy-self.hotspotzoomeffect, y-dy-self.hotspotzoomeffect, x+dy+self.hotspotzoomeffect, y-dy-self.hotspotzoomeffect)
				self.DrawLine(x-dy-self.hotspotzoomeffect, y+dy+self.hotspotzoomeffect, x+dy+self.hotspotzoomeffect, y+dy+self.hotspotzoomeffect)

				self.DrawLine(x-dy-self.hotspotzoomeffect, y-dy-self.hotspotzoomeffect, x-dy-self.hotspotzoomeffect, y+dy+self.hotspotzoomeffect)
				self.DrawLine(x+dy+self.hotspotzoomeffect, y-dy-self.hotspotzoomeffect, x+dy+self.hotspotzoomeffect, y+dy+self.hotspotzoomeffect)
			}
		}
	}

	// children
	for _, child := range self.GetChildren() {
		// adjust the clipping to the current child
		child.Repaint()
		rectSrc := sdl.Rect{0, 0, child.Width(), child.Height()}
		rectDst := sdl.Rect{child.X(), child.Y(), child.Width(), child.Height()}
		child.Surface().Blit(&rectSrc, self.Surface(), &rectDst)
	}
}

func (self *WorldmapWidget) SetLocationCallback(callback func(selected, hotspot string)) {
	self.hotspotcallback = callback
}

func (self *WorldmapWidget) Reset() {
	log.Debug("WorldmapWidget::Reset()")
	self.selected = ""
	self.hotspot = ""

	self.spotzoomeffect = 0
	if self.evtspotzoomeffect != nil {
		self.evtspotzoomeffect.StopRepeat()
		self.evtspotzoomeffect = nil
	}
	self.evtspotzoomeffect = sws.TimerAddEvent(time.Now(), 15*time.Millisecond, func(evt *sws.TimerEvent) {
		self.spotzoomeffect++
		if self.spotzoomeffect > 100 {
			evt.StopRepeat()
		}
		self.PostUpdate()
	})
}

func NewWorldmapWidget(w, h int32) *WorldmapWidget {
	corewidget := sws.NewCoreWidget(w, h)
	widget := &WorldmapWidget{
		CoreWidget: *corewidget,
		background: nil,
	}
	surface, err := sdl.CreateRGBSurface(0, w, h, 32, 0x00ff0000, 0x0000ff00, 0x000000ff, 0xff000000)
	if err != nil {
		panic(err)
	}
	widget.background = surface

	if image, err := global.LoadImageAsset("assets/ui/worldmap.jpg"); err == nil {
		dstw := w
		dsth := h
		if image.W*h > image.H*w {
			dsth = image.H * w / image.W
			widget.scale = float64(w) / float64(image.W)
		} else {
			dstw = image.W * h / image.H
			widget.scale = float64(h) / float64(image.H)
		}
		widget.xshift = (w - dstw) / 2
		widget.yshift = (h - dsth) / 2
		image.BlitScaled(&sdl.Rect{X: 0, Y: 0, W: image.W, H: image.H}, widget.background, &sdl.Rect{X: widget.xshift, Y: widget.yshift, W: dstw, H: dsth})
	} else {
		panic(err)
	}

	return widget
}
