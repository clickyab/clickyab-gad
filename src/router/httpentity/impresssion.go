package httpentity

import (
	"config"
	"crypto/sha1"
	"entity"
	"fmt"
	"models"
	"net"
	"net/http"
	"services/random"

	"github.com/mssola/user_agent"
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
	clientID int64
	attrs    map[string]string
	ip       net.IP
	os       *entity.OS
}

func (hi *httpImpression) Request() *http.Request {
	return hi.request
}

// MegaIMP return the random id of this imp object
func (hi *httpImpression) MegaIMP() string {
	return hi.mega
}

// ClientID is the key to identify client
func (hi *httpImpression) ClientID() int64 {
	if hi.clientID == 0 {
		copID := hi.request.URL.Query().Get("tid")
		if len(copID) < config.Config.Clickyab.CopLen {
			copID = httpCreateHash(config.Config.Clickyab.CopLen, []byte(hi.UserAgent()), []byte(hi.IP()))
		}
		hi.clientID = models.NewManager().CreateCookieProfile(copID, hi.IP()).ID
	}

	return hi.clientID
}

// IP return the client ip
func (hi *httpImpression) IP() net.IP {
	if hi.ip == nil {
		ra := hi.request.RemoteAddr
		if ip := hi.request.Header.Get(headerXForwardedFor); ip != "" {
			ra = ip
		} else if ip := hi.request.Header.Get(headerXRealIP); ip != "" {
			ra = ip
		} else {
			ra, _, _ = net.SplitHostPort(ra)
		}
		hi.ip = net.ParseIP(ra)
	}

	return hi.ip
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
	if hi.os == nil {
		os := findHTTPOS(user_agent.New(hi.UserAgent()))
		hi.os = &os
	}

	return *hi.os
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

func httpCreateHash(l int, data ...[]byte) string {
	h := sha1.New()
	for i := range data {
		_, _ = h.Write(data[i])
	}
	sum := fmt.Sprintf("%x", h.Sum(nil))
	if l >= len(sum) {
		l = len(sum)
	}
	if l < 1 {
		l = 1
	}
	return sum[:l]
}
