package models

// 坐标
type Location struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

func newLocation(lng, lat float64) *Location {
	return &Location{Lng: lng, Lat: lat}
}

func (this *Location) GeoJSON() []float64 {
	return []float64{this.Lng, this.Lat}
}
