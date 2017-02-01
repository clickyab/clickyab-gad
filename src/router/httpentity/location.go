package httpentity

import (
	"entity"
	"net/http"
)

type httpLocation struct {
}

func (hl *httpLocation) Country() entity.Country {
	return entity.Country{
		Valid: false,
	}
}
func (hl *httpLocation) Province() entity.Province {
	return entity.Province{
		Valid: false,
	}

}
func (hl *httpLocation) City() entity.City {
	return entity.City{
		Valid: false,
	}

}
func (hl *httpLocation) Hood() entity.Hood {
	return entity.Hood{
		Valid: false,
	}

}
func (hl *httpLocation) LatLon() entity.LatLon {
	return entity.LatLon{
		Valid: false,
	}

}

// NewHTTPLocation is the location based on request
func NewHTTPLocation(r *http.Request) entity.Location {
	// TODO actually load location
	return &httpLocation{}
}
