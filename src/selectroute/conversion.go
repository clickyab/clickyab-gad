package selectroute

import (
	"assert"
	"encoding/base64"
	"rabbit"
	"transport"

	"github.com/labstack/echo"
)

const MESSAGE = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQMAAAAl21bKAAAAA1BMVEUAAACnej3aAAAAAXRSTlMAQObYZgAAAApJREFUCNdjYAAAAAIAAeIhvDMAAAAASUVORK5CYII="

// TODO : send the conversion job and return an empty image
func (tc *selectController) conversion(c echo.Context) error {


	//click_id
	data, err := base64.StdEncoding.DecodeString(MESSAGE)
	assert.Nil(err)
	action_id:= c.QueryParam("action_id")
	clickID:= c.QueryParam("click_id")
	out := transport.Conversion{
		ConvID:   clickID,
		ActionID: action_id,
	}
	assert.Nil(rabbit.Publish("cy.conv", out))
	c.Response().Header().Set("Content-Type", "image/png")
	return c.String(200, string(data))

}
