package dctycoon

type LocationType struct {
	Metersquareprice    float64
	Internetfailrate    float64
	Electricityfailrate float64 // per day
	Temperatureaverage  float64 // in may?
	Electricitycost     float64 // per kwh
	Taxrate             float64
	Bankinterestrate    float64 // per month
}

var GlobalLocation string

var AvailableLocation = map[string]LocationType {
	"siliconvalley": LocationType{
		Metersquareprice:   862, //https://www.zillow.com/mountain-view-ca/home-values/
		Internetfailrate:   0.005,
		Electricityfailrate:0.005,
		Temperatureaverage: 16.11, //http://www.usclimatedata.com/climate/california/united-states/3174
		Electricitycost:    0.1534,
		Taxrate:            0.15,
		Bankinterestrate:   0.0025,
	},
}
