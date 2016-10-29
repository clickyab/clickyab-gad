package selector

import (
	"config"
	"errors"
	"filter"
	"middlewares"
	"modules"
	"mr"
	"net/http"
	"regexp"
	"selector"
	"strconv"

	"fmt"

	"redis"

	"time"

	"utils"

	"net"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

var (
	webSelector = selector.Mix(
		filter.CheckForSize,
		filter.CheckOS,
		filter.CheckWhiteList,
		filter.CheckBlackList,
		filter.CheckNetwork,
		filter.CheckCategory,
		filter.CheckCountry,
	)

	slotReg = regexp.MustCompile(`s\[(\d*)\]`)
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
	country, err := tc.FetchCountry(rd.IP)
	if err != nil {
		logrus.Warn(err)
	}
	//check if the website domain is valid
	if website.WDomain.Valid && website.WDomain.String != domain[0] {
		return errors.New("domain and public id mismatch")
	}

	slotPublic, sizeNumSlice := tc.slotSize(params)

	//call context
	m := selector.Context{
		RequestData:  *rd,
		WebsiteData:  *website,
		Size:         sizeNumSlice,
		SlotPublic:   slotPublic,
		Country2Info: *country,
	}
	x := selector.Apply(&m, selector.GetAdData(), webSelector, 3)

	minCap := findMinCap(m.CopID)
	fmt.Println(minCap)

	//let the game begin :)
	for s := range slotPublic {
		fmt.Println(s)
		i := 0
		/*for size := range x[sizeNumSlice[i]] {

		}*/
		i++
	}

	return c.JSON(http.StatusOK, x[7])
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
func (tc *selectController) FetchCountry(c net.IP) (*mr.Country2Info, error) {
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

func (tc *selectController) slotSize(params map[string][]string) ([]string, []int) {
	var size = make(map[string]string)
	var sizeNumSlice []int
	var slotPublic []string

	for key := range params {
		slice := slotReg.FindStringSubmatch(key)
		//fmt.Println(slice,len(slice))
		if len(slice) == 2 {
			slotPublic = append(slotPublic, slice[1])
			size[slice[1]] = params[key][0]
			//check for size
			SizeNum, _ := config.GetSize(size[slice[1]])
			sizeNumSlice = append(sizeNumSlice, SizeNum)

		}

	}
	return slotPublic, sizeNumSlice
}

func findMinCap(userKey string) int {
	//get user capping data
	result, err := aredis.GetAll(userKey, true, time.Hour)
	if err != nil {
		logrus.Warn("Hgetall error")
	}

	//find min capping for the user
	invMap := make(map[int]string, len(result))
	for k, v := range result {
		invMap[v] = k
	}

	//Sorting
	sortedKeys := make([]int, len(invMap))
	var i int = 0
	for k := range invMap {
		sortedKeys[i] = k
		i++
	}
	fmt.Println("Sorted keys")
	bubbleSorted := utils.BubbleSort(sortedKeys)
	var minCap int
	if len(bubbleSorted) == 0 {
		minCap = 0
	} else {
		minCap = bubbleSorted[0]
	}
	return minCap
}

func init() {

	modules.Register(&selectController{})
}
