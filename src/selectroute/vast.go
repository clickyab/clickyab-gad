package selectroute

import (
	"assert"
	"mr"
	"selector"

	"fmt"
	"strconv"

	"middlewares"

	"config"

	"bytes"

	"html/template"

	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
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

	rd, website, country, lenType, length, err := tc.getVastDataFromCtx(c)
	if err != nil {
		return err
	}
	webPublicID := website.WPubID
	slotSize, sizeNumSlice, vastSlotData := tc.slotSizeVast(webPublicID, length, *website)
	//call context
	m := selector.Context{
		RequestData: *rd,
		Website:     website,
		Size:        sizeNumSlice,
		Country:     country,
	}
	filteredAds := selector.Apply(&m, selector.GetAdData(), vastSelector)
	show := tc.makeShow(c, "vast", rd, filteredAds, sizeNumSlice, slotSize, website, true)

	var v = make([]vastAdTemplate, 0)
	for i := range sizeNumSlice {
		v = append(v, vastAdTemplate{
			Link:   template.HTML(fmt.Sprintf("<![CDATA[\n%s&pos=%s&type=%s&l=%s\n]]>", show[i], vastSlotData[i].Offset, vastSlotData[i].Type, lenType)),
			Offset: vastSlotData[i].Offset,
			Type:   vastSlotData[i].Type,
			Repeat: vastSlotData[i].Repeat,
		})
	}
	result := &bytes.Buffer{}

	assert.Nil(vastIndex.Execute(result, v))
	return c.XMLBlob(http.StatusOK, result.Bytes())
}

func (tc *selectController) slotSizeVast(websitePublicID int64, length map[string][]string, website mr.WebsiteData) (map[string]slotData, map[string]int, map[string]vastSlotData) {
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

func (tc *selectController) getVastDataFromCtx(c echo.Context) (*middlewares.RequestData, *mr.WebsiteData, *mr.CountryInfo, string, map[string][]string, error) {
	rd := middlewares.MustGetRequestData(c)

	publicParams := c.QueryParam("a")
	publicID, err := strconv.Atoi(publicParams)
	if err != nil {
		return nil, nil, nil, "", nil, c.HTML(400, "invalid request")
	}
	//fetch website and set in Context
	website, err := tc.fetchWebsite(publicID)
	if err != nil {
		return nil, nil, nil, "", nil, c.HTML(400, "invalid request")
	}
	country, err := tc.fetchCountry(rd.IP)
	if err != nil {
		logrus.Debug(err)
	}
	lenVast, vastCon := config.MakeVastLen(c.QueryParam("l"))
	return rd, website, country, lenVast, vastCon, nil
}

func (tc *selectController) slotSizeNormal(slotPublic []string, webID int64, sizeNumSlice map[string]int) (map[string]slotData, map[string]int) {
	slotPublicString := mr.Build(slotPublic)
	res, err := mr.NewManager().FetchSlots(slotPublicString, webID)
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

	insertedSlots := tc.insertNewSlots(webID, newSlots...)
	for i := range insertedSlots {
		answer[i] = slotData{
			ID:       insertedSlots[i],
			PublicID: i,
			SlotSize: sizeNumSlice[i],
		}
	}

	return answer, sizeNumSlice
}
