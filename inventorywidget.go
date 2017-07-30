package dctycoon

import (
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/sws"
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
	rootwindow      *sws.RootWidget
	mainwidget      *sws.MainWidget
	sub             []*supplier.SubInventory
	currentsub      *supplier.SubInventory
	menu            *sws.CoreWidget
	treeview        *sws.TreeViewWidget
	bottomsplitview *sws.SplitviewWidget
	title           *sws.LabelWidget
}

func (self *InventoryWidget) Show() {
	self.rootwindow.AddChild(self.mainwidget)
	self.rootwindow.SetFocus(self.mainwidget)
}

func (self *InventoryWidget) Hide() {
	self.rootwindow.RemoveChild(self.mainwidget)
	children := self.rootwindow.GetChildren()
	if len(children) > 0 {
		self.rootwindow.SetFocus(children[0])
	}
}
func (self *InventoryWidget) AddSubCategory(category *supplier.SubInventory,focus bool) {
	item := sws.NewTreeViewItem(category.Title,category.Icon,func() {
		self.SelectSubCategory(category)
	})

	self.treeview.AddItem(item)
	if focus {
		self.treeview.SetFocusOn(item)
	}
	self.sub = append(self.sub, category)
}

func (self *InventoryWidget) SelectSubCategory(category *supplier.SubInventory) {
	if self.currentsub != nil {
		self.menu.RemoveChild(self.currentsub.ButtonPanel)
	}
	self.title.SetText(category.Title)
	category.ButtonPanel.Move(450, 0)
	category.ButtonPanel.SetColor(0xffffffff)
	self.menu.AddChild(category.ButtonPanel)
	self.bottomsplitview.SetRightWidget(category.Widget)
	self.currentsub = category
	sws.PostUpdate()
}

func NewInventoryWidget(root *sws.RootWidget) *InventoryWidget {
	mainwidget := sws.NewMainWidget(650, 400, " Inventory Management ", true, true)
	widget := &InventoryWidget{
		rootwindow:      root,
		mainwidget:      mainwidget,
		sub:             make([]*supplier.SubInventory, 0),
		treeview:        sws.NewTreeViewWidget(),
		menu:            sws.NewCoreWidget(500, 50),
		bottomsplitview: sws.NewSplitviewWidget(200, 200, true),
	}
	mainwidget.SetCloseCallback(func() {
		widget.Hide()
	})

	sv := sws.NewSplitviewWidget(200, 200, false)
	sv.PlaceSplitBar(50)
	sv.SplitBarMovable(false)
	mainwidget.SetInnerWidget(sv)

	widget.menu.SetColor(0xffffffff)
	sv.SetLeftWidget(widget.menu)

	widget.bottomsplitview.PlaceSplitBar(220)
	widget.bottomsplitview.SplitBarMovable(false)
	sv.SetRightWidget(widget.bottomsplitview)

	widget.bottomsplitview.SetLeftWidget(widget.treeview)

	category := sws.NewLabelWidget(220, 50, "Category")
	category.SetColor(0xffffffff)
	category.SetCentered(true)
	widget.menu.AddChild(category)

	title := sws.NewLabelWidget(220, 50, "Title")
	title.SetColor(0xffffffff)
	title.Move(240, 0)
	widget.menu.AddChild(title)
	widget.title = title

	unallocated := supplier.NewUnallocatedInventorySub(supplier.GlobalInventory)
	widget.AddSubCategory(unallocated,true)
	widget.AddSubCategory(supplier.NewUnallocatedServerSub(supplier.GlobalInventory),false)
	widget.AddSubCategory(supplier.NewPoolSub(supplier.GlobalInventory),false)
	widget.AddSubCategory(supplier.NewOfferSub(supplier.GlobalInventory),false)

	widget.SelectSubCategory(unallocated)

	return widget
}
