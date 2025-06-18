package config

type Location struct {
	Name       string // "vadodara"
	Lat, Lon   float64
	PdfSlug    string   // "gujarat.pdf"
	RadarCodes []string // "baroda", fallback "ahmedabad"
}

var Locations = []Location{
	{Name: "vadodara", Lat: 22.3, Lon: 73.2, PdfSlug: "gujarat.pdf", RadarCodes: []string{"baroda", "ahmedabad"}},
}
