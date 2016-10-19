package selector

import (
	"middlewares"
	"modules"
	"mr"
	"net/http"

	"errors"
	"filter"
	"regexp"
	"selector"
	"strconv"

	"config"

	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

var (
	webSelector = selector.Mix(
		filter.CheckForSize,
		filter.CheckOS,
		filter.CheckWhiteList,
		filter.CheckNetwork,
		filter.CheckCategory,
		filter.CheckCountry,
	)
)

type selectController struct {
}

// Select function @todo
func (tc *selectController) Select(c echo.Context) error {
	rd := middlewares.MustGetRequestData(c)
	params := c.QueryParams()

	publicParams, ok := params["i"]
	if !ok {
		return c.HTML(400, "invalid request")
	}

	publicID, err := strconv.Atoi(publicParams[0])
	if err != nil {
		return c.HTML(400, "invalid request")
	}

	domain, ok := params["d"]
	if !ok {
		return c.HTML(400, "invalid request")
	}

	//fetch website and set in Context
	website, err := tc.FetchWebsite(publicID)
	if err != nil {
		return c.HTML(400, "invalid request")
	}

	country, err := tc.FetchCountry(rd.RealIP)
	if err != nil {
		logrus.Warn(err)
	}
	//check if the website domain is valid
	if website.WDomain.Valid && website.WDomain.String != domain[0] {
		return errors.New("domain and public id mismatch")
	}

	var size = make(map[string]string)
	var sizeNumSlice []int
	var slotPublic []string
	reg := regexp.MustCompile(`s\[(\d*)\]`)
	for key := range params {
		slice := reg.FindStringSubmatch(key)
		//fmt.Println(slice,len(slice))
		if len(slice) == 2 {
			slotPublic = append(slotPublic, slice[1])
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
		SlotPublic:   slotPublic,
		Country2Info: *country,
	}
	x := selector.Apply(&m, selector.GetAdData(), webSelector, 3)

	adIDBanner := GetAdID(x)

	adBanner, err := mr.NewManager().FetchSlotAd(mr.Build(slotPublic), mr.Build(adIDBanner))
	if err != nil {
		logrus.Info(err)
	}
	fmt.Println(adBanner)
	fmt.Println(len(adBanner))
	return c.JSON(http.StatusOK, adBanner)
}

//FetchWebsite website and set in Context
func (tc *selectController) FetchWebsite(publicID int) (*mr.WebsiteData, error) {
	website, err := mr.NewManager().FetchWebsite(publicID)
	if err != nil {
		return nil, err
	}
	return website, err
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

// GetAdID return ad ids as []string
func GetAdID(ad map[int][]mr.AdData) []string {
	var adIDBanner []string
	for _, adData := range ad {
		for adSliceData := range adData {

			adIDBanner = append(adIDBanner, strconv.FormatInt(adData[adSliceData].AdID, 10))
		}
	}
	return adIDBanner

}

func init() {

	modules.Register(&selectController{})
}
