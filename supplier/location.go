package supplier

type LocationType struct {
	Name                string
	Metersquareprice    float64
	Internetfailrate    float64
	Electricityfailrate float64 // per day
	Temperatureaverage  float64 // in may?
	Electricitycost     float64 // per kwh
	Taxrate             float64
	Bankinterestrate    float64 // per year
	Xmap                int32   // x worldmap
	Ymap                int32   // y worldmap
}

var GlobalLocation string

var AvailableLocation = map[string]*LocationType{
	"siliconvalley": &LocationType{
		Name:                "Silicon Valley",
		Metersquareprice:    862, //https://www.zillow.com/mountain-view-ca/home-values/
		Internetfailrate:    0.005,
		Electricityfailrate: 0.005,
		Temperatureaverage:  16.11, //http://www.usclimatedata.com/climate/california/united-states/3174
		Electricitycost:     0.1534,
		Taxrate:             0.15,
		Bankinterestrate:    0.03,
		Xmap:                772,
		Ymap:                1431,
	},
}
