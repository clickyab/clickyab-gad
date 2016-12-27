package fcgi

import (
	"errors"
	"net/http"
	"time"
)

type notFound struct {
}

func (notFound) ServeHTTP(http.ResponseWriter, *http.Request) (int, error) {
	return http.StatusNotFound, errors.New("Not forund")
}

// NewPHPFastCGIHandler return a simple handler
func NewPHPFastCGIHandler(root, path string, php string, rTimeout, sTimeout, dialTimeout time.Duration) Handler {
	r := Rule{
		Path: path,

		Address: "tcp://" + php,

		Ext:         ".php",
		SplitPath:   ".php",
		IndexFiles:  []string{"index.php"},
		ReadTimeout: rTimeout,
		SendTimeout: sTimeout,
		dialer: basicDialer{
			address: php,
			network: "tcp",
			timeout: dialTimeout,
		},
	}

	h := Handler{
		Next:    notFound{},
		Rules:   []Rule{r},
		Root:    root,
		AbsRoot: root,
		FileSys: http.Dir(root),
	}

	return h
}
