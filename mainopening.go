package dctycoon

import (
	"github.com/nzin/dctycoon/global"
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"
)

type MainOpening struct {
	sws.CoreWidget
	rootwindow *sws.RootWidget
	background *sdl.Surface
	gamemenu   *MainGameMenu
}

func (self *MainOpening) Repaint() {

	// background image
	self.FillRect(0, 0, self.Width(), self.Height(), 0xff000000)
	self.background.Blit(&sdl.Rect{X: 0, Y: 0, W: self.background.W, H: self.background.H}, self.Surface(), &sdl.Rect{X: 0, Y: 0, W: self.background.W, H: self.background.H})

	// children
	for _, child := range self.GetChildren() {
		// adjust the clipping to the current child
		child.Repaint()
		rectSrc := sdl.Rect{0, 0, child.Width(), child.Height()}
		rectDst := sdl.Rect{child.X(), child.Y(), child.Width(), child.Height()}
		child.Surface().Blit(&rectSrc, self.Surface(), &rectDst)
	}
}

func NewMainOpening(w, h int32, rootwindow *sws.RootWidget, gamemenu *MainGameMenu) *MainOpening {
	corewidget := sws.NewCoreWidget(w, h)
	widget := &MainOpening{
		CoreWidget: *corewidget,
		rootwindow: rootwindow,
		background: nil,
		gamemenu:   gamemenu,
	}
	surface, err := sdl.CreateRGBSurface(0, w, h, 32, 0x00ff0000, 0x0000ff00, 0x000000ff, 0xff000000)
	if err != nil {
		panic(err)
	}
	widget.background = surface

	if image, err := global.LoadImageAsset("assets/ui/opening.jpg"); err == nil {
		dstw := w
		dsth := h
		if image.W*h > image.H*w {
			dstw = image.W * h / image.H
			//	dsth = image.H * w / image.W
		} else {
			dsth = image.H * w / image.W
			//	dstw = image.W * h / image.H
		}
		image.BlitScaled(&sdl.Rect{X: 0, Y: 0, W: image.W, H: image.H}, widget.background, &sdl.Rect{X: (w - dstw) / 2, Y: (h - dsth) / 2, W: dstw, H: dsth})
	} else {
		panic(err)
	}

	gamemenu.Move((w-gamemenu.Width())/2, (h-gamemenu.Height())/2)
	widget.AddChild(gamemenu)

	return widget
}
