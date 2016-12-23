package routes

import (
	"assert"
	"encoding/base64"
	"rabbit"
	"transport"

	"gopkg.in/labstack/echo.v3"
)

const message = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQMAAAAl21bKAAAAA1BMVEUAAACnej3aAAAAAXRSTlMAQObYZgAAAApJREFUCNdjYAAAAAIAAeIhvDMAAAAASUVORK5CYII="

// send the conversion job and return an empty image
func (tc *selectController) conversion(c echo.Context) error {

	//click_id
	data, err := base64.StdEncoding.DecodeString(message)
	assert.Nil(err)
	actionID := c.QueryParam("action_id")
	clickID := c.QueryParam("click_id")
	out := transport.Conversion{
		ConvID:   clickID,
		ActionID: actionID,
	}
	assert.Nil(rabbit.Publish("cy.conv", out))
	c.Response().Header().Set("Content-Type", "image/png")
	return c.String(200, string(data))

}
