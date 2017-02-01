package web

import (
	"fmt"
	"net/http"
	"regexp"

	"entity"

	"github.com/julienschmidt/httprouter"
)

var slotReg = regexp.MustCompile(`^s\[(\d*)\]$`)

func extractSlice(imp entity.Impression) map[string]string {
	res := make(map[string]string)
	for key, value := range imp.Request().URL.Query() {
		if slice := slotReg.FindStringSubmatch(key); len(slice) == 2 {
			res[slice[1]] = value
		}
	}

	//if imp.OS().Valid && imp.OS().Mobile {
	//	slotPub := fmt.Sprintf("%d%s", imp.Source().ID(), webMobile)
	//	slotPublic = append(slotPublic, slotPub)
	//	sizeNumSlice[slotPub] = 8
	//}
	//return tc.slotSizeNormal(slotPublic, website.WID, sizeNumSlice)
}

func selectWeb(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
