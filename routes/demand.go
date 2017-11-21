package routes

import (
	"errors"
	"fmt"
	"net/http"

	"clickyab.com/gad/ip2location"
	"clickyab.com/gad/middlewares"
	"clickyab.com/gad/models"
	"clickyab.com/gad/models/selector"
	"clickyab.com/gad/utils"
	"github.com/clickyab/simple-rtb"
	"gopkg.in/labstack/echo.v3"
)

func (tc *selectController) selectDemandWebAd(c echo.Context, rd *middlewares.RequestData, e *srtb.BidRequest) error {
	rd, e, website, province, err := tc.getWebDataExchangeFromCtx(c, rd, e)
	if err != nil {
		return c.HTML(http.StatusBadRequest, err.Error())
	}
	slotSize, sizeNumSlice, trackIDs, bidfloors, err := tc.slotSizeWebExchange(e.Imp, *website)
	if err != nil {
		return c.HTML(http.StatusBadRequest, "slot size was wrong, reason : "+err.Error())
	}

	//call context
	m := selector.Context{
		RequestData: *rd,
		Website:     website,
		Size:        sizeNumSlice,
		Province:    province,
		BidFloors:   bidfloors,
	}
	var sel selector.FilterFunc
	if e.Site != nil {
		sel = webSelector
	} else {
		return c.HTML(http.StatusBadRequest, "not supported platform")
	}
	filteredAds := selector.Apply(&m, selector.GetAdData(), sel)
	show, ads := tc.makeShow(c, "sync", rd, filteredAds, nil, sizeNumSlice, slotSize, nil, website, false, bidfloors, false, true, floorDivDemand.Int64(), true)

	//substitute the webMobile slot if exists
	var bids []srtb.Bid
	for i := range ads {
		if ads[i] == nil {
			continue
		}

		h := utils.GetSizeByNumStringHeight(ads[i].AdSize)
		w := utils.GetSizeByNumStringWith(ads[i].AdSize)
		bids = append(bids, srtb.Bid{
			ID:       <-utils.ID,
			Height:   h,
			Width:    w,
			AdID:     fmt.Sprintf("%d", ads[i].AdID),
			ImpID:    trackIDs[ads[i].SlotPublicID],
			AdMarkup: fmt.Sprintf(`<iframe name="clickyab_frame" src="%s&px=${PIXEL_URL_IMAGE:B64}" marginwidth="0" marginheight="0" vspace="0" hspace="0" allowtransparency="true" scrolling="no" width="%d" height="%d" frameborder="0"></iframe>`, show[i], w, h),
			Price:    int64(float64(ads[i].WinnerBid) * ads[i].CTR * 10),
			WinURL:   "",
			Cat:      []string{},
		})
	}
	dm := srtb.BidResponse{
		ID:   e.ID,
		Bids: bids,
	}
	if len(dm.Bids) < 1 {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, dm)
}

// selectDemandWebAd function is the route that the real biding happens
func (tc *selectController) selectDemandAd(c echo.Context) error {
	rd := middlewares.MustGetRequestData(c)
	e := middlewares.MustExchangeGetRequestData(c)
	if e.App != nil && e.Site != nil {
		return c.HTML(http.StatusBadRequest, "wrong platform")
	}
	if e.Site != nil {
		return tc.selectDemandWebAd(c, rd, e)
	} // app platform selected
	//TODO implement later
	//return tc.selectDemandAppAd(c, rd, e)
	return c.HTML(http.StatusBadRequest, "wrong platform")

}

func (tc *selectController) getWebDataExchangeFromCtx(c echo.Context, rd *middlewares.RequestData, e *srtb.BidRequest) (*middlewares.RequestData, *srtb.BidRequest, *models.Website, int64, error) {
	name, userID, err := utils.GetSupplier(rd.SupplierKey)
	if err != nil {
		return nil, nil, nil, 0, fmt.Errorf("can not accept from supplier with key = %s", rd.SupplierKey)
	}
	website, err := tc.fetchWebsiteDomain(e.Site.Domain, name, userID)
	if err != nil {
		return nil, nil, nil, 0, err //errors.New("invalid request")
	}
	if !website.GetActive() {
		return nil, nil, nil, 0, errors.New("website is not active")
	}

	if !models.NewManager().IsUserActive(website.UserID) {
		return nil, nil, nil, 0, errors.New("user is banned")
	}

	var province int64
	if e.Device.Geo.Region.Valid {
		province = ip2location.GetProvinceIDByName(e.Device.Geo.Region.Name)
	}
	return rd, e, website, province, nil
}

//fetchWebsiteDomain website and check if the minimum floor is applied
func (tc *selectController) fetchWebsiteDomain(domain, supplier string, user int64) (*models.Website, error) {
	website, err := models.NewManager().FetchWebsiteByDomain(domain, supplier)
	if err != nil {
		website, err = models.NewManager().InsertWebsite(domain, supplier, user)
		if err != nil {
			return nil, err
		}
	}
	if website.WFloorCpm.Int64 < minCPMFloorWeb.Int64() {
		website.WFloorCpm.Int64 = minCPMFloorWeb.Int64()
	}
	return website, err
}

func (tc selectController) slotSizeWebExchange(imps []srtb.Impression, website models.Website) (map[string]*slotData, map[string]int, map[string]string, map[string]int64, error) {
	var sizeNumSlice = make(map[string]int)
	var slotPublics []string
	var trackIDs = make(map[string]string)
	var bidFloors = make(map[string]int64)
	var secureSlots = make(map[string]bool)
	//var attr = make(map[string]map[string]string)
	for i := range imps {
		size, err := utils.GetSize(fmt.Sprintf("%dx%d", imps[i].Banner.Width, imps[i].Banner.Height))
		slotPublic := fmt.Sprintf("%d%d%d", website.WPubID, size, i)
		//attr[slotPublic] = make(map[string]string)
		//for k, v := range imps[i].Attributes {
		//	attr[slotPublic][k] = v
		//}
		if err == nil {
			sizeNumSlice[slotPublic] = size
			secureSlots[slotPublic] = func() bool {
				if imps[i].Secure == 1 {
					return true
				}
				return false
			}()
			slotPublics = append(slotPublics, slotPublic)
			trackIDs[slotPublic] = imps[i].ID
			bidFloors[slotPublic] = int64(imps[i].BidFloor)
		}
	}
	if len(slotPublics) == 0 {
		return nil, nil, nil, bidFloors, errors.New("no supported i size")
	}
	all, size := tc.slotSizeNormal(slotPublics, website.WID, sizeNumSlice)
	for i := range all {
		all[i].Secure = secureSlots[i]
	}
	return all, size, trackIDs, bidFloors, nil
}
