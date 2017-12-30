package supplier

import (
	"fmt"

	"github.com/nzin/sws"
	//	"github.com/veandco/go-sdl2/sdl"
	"strconv"
)

type CartPageItemUi struct {
	sws.CoreWidget
	icon         *sws.LabelWidget
	desc         *sws.TextAreaWidget
	price        float64
	priceL       *sws.LabelWidget
	qty          int32
	qtyD         *sws.DropdownWidget
	delete       *sws.ButtonWidget
	total        *sws.LabelWidget
	totalchanged func()
}

func NewCartPageItemUi(icon, desc string, price float64, qty int32, totalcallback func()) *CartPageItemUi {
	choices := make([]string, 10)
	for i := 1; i <= 10; i++ {
		choices[i-1] = fmt.Sprintf("%d", i)
	}
	cartitem := &CartPageItemUi{
		CoreWidget:   *sws.NewCoreWidget(600, 100),
		icon:         sws.NewLabelWidget(100, 100, ""),
		desc:         sws.NewTextAreaWidget(150, 100, desc),
		price:        price,
		priceL:       sws.NewLabelWidget(100, 25, fmt.Sprintf("%.2f $", price)),
		qty:          qty,
		qtyD:         sws.NewDropdownWidget(50, 25, choices),
		delete:       sws.NewButtonWidget(100, 25, "remove"),
		total:        sws.NewLabelWidget(100, 25, fmt.Sprintf("%.2f $", price*float64(qty))),
		totalchanged: totalcallback,
	}
	cartitem.qtyD.SetActiveChoice(qty - 1)

	cartitem.SetColor(0xffffffff)
	cartitem.icon.SetImage(icon)
	cartitem.icon.SetColor(0xffffffff)

	cartitem.desc.Move(100, 0)
	cartitem.desc.SetDisabled(true)
	cartitem.desc.SetColor(0xffffffff)
	cartitem.priceL.Move(250, 0)
	cartitem.priceL.SetCentered(true)
	cartitem.priceL.SetColor(0xffffffff)
	cartitem.qtyD.Move(350, 0)
	cartitem.qtyD.SetColor(0xffffffff)
	cartitem.delete.Move(400, 0)
	cartitem.delete.SetColor(0xffffffff)
	cartitem.total.Move(500, 0)
	cartitem.total.SetCentered(true)
	cartitem.total.SetColor(0xffffffff)

	cartitem.AddChild(cartitem.icon)
	cartitem.AddChild(cartitem.desc)
	cartitem.AddChild(cartitem.priceL)
	cartitem.AddChild(cartitem.qtyD)
	cartitem.AddChild(cartitem.delete)
	cartitem.AddChild(cartitem.total)

	cartitem.qtyD.SetClicked(func() {
		if choice, err := strconv.Atoi(cartitem.qtyD.Choices[cartitem.qtyD.ActiveChoice]); err == nil {
			cartitem.qty = int32(choice)
			cartitem.total.SetText(fmt.Sprintf("%.2f $", price*float64(choice)))
			cartitem.totalchanged()
			cartitem.PostUpdate()
		}
	})

	return cartitem
}

//
// Cart page
//
// the cart inventory is stored into the GlobalInventory object
//
type CartPageWidget struct {
	sws.CoreWidget
	items       []*CartPageItemUi
	vbox        *sws.VBoxWidget
	grandTotalL *sws.LabelWidget
	grandTotal  *sws.LabelWidget
	buy         *sws.ButtonWidget
	inventory   *Inventory
}

func (self *CartPageWidget) SetBuyCallback(callback func()) {
	self.buy.SetClicked(callback)
}

func (self *CartPageWidget) Reset() {
	for _, i := range self.items {
		self.vbox.RemoveChild(i)
	}
	self.items = make([]*CartPageItemUi, 0)
	self.inventory.Cart = make([]*CartItem, 0)
	self.grandTotal.SetText("0 $")
}

func (self *CartPageWidget) AddItem(productitem int32, conf *ServerConf, unitprice float64, nb int32) {
	item := &CartItem{
		Typeitem:   productitem,
		Serverconf: conf,
		Unitprice:  unitprice,
		Nb:         nb,
	}
	self.inventory.Cart = append(self.inventory.Cart, item)
	var ui *CartPageItemUi
	switch productitem {
	case PRODUCT_SERVER:
		ramSizeText := fmt.Sprintf("%d Mo", conf.NbSlotRam*conf.RamSize)
		if conf.NbSlotRam*conf.RamSize >= 2048 {
			ramSizeText = fmt.Sprintf("%d Go", conf.NbSlotRam*conf.RamSize/1024)
		}

		ui = NewCartPageItemUi("resources/"+conf.ConfType.ServerSprite+"0.png",
			fmt.Sprintf("%dx %d cores\n%s RAM\n%d disks", conf.NbProcessors, conf.NbCore, ramSizeText, conf.NbDisks),
			unitprice,
			nb, func() {
				var totalprice float64
				item.Nb = ui.qty
				for _, itemui := range self.items {
					totalprice += itemui.price * float64(itemui.qty)
				}
				self.grandTotal.SetText(fmt.Sprintf("%.2f $", totalprice))
			})
	case PRODUCT_RACK:
		ui = NewCartPageItemUi("resources/rack0.100.png",
			"Rack",
			unitprice,
			nb, func() {
				var totalprice float64
				item.Nb = ui.qty
				for _, itemui := range self.items {
					totalprice += itemui.price * float64(itemui.qty)
				}
				self.grandTotal.SetText(fmt.Sprintf("%.2f $", totalprice))
			})
	case PRODUCT_AC:
		ui = NewCartPageItemUi("resources/ac0.100.png",
			"Air Climatizer",
			unitprice,
			nb, func() {
				var totalprice float64
				item.Nb = ui.qty
				for _, itemui := range self.items {
					totalprice += itemui.price * float64(itemui.qty)
				}
				self.grandTotal.SetText(fmt.Sprintf("%.2f $", totalprice))
			})
	case PRODUCT_GENERATOR:
		ui = NewCartPageItemUi("resources/generator0.100.png",
			"Generator",
			unitprice,
			nb, func() {
				var totalprice float64
				item.Nb = ui.qty
				for _, itemui := range self.items {
					totalprice += itemui.price * float64(itemui.qty)
				}
				self.grandTotal.SetText(fmt.Sprintf("%.2f $", totalprice))
			})
	}

	self.items = append(self.items, ui)
	self.vbox.AddChild(ui)
	ui.delete.SetClicked(func() {
		self.DeleteItem(item)
	})
	self.grandTotal.Move(500, 80+100*int32(len(self.items)))
	self.grandTotalL.Move(450, 80+100*int32(len(self.items)))
	self.buy.Move(500, 120+100*int32(len(self.items)))
	var totalprice float64
	for _, item := range self.items {
		totalprice += item.price * float64(item.qty)
	}
	self.grandTotal.SetText(fmt.Sprintf("%.2f $", totalprice))
	self.Resize(600, 150+100*int32(len(self.items)))
}

func (self *CartPageWidget) DeleteItem(cartitem *CartItem) {
	pos := -1
	for i, v := range self.inventory.Cart {
		if v == cartitem {
			pos = i
		}
	}

	self.vbox.RemoveChild(self.items[pos])
	self.inventory.Cart = append(self.inventory.Cart[:pos], self.inventory.Cart[pos+1:]...)
	self.items = append(self.items[:pos], self.items[pos+1:]...)
	if len(self.items) == 0 {
		self.grandTotal.Move(500, 180)
		self.grandTotal.SetText("0 $")
		self.grandTotalL.Move(450, 180)
		self.buy.Move(500, 120)
		self.Resize(600, 250)
	} else {
		var totalprice float64
		for _, item := range self.items {
			totalprice += item.price * float64(item.qty)
		}
		self.grandTotal.Move(500, 80+100*int32(len(self.items)))
		self.grandTotal.SetText(fmt.Sprintf("%.2f $", totalprice))
		self.grandTotalL.Move(450, 80+100*int32(len(self.items)))
		self.buy.Move(500, 120+100*int32(len(self.items)))
		self.Resize(600, 150+100*int32(len(self.items)))
	}
}

func NewCartPageWidget(width, height int32, inventory *Inventory) *CartPageWidget {
	cartpage := &CartPageWidget{
		CoreWidget: *sws.NewCoreWidget(width, height),
		items:      make([]*CartPageItemUi, 0),
		vbox:       sws.NewVBoxWidget(600, 0),
		inventory:  inventory,
	}
	cartpage.SetColor(0xffffffff)
	title := sws.NewLabelWidget(200, 30, "Shopping Cart")
	title.SetColor(0xffffffff)
	title.SetFont(sws.LatoRegular20)
	title.Move(20, 0)
	title.SetCentered(false)
	cartpage.AddChild(title)

	hProduct := sws.NewLabelWidget(250, 25, "Product")
	hProduct.Move(0, 55)
	cartpage.AddChild(hProduct)

	hPrice := sws.NewLabelWidget(100, 25, "Unit price")
	hPrice.Move(250, 55)
	cartpage.AddChild(hPrice)

	hQty := sws.NewLabelWidget(50, 25, "Qty")
	hQty.Move(350, 55)
	cartpage.AddChild(hQty)

	hRemove := sws.NewLabelWidget(100, 25, "Remove")
	hRemove.Move(400, 55)
	cartpage.AddChild(hRemove)

	hTotal := sws.NewLabelWidget(100, 25, "Total price")
	hTotal.Move(500, 55)
	cartpage.AddChild(hTotal)

	buy := sws.NewButtonWidget(100, 25, "Buy")
	buy.SetColor(0xffffffff)
	buy.Move(500, 120)
	cartpage.AddChild(buy)
	cartpage.buy = buy

	empty := sws.NewLabelWidget(600, 100, "Your shopping cart is empty")
	empty.SetColor(0xffffffff)
	empty.SetCentered(true)
	empty.Move(0, 80)
	cartpage.AddChild(empty)

	cartpage.vbox.Move(0, 80)
	cartpage.AddChild(cartpage.vbox)

	grandTotalL := sws.NewLabelWidget(50, 25, "Total:")
	grandTotalL.SetColor(0xffffffff)
	grandTotalL.Move(450, 180)
	cartpage.AddChild(grandTotalL)
	cartpage.grandTotalL = grandTotalL

	grandTotal := sws.NewLabelWidget(100, 25, "0 $")
	grandTotal.SetColor(0xffffffff)
	grandTotal.SetCentered(true)
	grandTotal.Move(500, 180)
	cartpage.AddChild(grandTotal)
	cartpage.grandTotal = grandTotal

	cartpage.Resize(600, 250)

	return cartpage
}
