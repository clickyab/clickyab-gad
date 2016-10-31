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

	"net"

	"redis"
	"time"

	"sort"
	"transport"
	"utils"

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
	redisUserHashKey := fmt.Sprintf("%s%s%s%s%s", transport.USER_CAPPING, transport.DELIMITER, m.CopID, transport.DELIMITER, time.Now().Format("060102"))
	//let the game begin :)
	/*for s := range slotPublic {
	fmt.Println(s)
	i := 0
	*/ /*for size := range x[sizeNumSlice[i]] {

	}*/ /*
		i++
	}*/
	return c.JSON(http.StatusOK, x[7])
	results, _ := aredis.HGetAll(redisUserHashKey, true, 72*time.Hour)
	for i := range sizeNumSlice {
		for ad := range x[sizeNumSlice[i]] {
			if view, ok := results[fmt.Sprintf("%s%s%d", transport.CAMPAIGN, transport.DELIMITER, x[sizeNumSlice[i]][ad].CpID)]; ok {
				if x[sizeNumSlice[i]][ad].CpFrequency <= 0 {
					x[sizeNumSlice[i]][ad].CpFrequency = 2
				}
				x[sizeNumSlice[i]][ad].Capping = view / x[sizeNumSlice[i]][ad].CpFrequency
			} else {
				x[sizeNumSlice[i]][ad].Capping = 0
			}
		}
		sortCap := mr.ByCapping(x[sizeNumSlice[i]])
		sort.Sort(sortCap)
		x[sizeNumSlice[i]] = []mr.AdData(sortCap)
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

/*func findMinCap(userKey string) (map[string]string, error) {
	//get user capping data
	result, err := aredis.HGetAll(userKey, true, 72*time.Hour)
	if err != nil {
		return nil, errors.New("error")
	}
	for i:= range result{
		fmt.Println(result[i])
	}

	return result, nil
}*/

func CalculateCtr(cpID int64, adID int64, wID int64, slotPublicID int64) (float64, string) {
	day := 2
	final := make(map[string]int)
	for c := range config.Config.CtrConst {
		var key string
		switch config.Config.CtrConst[c] {
		case transport.AD_SLOT:

			key = fmt.Sprintf("%s%s%d%s%d%s",
				transport.AD_SLOT,
				transport.DELIMITER,
				adID, transport.DELIMITER,
				slotPublicID, transport.DELIMITER)

		case transport.CAMPAIGN:

			key = fmt.Sprintf("%s%s%d%s",
				transport.CAMPAIGN,
				transport.DELIMITER,
				cpID, transport.DELIMITER)

		case transport.ADVERTISE:

			key = fmt.Sprintf("%s%s%d%s",
				transport.ADVERTISE,
				transport.DELIMITER,
				adID, transport.DELIMITER)

		case transport.SLOT:

			fmt.Sprintf("%s%s%d",
				transport.SLOT,
				transport.DELIMITER,
				slotPublicID,
			)

		case transport.WEBSITE:

			key = fmt.Sprintf("%s%s%d",
				transport.WEBSITE,
				transport.DELIMITER,
				wID,
			)

		case transport.AD_WEBSITE:

			key = fmt.Sprintf("%s%s%d%s%d%s",
				transport.AD_WEBSITE,
				transport.DELIMITER,
				adID,
				transport.DELIMITER,
				wID,
				transport.DELIMITER,
			)

		case transport.CAMPAIGN_SLOT:

			key = fmt.Sprintf("%s%s%d%s%d%s",
				transport.CAMPAIGN_SLOT,
				transport.DELIMITER,
				cpID,
				transport.DELIMITER,
				slotPublicID,
				transport.DELIMITER,
			)

		}
		result, err := aredis.SumHMGetField(key, day, "i", "c")
		if err != nil || result["c"] == 0 || result["i"] < config.Config.MinImp {
			final[config.Config.CtrConst[c]] = 0
		} else {
			return utils.Ctr(result["i"], result["c"]), config.Config.CtrConst[c]
		}
	}
	return config.Config.DefaultCTR, "default"
}

func init() {

	modules.Register(&selectController{})
}
