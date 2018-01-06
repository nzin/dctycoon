package dctycoon

import (
	"fmt"

	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/dctycoon/global"
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
	configurepage     *supplier.ServerPageConfigureWidget
	trend             *supplier.Trend
	timer             *timer.GameTimer
	inventory         *supplier.Inventory
	ledger            *accounting.Ledger
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

func NewMainSupplierWidget(root *sws.RootWidget) *MainSupplierWidget {
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
		configurepage:     supplier.NewServerPageConfigureWidget(480, 700),
		trend:             nil,
		timer:             nil,
		inventory:         nil,
		ledger:            nil,
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
	if icon, err := global.LoadImageAsset("assets/ui/deal.small2.png"); err == nil {
		widgeticon.SetImageSurface(icon)
	}
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
	if icon, err := global.LoadImageAsset("assets/ui/cart.small.png"); err == nil {
		cart.SetImageSurface(icon)
	}
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
	widget.configurepage.Move(120, 160)

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
		serverpage.RemoveChild(widget.content)
		//serverpage.RemoveChild(banners)
		widget.configurepage.SetConfType(widget.trend, "T1000", widget.timer.CurrentTime)
		widget.content = widget.configurepage
		serverpage.AddChild(widget.configurepage)
		scrollwidgetshop.SetHorizontalPosition(0)
		scrollwidgetshop.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})

	rackpage.SetConfigureRack1Callback(func() {
		serverpage.RemoveChild(widget.content)
		//serverpage.RemoveChild(banners)
		widget.configurepage.SetConfType(widget.trend, "R100", widget.timer.CurrentTime)
		widget.content = widget.configurepage
		serverpage.AddChild(widget.configurepage)
		scrollwidgetshop.SetHorizontalPosition(0)
		scrollwidgetshop.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})

	rackpage.SetConfigureRack2Callback(func() {
		serverpage.RemoveChild(widget.content)
		//serverpage.RemoveChild(banners)
		widget.configurepage.SetConfType(widget.trend, "R200", widget.timer.CurrentTime)
		widget.content = widget.configurepage
		serverpage.AddChild(widget.configurepage)
		scrollwidgetshop.SetHorizontalPosition(0)
		scrollwidgetshop.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})

	rackpage.SetConfigureRack4Callback(func() {
		serverpage.RemoveChild(widget.content)
		//serverpage.RemoveChild(banners)
		widget.configurepage.SetConfType(widget.trend, "R400", widget.timer.CurrentTime)
		widget.content = widget.configurepage
		serverpage.AddChild(widget.configurepage)
		scrollwidgetshop.SetHorizontalPosition(0)
		scrollwidgetshop.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})

	rackpage.SetConfigureRack6Callback(func() {
		serverpage.RemoveChild(widget.content)
		//serverpage.RemoveChild(banners)
		widget.configurepage.SetConfType(widget.trend, "R600", widget.timer.CurrentTime)
		widget.content = widget.configurepage
		serverpage.AddChild(widget.configurepage)
		scrollwidgetshop.SetHorizontalPosition(0)
		scrollwidgetshop.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})

	bladepage.SetConfigureBlade1Callback(func() {
		serverpage.RemoveChild(widget.content)
		//serverpage.RemoveChild(banners)
		widget.configurepage.SetConfType(widget.trend, "B100", widget.timer.CurrentTime)
		widget.content = widget.configurepage
		serverpage.AddChild(widget.configurepage)
		scrollwidgetshop.SetHorizontalPosition(0)
		scrollwidgetshop.SetVerticalPosition(0)
		//		sws.PostUpdate()
	})

	bladepage.SetConfigureBlade2Callback(func() {
		serverpage.RemoveChild(widget.content)
		//serverpage.RemoveChild(banners)
		widget.configurepage.SetConfType(widget.trend, "B200", widget.timer.CurrentTime)
		widget.content = widget.configurepage
		serverpage.AddChild(widget.configurepage)
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
	widget.configurepage.SetAddCartCallback(func() {
		sv.SetRightWidget(scrollwidgetcart)
		widget.cartpage.AddItem(widget.configurepage.GetProductType(), widget.configurepage.GetConf(), widget.configurepage.GetUnitPrice(), widget.configurepage.GetNbUnit())
		scrollwidgetcart.Resize(scrollwidgetcart.Width(), scrollwidgetcart.Height())
		//		sws.PostUpdate()
	})

	widget.cartpage.SetBuyCallback(func() {
		var totalprice float64
		for _, item := range widget.inventory.Cart {
			totalprice += item.Unitprice * float64(item.Nb)
		}
		accounts := widget.ledger.GetYearAccount(widget.timer.CurrentTime.Year())
		bankAccount := accounts["51"]
		if bankAccount < totalprice {
			// show modal window
			iconsurface, _ := global.LoadImageAsset("assets/ui/icon-triangular-big.png")
			sws.ShowModalErrorSurfaceicon(widget.rootwindow, "Not enough funds", iconsurface, fmt.Sprintf("You cannot buy for %.2f $ of goods: your bank account is currently credited of %.2f $!", totalprice, bankAccount), nil)
		} else {
			// we buy
			for _, item := range widget.inventory.Cart {
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
				widget.ledger.BuyProduct(desc, widget.timer.CurrentTime, item.Unitprice*float64(item.Nb))
			}
			widget.inventory.BuyCart(widget.timer.CurrentTime)
			// we reset the cart
			widget.cartpage.Reset()
			sv.SetRightWidget(scrollwidgettrack)
		}
	})

	widget.trackpage = supplier.NewTrackPageWidget(600, 850)
	scrollwidgettrack.SetInnerWidget(widget.trackpage)

	track.SetClicked(func() {
		sv.SetRightWidget(scrollwidgettrack)
		scrollwidgetcart.SetHorizontalPosition(0)
		scrollwidgetcart.SetVerticalPosition(0)
	})

	return widget
}

func (self *MainSupplierWidget) SetGame(timer *timer.GameTimer, inventory *supplier.Inventory, ledger *accounting.Ledger, trend *supplier.Trend) {
	self.trend = trend
	self.timer = timer
	self.inventory = inventory
	self.ledger = ledger
	self.configurepage.SetGame(trend)
	self.cartpage.SetGame(inventory)
	self.trackpage.SetGame(inventory, timer.CurrentTime)
}
