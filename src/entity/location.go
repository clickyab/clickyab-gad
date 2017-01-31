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
	Country() Country
	Province() Province
	City() City
	Hood() Hood
	LatLon() LatLon
}
