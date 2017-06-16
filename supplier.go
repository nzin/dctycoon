package dctycoon

import(
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl"
)
//
// top 'menu': [Deal] <shop>  <in transit> <support> <cart>
//
// shop: a gauche: tower, rack server, blade, a droite: | tower | rack server | blade (+pub au dessus)
//
// in transit: orders in transit (list)
//
// suppport: bought last 3 years
//
// cart currently in cart, not paid

type Supplier struct {
	rootwindow   *sws.SWS_RootWidget 
	mainwidget   *sws.SWS_MainWidget
	scrollwidget *sws.SWS_ScrollWidget
	serverpage   *sws.SWS_CoreWidget
}

func (self *Supplier) Show() {
	self.rootwindow.AddChild(self.mainwidget)
}

func (self *Supplier) Hide() {
	self.rootwindow.RemoveChild(self.mainwidget)
}

func CreateSupplier(root *sws.SWS_RootWidget) *Supplier {
	mainwidget := sws.CreateMainWidget(600,400," Your DEAL supplier",true,true)
	scrollwidget := sws.CreateScrollWidget(600,350)
	supplier := &Supplier{
		rootwindow: root,
		mainwidget: mainwidget,
		scrollwidget: scrollwidget,
	}
	mainwidget.SetCloseCallback(func() {
		supplier.Hide()
	})
	sv := sws.CreateSplitviewWidget(200,200,false)
	sv.PlaceSplitBar(50)
	sv.SplitBarMovable(false)
	mainwidget.SetInnerWidget(sv)
	
	// banner
	banner:=sws.CreateCoreWidget(600,50)
	banner.SetColor(0xff0684dc)
	sv.SetLeftWidget(banner)
	suppliericon:=sws.CreateLabel(100,50,"")
	suppliericon.SetColor(0xff0684dc)
	if img,err := img.Load("resources/deal.small2.png"); err==nil {
		suppliericon.SetImage(img)
	}
	banner.AddChild(suppliericon)

	shop:=sws.CreateFlatButtonWidget(100,50,"Shop")
	shop.SetColor(0xff0684dc)
	shop.SetTextColor(sdl.Color{255,255,255,255})
	shop.Move(100,0)
	banner.AddChild(shop)
	
	ups:=sws.CreateFlatButtonWidget(100,50,"Tracking")
	ups.SetColor(0xff0684dc)
	ups.SetTextColor(sdl.Color{255,255,255,255})
	ups.Move(200,0)
	banner.AddChild(ups)
	
	support:=sws.CreateFlatButtonWidget(100,50,"Support")
	support.SetTextColor(sdl.Color{255,255,255,255})
	support.SetColor(0xff0684dc)
	support.Move(300,0)
	banner.AddChild(support)
	
	cart:=sws.CreateFlatButtonWidget(100,50,"")
	cart.SetColor(0xff0684dc)
	if img,err := img.Load("resources/cart.small.png"); err==nil {
		cart.SetImage(img)
	}
	cart.Move(400,0)
	banner.AddChild(cart)
	
	sv.SetRightWidget(scrollwidget)

	// server page
	serverpage:=sws.CreateCoreWidget(600,350)
	supplier.serverpage=serverpage
	servertitle:=sws.CreateLabel(100,40,"DEAL Servers")
	servertitle.Move(40,0)
	serverpage.AddChild(servertitle)

	towerservers:=sws.CreateFlatButtonWidget(120,20,"> Tower Servers")
	towerservers.SetTextColor(sdl.Color{0x06,0x84,0xdc,0xff})
	towerservers.SetCentered(false)
	towerservers.Move(0,40)
	serverpage.AddChild(towerservers)
	
	rackservers:=sws.CreateFlatButtonWidget(120,20,"> Rack Servers")
	rackservers.SetTextColor(sdl.Color{0x06,0x84,0xdc,0xff})
	rackservers.SetCentered(false)
	rackservers.Move(0,60)
	serverpage.AddChild(rackservers)
	
	bladeservers:=sws.CreateFlatButtonWidget(120,20,"> Blade Servers")
	bladeservers.SetTextColor(sdl.Color{0x06,0x84,0xdc,0xff})
	bladeservers.SetCentered(false)
	bladeservers.Move(0,80)
	serverpage.AddChild(bladeservers)
	
	scrollwidget.SetInnerWidget(serverpage)

	// content
	content := sws.CreateCoreWidget(480,300)
	content.SetColor(0xffffffff)
	content.Move(120,40)

	serverpage.AddChild(content)
	

	
	return supplier
}
