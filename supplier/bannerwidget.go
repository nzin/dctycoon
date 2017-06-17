package supplier

import(
	"github.com/nzin/sws"
	"time"
	"strconv"
)

type BannerWidget struct {
	sws.SWS_CoreWidget
	banners       []*sws.SWS_Label
	currentBanner int32
	te            *sws.TimerEvent
}

func CreateBannerWidget(width,height int32) *BannerWidget {
	widget:=&BannerWidget{
		SWS_CoreWidget: *sws.CreateCoreWidget(width,height),
		banners:        []*sws.SWS_Label{},
		currentBanner:  0,
	}
	widget.SetColor(0xffffffff)

	for i:=1;i<=3;i++ {
		banner:=sws.CreateLabel(width,100,"")
		banner.SetImage("resources/banner"+strconv.Itoa(i)+".png")
		banner.SetColor(0xffffffff)
		banner.Move((width-400)/2,(height-100)/2)
		widget.banners=append(widget.banners,banner)
	}
        widget.AddChild(widget.banners[0])

	widget.te = sws.TimerAddEvent(time.Now(), 6000*time.Millisecond, func() {
        	widget.RemoveChild(widget.banners[widget.currentBanner])
		widget.currentBanner++
		if widget.currentBanner>=int32(len(widget.banners)) { widget.currentBanner=0 }
        	widget.AddChild(widget.banners[widget.currentBanner])
		sws.PostUpdate()
	})

	return widget
}

