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

	"transport"
	"utils"

	"sort"

	"assert"

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

	slotSize := tc.slotSize2(params)

	// TODO : err check?
	var userMinView int

	results, _ := aredis.HGetAll(redisUserHashKey, true, 72*time.Hour)
	for i := range sizeNumSlice {
		for ad := range x[sizeNumSlice[i]] {
			view := results[fmt.Sprintf("%s%s%d", transport.CAMPAIGN, transport.DELIMITER, x[sizeNumSlice[i]][ad].CpID)]
			if x[sizeNumSlice[i]][ad].CpFrequency <= 0 {
				// TODO : use default freq from config
				x[sizeNumSlice[i]][ad].CpFrequency = 2
			}
			x[sizeNumSlice[i]][ad].Capping = mr.NewCapping(c, x[sizeNumSlice[i]][ad].CpID, view, x[sizeNumSlice[i]][ad].CpFrequency)
			if userMinView == 0 {
				userMinView = view
			} else if view > 0 && userMinView > view {
				userMinView = view
			}
		}
		sortCap := mr.ByCapping(x[sizeNumSlice[i]])
		sort.Sort(sortCap)
		x[sizeNumSlice[i]] = []mr.MinAdData(sortCap)
	}

	fmt.Println("User min view ", userMinView)
	// TODO {@mahm0ud22} check minimum CpmFloor for the entire site
	var winnerAd = make(map[string]*mr.MinAdData)
	var minCapFloor int
	for slotID := range slotSize {
		var exceedFloor []*mr.MinAdData
		minCapFloor = 0
		for ad := range x[slotSize[slotID]] {

			x[slotSize[slotID]][ad].CTR, _ = CalculateCtr(x[slotSize[slotID]][ad].CpID, x[slotSize[slotID]][ad].AdID, website.WID, slotID)
			x[slotSize[slotID]][ad].CPM = utils.Cpm(x[slotSize[slotID]][ad].CpMaxbid, x[slotSize[slotID]][ad].CTR)
			//exceed cpm floor
			if x[slotSize[slotID]][ad].CPM >= website.WFloorCpm.Int64 {
				if len(exceedFloor) == 0 {
					minCapFloor = x[slotSize[slotID]][ad].Capping.GetCapping()
				}

				//minimum capping
				if x[slotSize[slotID]][ad].Capping.GetCapping() <= minCapFloor && x[slotSize[slotID]][ad].WinnerBid == 0 {
					exceedFloor = append(exceedFloor, &x[slotSize[slotID]][ad])
				}
			}
		}
		if len(exceedFloor) == 0 {
			// TODO : send a warning, log it or anything else:)
			continue
		}
		ef := mr.ByCPM(exceedFloor)
		sort.Sort(ef)
		exceedFloor = []*mr.MinAdData(ef)

		var secondCPM = website.WFloorCpm.Int64
		if len(exceedFloor) > 1 && exceedFloor[0].Capping.GetSelected() == exceedFloor[1].Capping.GetSelected() {
			secondCPM = exceedFloor[1].CPM
		}

		exceedFloor[0].WinnerBid = utils.WinnerBid(secondCPM, exceedFloor[0].CTR)
		exceedFloor[0].Capping.IncView(1)
		winnerAd[slotID] = exceedFloor[0]

		_, err := aredis.IncHash(
			redisUserHashKey,
			fmt.Sprintf("%s%s%d", transport.CAMPAIGN, transport.DELIMITER, exceedFloor[0].CpID),
			1,
			true,
			config.Config.Redis.DailyCapExpireTime,
		)

		assert.Nil(err)
		// TODO {fzerorubigd} : Can we check for inner capping increase?

	}

	// add mega imp
	ip, err := utils.IP2long(rd.IP)
	assert.Nil(err)
	tmp := []interface{}{
		"IP",
		ip,
		"UA",
		rd.UserAgent,
		"WS",
		website.WID,
		"T",
		time.Now().Unix(),
	}

	for i := range winnerAd {
		tmp = append(tmp, fmt.Sprintf("ad_%d", winnerAd[i].AdID), winnerAd[i].WinnerBid)
	}

	assert.Nil(aredis.HMSet(
		"mega_"+rd.MegaImp, true, time.Hour,
		tmp...,
	))

	return c.JSON(http.StatusOK, winnerAd)
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

//must be checked after connect database
func (tc *selectController) slotSize2(params map[string][]string) map[string]int {

	var size = make(map[string]int)
	//var realSize int
	for key := range params {
		slice := slotReg.FindStringSubmatch(key)

		if len(slice) == 2 {
			//check for size
			SizeNum, _ := config.GetSize(params[key][0])
			size[string(slice[1])] = SizeNum
		}

	}

	return size
}

func CalculateCtr(cpID int64, adID int64, wID int64, slotPublicID string) (float64, string) {
	day := 2
	final := make(map[string]int)
	for c := range config.Config.CtrConst {
		var key string
		switch config.Config.CtrConst[c] {
		case transport.AD_SLOT:

			key = fmt.Sprintf("%s%s%d%s%s%s",
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

			fmt.Sprintf("%s%s%s",
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

			key = fmt.Sprintf("%s%s%d%s%s%s",
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
