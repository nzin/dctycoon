package supplier

import (
	"strconv"
	"time"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/sws"
)

//
// ad widget
//
type BannerWidget struct {
	sws.CoreWidget
	banners       []*sws.LabelWidget
	currentBanner int32
	te            *sws.TimerEvent
}

func NewBannerWidget(width, height int32) *BannerWidget {
	widget := &BannerWidget{
		CoreWidget:    *sws.NewCoreWidget(width, height),
		banners:       []*sws.LabelWidget{},
		currentBanner: 0,
	}
	widget.SetColor(0xffffffff)

	for i := 1; i <= 3; i++ {
		banner := sws.NewLabelWidget(width, 100, "")
		if img, err := global.LoadImageAsset("assets/ui/banner" + strconv.Itoa(i) + ".png"); err == nil {
			banner.SetImageSurface(img)
		}
		banner.SetColor(0xffffffff)
		banner.Move((width-400)/2, (height-100)/2)
		widget.banners = append(widget.banners, banner)
	}
	widget.AddChild(widget.banners[0])

	widget.te = sws.TimerAddEvent(time.Now(), 6000*time.Millisecond, func(evt *sws.TimerEvent) {
		widget.RemoveChild(widget.banners[widget.currentBanner])
		widget.currentBanner++
		if widget.currentBanner >= int32(len(widget.banners)) {
			widget.currentBanner = 0
		}
		widget.AddChild(widget.banners[widget.currentBanner])
	})

	return widget
}
