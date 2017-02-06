package httpentity

import (
	"entity"
	"models"
	"net"
	"net/http"
)

type httpLocation struct {
	ip       net.IP
	location *models.IP2Location
	err      error
	country  *entity.Country
}

func (hl *httpLocation) loadLoaction() {
	if hl.err == nil && hl.location == nil {
		hl.location, hl.err = models.NewManager().GetLocation(hl.ip)
	}
}

func (hl *httpLocation) Country() entity.Country {
	hl.loadLoaction()
	// TODO : Currently the only country supported is IRAN, So take it easy now, but watch for it later
	if hl.err == nil {
		if hl.country == nil {
			hl.country = &entity.Country{
				Valid: true,
				ID:    1,
				ISO:   "ir",
				Name:  "Iran",
			}
		}

		return *hl.country
	}
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
