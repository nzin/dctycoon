package supplier

import(
	"github.com/nzin/sws"
	"strconv"
	"time"
	"math"
)

//
// Page Shop>>Explore>>xxx servers>>Configure
//
type ServerPageConfigureWidget struct {
	sws.SWS_CoreWidget
	title         *sws.SWS_Label
	buybutton     *sws.SWS_ButtonWidget
	configureicon *sws.SWS_Label
	nbproc        *sws.SWS_DropdownWidget
	nbcoretitle   *sws.SWS_Label
	nbcores       *sws.SWS_DropdownWidget
	nbcorechoice  []int32
	vttitle       *sws.SWS_Label
	ddsizechoice  []int32
	nbdisk        *sws.SWS_DropdownWidget
	disksize      *sws.SWS_DropdownWidget
	nbram         *sws.SWS_DropdownWidget
	ramsizechoice []int32
	ramsize       *sws.SWS_DropdownWidget
	pricevalue    *sws.SWS_Label
	howmany       *sws.SWS_DropdownWidget
	nbunits       int32
	pricetotal    *sws.SWS_Label
	conftype      *ServerConfType
	conf          *ServerConf
	today         time.Time
}

func (self *ServerPageConfigureWidget) SetBuyCallback(callback func()) {
	self.buybutton.SetClicked(callback)
}

func (self *ServerPageConfigureWidget) SetConfType(conftypename string, today time.Time) {
	for i,c := range(AvailableConfs) {
		if c.ServerName==conftypename {
			self.conftype=&AvailableConfs[i]
		}
	}
	if self.conftype==nil { return }
	
	self.today=today
	self.configureicon.SetImage("resources/"+self.conftype.ServerSprite+"0.png")
	
	vt:=false
	if Trends.Vt.CurrentValue(self.today)>0 {
		vt=true
	}
	// conf
	self.conf = &ServerConf{
		NbProcessors:self.conftype.NbProcessors[0],
		NbCore:      1,
		VtSupport:   vt,
		NbDisks:     self.conftype.NbDisks[0],
		NbSlotRam:   self.conftype.NbSlotRam[0],
		DiskSize:    Trends.Disksize.CurrentValue(self.today)/4,
		RamSize:     Trends.Ramsize.CurrentValue(self.today)/8,
		ConfType:    self.conftype,
		PricePaid:   0,
	}
	//////// configuration
	// processors
	var nbprocs []string
	maxcores:=Trends.Corepercpu.CurrentValue(self.today)
	for i:=self.conftype.NbProcessors[0];i<=self.conftype.NbProcessors[1];i++ {
		nbprocs=append(nbprocs,strconv.Itoa(int(i)))
	}
	self.nbproc.SetChoices(nbprocs)
	self.nbcoretitle.SetText("Deal offers you the best processor in the category, equipped up to "+strconv.Itoa(int(maxcores))+" cores")
	// vt
	if (vt) {
		self.vttitle.SetText("These processors are VT equipped")
	}
	// nb cores
	var choices []string=make([]string,0,3)
	if (maxcores<=4) {
		self.nbcorechoice=[]int32{1,2,4}
		if (maxcores==1) { choices=[]string{"Altium"} }
		if (maxcores<=2) { choices=[]string{"Altium","Altium Duo"} }
		if (maxcores<=4) { choices=[]string{"Altium","Altium Duo","Altium Quad"} }
	} else {
		for i:=maxcores-4;i<=maxcores;i+=2 {
			choices=append(choices,"Altium "+strconv.Itoa(int(i))+" cores")
		}
		self.conf.NbCore=maxcores-4
		self.nbcorechoice=[]int32{maxcores-4,maxcores-2,maxcores}
	}
	self.nbcores.SetChoices(choices)
	
	// nb disks
	var nbdisks []string
	for i:=self.conftype.NbDisks[0];i<=self.conftype.NbDisks[1];i++ {
		nbdisks=append(nbdisks,strconv.Itoa(int(i)))
	}
	self.nbdisk.SetChoices(nbdisks)
	
	// disk size
	maxsize:=Trends.Disksize.CurrentValue(self.today)
	self.ddsizechoice=[]int32{maxsize/4,maxsize/2,maxsize}
	var ddsize=make([]string,3)
	if maxsize>8000000 {
		ddsize[0]=strconv.Itoa(int(maxsize/4000000))+" To"
		ddsize[1]=strconv.Itoa(int(maxsize/2000000))+" To"
		ddsize[2]=strconv.Itoa(int(maxsize/1000000))+" To"
	} else if maxsize>8000 {
		ddsize[0]=strconv.Itoa(int(maxsize/4000))+" Go"
		ddsize[1]=strconv.Itoa(int(maxsize/2000))+" Go"
		ddsize[2]=strconv.Itoa(int(maxsize/1000))+" Go"
	} else {
		ddsize[0]=strconv.Itoa(int(maxsize/4))+" Mo"
		ddsize[1]=strconv.Itoa(int(maxsize/2))+" Mo"
		ddsize[2]=strconv.Itoa(int(maxsize/1))+" Mo"
	}
	self.disksize.SetChoices(ddsize)

	// nb slot ram
	var nbrams []string
	for i:=self.conftype.NbSlotRam[0];i<=self.conftype.NbSlotRam[1];i++ {
		nbrams=append(nbrams,strconv.Itoa(int(i)))
	}
	self.nbram.SetChoices(nbrams)
	
	// disk size
	maxramsize:=Trends.Ramsize.CurrentValue(self.today)
	self.ramsizechoice=[]int32{maxramsize/8,maxramsize/4,maxramsize/2,maxramsize}
	var ramsize=make([]string,4)
	if maxramsize>16000 {
		ramsize[0]=strconv.Itoa(int(maxramsize/8000))+" Go"
		ramsize[1]=strconv.Itoa(int(maxramsize/4000))+" Go"
		ramsize[2]=strconv.Itoa(int(maxramsize/2000))+" Go"
		ramsize[3]=strconv.Itoa(int(maxramsize/1000))+" Go"
	} else {
		ramsize[0]=strconv.Itoa(int(maxramsize/4))+" Mo"
		ramsize[1]=strconv.Itoa(int(maxramsize/4))+" Mo"
		ramsize[2]=strconv.Itoa(int(maxramsize/2))+" Mo"
		ramsize[3]=strconv.Itoa(int(maxramsize))+" Mo"
	}
	self.ramsize.SetChoices(ramsize)

	// price
	self.conf.PricePaid=math.Floor(self.conf.Price(self.today))
	self.pricevalue.SetText(strconv.FormatFloat(self.conf.PricePaid,'f',0,64))
	
	// how many
	self.howmany.SetChoices([]string{"1","2","3","4","5","6","7","8","9","10"})
	self.nbunits=1
	
	// price total
	self.pricetotal.SetText(strconv.FormatFloat(self.conf.PricePaid*float64(self.nbunits),'f',0,64)) 
	//////// configuration
	sws.PostUpdate()
}

func (self *ServerPageConfigureWidget) GetConf() *ServerConf {
	return self.conf
}

func CreateServerPageConfigureWidget(width,height int32) *ServerPageConfigureWidget {
	serverpageconfigure:=&ServerPageConfigureWidget{
		SWS_CoreWidget: *sws.CreateCoreWidget(width,height),
	}
	serverpageconfigure.SetColor(0xffffffff)
	
        title:=sws.CreateLabel(200,20,"Configure Server")
        title.SetFont(sws.LatoRegular20)
	title.SetColor(0xffffffff)
        title.Move(20,0)
        title.SetCentered(false)
        serverpageconfigure.AddChild(title)
        serverpageconfigure.title=title

	configureIcon:=sws.CreateLabel(150,100,"")
	configureIcon.SetColor(0xffffffff)
        configureIcon.SetCentered(true)
	configureIcon.Move(0,20)
        serverpageconfigure.AddChild(configureIcon)
        serverpageconfigure.configureicon=configureIcon


	// nb processors
	nbproctitle:=sws.CreateLabel(150,25,"Nb Processors:")
	nbproctitle.SetColor(0xffffffff)
	nbproctitle.Move(0,140)
        serverpageconfigure.AddChild(nbproctitle)
	
	nbproc:=sws.CreateDropdownWidget(75,25,[]string{})
	nbproc.SetColor(0xffffffff)
	nbproc.Move(150,140)
        serverpageconfigure.AddChild(nbproc)
	serverpageconfigure.nbproc=nbproc
	nbproc.SetClicked(func() {
		if choice,err:=strconv.Atoi(nbproc.Choices[nbproc.ActiveChoice]); err==nil {
			serverpageconfigure.conf.NbProcessors=int32(choice)
		}
		serverpageconfigure.conf.PricePaid=math.Floor(serverpageconfigure.conf.Price(serverpageconfigure.today))
		serverpageconfigure.pricevalue.SetText(strconv.FormatFloat(serverpageconfigure.conf.PricePaid,'f',0,64))
		serverpageconfigure.pricetotal.SetText(strconv.FormatFloat(serverpageconfigure.conf.PricePaid*float64(serverpageconfigure.nbunits),'f',0,64))
	})

	nbcoretitle:=sws.CreateLabel(430,20,"")
	nbcoretitle.SetFont(sws.LatoRegular12)
	nbcoretitle.SetColor(0xffffffff)
	nbcoretitle.Move(0,165)
        serverpageconfigure.AddChild(nbcoretitle)
        serverpageconfigure.nbcoretitle=nbcoretitle
	
	vttitle:=sws.CreateLabel(430,20,"")
	vttitle.SetFont(sws.LatoRegular12)
	vttitle.SetColor(0xffffffff)
	vttitle.Move(0,185)
        serverpageconfigure.AddChild(vttitle)
        serverpageconfigure.vttitle=vttitle
	

	// processor type
	processor:=sws.CreateLabel(150,25,"Processor")
	processor.SetColor(0xffffffff)
	processor.Move(0,205)
	serverpageconfigure.AddChild(processor)

	nbcores:=sws.CreateDropdownWidget(150,25,[]string{})
	nbcores.SetColor(0xffffffff)
	nbcores.Move(150,205)
        serverpageconfigure.AddChild(nbcores)
	serverpageconfigure.nbcores=nbcores
	nbcores.SetClicked(func() {
		serverpageconfigure.conf.NbCore=serverpageconfigure.nbcorechoice[nbcores.ActiveChoice]
		serverpageconfigure.conf.PricePaid=math.Floor(serverpageconfigure.conf.Price(serverpageconfigure.today))
		serverpageconfigure.pricevalue.SetText(strconv.FormatFloat(serverpageconfigure.conf.PricePaid,'f',0,64))
		serverpageconfigure.pricetotal.SetText(strconv.FormatFloat(serverpageconfigure.conf.PricePaid*float64(serverpageconfigure.nbunits),'f',0,64))
	})


	// nb disks
	nbdisktitle:=sws.CreateLabel(150,25,"Nb disks")
	nbdisktitle.SetColor(0xffffffff)
	nbdisktitle.Move(0,230)
	serverpageconfigure.AddChild(nbdisktitle)

	nbdisk:=sws.CreateDropdownWidget(75,25,[]string{})
	nbdisk.SetColor(0xffffffff)
	nbdisk.Move(150,230)
        serverpageconfigure.AddChild(nbdisk)
	serverpageconfigure.nbdisk=nbdisk
	nbdisk.SetClicked(func() {
		if choice,err:=strconv.Atoi(nbdisk.Choices[nbdisk.ActiveChoice]); err==nil {
			serverpageconfigure.conf.NbDisks=int32(choice)
		}
		serverpageconfigure.conf.PricePaid=math.Floor(serverpageconfigure.conf.Price(serverpageconfigure.today))
		serverpageconfigure.pricevalue.SetText(strconv.FormatFloat(serverpageconfigure.conf.PricePaid,'f',0,64))
		serverpageconfigure.pricetotal.SetText(strconv.FormatFloat(serverpageconfigure.conf.PricePaid*float64(serverpageconfigure.nbunits),'f',0,64))
	})

	
	// disk size
	disksizetitle:=sws.CreateLabel(150,25,"Disks size")
	disksizetitle.SetColor(0xffffffff)
	disksizetitle.Move(0,255)
	serverpageconfigure.AddChild(disksizetitle)

	disksize:=sws.CreateDropdownWidget(75,25,[]string{})
	disksize.SetColor(0xffffffff)
	disksize.Move(150,255)
        serverpageconfigure.AddChild(disksize)
	serverpageconfigure.disksize=disksize
	disksize.SetClicked(func() {
		serverpageconfigure.conf.DiskSize=serverpageconfigure.ddsizechoice[disksize.ActiveChoice]
		serverpageconfigure.conf.PricePaid=math.Floor(serverpageconfigure.conf.Price(serverpageconfigure.today))
		serverpageconfigure.pricevalue.SetText(strconv.FormatFloat(serverpageconfigure.conf.PricePaid,'f',0,64))
		serverpageconfigure.pricetotal.SetText(strconv.FormatFloat(serverpageconfigure.conf.PricePaid*float64(serverpageconfigure.nbunits),'f',0,64))
	})

	
	// nb ram
	nbramtitle:=sws.CreateLabel(150,25,"Nb SDRAM DIMM")
	nbramtitle.SetColor(0xffffffff)
	nbramtitle.Move(0,280)
	serverpageconfigure.AddChild(nbramtitle)

	nbram:=sws.CreateDropdownWidget(75,25,[]string{})
	nbram.SetColor(0xffffffff)
	nbram.Move(150,280)
        serverpageconfigure.AddChild(nbram)
	serverpageconfigure.nbram=nbram
	nbram.SetClicked(func() {
		if choice,err:=strconv.Atoi(nbram.Choices[nbram.ActiveChoice]); err==nil {
			serverpageconfigure.conf.NbSlotRam=int32(choice)
		}
		serverpageconfigure.conf.PricePaid=math.Floor(serverpageconfigure.conf.Price(serverpageconfigure.today))
		serverpageconfigure.pricevalue.SetText(strconv.FormatFloat(serverpageconfigure.conf.PricePaid,'f',0,64))
		serverpageconfigure.pricetotal.SetText(strconv.FormatFloat(serverpageconfigure.conf.PricePaid*float64(serverpageconfigure.nbunits),'f',0,64))
	})

	
	// ram size
	ramsizetitle:=sws.CreateLabel(150,25,"SDRAM size")
	ramsizetitle.SetColor(0xffffffff)
	ramsizetitle.Move(0,305)
	serverpageconfigure.AddChild(ramsizetitle)

	ramsize:=sws.CreateDropdownWidget(75,25,[]string{})
	ramsize.SetColor(0xffffffff)
	ramsize.Move(150,305)
        serverpageconfigure.AddChild(ramsize)
	serverpageconfigure.ramsize=ramsize
	ramsize.SetClicked(func() {
		serverpageconfigure.conf.RamSize=serverpageconfigure.ramsizechoice[ramsize.ActiveChoice]
		serverpageconfigure.conf.PricePaid=math.Floor(serverpageconfigure.conf.Price(serverpageconfigure.today))
		serverpageconfigure.pricevalue.SetText(strconv.FormatFloat(serverpageconfigure.conf.PricePaid,'f',0,64))
		serverpageconfigure.pricetotal.SetText(strconv.FormatFloat(serverpageconfigure.conf.PricePaid*float64(serverpageconfigure.nbunits),'f',0,64))
	})

	
	// shipping 
	shipping:=sws.CreateLabel(200,20,"Ships in 3-5 business days")
	shipping.SetFont(sws.LatoRegular12)
	shipping.SetColor(0xffffffff)
	shipping.Move(150,330)
        serverpageconfigure.AddChild(shipping)


	// price 
	priceText:=sws.CreateLabel(100,25,"Unit Price")
	priceText.SetColor(0xffffffff)
	priceText.Move(0,355)
        serverpageconfigure.AddChild(priceText)

	priceValue:=sws.CreateLabel(100,25,"0")
	priceValue.SetColor(0xffffffff)
	priceValue.Move(150,355)
        serverpageconfigure.AddChild(priceValue)
	serverpageconfigure.pricevalue=priceValue


	// how many
	nbunittitle:=sws.CreateLabel(150,25,"Nb Units")
	nbunittitle.SetColor(0xffffffff)
	nbunittitle.Move(0,380)
	serverpageconfigure.AddChild(nbunittitle)

	nbunits:=sws.CreateDropdownWidget(75,25,[]string{})
	nbunits.SetColor(0xffffffff)
	nbunits.Move(150,380)
        serverpageconfigure.AddChild(nbunits)
	serverpageconfigure.howmany=nbunits
	nbunits.SetClicked(func() {
		if choice,err:=strconv.Atoi(nbunits.Choices[nbunits.ActiveChoice]); err==nil {
			serverpageconfigure.nbunits=int32(choice)
		}
		serverpageconfigure.conf.PricePaid=math.Floor(serverpageconfigure.conf.Price(serverpageconfigure.today))
		serverpageconfigure.pricevalue.SetText(strconv.FormatFloat(serverpageconfigure.conf.PricePaid,'f',0,64))
		serverpageconfigure.pricetotal.SetText(strconv.FormatFloat(serverpageconfigure.conf.PricePaid*float64(serverpageconfigure.nbunits),'f',0,64))
	})

	
	// price total
	pricetotalText:=sws.CreateLabel(100,25,"Final Price")
	pricetotalText.SetColor(0xffffffff)
	pricetotalText.Move(0,405)
        serverpageconfigure.AddChild(pricetotalText)

	pricetotalValue:=sws.CreateLabel(100,25,"0")
	pricetotalValue.SetColor(0xffffffff)
	pricetotalValue.Move(150,405)
        serverpageconfigure.AddChild(pricetotalValue)
	serverpageconfigure.pricetotal=pricetotalValue


	
	// buy button
	buyButton:=sws.CreateButtonWidget(100,25,"Buy >")
	buyButton.SetColor(0xffffffff)
	buyButton.Move(150,430)
	serverpageconfigure.AddChild(buyButton)
	serverpageconfigure.buybutton=buyButton


	return serverpageconfigure
}

