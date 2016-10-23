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

	"sort"
	"utils"

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
	adIDBanner := GetAdID(x)
	adBanner, _ := mr.NewManager().FetchSlotAd(mr.Build(slotPublic), mr.Build(adIDBanner))
	tc.AddCTR(adBanner, x)
	filteredAdd := tc.CpmFloor(x, *website)
	return c.JSON(http.StatusOK, filteredAdd[7])
}

//AddCTR add ctr from slot_ad
func (tc *selectController) AddCTR(adBanner []mr.SlotData, x map[int][]mr.AdData) {
	for i := range adBanner {
		for r := range x[adBanner[i].SlotSize] {
			if x[adBanner[i].SlotSize][r].AdID == adBanner[i].AdID {
				x[adBanner[i].SlotSize][r].CTR = utils.Ctr(adBanner[i].SLAImps, adBanner[i].SLAClicks)
			}

		}

	}
}

//CpmFloor create slice ads filter cpm > cpm floor and sort by cpm
func (tc *selectController) CpmFloor(x map[int][]mr.AdData, website mr.WebsiteData) map[int][]mr.AdData {
	filteredAdd := make(map[int][]mr.AdData)
	for size := range x {
		for ad := range x[size] {
			//if ads not have slot_ad ctr get ctr of ad_ctr
			if x[size][ad].AdCtr != 0 && x[size][ad].CTR == config.Config.DefaultCTR {
				x[size][ad].CTR = x[size][ad].AdCtr
			}
			//callculate cpm
			x[size][ad].CPM = utils.Cpm(x[size][ad].CpMaxbid, x[size][ad].CTR)

			if x[size][ad].CPM >= website.WFloorCpm.Int64 {
				filteredAdd[size] = append(filteredAdd[size], x[size][ad])
			}

		}
		//sort by cpm
		sort.Sort(mr.ByCPM(filteredAdd[size]))
	}
	return filteredAdd
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

func (tc *selectController) slotSize(params map[string][]string) ([]string, []int) {
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
	return slotPublic, sizeNumSlice
}

func init() {

	modules.Register(&selectController{})
}
