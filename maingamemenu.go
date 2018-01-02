package dctycoon

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/nzin/sws"
)

//
// MainGameMenu is the main game window (new/load/save/quit...)
type MainGameMenu struct {
	sws.CoreWidget
	game         *Game
	newbutton    *sws.ButtonWidget
	loadbutton   *sws.ButtonWidget
	savebutton   *sws.ButtonWidget
	cancelbutton *sws.ButtonWidget
	quitbutton   *sws.ButtonWidget
	loadwidget   *MainGameMenuLoad
}

func NewMainGameMenu(game *Game, root *sws.RootWidget, quit *bool) *MainGameMenu {
	corewidget := sws.NewCoreWidget(600, 350)
	widget := &MainGameMenu{
		CoreWidget:   *corewidget,
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
	})

	widget.quitbutton.Move(150, 170)
	widget.AddChild(widget.quitbutton)
	widget.quitbutton.SetClicked(func() {
		*quit = true
	})

	widget.loadwidget.Move((root.Width()-widget.loadwidget.Width())/2, (root.Height()-widget.loadwidget.Height())/2)
	widget.loadwidget.SetCancelCallback(func() {
		fmt.Println("cancel callback")
		root.RemoveChild(widget.loadwidget)
	})

	widget.loadwidget.SetLoadCallback(func(filename string) {
		root.RemoveChild(widget.loadwidget)
		game.LoadGame(filename)
	})

	return widget
}

func (self *MainGameMenu) ShowSave() {
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
	sws.CoreWidget
	vbox         *sws.VBoxWidget
	scroll       *sws.ScrollWidget
	cancelbutton *sws.ButtonWidget
	loadcallback func(filename string)
}

func NewMainGameMenuLoad() *MainGameMenuLoad {
	corewidget := sws.NewCoreWidget(400, 500)
	widget := &MainGameMenuLoad{
		CoreWidget:   *corewidget,
		vbox:         sws.NewVBoxWidget(360, 0),
		scroll:       sws.NewScrollWidget(360, 400),
		cancelbutton: sws.NewButtonWidget(100, 25, "Cancel"),
	}

	widget.scroll.SetInnerWidget(widget.vbox)
	widget.scroll.Move(20, 50)
	widget.scroll.ShowHorizontalScrollbar(false)
	widget.AddChild(widget.scroll)

	widget.cancelbutton.Move(280, 465)
	widget.AddChild(widget.cancelbutton)

	return widget
}

func (self *MainGameMenuLoad) Loadfiles() {
	for _, c := range self.vbox.GetChildren() {
		self.vbox.RemoveChild(c)
	}

	// check in the working directory all files in ".map"
	files, err := ioutil.ReadDir(".")
	if err == nil {
		for _, f := range files {
			filename := f.Name()
			if strings.HasSuffix(filename, ".map") {
				loadbutton := sws.NewButtonWidget(360, 30, filename[:len(filename)-4])
				loadbutton.SetClicked(func() {
					self.loadcallback(filename)
				})
				self.vbox.AddChild(loadbutton)
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
