package dctycoon

import (
	"fmt"

	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/dctycoon/timer"
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"
)

//
// top 'menu': [Deal] <shop>  <in transit> <support> <cart>
//
// shop: to the left: tower, rack server, blade, to the right: | tower | rack server | blade (+pub au dessus)
//
// in transit: orders in transit (list)
//
// suppport: bought last 3 years
//
// cart currently in cart, not paid

type MainSupplierWidget struct {
	rootwindow        *sws.RootWidget
	mainwidget        *sws.MainWidget
	splitviewwidget   *sws.SplitviewWidget
	scrollwidgetshop  *sws.ScrollWidget
	scrollwidgetcart  *sws.ScrollWidget
	scrollwidgettrack *sws.ScrollWidget
	serverpage        *supplier.ServerPageWidget
	cartpage          *supplier.CartPageWidget
	trackpage         *supplier.TrackPageWidget
	content           sws.Widget
	bannerwidget      *supplier.BannerWidget
	explorewidget     *supplier.ServerPageExploreWidget
}

func (self *MainSupplierWidget) Show() {
	self.rootwindow.AddChild(self.mainwidget)
	self.rootwindow.SetFocus(self.mainwidget)
	self.splitviewwidget.SetRightWidget(self.scrollwidgetshop)
	self.serverpage.RemoveChild(self.content)
	self.serverpage.AddChild(self.bannerwidget)
	self.content = self.explorewidget
	self.serverpage.AddChild(self.explorewidget)
	self.scrollwidgetshop.SetHorizontalPosition(0)
	self.scrollwidgetshop.SetVerticalPosition(0)
}

func (self *MainSupplierWidget) Hide() {
	self.rootwindow.RemoveChild(self.mainwidget)
	children := self.rootwindow.GetChildren()
	if len(children) > 0 {
		self.rootwindow.SetFocus(children[len(children)-1])
	}
}

func NewMainSupplierWidget(trend *supplier.Trend, root *sws.RootWidget) *MainSupplierWidget {
	mainwidget := sws.NewMainWidget(650, 400, " Your DEAL supplier", true, true)
	scrollwidgetshop := sws.NewScrollWidget(600, 550)
	scrollwidgetshop.SetColor(0xffffffff)
	scrollwidgetcart := sws.NewScrollWidget(600, 550)
	scrollwidgetcart.SetColor(0xffffffff)
	scrollwidgettrack := sws.NewScrollWidget(600, 550)
	scrollwidgettrack.SetColor(0xffffffff)
	sv := sws.NewSplitviewWidget(200, 200, false)
	sv.PlaceSplitBar(50)
	sv.SplitBarMovable(false)

	widget := &MainSupplierWidget{
		rootwindow:        root,
		mainwidget:        mainwidget,
		scrollwidgetshop:  scrollwidgetshop,
		scrollwidgetcart:  scrollwidgetcart,
		scrollwidgettrack: scrollwidgettrack,
		splitviewwidget:   sv,
	}
	mainwidget.SetCloseCallback(func() {
		widget.Hide()
	})
	mainwidget.SetInnerWidget(sv)

	// banner
	banner := sws.NewCoreWidget(600, 50)
	banner.SetColor(0xff0684dc)
	sv.SetLeftWidget(banner)
	widgeticon := sws.NewLabelWidget(100, 50, "")
	widgeticon.SetColor(0xff0684dc)
	widgeticon.SetImage("resources/deal.small2.png")
	banner.AddChild(widgeticon)

	shop := sws.NewFlatButtonWidget(100, 50, "Shop")
	shop.SetColor(0xff0684dc)
	shop.SetTextColor(sdl.Color{255, 255, 255, 255})
	shop.Move(100, 0)
	banner.AddChild(shop)

	track := sws.NewFlatButtonWidget(100, 50, "Tracking")
	track.SetColor(0xff0684dc)
	track.SetTextColor(sdl.Color{255, 255, 255, 255})
	track.Move(200, 0)
	banner.AddChild(track)

	support := sws.NewFlatButtonWidget(100, 50, "Support")
	support.SetTextColor(sdl.Color{255, 255, 255, 255})
	support.SetColor(0xff0684dc)
	support.Move(300, 0)
	banner.AddChild(support)

	cart := sws.NewFlatButtonWidget(100, 50, "")
	cart.SetColor(0xff0684dc)
	cart.SetImage("resources/cart.small.png")
	cart.Move(400, 0)
	banner.AddChild(cart)

	sv.SetRightWidget(scrollwidgetshop)

	// server page
	serverpage := supplier.NewServerPageWidget(600, 850)
	widget.serverpage = serverpage
	scrollwidgetshop.SetInnerWidget(serverpage)

	// content
	banners := supplier.NewBannerWidget(480, 120)
	banners.Move(120, 40)
	serverpage.AddChild(banners)
	widget.bannerwidget = banners

	explore := supplier.NewServerPageExploreWidget(480, 700)
	explore.Move(120, 160)
	serverpage.AddChild(explore)
	widget.content = explore
	widget.explorewidget = explore

	towerpage := supplier.NewServerPageTowerWidget(480, 700)
	towerpage.Move(120, 160)

	rackpage := supplier.NewServerPageRackWidget(480, 700)
	rackpage.Move(120, 160)

	bladepage := supplier.NewServerPageBladeWidget(480, 700)
	bladepage.Move(120, 160)

	// configure
	configurepage := supplier.NewServerPageConfigureWidget(trend, 480, 700)
	configurepage.Move(120, 160)

	// buttons callback

	shop.SetClicked(func() {
		sv.SetRightWidget(scrollwidgetshop)
		serverpage.RemoveChild(widget.content)
		serverpage.AddChild(banners)
		widget.content = explore
		serverpage.AddChild(explore)
		scrollwidgetshop.SetHorizontalPosition(0)
		scrollwidgetshop.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})

	serverpage.SetTowerCallback(func() {
		serverpage.RemoveChild(widget.content)
		serverpage.AddChild(banners)
		serverpage.AddChild(towerpage)
		scrollwidgetshop.SetHorizontalPosition(0)
		scrollwidgetshop.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})

	serverpage.SetRackServerCallback(func() {
		serverpage.RemoveChild(widget.content)
		serverpage.AddChild(banners)
		widget.content = rackpage
		serverpage.AddChild(rackpage)
		scrollwidgetshop.SetHorizontalPosition(0)
		scrollwidgetshop.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})

	serverpage.SetBladeCallback(func() {
		serverpage.RemoveChild(widget.content)
		serverpage.AddChild(banners)
		widget.content = bladepage
		serverpage.AddChild(bladepage)
		scrollwidgetshop.SetHorizontalPosition(0)
		scrollwidgetshop.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})

	serverpage.SetAcCallback(func() {
		sv.SetRightWidget(scrollwidgetcart)
		widget.cartpage.AddItem(supplier.PRODUCT_AC, nil, 2000, 1)
		scrollwidgetcart.Resize(scrollwidgetcart.Width(), scrollwidgetcart.Height())
		//		sws.PostUpdate()
	})

	serverpage.SetRackCallback(func() {
		sv.SetRightWidget(scrollwidgetcart)
		widget.cartpage.AddItem(supplier.PRODUCT_RACK, nil, 500, 1)
		scrollwidgetcart.Resize(scrollwidgetcart.Width(), scrollwidgetcart.Height())
		//		sws.PostUpdate()
	})

	serverpage.SetGeneratorCallback(func() {
		sv.SetRightWidget(scrollwidgetcart)
		widget.cartpage.AddItem(supplier.PRODUCT_GENERATOR, nil, 3000, 1)
		scrollwidgetcart.Resize(scrollwidgetcart.Width(), scrollwidgetcart.Height())
		//		sws.PostUpdate()
	})

	explore.SetTowerCallback(func() {
		serverpage.RemoveChild(widget.content)
		serverpage.AddChild(banners)
		widget.content = towerpage
		serverpage.AddChild(towerpage)
		scrollwidgetshop.SetHorizontalPosition(0)
		scrollwidgetshop.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})

	explore.SetRackServerCallback(func() {
		serverpage.RemoveChild(widget.content)
		serverpage.AddChild(banners)
		widget.content = rackpage
		serverpage.AddChild(rackpage)
		scrollwidgetshop.SetHorizontalPosition(0)
		scrollwidgetshop.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})

	explore.SetBladeCallback(func() {
		serverpage.RemoveChild(widget.content)
		serverpage.AddChild(banners)
		widget.content = bladepage
		serverpage.AddChild(bladepage)
		scrollwidgetshop.SetHorizontalPosition(0)
		scrollwidgetshop.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})

	explore.SetAcCallback(func() {
		sv.SetRightWidget(scrollwidgetcart)
		widget.cartpage.AddItem(supplier.PRODUCT_AC, nil, 2000, 1)
		scrollwidgetcart.Resize(scrollwidgetcart.Width(), scrollwidgetcart.Height())
		//		sws.PostUpdate()
	})

	explore.SetRackCallback(func() {
		sv.SetRightWidget(scrollwidgetcart)
		widget.cartpage.AddItem(supplier.PRODUCT_RACK, nil, 500, 1)
		scrollwidgetcart.Resize(scrollwidgetcart.Width(), scrollwidgetcart.Height())
		//		sws.PostUpdate()
	})

	explore.SetGeneratorCallback(func() {
		sv.SetRightWidget(scrollwidgetcart)
		widget.cartpage.AddItem(supplier.PRODUCT_GENERATOR, nil, 3000, 1)
		scrollwidgetcart.Resize(scrollwidgetcart.Width(), scrollwidgetcart.Height())
		//		sws.PostUpdate()
	})

	// callback configure
	towerpage.SetConfigureTower1Callback(func() {
		now := timer.GlobalGameTimer.CurrentTime
		serverpage.RemoveChild(widget.content)
		//serverpage.RemoveChild(banners)
		configurepage.SetConfType(trend, "T1000", now)
		widget.content = configurepage
		serverpage.AddChild(configurepage)
		scrollwidgetshop.SetHorizontalPosition(0)
		scrollwidgetshop.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})

	rackpage.SetConfigureRack1Callback(func() {
		now := timer.GlobalGameTimer.CurrentTime
		serverpage.RemoveChild(widget.content)
		//serverpage.RemoveChild(banners)
		configurepage.SetConfType(trend, "R100", now)
		widget.content = configurepage
		serverpage.AddChild(configurepage)
		scrollwidgetshop.SetHorizontalPosition(0)
		scrollwidgetshop.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})

	rackpage.SetConfigureRack2Callback(func() {
		now := timer.GlobalGameTimer.CurrentTime
		serverpage.RemoveChild(widget.content)
		//serverpage.RemoveChild(banners)
		configurepage.SetConfType(trend, "R200", now)
		widget.content = configurepage
		serverpage.AddChild(configurepage)
		scrollwidgetshop.SetHorizontalPosition(0)
		scrollwidgetshop.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})

	rackpage.SetConfigureRack4Callback(func() {
		now := timer.GlobalGameTimer.CurrentTime
		serverpage.RemoveChild(widget.content)
		//serverpage.RemoveChild(banners)
		configurepage.SetConfType(trend, "R400", now)
		widget.content = configurepage
		serverpage.AddChild(configurepage)
		scrollwidgetshop.SetHorizontalPosition(0)
		scrollwidgetshop.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})

	rackpage.SetConfigureRack6Callback(func() {
		now := timer.GlobalGameTimer.CurrentTime
		serverpage.RemoveChild(widget.content)
		//serverpage.RemoveChild(banners)
		configurepage.SetConfType(trend, "R600", now)
		widget.content = configurepage
		serverpage.AddChild(configurepage)
		scrollwidgetshop.SetHorizontalPosition(0)
		scrollwidgetshop.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})

	bladepage.SetConfigureBlade1Callback(func() {
		now := timer.GlobalGameTimer.CurrentTime
		serverpage.RemoveChild(widget.content)
		//serverpage.RemoveChild(banners)
		configurepage.SetConfType(trend, "B100", now)
		widget.content = configurepage
		serverpage.AddChild(configurepage)
		scrollwidgetshop.SetHorizontalPosition(0)
		scrollwidgetshop.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})

	bladepage.SetConfigureBlade2Callback(func() {
		now := timer.GlobalGameTimer.CurrentTime
		serverpage.RemoveChild(widget.content)
		//serverpage.RemoveChild(banners)
		configurepage.SetConfType(trend, "B200", now)
		widget.content = configurepage
		serverpage.AddChild(configurepage)
		scrollwidgetshop.SetHorizontalPosition(0)
		scrollwidgetshop.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})

	widget.cartpage = supplier.NewCartPageWidget(600, 850)
	scrollwidgetcart.SetInnerWidget(widget.cartpage)

	cart.SetClicked(func() {
		sv.SetRightWidget(scrollwidgetcart)
		scrollwidgetcart.SetHorizontalPosition(0)
		scrollwidgetcart.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})
	configurepage.SetAddCartCallback(func() {
		sv.SetRightWidget(scrollwidgetcart)
		widget.cartpage.AddItem(configurepage.GetProductType(), configurepage.GetConf(), configurepage.GetUnitPrice(), configurepage.GetNbUnit())
		scrollwidgetcart.Resize(scrollwidgetcart.Width(), scrollwidgetcart.Height())
		//		sws.PostUpdate()
	})

	widget.cartpage.SetBuyCallback(func() {
		var totalprice float64
		for _, item := range supplier.GlobalInventory.Cart {
			totalprice += item.Unitprice * float64(item.Nb)
		}
		accounts := accounting.GlobalLedger.GetYearAccount(timer.GlobalGameTimer.CurrentTime.Year())
		bankAccount := accounts["51"]
		if bankAccount < totalprice {
			// show modal window
			sws.ShowModalError(widget.rootwindow, "Not enough funds", "resources/icon-triangular-big.png", fmt.Sprintf("You cannot buy for %.2f $ of goods: your bank account is currently credited of %.2f $!", totalprice, bankAccount), nil)
		} else {
			// we buy
			for _, item := range supplier.GlobalInventory.Cart {
				var desc string
				switch item.Typeitem {
				case supplier.PRODUCT_SERVER:
					desc = fmt.Sprintf("%dx %s", item.Nb, item.Serverconf.ConfType.ServerName)
				case supplier.PRODUCT_RACK:
					desc = fmt.Sprintf("%dx Rack", item.Nb)
				case supplier.PRODUCT_AC:
					desc = fmt.Sprintf("%dx AC", item.Nb)
				case supplier.PRODUCT_GENERATOR:
					desc = fmt.Sprintf("%dx Generator", item.Nb)
				}
				accounting.GlobalLedger.BuyProduct(desc, timer.GlobalGameTimer.CurrentTime, item.Unitprice*float64(item.Nb))
			}
			supplier.GlobalInventory.BuyCart(timer.GlobalGameTimer.CurrentTime)
			// we reset the cart
			widget.cartpage.Reset()
			sv.SetRightWidget(scrollwidgettrack)
		}
	})

	widget.trackpage = supplier.NewTrackPageWidget(600, 850, supplier.GlobalInventory)
	scrollwidgettrack.SetInnerWidget(widget.trackpage)

	track.SetClicked(func() {
		sv.SetRightWidget(scrollwidgettrack)
		scrollwidgetcart.SetHorizontalPosition(0)
		scrollwidgetcart.SetVerticalPosition(0)
	})

	return widget
}
