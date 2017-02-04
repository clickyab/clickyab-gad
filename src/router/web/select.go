package web

import (
	"entity"
	"fmt"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
)

var slotReg = regexp.MustCompile(`^s\[(\d*)\]$`)

func extractSlice(imp entity.Impression) map[string]string {
	res := make(map[string]string)
	for key, value := range imp.Request().URL.Query() {
		if slice := slotReg.FindStringSubmatch(key); len(slice) == 2 {
			res[slice[1]] = value[0]
		}
	}

	if imp.OS().Valid && imp.OS().Mobile {
		slotPub := fmt.Sprintf("%d1000", imp.Source().ID())
		res[slotPub] = "300x250"
	}
	return res
}

func selectWeb(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
