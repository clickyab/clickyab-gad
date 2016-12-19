package routes

import (
	"selector"

	"gopkg.in/labstack/echo.v3"
)

func (tc *selectController) applyAd(c echo.Context) error {
	params := c.QueryParams()

	rd, website, country, err := tc.getWebDataFromCtx(c)
	if err != nil {
		return err
	}
	_, sizeNumSlice := tc.slotSizeWeb(params, *website, false)
	//call context
	_ = selector.Context{
		RequestData: *rd,
		Website:     website,
		Size:        sizeNumSlice,
		Country:     country,
	}
	// selector.Apply(&m, selector.GetAdData(), webSelector)
	return nil
}
