package routes

import (
	"assert"
	"bytes"
	"config"
	"fmt"
	"html/template"
	"middlewares"
	"mr"
	"net/http"
	"redis"
	"selector"
	"strconv"
	"transport"
	"utils"

	"errors"

	"github.com/Sirupsen/logrus"
	"gopkg.in/labstack/echo.v3"
)

type vastAdTemplate struct {
	Link     template.HTML
	Repeat   string
	Offset   string
	Type     string
	PublicID string
	Len      string
}

// Select function is the route that the real biding happen
func (tc *selectController) selectVastAd(c echo.Context) error {

	rd, website, province, lenType, length, err := tc.getVastDataFromCtx(c)
	if err != nil {
		return c.HTML(http.StatusBadRequest, err.Error())
	}
	webPublicID := website.WPubID
	slotSize, sizeNumSlice, vastSlotData := tc.slotSizeVast(webPublicID, length, *website)
	// TODO : Move this to slotSizeVast func
	for i := range slotSize {
		slotSize[i].ExtraParam = map[string]string{
			"pos":  vastSlotData[i].Offset,
			"type": vastSlotData[i].Type,
			"l":    lenType,
		}
	}
	//call context
	m := selector.Context{
		RequestData: *rd,
		Website:     website,
		Size:        sizeNumSlice,
		Province:    province,
	}
	filteredAds := selector.Apply(&m, selector.GetAdData(), vastSelector)
	show := tc.makeShow(c, "vast", rd, filteredAds, sizeNumSlice, slotSize, website, true)

	var v = make([]vastAdTemplate, 0)
	for i := range sizeNumSlice {
		v = append(v, vastAdTemplate{
			Link:   template.HTML(fmt.Sprintf("<![CDATA[\n%s\n]]>", show[i])),
			Offset: vastSlotData[i].Offset,
			Type:   vastSlotData[i].Type,
			Repeat: vastSlotData[i].Repeat,
		})
	}
	result := &bytes.Buffer{}

	assert.Nil(vastIndex.Execute(result, v))
	return c.XMLBlob(http.StatusOK, result.Bytes())
}

func (tc *selectController) slotSizeVast(websitePublicID int64, length map[string][]string, website mr.Website) (map[string]*slotData, map[string]int, map[string]vastSlotData) {
	var sizeNumSlice = make(map[string]int)
	var slotPublic []string
	var vastSlot = make(map[string]vastSlotData)
	var i int
	i = 0
	for m := range length {
		i++
		lenType := length[m][0]
		pub := fmt.Sprintf("%d%s", websitePublicID, length[m][1])
		sizeNumSlice[pub] = config.VastNonLinearSize
		if lenType == "linear" {
			sizeNumSlice[pub] = config.VastLinearSize
		}
		slotPublic = append(slotPublic, pub)
		vastSlot[pub] = vastSlotData{
			Offset: m,
			Repeat: length[m][2],
			Type:   lenType,
		}
	}
	finalSlotData, finalSizeNumSlice := tc.slotSizeNormal(slotPublic, website.WID, sizeNumSlice)
	return finalSlotData, finalSizeNumSlice, vastSlot

}

func (tc *selectController) getVastDataFromCtx(c echo.Context) (*middlewares.RequestData, *mr.Website, *mr.Province, string, map[string][]string, error) {
	rd := middlewares.MustGetRequestData(c)

	publicID, err := strconv.ParseInt(c.QueryParam("a"), 10, 0)
	if err != nil {
		return nil, nil, nil, "", nil, errors.New("invalid request")
	}
	//fetch website and set in Context
	website, err := tc.fetchWebsite(publicID)
	if err != nil {
		return nil, nil, nil, "", nil, errors.New("invalid request")
	}

	if !mr.NewManager().IsUserActive(website.UserID) {
		return nil, nil, nil, "", nil, errors.New("user is banned")
	}

	province, err := tc.fetchProvince(rd.IP, c.Request().Header.Get("Cf-Ipcountry"))
	if err != nil {
		logrus.Debug(err)
	}
	lenVast, vastCon := config.MakeVastLen(c.QueryParam("l"))
	return rd, website, province, lenVast, vastCon, nil
}

func (tc *selectController) slotSizeNormal(slotPublic []string, webID int64, sizeNumSlice map[string]int) (map[string]*slotData, map[string]int) {
	slotPublicString := mr.Build(slotPublic)
	res, err := mr.NewManager().FetchSlots(slotPublicString, webID)
	assert.Nil(err)

	answer := make(map[string]*slotData)
	var newSlots []int64
	for i := range slotPublic {
		if _, ok := answer[slotPublic[i]]; ok {
			continue
		}
		for j := range res {
			if fmt.Sprintf("%d", res[j].PublicID) == slotPublic[i] {
				answer[slotPublic[i]] = &slotData{
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

	insertedSlots := tc.insertNewSlots(webID, newSlots...)
	for i := range insertedSlots {
		answer[i] = &slotData{
			ID:       insertedSlots[i],
			PublicID: i,
			SlotSize: sizeNumSlice[i],
		}
	}

	for i := range answer {
		result, err := aredis.SumHMGetField(transport.KeyGenDaily(transport.SLOT, strconv.FormatInt(answer[i].ID, 10)), config.Config.Redis.Days, "i", "c")
		if err != nil || result["c"] == 0 || result["i"] < config.Config.Clickyab.MinImp {
			answer[i].Ctr = config.Config.Clickyab.DefaultCTR
		} else {
			answer[i].Ctr = utils.Ctr(result["i"], result["c"])
		}
	}

	return answer, sizeNumSlice
}
