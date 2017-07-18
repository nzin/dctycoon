package dctycoon

import(
	"github.com/nzin/sws"
	"github.com/nzin/dctycoon/supplier"
//	"github.com/veandco/go-sdl2/sdl"
)

//
// We have to:
// - list all categories (ac, generator, rack, servers)
// - servers: list non attributed servers (or be able to filter? attribute, type, subtype, ...)
//    -> a la gmail? (faire des checkboxwidget)
// - see/build pools (Hardware / VPS)
// - see/build offers
// - see/build contract?
//
// tabwidget?
// upper: title, + buttons
//
type InventoryWidget struct {
	rootwindow        *sws.SWS_RootWidget 
	mainwidget        *sws.SWS_MainWidget
	sub               []*supplier.SubInventory
	currentsub        *supplier.SubInventory
	buttonsub         map[*supplier.SubInventory]*sws.SWS_FlatButtonWidget
	menu              *sws.SWS_CoreWidget
	vbox              *sws.SWS_VBoxWidget
	bottomsplitview   *sws.SWS_SplitviewWidget
	title             *sws.SWS_Label
}

func (self *InventoryWidget) Show() {
	self.rootwindow.AddChild(self.mainwidget)
	self.rootwindow.SetFocus(self.mainwidget)
}

func (self *InventoryWidget) Hide() {
	self.rootwindow.RemoveChild(self.mainwidget)
	children:=self.rootwindow.GetChildren()
	if len(children)>0 {
		self.rootwindow.SetFocus(children[0])
	}
}
func (self *InventoryWidget) AddSubCategory(category *supplier.SubInventory) {
	button:=sws.CreateFlatButtonWidget(220,40,category.Title)
	button.SetClicked(func() {
		self.SelectSubCategory(category)
	})
	self.buttonsub[category]=button
	button.SetImage(category.Icon)
	button.AlignImageLeft(true)
	button.SetCentered(false)
	
	self.vbox.AddChild(button)
	self.sub=append(self.sub,category)
}

func (self *InventoryWidget) SelectSubCategory(category *supplier.SubInventory) {
	if self.currentsub!=nil {
		self.buttonsub[self.currentsub].SetColor(0xffdddddd)
		self.menu.RemoveChild(self.currentsub.ButtonPanel)
	}
	self.title.SetText(category.Title)
	category.ButtonPanel.Move(450,0)
	category.ButtonPanel.SetColor(0xffffffff)
	self.menu.AddChild(category.ButtonPanel)
	self.bottomsplitview.SetRightWidget(category.Widget)
	self.currentsub=category
	self.buttonsub[self.currentsub].SetColor(0xffcccccc)
	sws.PostUpdate()
}

func NewInventoryWidget(root *sws.SWS_RootWidget) *InventoryWidget {
	mainwidget := sws.CreateMainWidget(650,400," Inventory Management ",true,true)
	widget := &InventoryWidget{
		rootwindow: root,
		mainwidget: mainwidget,
		sub: make([]*supplier.SubInventory,0),
		vbox: sws.CreateVBoxWidget(200,100),
		menu: sws.CreateCoreWidget(500,50),
		bottomsplitview: sws.CreateSplitviewWidget(200,200,true),
		buttonsub: make(map[*supplier.SubInventory]*sws.SWS_FlatButtonWidget),
	}
	mainwidget.SetCloseCallback(func() {
		widget.Hide()
	})
	
	sv := sws.CreateSplitviewWidget(200,200,false)
	sv.PlaceSplitBar(50)
	sv.SplitBarMovable(false)
	mainwidget.SetInnerWidget(sv)
	
	widget.menu.SetColor(0xffffffff)
	sv.SetLeftWidget(widget.menu)
	
	widget.bottomsplitview.PlaceSplitBar(220)
	widget.bottomsplitview.SplitBarMovable(false)
	sv.SetRightWidget(widget.bottomsplitview)
	
	widget.bottomsplitview.SetLeftWidget(widget.vbox)
	
	category:= sws.CreateLabel(220,50,"Category")
	category.SetColor(0xffffffff)
	category.SetCentered(true)
	widget.menu.AddChild(category)
	
	title:= sws.CreateLabel(220,50,"Title")
	title.SetColor(0xffffffff)
	title.Move(240,0)
	widget.menu.AddChild(title)
	widget.title=title
	
	unallocated:=supplier.NewUnallocatedInventorySub(supplier.GlobalInventory)
	widget.AddSubCategory(unallocated)
	widget.AddSubCategory(supplier.NewUnallocatedServerSub(supplier.GlobalInventory))
	widget.AddSubCategory(supplier.NewPoolSub(supplier.GlobalInventory))
	widget.AddSubCategory(supplier.NewOfferSub(supplier.GlobalInventory))
	
	widget.SelectSubCategory(unallocated)

	return widget
}
