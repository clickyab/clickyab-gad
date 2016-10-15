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

	//fetch website and set in Context
	website, err := mr.NewManager().FetchWebsite(publicID)
	if err != nil {
		logrus.Fatal(err)
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

	rd := middlewares.MustGetRequestData(c)

	//call context
	m := selector.Context{
		RequestData: *rd,
		WebsiteData: *website,
		Size:        sizeNumSlice,
	}
	x := selector.Apply(&m, selector.GetAdData(), selector.Mix(filter.CheckForSize, filter.CheckOS, filter.CheckWhiteList, filter.CheckNetwork), 3)
	fmt.Println(len(x))
	return c.JSON(http.StatusOK, x)
}

// Routes function @todo
func (tc *selectController) Routes(e *echo.Echo, _ string) {
	e.Get("/select", tc.Select)
}

func init() {

	modules.Register(&selectController{})
}
