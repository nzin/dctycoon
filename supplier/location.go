package supplier

type LocationType struct {
	Name                string
	Metersquareprice    float64 // renting price per month
	Internetfailrate    float64 // per day
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
		Name:                "Mountain View",
		Metersquareprice:    30, // https://www.zillow.com/mountain-view-ca/home-values/ ?
		Internetfailrate:    0.005,
		Electricityfailrate: 0.010,
		Temperatureaverage:  16.11,  // http://www.usclimatedata.com/climate/california/united-states/3174
		Electricitycost:     0.1534, // https://en.wikipedia.org/wiki/Electricity_pricing
		Taxrate:             0.12,
		Bankinterestrate:    0.03,
		Xmap:                772,
		Ymap:                1431,
	},
	"newyork": &LocationType{
		Name:                "New York",
		Metersquareprice:    26, // https://www.numbeo.com/cost-of-living/city_price_rankings?itemId=100 x2
		Internetfailrate:    0.003,
		Electricityfailrate: 0.005,
		Temperatureaverage:  12.77, // https://www.usclimatedata.com/climate/new-york/united-states/3202
		Electricitycost:     0.13,  // https://en.wikipedia.org/wiki/Electricity_pricing
		Taxrate:             0.10,
		Bankinterestrate:    0.03,
		Xmap:                1228,
		Ymap:                1388,
	},
	"london": &LocationType{
		Name:                "London",
		Metersquareprice:    36, // https://www.numbeo.com/cost-of-living/city_price_rankings?itemId=100 x2
		Internetfailrate:    0.002,
		Electricityfailrate: 0.004,
		Temperatureaverage:  10.4, // http://www.london.climatemps.com/temperatures.php
		Electricitycost:     0.22, // https://en.wikipedia.org/wiki/Electricity_pricing
		Taxrate:             0.19, // https://en.wikipedia.org/wiki/List_of_countries_by_tax_rates
		Bankinterestrate:    0.03,
		Xmap:                1933,
		Ymap:                1233,
	},
	"lulea": &LocationType{
		Name:                "Luleå",
		Metersquareprice:    1,
		Internetfailrate:    0.002,
		Electricityfailrate: 0.003,
		Temperatureaverage:  0,      //
		Electricitycost:     0.0833, // https://en.wikipedia.org/wiki/Electricity_pricing
		Taxrate:             0.22,   // https://en.wikipedia.org/wiki/List_of_countries_by_tax_rates
		Bankinterestrate:    0.03,
		Xmap:                2146,
		Ymap:                1003,
	},
	"shanghai": &LocationType{
		Name:                "Shanghai",
		Metersquareprice:    30, // https://www.numbeo.com/cost-of-living/city_price_rankings?itemId=100 x2
		Internetfailrate:    0.002,
		Electricityfailrate: 0.003,
		Temperatureaverage:  22,
		Electricitycost:     0.042, // https://en.wikipedia.org/wiki/Electricity_pricing
		Taxrate:             0.25,  // https://en.wikipedia.org/wiki/List_of_countries_by_tax_rates
		Bankinterestrate:    0.04,
		Xmap:                3085,
		Ymap:                1500,
	},
	"sydney": &LocationType{
		Name:                "Sydney",
		Metersquareprice:    28, // https://www.numbeo.com/cost-of-living/city_price_rankings?itemId=100 x 2
		Internetfailrate:    0.002,
		Electricityfailrate: 0.003,
		Temperatureaverage:  23,   //
		Electricitycost:     0.23, // https://en.wikipedia.org/wiki/Electricity_pricing
		Taxrate:             0.29, // https://en.wikipedia.org/wiki/List_of_countries_by_tax_rates
		Bankinterestrate:    0.04,
		Xmap:                3369,
		Ymap:                2183,
	},
	"bangalore": &LocationType{
		Name:                "Bangalore",
		Metersquareprice:    3.4, // https://www.numbeo.com/cost-of-living/city_price_rankings?itemId=100 x 2
		Internetfailrate:    0.003,
		Electricityfailrate: 0.004,
		Temperatureaverage:  30,
		Electricitycost:     0.07, // https://en.wikipedia.org/wiki/Electricity_pricing
		Taxrate:             0.08, // https://en.wikipedia.org/wiki/List_of_countries_by_tax_rates
		Bankinterestrate:    0.04,
		Xmap:                2673,
		Ymap:                1706,
	},
	"moscow": &LocationType{
		Name:                "Moscow",
		Metersquareprice:    11, // https://www.numbeo.com/cost-of-living/city_price_rankings?itemId=100 x 2
		Internetfailrate:    0.003,
		Electricityfailrate: 0.005,
		Temperatureaverage:  9,    // ~ https://www.worldweatheronline.com/moscow-weather-averages/moscow-city/ru.aspx
		Electricitycost:     0.08, // https://en.wikipedia.org/wiki/Electricity_pricing
		Taxrate:             0.20, // https://en.wikipedia.org/wiki/List_of_countries_by_tax_rates
		Bankinterestrate:    0.06,
		Xmap:                2294,
		Ymap:                1172,
	},
}
