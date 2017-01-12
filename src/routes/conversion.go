package routes

import (
	"assert"
	"encoding/base64"
	"net/http"
	"rabbit"
	"time"
	"transport"

	"gopkg.in/labstack/echo.v3"
)

// send the conversion job and return an empty image
func (tc *selectController) conversion(c echo.Context) error {
	data, err := base64.StdEncoding.DecodeString(message)
	assert.Nil(err)
	actionID := c.QueryParam("action_id")
	clickID := c.QueryParam("click_id")
	out := transport.Conversion{
		ConvID:   clickID,
		ActionID: actionID,
	}
	rabbit.MustPublishAfter(out, time.Minute)
	c.Response().Header().Set("Content-Type", "image/png")
	return c.String(http.StatusOK, string(data))

}
