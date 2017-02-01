package httpentity

import (
	"entity"
	"net"
	"net/http"
	"services/random"
)

const (
	headerXForwardedFor = "X-Forwarded-For"
	headerXRealIP       = "X-Real-IP"
)

type httpImpression struct {
	request  *http.Request
	mega     string
	source   entity.Publisher
	location entity.Location
	attrs    map[string]string
}

func (hi *httpImpression) Request() *http.Request {
	return hi.request
}

// MegaIMP return the random id of this imp object
func (hi *httpImpression) MegaIMP() string {
	return hi.mega
}

// ClientID is the key to identify client
func (*httpImpression) ClientID() int64 {

}

// IP return the client ip
func (hi *httpImpression) IP() net.IP {
	ra := hi.request.RemoteAddr
	if ip := hi.request.Header.Get(headerXForwardedFor); ip != "" {
		ra = ip
	} else if ip := hi.request.Header.Get(headerXRealIP); ip != "" {
		ra = ip
	} else {
		ra, _, _ = net.SplitHostPort(ra)
	}
	return net.ParseIP(ra)
}

// UserAgent return the client user agent
func (hi *httpImpression) UserAgent() string {
	return hi.request.UserAgent()
}

// Source return the publisher that this client is going into system from that
func (hi *httpImpression) Source() entity.Publisher {
	return hi.source
}

// Location of the request
func (hi *httpImpression) Location() entity.Location {
	if hi.location == nil {
		hi.location = NewHTTPLocation(hi.request)
	}
	return hi.location
}

func (hi *httpImpression) OS() entity.OS {
	// TODO : detect the os here
	return entity.OS{
		Valid: false,
	}
}

// Attributes is the generic attribute system
func (*httpImpression) Attributes(entity.ImpressionAttributes) interface{} {
	// TODO : implement if needed
	return nil
}

// NewHTTPImpression return an impression object for this session
func NewHTTPImpression(source entity.Publisher, r *http.Request) entity.Impression {
	return &httpImpression{
		request: r,
		mega:    <-random.ID,
		source:  source,
	}
}
