package dctycoon

import (
	"io/ioutil"
	"strings"

	"github.com/nzin/sws"
)

//
// StaticWidget is a basic bezeled widget
type StaticWidget struct {
	sws.CoreWidget
}

func NewStaticWidget(h, w int32) *StaticWidget {
	corewidget := sws.NewCoreWidget(h, w)
	return &StaticWidget{
		CoreWidget: *corewidget,
	}
}

func (self *StaticWidget) Repaint() {
	self.CoreWidget.Repaint()

	// bezel
	self.SetDrawColor(0x88, 0x88, 0x88, 0xff)
	self.DrawLine(0, 0, 0, self.Height()-1)
	self.DrawLine(0, self.Height()-1, self.Width()-1, self.Height()-1)
	self.DrawLine(self.Width()-1, self.Height()-1, self.Width()-1, 0)
	self.DrawLine(self.Width()-1, 0, 0, 0)

	self.SetDrawColor(0xff, 0xff, 0xff, 0xff)
	self.DrawLine(1, 1, 1, self.Height()-1)
	self.DrawLine(self.Width()-2, 1, 2, 1)

	self.SetDrawColor(0x88, 0x88, 0x88, 0xff)
	self.DrawLine(1, self.Height()-2, self.Width()-2, self.Height()-2)
	self.DrawLine(self.Width()-2, self.Height()-2, self.Width()-2, 2)

	// bezel interior
	self.SetDrawColor(0xdd, 0xdd, 0xdd, 0xff)
	self.DrawLine(2, 2, self.Width()-3, 2)
	self.DrawLine(self.Width()-3, 2, self.Width()-3, self.Height()-3)
	self.DrawLine(self.Width()-3, self.Height()-3, 2, self.Height()-3)
	self.DrawLine(2, self.Height()-3, 2, 2)

	self.SetDrawColor(0xdd, 0xdd, 0xdd, 0xff)
	self.DrawLine(3, 3, self.Width()-4, 3)
	self.DrawLine(self.Width()-4, 3, self.Width()-4, self.Height()-4)
	self.DrawLine(self.Width()-4, self.Height()-4, 3, self.Height()-4)
	self.DrawLine(3, self.Height()-4, 3, 3)

	self.SetDrawColor(0xbb, 0xbb, 0xbb, 0xff)
	self.DrawLine(4, 4, self.Width()-5, 4)
	self.DrawLine(self.Width()-5, 4, self.Width()-5, self.Height()-5)
	self.DrawLine(self.Width()-5, self.Height()-5, 4, self.Height()-5)
	self.DrawLine(4, self.Height()-5, 4, 4)

	self.SetDrawColor(0x88, 0x88, 0x88, 0xff)
	self.DrawLine(5, 5, self.Width()-6, 5)
	self.DrawLine(self.Width()-6, 5, self.Width()-6, self.Height()-6)
	self.DrawLine(self.Width()-6, self.Height()-6, 4, self.Height()-6)
	self.DrawLine(5, self.Height()-6, 5, 5)
}

//
// MainGameMenu is the main game window (new/load/save/quit...)
type MainGameMenu struct {
	StaticWidget
	game         *Game
	newbutton    *sws.ButtonWidget
	loadbutton   *sws.ButtonWidget
	savebutton   *sws.ButtonWidget
	cancelbutton *sws.ButtonWidget
	quitbutton   *sws.ButtonWidget
	loadwidget   *MainGameMenuLoad
}

func NewMainGameMenu(game *Game, root *sws.RootWidget, quit *bool) *MainGameMenu {
	corewidget := NewStaticWidget(600, 250)
	widget := &MainGameMenu{
		StaticWidget: *corewidget,
		game:         game,
		newbutton:    sws.NewButtonWidget(300, 50, "New Game"),
		loadbutton:   sws.NewButtonWidget(300, 50, "Load Game"),
		savebutton:   sws.NewButtonWidget(300, 50, "Save Game"),
		quitbutton:   sws.NewButtonWidget(300, 50, "Quit"),
		cancelbutton: sws.NewButtonWidget(150, 30, "Cancel"),
		loadwidget:   NewMainGameMenuLoad(),
	}

	widget.newbutton.Move(150, 50)
	widget.AddChild(widget.newbutton)

	widget.loadbutton.Move(150, 110)
	widget.AddChild(widget.loadbutton)
	widget.loadbutton.SetClicked(func() {
		widget.loadwidget.Loadfiles()
		root.AddChild(widget.loadwidget)
		root.SetFocus(widget.loadwidget.listwidget)
	})

	widget.quitbutton.Move(150, 170)
	widget.AddChild(widget.quitbutton)
	widget.quitbutton.SetClicked(func() {
		*quit = true
	})

	widget.loadwidget.Move((root.Width()-widget.loadwidget.Width())/2, (root.Height()-widget.loadwidget.Height())/2)
	widget.loadwidget.SetCancelCallback(func() {
		root.RemoveChild(widget.loadwidget)
	})

	widget.loadwidget.SetLoadCallback(func(filename string) {
		root.RemoveChild(widget.loadwidget)
		game.LoadGame(filename)
	})

	return widget
}

func (self *MainGameMenu) ShowSave() {
	self.Resize(600, 350)
	self.newbutton.Move(150, 50)
	self.loadbutton.Move(150, 110)
	self.savebutton.Move(150, 170)
	self.AddChild(self.savebutton)
	self.quitbutton.Move(150, 230)
	self.cancelbutton.Move(300, 290)
	self.AddChild(self.cancelbutton)
}

func (self *MainGameMenu) SetCancelCallback(callback func()) {
	self.cancelbutton.SetClicked(callback)
}

//
// MainGameMenuLoad is the load game window.
// Used by MainGameMenu
type MainGameMenuLoad struct {
	StaticWidget
	listwidget   *sws.ListWidget
	loadbutton   *sws.ButtonWidget
	cancelbutton *sws.ButtonWidget
	loadcallback func(filename string)
}

func NewMainGameMenuLoad() *MainGameMenuLoad {
	corewidget := NewStaticWidget(400, 500)
	widget := &MainGameMenuLoad{
		StaticWidget: *corewidget,
		listwidget:   sws.NewListWidget(360, 400),
		cancelbutton: sws.NewButtonWidget(100, 25, "Cancel"),
		loadbutton:   sws.NewButtonWidget(100, 25, "Load"),
	}

	widget.listwidget.Move(20, 50)
	widget.AddChild(widget.listwidget)

	widget.cancelbutton.Move(280, 460)
	widget.AddChild(widget.cancelbutton)

	widget.loadbutton.Move(20, 460)
	widget.AddChild(widget.loadbutton)
	widget.loadbutton.SetClicked(func() {
		currentitem := widget.listwidget.GetCurrentItem()
		if currentitem != nil {
			widget.loadcallback(currentitem.GetText() + ".map")
		}
	})

	return widget
}

func (self *MainGameMenuLoad) Loadfiles() {
	for _, c := range self.listwidget.GetItems() {
		self.listwidget.RemoveItem(c)
	}

	// check in the working directory all files in ".map"
	files, err := ioutil.ReadDir(".")
	if err == nil {
		for _, f := range files {
			filename := f.Name()
			if strings.HasSuffix(filename, ".map") {
				self.listwidget.AddItem(filename[:len(filename)-4])
			}
		}
	}
}

func (self *MainGameMenuLoad) SetLoadCallback(callback func(filename string)) {
	self.loadcallback = callback
}

func (self *MainGameMenuLoad) SetCancelCallback(callback func()) {
	self.cancelbutton.SetClicked(callback)
}
