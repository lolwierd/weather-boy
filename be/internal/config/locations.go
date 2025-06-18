package config

type Location struct {
	Name       string
	Lat, Lon   float64
	PdfSlug    string
	RadarCodes []string
}

// Locations lists the supported cities for Weather Boy.
var Locations = []Location{
	{Name: "vadodara", Lat: 22.30, Lon: 73.20, PdfSlug: "gujarat.pdf", RadarCodes: []string{"baroda", "ahmedabad"}},
	{Name: "mumbai", Lat: 19.08, Lon: 72.88, PdfSlug: "maharashtra.pdf", RadarCodes: []string{"mumbai"}},
	{Name: "thane", Lat: 19.22, Lon: 72.97, PdfSlug: "maharashtra.pdf", RadarCodes: []string{"mumbai"}},
	{Name: "pune", Lat: 18.52, Lon: 73.85, PdfSlug: "maharashtra.pdf", RadarCodes: []string{"mumbai"}},
}
