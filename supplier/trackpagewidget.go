package supplier

import(
	"fmt"
	"github.com/nzin/sws"
	"time"
)

type TrackPageItemUi struct {
	sws.SWS_CoreWidget
	icon     *sws.SWS_Label
	desc     *sws.SWS_TextAreaWidget
	delivery *sws.SWS_Label
}

func NewTrackPageItemUi(icon, desc string, deliveryDate time.Time) *TrackPageItemUi {
	trackitem:=&TrackPageItemUi{
		SWS_CoreWidget: *sws.CreateCoreWidget(600,100),
		icon: sws.CreateLabel(100,100,""),
		desc: sws.CreateTextAreaWidget(150,100,desc),
		delivery: sws.CreateLabel(100,100,fmt.Sprintf("%d / %d / %d",deliveryDate.Day(),deliveryDate.Month(),deliveryDate.Year())),
	}

	trackitem.SetColor(0xffffffff)
	trackitem.icon.SetImage(icon)
	trackitem.icon.SetColor(0xffffffff)
	
	trackitem.desc.Move(100,0)
	trackitem.desc.SetReadonly(true)
	trackitem.desc.SetColor(0xffffffff)
	trackitem.delivery.Move(250,0)
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
	sws.SWS_CoreWidget
	vbox        *sws.SWS_VBoxWidget
	intransit   map[*InventoryItem]*TrackPageItemUi
}

func (self *TrackPageWidget) ItemInTransit(item *InventoryItem) {
	desc:=""
	icon:=""
	if item.Typeitem==PRODUCT_SERVER{
		ramSizeText:=fmt.Sprintf("%d Mo",item.Serverconf.NbSlotRam*item.Serverconf.RamSize)
                if (item.Serverconf.NbSlotRam*item.Serverconf.RamSize>=2048) {
                        ramSizeText=fmt.Sprintf("%d Go",item.Serverconf.NbSlotRam*item.Serverconf.RamSize/1024)
                }
		icon="resources/"+item.Serverconf.ConfType.ServerSprite+"0.png"
		desc=fmt.Sprintf("%dx %d cores\n%s RAM\n%d disks",item.Serverconf.NbProcessors,item.Serverconf.NbCore,ramSizeText,item.Serverconf.NbDisks)
	}
	self.intransit[item]=NewTrackPageItemUi(icon, desc, item.Deliverydate)
	self.vbox.AddChild(self.intransit[item])
	sws.PostUpdate()
}

func (self *TrackPageWidget) ItemInStock(item *InventoryItem) {
	self.vbox.RemoveChild(self.intransit[item])
	delete(self.intransit,item)
	sws.PostUpdate()
}

func (self *TrackPageWidget) ItemRemoveFromStock(*InventoryItem) {
}

func (self *TrackPageWidget) ItemInstalled(*InventoryItem) {
}

func (self *TrackPageWidget) ItemUninstalled(*InventoryItem) {
}


func NewTrackPageWidget(width,height int32,inventory *Inventory) *TrackPageWidget {
	trackpage:=&TrackPageWidget{
		SWS_CoreWidget: *sws.CreateCoreWidget(width,height),
		vbox: sws.CreateVBoxWidget(600,0),
		intransit: make(map[*InventoryItem]*TrackPageItemUi),
	}
	inventory.AddSubscriber(trackpage)
	trackpage.SetColor(0xffffffff)
	title:=sws.CreateLabel(200,30,"Product Tracking")
	title.SetColor(0xffffffff)
	title.SetFont(sws.LatoRegular20)
	title.Move(20,0)
	title.SetCentered(false)
	trackpage.AddChild(title)

	
	hProduct:=sws.CreateLabel(250,25,"Product")
	hProduct.Move(0,55)
	trackpage.AddChild(hProduct)

	hDelivery:=sws.CreateLabel(100,25,"Delivery date")
	hDelivery.Move(250,55)
	trackpage.AddChild(hDelivery)
	
	empty:=sws.CreateLabel(600,100,"You don't have anything to track")
	empty.SetColor(0xffffffff)
	empty.SetCentered(true)
	empty.Move(0,80)
	trackpage.AddChild(empty)
	
	trackpage.vbox.Move(0,80)
	trackpage.AddChild(trackpage.vbox)
	
	trackpage.Resize(600,250)

	return trackpage
}

