package entity

// Country is the country object
type Country struct {
	Valid bool
	ID    int64
	Name  string
	ISO   string
}

// Province of the request
type Province struct {
	Valid bool
	ID    int64
	Name  string
}

// City if available
type City struct {
	Valid bool
	ID    int64
	Name  string
}

// Hood id
type Hood struct {
	Valid bool
	ID    int64
	Name  string
}

// LatLon is the latitude longitude
type LatLon struct {
	Valid    bool
	Lat, Lon float64
}

// Location is the location provider
type Location interface {
	// Country get the country if available
	Country() Country
	// Province get the province of request if available
	Province() Province
	// City return the city if available
	City() City
	// Hood return the hood if any
	Hood() Hood
	// LatLon return the latitude longitude if any
	LatLon() LatLon
}
