package dctycoon

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/supplier"
	"github.com/veandco/go-sdl2/sdl"

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
	root          *sws.RootWidget
	game          *Game
	newbutton     *sws.ButtonWidget
	loadbutton    *sws.ButtonWidget
	savebutton    *sws.ButtonWidget
	cancelbutton  *sws.ButtonWidget
	quitbutton    *sws.ButtonWidget
	loadwidget    *MainGameMenuLoad
	savewidget    *MainGameMenuSave
	newwidget     *MainGameMenuNew
	previousSpeed int
}

func NewMainGameMenu(game *Game, root *sws.RootWidget, quit *bool) *MainGameMenu {
	corewidget := NewStaticWidget(600, 250)
	widget := &MainGameMenu{
		StaticWidget:  *corewidget,
		root:          root,
		game:          game,
		newbutton:     sws.NewButtonWidget(300, 50, "New Game"),
		loadbutton:    sws.NewButtonWidget(300, 50, "Load Game"),
		savebutton:    sws.NewButtonWidget(300, 50, "Save Game"),
		quitbutton:    sws.NewButtonWidget(300, 50, "Quit"),
		cancelbutton:  sws.NewButtonWidget(150, 30, "Cancel"),
		loadwidget:    NewMainGameMenuLoad(),
		savewidget:    NewMainGameMenuSave(),
		newwidget:     NewMainGameMenuNew(root),
		previousSpeed: SPEED_STOP,
	}

	widget.newbutton.Move(150, 50)
	widget.AddChild(widget.newbutton)
	widget.newbutton.SetClicked(func() {
		root.AddChild(widget.newwidget)
		widget.newwidget.Reset()
		root.SetModal(widget.newwidget)
	})
	widget.newwidget.SetCancelCallback(func() {
		root.RemoveChild(widget.newwidget)
		root.SetModal(widget)
	})
	widget.newwidget.SetNewGameCallback(func(location string, difficulty int32, companyname string) {
		root.RemoveChild(widget.newwidget)
		game.InitGame(location, difficulty, companyname)
	})

	widget.loadbutton.Move(150, 110)
	widget.AddChild(widget.loadbutton)
	widget.loadbutton.SetClicked(func() {
		widget.loadwidget.Loadfiles()
		root.AddChild(widget.loadwidget)
		root.SetFocus(widget.loadwidget.listwidget)
		root.SetModal(widget.loadwidget)
	})

	widget.quitbutton.Move(150, 170)
	widget.AddChild(widget.quitbutton)
	widget.quitbutton.SetClicked(func() {
		*quit = true
	})

	widget.loadwidget.Move((root.Width()-widget.loadwidget.Width())/2, (root.Height()-widget.loadwidget.Height())/2)
	widget.loadwidget.SetCancelCallback(func() {
		root.RemoveChild(widget.loadwidget)
		root.SetModal(widget)
	})

	widget.loadwidget.SetLoadCallback(func(filename string) {
		root.RemoveChild(widget.loadwidget)
		game.LoadGame(filename)
	})

	widget.savewidget.Move((root.Width()-widget.savewidget.Width())/2, (root.Height()-widget.savewidget.Height())/2)
	widget.savebutton.SetClicked(func() {
		widget.savewidget.Loadfiles()
		root.AddChild(widget.savewidget)
		root.SetFocus(widget.savewidget.filenameinput)
		root.SetModal(widget.savewidget)
	})

	widget.savewidget.SetCancelCallback(func() {
		root.RemoveChild(widget.savewidget)
		root.SetModal(widget)
	})

	widget.savewidget.SetSaveCallback(func(filename string) {
		root.RemoveChild(widget.savewidget)
		game.SaveGame(filename)
		root.SetModal(widget)
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

	self.previousSpeed = self.game.GetCurrentSpeed()
	self.game.ChangeGameSpeed(SPEED_STOP)
	self.root.AddChild(self)
	self.root.SetModal(self)
	self.SetCancelCallback(func() {
		self.root.RemoveChild(self)
		self.game.ChangeGameSpeed(self.previousSpeed)
	})
}

func (self *MainGameMenu) SetCancelCallback(callback func()) {
	self.cancelbutton.SetClicked(callback)
}

//
// MainGameMenuLoad is the load game window.
// Used by MainGameMenu
type MainGameMenuLoad struct {
	StaticWidget
	loadlabel    *sws.LabelWidget
	hr           *sws.Hr
	listwidget   *sws.ListWidget
	loadbutton   *sws.ButtonWidget
	cancelbutton *sws.ButtonWidget
	loadcallback func(filename string)
}

func NewMainGameMenuLoad() *MainGameMenuLoad {
	corewidget := NewStaticWidget(400, 500)
	widget := &MainGameMenuLoad{
		StaticWidget: *corewidget,
		loadlabel:    sws.NewLabelWidget(360, 25, "Load game"),
		hr:           sws.NewHr(360),
		listwidget:   sws.NewListWidget(360, 380),
		cancelbutton: sws.NewButtonWidget(100, 25, "Cancel"),
		loadbutton:   sws.NewButtonWidget(100, 25, "Load"),
	}

	widget.loadlabel.Move(20, 20)
	widget.loadlabel.SetCentered(true)
	widget.AddChild(widget.loadlabel)

	widget.hr.Move(20, 50)
	widget.AddChild(widget.hr)

	widget.listwidget.Move(20, 70)
	widget.AddChild(widget.listwidget)

	widget.cancelbutton.Move(20, 460)
	widget.AddChild(widget.cancelbutton)

	widget.loadbutton.Move(280, 460)
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
	self.listwidget.RemoveAllItems()

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

//
// MainGameMenuSave is the save game window.
// Used by MainGameMenu
type MainGameMenuSave struct {
	StaticWidget
	savelabel     *sws.LabelWidget
	hr            *sws.Hr
	filenamelabel *sws.LabelWidget
	filenameinput *sws.InputWidget
	listwidget    *sws.ListWidget
	savebutton    *sws.ButtonWidget
	cancelbutton  *sws.ButtonWidget
	savecallback  func(filename string)
}

func NewMainGameMenuSave() *MainGameMenuSave {
	corewidget := NewStaticWidget(400, 500)
	widget := &MainGameMenuSave{
		StaticWidget:  *corewidget,
		savelabel:     sws.NewLabelWidget(360, 25, "Save game"),
		hr:            sws.NewHr(360),
		filenamelabel: sws.NewLabelWidget(100, 25, "Filename:"),
		filenameinput: sws.NewInputWidget(260, 25, ""),
		listwidget:    sws.NewListWidget(360, 350),
		cancelbutton:  sws.NewButtonWidget(100, 25, "Cancel"),
		savebutton:    sws.NewButtonWidget(100, 25, "Save"),
	}

	widget.savelabel.Move(20, 20)
	widget.savelabel.SetCentered(true)
	widget.AddChild(widget.savelabel)

	widget.hr.Move(20, 50)
	widget.AddChild(widget.hr)

	widget.filenamelabel.Move(20, 65)
	widget.AddChild(widget.filenamelabel)

	widget.filenameinput.Move(120, 65)
	widget.AddChild(widget.filenameinput)

	widget.listwidget.Move(20, 100)
	widget.AddChild(widget.listwidget)
	widget.listwidget.SetCallbackValueChanged(func() {
		currentitem := widget.listwidget.GetCurrentItem()
		if currentitem != nil {
			widget.filenameinput.SetText(currentitem.GetText())
		}
	})

	widget.cancelbutton.Move(20, 460)
	widget.AddChild(widget.cancelbutton)

	widget.savebutton.Move(280, 460)
	widget.AddChild(widget.savebutton)
	widget.savebutton.SetClicked(func() {
		filename := widget.filenameinput.GetText()
		if filename != "" {
			widget.savecallback(filename + ".map")
		}
	})

	return widget
}

func (self *MainGameMenuSave) Loadfiles() {
	self.listwidget.RemoveAllItems()

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

func (self *MainGameMenuSave) SetSaveCallback(callback func(filename string)) {
	self.savecallback = callback
}

func (self *MainGameMenuSave) SetCancelCallback(callback func()) {
	self.cancelbutton.SetClicked(callback)
}

//
// MainGameMenuNew is the 'new' game window.
// Used by MainGameMenu
type MainGameMenuNew struct {
	StaticWidget
	root            *sws.RootWidget
	worldmap        *WorldmapWidget
	locationname    *sws.LabelWidget
	currentlocation string
	bankrate        *sws.LabelWidget
	taxrate         *sws.LabelWidget
	temperature     *sws.LabelWidget
	rentingprice    *sws.LabelWidget
	difficulty      *sws.DropdownWidget
	companyname     *sws.InputWidget
	createbutton    *sws.ButtonWidget
	cancelbutton    *sws.ButtonWidget
}

func (self *MainGameMenuNew) Reset() {
	self.difficulty.SetActiveChoice(0)
	self.worldmap.Reset()
	self.SetLocation("", "")
}

func (self *MainGameMenuNew) SetLocation(selected, hotspot string) {
	var location *supplier.LocationType
	if hotspot != "" {
		location = supplier.AvailableLocation[hotspot]
		self.currentlocation = hotspot
	} else {
		location = supplier.AvailableLocation[selected]
		self.currentlocation = selected
	}
	if location != nil {
		self.locationname.SetText(location.Name)
		self.bankrate.SetText(fmt.Sprintf("%.2f %%", location.Bankinterestrate*100))
		self.taxrate.SetText(fmt.Sprintf("%.2f %%", location.Taxrate*100))
		self.temperature.SetText(fmt.Sprintf("%.1f °C", location.Temperatureaverage))
		self.rentingprice.SetText(fmt.Sprintf("%.0f $/m²", location.Metersquareprice))
	} else {
		self.locationname.SetText("")
		self.bankrate.SetText("- %")
		self.taxrate.SetText("- %")
		self.temperature.SetText("- °C")
		self.rentingprice.SetText("- $/m²")
	}
}

func (self *MainGameMenuNew) SetCancelCallback(callback func()) {
	self.cancelbutton.SetClicked(callback)
}

func (self *MainGameMenuNew) SetNewGameCallback(callback func(location string, difficulty int32, companyname string)) {
	self.createbutton.SetClicked(func() {
		if self.currentlocation == "" {
			iconsurface, _ := global.LoadImageAsset("assets/ui/icon-triangular-big.png")
			sws.ShowModalErrorSurfaceicon(self.root, "No location selected", iconsurface, "You must select a location where you want to be based", nil)
			return
		}
		if self.companyname.GetText() == "" {
			iconsurface, _ := global.LoadImageAsset("assets/ui/icon-triangular-big.png")
			sws.ShowModalErrorSurfaceicon(self.root, "No company name", iconsurface, "You must give your company a name", nil)
			return
		}
		callback(self.currentlocation, self.difficulty.ActiveChoice, self.companyname.GetText())
	})
}

func NewMainGameMenuNew(root *sws.RootWidget) *MainGameMenuNew {
	corewidget := NewStaticWidget(root.Width()-200, root.Height()-200)
	widget := &MainGameMenuNew{
		StaticWidget: *corewidget,
		root:         root,
		worldmap:     NewWorldmapWidget(root.Width()-400, root.Height()-200),
		locationname: sws.NewLabelWidget(100, 25, ""),
		bankrate:     sws.NewLabelWidget(100, 25, ""),
		taxrate:      sws.NewLabelWidget(100, 25, ""),
		temperature:  sws.NewLabelWidget(100, 25, ""),
		rentingprice: sws.NewLabelWidget(100, 25, ""),
		difficulty:   sws.NewDropdownWidget(100, 25, []string{"Easy", "Medium", "Hard"}),
		companyname:  sws.NewInputWidget(100, 25, "noname"),
		createbutton: sws.NewButtonWidget(100, 40, "Create"),
		cancelbutton: sws.NewButtonWidget(100, 40, "Cancel"),
	}
	widget.SetColor(0xff000000)

	widget.AddChild(widget.worldmap)
	widget.Move(100, 100)
	widget.worldmap.SetLocationCallback(func(selected, hotspot string) {
		widget.SetLocation(selected, hotspot)
	})

	selectyourlocation := sws.NewLabelWidget(300, 30, "Select your location on the map:")
	selectyourlocation.SetFont(sws.LatoRegular20)
	selectyourlocation.SetColor(0x00000000)
	selectyourlocation.SetTextColor(sdl.Color{0xff, 0xff, 0xff, 0xff})
	selectyourlocation.Move(20, 20)
	widget.AddChild(selectyourlocation)

	name := sws.NewLabelWidget(130, 25, "Location:")
	name.Move(root.Width()-450, 25)
	name.SetTextColor(sdl.Color{0xff, 0xff, 0xff, 0xff})
	name.SetColor(0xff000000)
	widget.AddChild(name)
	widget.locationname.Move(root.Width()-310, 25)
	widget.locationname.SetTextColor(sdl.Color{0xff, 0xff, 0xff, 0xff})
	widget.locationname.SetColor(0xff000000)
	widget.AddChild(widget.locationname)

	bankrate := sws.NewLabelWidget(130, 25, "Bank rate:")
	bankrate.Move(root.Width()-450, 50)
	bankrate.SetTextColor(sdl.Color{0xff, 0xff, 0xff, 0xff})
	bankrate.SetColor(0xff000000)
	widget.AddChild(bankrate)
	widget.bankrate.Move(root.Width()-310, 50)
	widget.bankrate.SetTextColor(sdl.Color{0xff, 0xff, 0xff, 0xff})
	widget.bankrate.SetColor(0xff000000)
	widget.AddChild(widget.bankrate)

	taxrate := sws.NewLabelWidget(130, 25, "Tax rate:")
	taxrate.Move(root.Width()-450, 75)
	taxrate.SetTextColor(sdl.Color{0xff, 0xff, 0xff, 0xff})
	taxrate.SetColor(0xff000000)
	widget.AddChild(taxrate)
	widget.taxrate.Move(root.Width()-310, 75)
	widget.taxrate.SetTextColor(sdl.Color{0xff, 0xff, 0xff, 0xff})
	widget.taxrate.SetColor(0xff000000)
	widget.AddChild(widget.taxrate)

	temperature := sws.NewLabelWidget(130, 25, "Avg temperature:")
	temperature.Move(root.Width()-450, 100)
	temperature.SetTextColor(sdl.Color{0xff, 0xff, 0xff, 0xff})
	temperature.SetColor(0xff000000)
	widget.AddChild(temperature)
	widget.temperature.Move(root.Width()-310, 100)
	widget.temperature.SetTextColor(sdl.Color{0xff, 0xff, 0xff, 0xff})
	widget.temperature.SetColor(0xff000000)
	widget.AddChild(widget.temperature)

	rentingprice := sws.NewLabelWidget(130, 25, "Renting price:")
	rentingprice.Move(root.Width()-450, 125)
	rentingprice.SetTextColor(sdl.Color{0xff, 0xff, 0xff, 0xff})
	rentingprice.SetColor(0xff000000)
	widget.AddChild(rentingprice)
	widget.rentingprice.Move(root.Width()-310, 125)
	widget.rentingprice.SetTextColor(sdl.Color{0xff, 0xff, 0xff, 0xff})
	widget.rentingprice.SetColor(0xff000000)
	widget.AddChild(widget.rentingprice)

	hr := sws.NewHr(160)
	hr.Move(root.Width()-390, 160)
	widget.AddChild(hr)

	difficulty := sws.NewLabelWidget(130, 25, "Difficulty: ")
	difficulty.Move(root.Width()-450, 175)
	difficulty.SetTextColor(sdl.Color{0xff, 0xff, 0xff, 0xff})
	difficulty.SetColor(0xff000000)
	widget.AddChild(difficulty)
	widget.difficulty.Move(root.Width()-310, 175)
	widget.difficulty.SetColor(0xff000000)
	widget.AddChild(widget.difficulty)

	companyname := sws.NewLabelWidget(130, 25, "Company name: ")
	companyname.Move(root.Width()-450, 200)
	companyname.SetTextColor(sdl.Color{0xff, 0xff, 0xff, 0xff})
	companyname.SetColor(0xff000000)
	widget.AddChild(companyname)
	widget.companyname.Move(root.Width()-310, 200)
	widget.companyname.SetColor(0xff000000)
	widget.AddChild(widget.companyname)

	widget.createbutton.Move(root.Width()-350, root.Height()-280)
	widget.createbutton.SetColor(0xff000000)
	widget.AddChild(widget.createbutton)

	widget.cancelbutton.Move(20, root.Height()-280)
	widget.cancelbutton.SetColor(0xff000000)
	widget.AddChild(widget.cancelbutton)

	return widget
}
