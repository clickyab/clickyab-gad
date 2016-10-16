package selector

import (
	"middlewares"
	"modules"
	"mr"
	"net/http"

	"errors"
	"filter"
	"fmt"
	"regexp"
	"selector"
	"strconv"

	"config"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

type selectController struct {
}

// Select functioon @todo
func (tc *selectController) Select(c echo.Context) error {

	params := c.QueryParams()

	publicParams, ok := params["i"]
	if !ok {
		return errors.New("params i not found")
	}
	publicID, err := strconv.Atoi(publicParams[0])
	if err != nil {
		return errors.New("public_id not found")
	}
	domain, ok := params["d"]
	if !ok {
		return errors.New("domain not found")
	}
	rd := middlewares.MustGetRequestData(c)
	//fetch website and set in Context
	website := tc.FetchWebsite(publicID)
	country, err := tc.FetchCountry(rd.RealIP)
	if err != nil {
		logrus.Info(err)
	}
	//check if the website domain is valid
	if website.WDomain.Valid && website.WDomain.String != domain[0] {
		return errors.New("domain and public id mismatch")
	}

	var size = make(map[string]string)
	var sizeNumSlice []int
	reg := regexp.MustCompile(`s\[(\d*)\]`)
	for key := range params {
		slice := reg.FindStringSubmatch(key)
		//fmt.Println(slice,len(slice))
		if len(slice) == 2 {
			size[slice[1]] = params[key][0]
			//check for size
			SizeNum, _ := config.GetSize(size[slice[1]])
			sizeNumSlice = append(sizeNumSlice, SizeNum)

		}

	}

	//call context
	m := selector.Context{
		RequestData:  *rd,
		WebsiteData:  *website,
		Size:         sizeNumSlice,
		Country2Info: *country,
	}
	x := selector.Apply(&m, selector.GetAdData(), selector.Mix(filter.CheckForSize, filter.CheckOS, filter.CheckWhiteList, filter.CheckNetwork, filter.CheckCategory, filter.CheckCountry), 3)
	fmt.Println(len(x))
	return c.JSON(http.StatusOK, x)
}

//FetchWebsite website and set in Context
func (tc *selectController) FetchWebsite(publicID int) *mr.WebsiteData {
	website, err := mr.NewManager().FetchWebsite(publicID)
	if err != nil {
		logrus.Fatal(err)
	}
	return website
}

//FetchCountry find country and set context
func (tc *selectController) FetchCountry(c string) (*mr.Country2Info, error) {
	var country mr.Country2Info
	ip, err := mr.NewManager().GetLocation(c)
	if err != nil || !ip.CountryName.Valid {
		return &country, errors.New("Country not found")
	}
	country, err = mr.NewManager().ConvertCountry2Info(ip.CountryName.String)
	if err != nil {
		return &country, errors.New("Country not found")
	}
	return &country, nil

}

// Routes function @todo
func (tc *selectController) Routes(e *echo.Echo, _ string) {
	e.Get("/select", tc.Select)
}

func init() {

	modules.Register(&selectController{})
}
