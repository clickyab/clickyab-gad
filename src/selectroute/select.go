package selectroute

import (
	"assert"
	"config"
	"encoding/json"
	"errors"
	"filter"
	"fmt"
	"middlewares"
	"modules"
	"mr"
	"net"
	"redis"
	"regexp"
	"selector"
	"sort"
	"strconv"
	"time"
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

// SlotData is the single slot data in database
type slotData struct {
	SlotSize int
	ID       int64
	PublicID string
}

// Select function is the route that the real biding happen
func (tc *selectController) selectAd(c echo.Context) error {
	params := c.QueryParams()

	rd, website, country, err := tc.getDataFromCtx(c)
	if err != nil {
		return err
	}
	slotSize, sizeNumSlice := tc.slotSize(params, website.WID)
	//call context
	m := selector.Context{
		RequestData: *rd,
		Website:     website,
		Size:        sizeNumSlice,
		Country:     country,
	}
	filteredAds := selector.Apply(&m, selector.GetAdData(), webSelector)
	filteredAds = getCapping(c, rd.CopID, sizeNumSlice, filteredAds)

	var (
		winnerAd = make(map[string]*mr.MinAdData)
		show     = make(map[string]string)
		video    bool // once set, never unset it again
	)
	// TODO : must loop over this values, from lowest data to highest. the size with less ad count must be in higher priority
	for slotID := range slotSize {
		var exceedFloor []*mr.MinAdData
		minCapFloor := 0
		for _, adData := range filteredAds[slotSize[slotID].SlotSize] {
			tc.doBid(adData, website, slotSize[slotID].ID, &exceedFloor, video, &minCapFloor)
		}
		if len(exceedFloor) == 0 {
			// TODO : send a warning, log it or anything else:)
			continue
		}
		ef := mr.ByCPM(exceedFloor)
		sort.Sort(ef)
		exceedFloor = []*mr.MinAdData(ef)

		secondCPM := tc.getSecondCPM(website.WFloorCpm.Int64, exceedFloor)
		exceedFloor[0].WinnerBid = utils.WinnerBid(secondCPM, exceedFloor[0].CTR)
		exceedFloor[0].Capping.IncView(1)
		winnerAd[slotID] = exceedFloor[0]
		video = video || exceedFloor[0].AdType == config.AdTypeVideo
		show[slotID] = fmt.Sprintf("%s://%s/%s/%s/%d?tid=%s&ref=%s&s=%d", rd.Proto, rd.URL, "show", rd.MegaImp, exceedFloor[0].AdID, rd.TID, rd.Parent, slotSize[slotID].ID)

		assert.Nil(storeCapping(m.CopID, exceedFloor[0].CampaignID))
		// TODO {fzerorubigd} : Can we check for inner capping increase?

	}
	err = tc.addMegaKey(rd, website, winnerAd)
	assert.Nil(err)
	b, _ := json.MarshalIndent(show, "\t", "\t")
	result := "renderFarm(" + string(b) + ");"
	return c.HTML(200, result)
}

func (tc *selectController) doBid(adData *mr.MinAdData, website *mr.WebsiteData, slotID int64, exceedFloor *[]*mr.MinAdData, video bool, minCapFloor *int) {
	adData.CTR, _ = tc.calculateCTR(
		adData.CampaignID,
		adData.AdID,
		website.WID,
		slotID,
	)
	adData.CPM = utils.Cpm(adData.CampaignMaxBid, adData.CTR)
	//exceed cpm floor
	if adData.CPM >= website.WFloorCpm.Int64 && (!video || adData.AdType != config.AdTypeVideo) {
		if len(*exceedFloor) == 0 {
			*minCapFloor = adData.Capping.GetCapping()
		}

		//minimum capping
		if adData.Capping.GetCapping() <= *minCapFloor && adData.WinnerBid == 0 {
			*exceedFloor = append(*exceedFloor, adData)

		}
	}
}

func (tc *selectController) getSecondCPM(floorCPM int64, exceedFloor []*mr.MinAdData) int64 {
	var secondCPM = floorCPM
	if len(exceedFloor) > 1 && exceedFloor[0].Capping.GetSelected() == exceedFloor[1].Capping.GetSelected() {
		secondCPM = exceedFloor[1].CPM
	}

	return secondCPM
}

func (tc *selectController) addMegaKey(rd *middlewares.RequestData, website *mr.WebsiteData, winnerAd map[string]*mr.MinAdData) error {
	// add mega imp
	ip, err := utils.IP2long(rd.IP)
	if err != nil {
		return err
	}
	// TODO : get interface from redis?
	tmp := []interface{}{
		"IP",
		fmt.Sprintf("%d", ip),
		"UA",
		rd.UserAgent,
		"WS",
		fmt.Sprintf("%d", website.WID),
		"T",
		fmt.Sprintf("%d", time.Now().Unix()),
	}

	for i := range winnerAd {
		tmp = append(tmp, fmt.Sprintf("ad_%d", winnerAd[i].AdID), fmt.Sprintf("%d", winnerAd[i].WinnerBid))
	}

	return aredis.HMSet(
		fmt.Sprintf("%s%s%s", transport.MEGA, transport.DELIMITER, rd.MegaImp), true, time.Hour,
		tmp...,
	)
}

func (tc *selectController) getDataFromCtx(c echo.Context) (*middlewares.RequestData, *mr.WebsiteData, *mr.CountryInfo, error) {
	rd := middlewares.MustGetRequestData(c)

	params := c.QueryParams()
	publicParams, ok := params["i"]
	if !ok {
		return nil, nil, nil, c.HTML(400, "invalid request")
	}
	publicID, err := strconv.Atoi(publicParams[0])
	if err != nil {
		return nil, nil, nil, c.HTML(400, "invalid request")
	}
	domain, ok := params["d"]
	if !ok {
		return nil, nil, nil, c.HTML(400, "invalid request")
	}
	//fetch website and set in Context
	website, err := tc.fetchWebsite(publicID)
	if err != nil {
		return nil, nil, nil, c.HTML(400, "invalid request")
	}
	country, err := tc.fetchCountry(rd.IP)
	if err != nil {
		logrus.Warn(err)
	}
	//check if the website domain is valid
	if website.WDomain.Valid && website.WDomain.String != domain[0] {
		return nil, nil, nil, errors.New("domain and public id mismatch")
	}

	return rd, website, country, nil
}

//FetchWebsite website and check if the minimum floor is applied
func (selectController) fetchWebsite(publicID int) (*mr.WebsiteData, error) {
	website, err := mr.NewManager().FetchWebsite(publicID)
	if err != nil {
		return nil, err
	}
	if website.WFloorCpm.Int64 < config.Config.Clickyab.MinCPMFloor {
		website.WFloorCpm.Int64 = config.Config.Clickyab.MinCPMFloor
	}
	return website, err
}

//FetchCountry find country and set context
func (selectController) fetchCountry(c net.IP) (*mr.CountryInfo, error) {
	var country mr.CountryInfo
	ip, err := mr.NewManager().GetLocation(c)
	if err != nil || !ip.CountryName.Valid {
		return nil, errors.New("Country not found")
	}
	country, err = mr.NewManager().ConvertCountry2Info(ip.CountryName.String)
	if err != nil {
		return nil, errors.New("Country not found")
	}
	return &country, nil

}

func (selectController) slotSize(params map[string][]string, wID int64) (map[string]slotData, map[string]int) {
	var size = make(map[string]string)
	var sizeNumSlice map[string]int
	var slotPublic []string

	for key := range params {
		slice := slotReg.FindStringSubmatch(key)

		//fmt.Println(slice,len(slice))
		if len(slice) == 2 {

			slotPublic = append(slotPublic, slice[1])
			size[slice[1]] = params[key][0]
			//check for size
			//size[slice[1]] = strings.Trim(size[slice[1]], "a")
			SizeNum, _ := config.GetSize(size[slice[1]])
			sizeNumSlice[slice[1]] = SizeNum

		}

	}

	//query to fetch slot ID
	slotPublicString := mr.Build(slotPublic)
	res, err := mr.NewManager().FetchSlots(slotPublicString, wID)
	assert.Nil(err)

	answer := make(map[string]slotData)
	var newSlots []int64
	for i := range slotPublic {
		if _, ok := answer[slotPublic[i]]; ok {
			continue
		}
		for j := range res {
			if fmt.Sprintf("%d", res[j].PublicID) == slotPublic[i] {
				answer[slotPublic[i]] = slotData{
					ID:       res[j].ID,
					PublicID: slotPublic[i],
					SlotSize: sizeNumSlice[slotPublic[i]],
				}
				break
			}
		}
		if _, ok := answer[slotPublic[i]]; !ok {
			s, err := strconv.ParseInt(slotPublic[i], 10, 0)
			if err == nil {
				newSlots = append(newSlots, s)
			}
		}
	}

	if len(newSlots) > 0 {
		insertedSlots, err := mr.NewManager().InsertSlots(wID, newSlots...)
		if err == nil {
			for i := range insertedSlots {
				p := fmt.Sprintf("%d", insertedSlots[i].PublicID)
				answer[p] = slotData{
					ID:       insertedSlots[i].ID,
					PublicID: p,
					SlotSize: sizeNumSlice[p],
				}
			}
		}
	}

	return answer, sizeNumSlice
}

// CalculateCtr calculate ctr
func (selectController) calculateCTR(cpID int64, adID int64, wID int64, slotID int64) (float64, string) {
	day := 2
	final := make(map[string]int)
	for c := range config.Config.Clickyab.CTRConst {
		key := bestCTRKey(c, adID, slotID, cpID, wID)
		result, err := aredis.SumHMGetField(key, day, "i", "c")
		if err != nil || result["c"] == 0 || result["i"] < config.Config.Clickyab.MinImp {
			final[config.Config.Clickyab.CTRConst[c]] = 0
		} else {
			return utils.Ctr(result["i"], result["c"]), config.Config.Clickyab.CTRConst[c]
		}
	}
	return config.Config.Clickyab.DefaultCTR, "default"
}

func bestCTRKey(c int, adID int64, slotID int64, cpID int64, wID int64) string {
	var key string
	switch config.Config.Clickyab.CTRConst[c] {
	case transport.AD_SLOT:

		key = fmt.Sprintf("%s%s%d%s%d%s",
			transport.AD_SLOT,
			transport.DELIMITER,
			adID, transport.DELIMITER,
			slotID, transport.DELIMITER)

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
			slotID,
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
			slotID,
			transport.DELIMITER,
		)

	}
	return key
}

func init() {
	modules.Register(&selectController{})
}
