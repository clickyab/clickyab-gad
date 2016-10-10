package selector

import (
	"fmt"
	"middlewares"
	"modules"
	"mr"
	"net/http"

	"banners"
	"errors"
	"filter"
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"regexp"
	"selector"
	"strconv"
)

type selectController struct {
}

func filterNonApp(c *selector.Context, in mr.AdData) bool {
	return in.CpType == 0
}

type size []int

func (tc *selectController) Select(c echo.Context) error {

	params := c.QueryParams()

	public_params, ok := params["i"]
	if !ok {
		return errors.New("params i not found")
	}
	public_id, err := strconv.Atoi(public_params[0])
	if err != nil {
		return errors.New("public_id not found")
	}

	////fetch website and set in Context
	website, err := mr.NewManager().FetchWebsite(public_id)
	if err != nil {
		logrus.Fatal(err)
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
			SizeNum, _ := banners.GetSize(size[slice[1]])
			sizeNumSlice = append(sizeNumSlice, SizeNum)

		}

	}

	rd := c.Get("RequestData").(*middlewares.RequestData)

	//Fetch Category



	//call context
	m := selector.Context{
		RequestData: *rd,
		WebsiteData: *website,
		Size:        sizeNumSlice,
	}
	x := selector.Apply(&m, selector.GetAdData(), selector.Mix(filter.CheckForSize, filter.CheckOS, filter.CheckWhiteList, filter.CheckNetwork), 3)
	fmt.Println(len(x))
	return c.JSON(http.StatusOK, 1)
}

func (tc *selectController) Routes(e *echo.Echo, _ string) {
	e.Get("/select", tc.Select)
}

func init() {

	modules.Register(&selectController{})
}
