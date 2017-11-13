package routes

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/clickyab/services/assert"

	"strconv"

	"clickyab.com/gad/mr"

	"gopkg.in/labstack/echo.v3"
)

// send the conversion job and return an empty image
func (tc *selectController) conversion(c echo.Context) error {
	today := time.Now()
	yesterday := time.Now().AddDate(0, 0, -1)
	data, err := base64.StdEncoding.DecodeString(message)
	assert.Nil(err)
	actionID := c.QueryParam("action_id")
	impID := c.QueryParam("imp_id")
	impIDINT, err := strconv.ParseInt(impID, 10, 64)
	if err != nil {
		return err
	}
	go func() {
		//query database for current impression
		imp, err := mr.NewManager().FindImpByIDDate(impIDINT, today.Format("20060102"))
		if err != nil {
			imp, err = mr.NewManager().FindImpByIDDate(impIDINT, yesterday.Format("20060102"))
			if err != nil {
				return
			}
		}
		//apply the conversion query
		err = mr.NewManager().InsertConversion(actionID, imp)
		if err != nil {
			return
		}
	}()
	c.Response().Header().Set("Content-Type", "image/png")
	return c.String(http.StatusOK, string(data))

}
