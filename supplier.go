package dctycoon

import(
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/nzin/dctycoon/supplier"
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
	serverpage   *supplier.ServerPageWidget
	content      sws.SWS_Widget
}

func (self *Supplier) Show() {
	self.rootwindow.AddChild(self.mainwidget)
}

func (self *Supplier) Hide() {
	self.rootwindow.RemoveChild(self.mainwidget)
	children:=self.rootwindow.GetChildren()
	if len(children)>0 {
		self.rootwindow.SetFocus(children[0])
	}
}

func CreateSupplier(root *sws.SWS_RootWidget) *Supplier {
	mainwidget := sws.CreateMainWidget(600,400," Your DEAL supplier",true,true)
	scrollwidget := sws.CreateScrollWidget(600,550)
	widget := &Supplier{
		rootwindow: root,
		mainwidget: mainwidget,
		scrollwidget: scrollwidget,
	}
	mainwidget.SetCloseCallback(func() {
		widget.Hide()
	})
	sv := sws.CreateSplitviewWidget(200,200,false)
	sv.PlaceSplitBar(50)
	sv.SplitBarMovable(false)
	mainwidget.SetInnerWidget(sv)
	
	// banner
	banner:=sws.CreateCoreWidget(600,50)
	banner.SetColor(0xff0684dc)
	sv.SetLeftWidget(banner)
	widgeticon:=sws.CreateLabel(100,50,"")
	widgeticon.SetColor(0xff0684dc)
	widgeticon.SetImage("resources/deal.small2.png")
	banner.AddChild(widgeticon)

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
	cart.SetImage("resources/cart.small.png")
	cart.Move(400,0)
	banner.AddChild(cart)
	
	sv.SetRightWidget(scrollwidget)

	// server page
	serverpage:=supplier.CreateServerPageWidget(600,550)
	widget.serverpage=serverpage
	scrollwidget.SetInnerWidget(serverpage)
	
	// content
	banners:=supplier.CreateBannerWidget(480,120)
	banners.Move(120,40)
	serverpage.AddChild(banners)
	
	explore:=supplier.CreateServerPageExploreWidget(480,500)
	explore.Move(120,160)
	serverpage.AddChild(explore)
	widget.content=explore
	
	towerpage:=supplier.CreateServerPageTowerWidget(480,500)
	towerpage.Move(120,160)
	
	rackpage:=supplier.CreateServerPageRackWidget(480,500)
	rackpage.Move(120,160)
	
	bladepage:=supplier.CreateServerPageBladeWidget(480,500)
	bladepage.Move(120,160)
	
	// configure
	
	// buttons callback
	
	shop.SetClicked(func() {
		serverpage.RemoveChild(widget.content)
		serverpage.AddChild(explore)
		sws.PostUpdate()
	})

	serverpage.SetTowerCallback(func() {
		serverpage.RemoveChild(widget.content)
		serverpage.AddChild(towerpage)
		sws.PostUpdate()
	})
	
	serverpage.SetRackCallback(func() {
		serverpage.RemoveChild(widget.content)
		serverpage.AddChild(rackpage)
		sws.PostUpdate()
	})
	
	serverpage.SetBladeCallback(func() {
		serverpage.RemoveChild(widget.content)
		serverpage.AddChild(bladepage)
		sws.PostUpdate()
	})
	
	explore.SetTowerCallback(func() {
		serverpage.RemoveChild(widget.content)
		serverpage.AddChild(towerpage)
		sws.PostUpdate()
	})
	
	explore.SetRackCallback(func() {
		serverpage.RemoveChild(widget.content)
		serverpage.AddChild(rackpage)
		sws.PostUpdate()
	})
	
	explore.SetBladeCallback(func() {
		serverpage.RemoveChild(widget.content)
		serverpage.AddChild(bladepage)
		sws.PostUpdate()
	})
	
	return widget
}
