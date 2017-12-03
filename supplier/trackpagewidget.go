package supplier

import (
	"fmt"
	"time"

	"github.com/nzin/sws"
)

type TrackPageItemUi struct {
	sws.CoreWidget
	icon     *sws.LabelWidget
	desc     *sws.TextAreaWidget
	delivery *sws.LabelWidget
}

func NewTrackPageItemUi(icon, desc string, deliveryDate time.Time) *TrackPageItemUi {
	trackitem := &TrackPageItemUi{
		CoreWidget: *sws.NewCoreWidget(600, 100),
		icon:       sws.NewLabelWidget(100, 100, ""),
		desc:       sws.NewTextAreaWidget(150, 100, desc),
		delivery:   sws.NewLabelWidget(100, 100, fmt.Sprintf("%d / %d / %d", deliveryDate.Day(), deliveryDate.Month(), deliveryDate.Year())),
	}

	trackitem.SetColor(0xffffffff)
	trackitem.icon.SetImage(icon)
	trackitem.icon.SetColor(0xffffffff)

	trackitem.desc.Move(100, 0)
	trackitem.desc.SetReadonly(true)
	trackitem.desc.SetColor(0xffffffff)
	trackitem.delivery.Move(250, 0)
	trackitem.delivery.SetCentered(true)
	trackitem.delivery.SetColor(0xffffffff)

	trackitem.AddChild(trackitem.icon)
	trackitem.AddChild(trackitem.desc)
	trackitem.AddChild(trackitem.delivery)

	return trackitem
}

//
// Track page
//
// the track inventory is stored into the GlobalInventory object
//
type TrackPageWidget struct {
	sws.CoreWidget
	vbox      *sws.VBoxWidget
	intransit map[*InventoryItem]*TrackPageItemUi
}

func (self *TrackPageWidget) ItemInTransit(item *InventoryItem) {
	desc := ""
	icon := ""
	switch item.Typeitem {
	case PRODUCT_SERVER:
		ramSizeText := fmt.Sprintf("%d Mo", item.Serverconf.NbSlotRam*item.Serverconf.RamSize)
		if item.Serverconf.NbSlotRam*item.Serverconf.RamSize >= 2048 {
			ramSizeText = fmt.Sprintf("%d Go", item.Serverconf.NbSlotRam*item.Serverconf.RamSize/1024)
		}
		icon = "resources/" + item.Serverconf.ConfType.ServerSprite + "0.png"
		desc = fmt.Sprintf("%dx %d cores\n%s RAM\n%d disks", item.Serverconf.NbProcessors, item.Serverconf.NbCore, ramSizeText, item.Serverconf.NbDisks)
	case PRODUCT_AC:
		icon = "resources/ac0.100.png"
		desc = "Air climatiser"
	case PRODUCT_RACK:
		icon = "resources/rack0.100.png"
		desc = "Rack chassis"
	case PRODUCT_GENERATOR:
		icon = "resources/generator0.100.png"
		desc = "Generator"
	}
	self.intransit[item] = NewTrackPageItemUi(icon, desc, item.Deliverydate)
	self.vbox.AddChild(self.intransit[item])
	self.Resize(600, 80+self.vbox.Height())
}

func (self *TrackPageWidget) ItemInStock(item *InventoryItem) {
	self.vbox.RemoveChild(self.intransit[item])
	self.Resize(600, 80+self.vbox.Height())
	delete(self.intransit, item)
}

func (self *TrackPageWidget) ItemRemoveFromStock(*InventoryItem) {
}

func (self *TrackPageWidget) ItemInstalled(*InventoryItem) {
}

func (self *TrackPageWidget) ItemUninstalled(*InventoryItem) {
}

func (self *TrackPageWidget) ItemChangedPool(*InventoryItem) {
}

func NewTrackPageWidget(width, height int32, inventory *Inventory) *TrackPageWidget {
	trackpage := &TrackPageWidget{
		CoreWidget: *sws.NewCoreWidget(width, height),
		vbox:       sws.NewVBoxWidget(600, 0),
		intransit:  make(map[*InventoryItem]*TrackPageItemUi),
	}
	inventory.AddInventorySubscriber(trackpage)
	trackpage.SetColor(0xffffffff)
	title := sws.NewLabelWidget(200, 30, "Product Tracking")
	title.SetColor(0xffffffff)
	title.SetFont(sws.LatoRegular20)
	title.Move(20, 0)
	title.SetCentered(false)
	trackpage.AddChild(title)

	hProduct := sws.NewLabelWidget(250, 25, "Product")
	hProduct.Move(0, 55)
	trackpage.AddChild(hProduct)

	hDelivery := sws.NewLabelWidget(100, 25, "Delivery date")
	hDelivery.Move(250, 55)
	trackpage.AddChild(hDelivery)

	empty := sws.NewLabelWidget(600, 100, "You don't have anything to track")
	empty.SetColor(0xffffffff)
	empty.SetCentered(true)
	empty.Move(0, 80)
	trackpage.AddChild(empty)

	trackpage.vbox.Move(0, 80)
	trackpage.AddChild(trackpage.vbox)

	trackpage.Resize(600, 250)

	return trackpage
}
