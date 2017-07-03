package supplier

import(
	"github.com/nzin/sws"
//	"github.com/veandco/go-sdl2/sdl"
)

//
// Page Shop
//
type CartPageWidget struct {
	sws.SWS_CoreWidget
}

func CreateCartPageWidget(width,height int32) *CartPageWidget {
	cartpage:=&CartPageWidget{
		SWS_CoreWidget: *sws.CreateCoreWidget(width,height),
	}

	return cartpage
}

