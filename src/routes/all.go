package routes

// AllData return all data required to render the all routes
// TODO : Rename this
//type AllData struct {
//	Website  []*mr.Website
//	Province []*mr.Province
//	//Campaign *[]mr.Campaign
//	Size map[string]int
//	Vast bool
//	Data []*mr.AdData
//	Len  int
//}
//
//var allFiter = map[string]selector.FilterFunc{
//	"isWebNetwork":  filter.IsWebNetwork,
//	"webSize":       filter.CheckWebSize,
//	"appSize":       filter.CheckAppSize,
//	"vastSize":      filter.CheckVastSize,
//	"os":            filter.CheckOS,
//	"whiteList":     filter.CheckWhiteList,
//	"blackList":     filter.CheckWebBlackList,
//	"webCategory":   filter.CheckWebCategory,
//	"checkProvince": filter.CheckProvince,
//	"isWebMobile":   filter.IsWebMobile,
//	"notWebMobile":  filter.IsNotWebMobile,
//	"checkCampaign": filter.CheckCampaign,
//	"webMobileSize": filter.CheckWebMobileSize,
//	"appBlackList":  filter.CheckAppBlackList,
//	"appWhiteList":  filter.CheckAppWhiteList,
//	"appCategory":   filter.CheckAppCategory,
//	"appBrand":      filter.CheckAppBrand,
//	"appHood":       filter.CheckAppHood,
//	"appProvider":   filter.CheckProvder,
//	"appAreaInGlob": filter.CheckAppAreaInGlob,
//}
//
//// Ints returns a unique subset of the int slice provided.
//func UniqueStr(input []string) []string {
//	u := make([]string, 0, len(input))
//	m := make(map[string]bool)
//
//	for _, val := range input {
//		if _, ok := m[val]; !ok {
//			m[val] = true
//			u = append(u, val)
//		}
//	}
//
//	return u
//}
//
//func (tc *selectController) allAds(c echo.Context) error {
//	w := c.QueryParam("w")
//	p := c.QueryParam("p")
//	v := c.QueryParam("v")
//	cam := c.QueryParam("cam")
//	s := c.QueryParam("s")
//	ab := c.QueryParam("ab")
//	wa := c.QueryParam("wa")
//	ca := c.QueryParam("ca")
//	b := c.QueryParam("b")
//	pr := c.QueryParam("pr")
//	h := c.QueryParam("h")
//	ar := c.QueryParam("ar")
//	webMobile := c.QueryParam("webmobile")
//
//	//cat := c.QueryParam("cat")
//	var campaign int64
//	var province int64
//	var resFilter = []selector.FilterFunc{}
//	var fltrString []string
//
//	var sizeNumSlice = make(map[string]int)
//	var website *mr.Website
//	var app *mr.App
//	var err error
//	var vv bool
//	rd := middlewares.MustGetRequestData(c)
//	m := selector.Context{}
//
//	if v != "" || v == "on" {
//		fltrString = append(fltrString,
//			"isWebNetwork",
//		)
//		vv = true
//	} else {
//		if s != "" {
//			ss := strings.Split(s, ",")
//			var strin string
//			for _, sss := range ss {
//				size, err := strconv.Atoi(sss)
//				if err == nil {
//					strin = fmt.Sprintf("1jhgy%d", rand.Intn(200))
//					sizeNumSlice[strin] = size
//				}
//			}
//			if len(sizeNumSlice) > 0 {
//				fltrString = append(fltrString,
//					"webSize",
//				)
//			}
//		}
//	}
//	if w != "" {
//		ww, err := strconv.ParseInt(w, 10, 0)
//		if err != nil {
//			website, err = mr.NewManager().FetchWebsite(ww)
//			if err == nil {
//				fltrString = append(fltrString, "whiteList", "blackList")
//			}
//		}
//	}
//	if ab != "" {
//		abab, err := strconv.ParseInt(w, 10, 0)
//		if err == nil {
//			app, err = mr.NewManager().GetAppByID(abab)
//			if err == nil {
//				fltrString = append(fltrString, "appBlackList")
//			}
//		}
//
//	}
//	if wa != "" {
//		wawa, err := strconv.ParseInt(wa, 10, 0)
//		if err == nil {
//			app, err = mr.NewManager().GetAppByID(wawa)
//			if err == nil {
//				fltrString = append(fltrString, "appWhiteList")
//			}
//		}
//
//	}
//	if ca != "" {
//		caca, err := strconv.ParseInt(ca, 10, 0)
//		if err == nil {
//			//todo
//			if true {
//				m.Website.WCategories = mr.SharpArray(fmt.Sprintf("#%d#", caca))
//				fltrString = append(fltrString, "webCategory")
//			}
//			//todo
//			if false {
//				m.App.Appcat = mr.SharpArray(fmt.Sprintf("#%d#", caca))
//				fltrString = append(fltrString, "appCategory")
//			}
//		}
//	}
//	if b != "" {
//		bb, err := strconv.ParseInt(b, 10, 0)
//		if err == nil {
//			m.PhoneData.BrandID = bb
//			fltrString = append(fltrString, "appBrand")
//		}
//	}
//	if pr != "" {
//		prpr, err := strconv.ParseInt(pr, 10, 0)
//		if err == nil {
//			m.PhoneData.NetworkID = prpr
//			fltrString = append(fltrString, "appProvider")
//		}
//	}
//	if h != "" {
//		hh, err := strconv.ParseInt(h, 10, 0)
//		if err == nil {
//			m.CellLocation.NeighborhoodsID = hh
//			fltrString = append(fltrString, "appHood")
//		}
//	}
//	if ar != "" {
//		m.CellLocation.Location = "asd"
//		fltrString = append(fltrString, "appAreaInGlob")
//	}
//	if cam != "" {
//
//		campaign, err = strconv.ParseInt(s, 10, 0)
//		if err == nil {
//			fltrString = append(fltrString, "checkCampaign")
//		}
//	} else {
//		campaign = 0
//	}
//	if p != "" {
//		i64, err := strconv.ParseInt(s, 10, 0)
//		if err == nil {
//			province, err = mr.NewManager().ConvertProvinceID2Info(i64)
//			if err == nil {
//				fltrString = append(fltrString,
//					"checkCampaign",
//				)
//			}
//		}
//
//	}
//
//	//check webmobile filter
//	if webMobile == "" || webMobile == "all" {
//		fltrString = append(fltrString,
//			"isWebNetwork",
//		)
//	} else if webMobile == "on" {
//		fltrString = append(fltrString,
//			"isWebNetwork",
//			"isWebMobile",
//			"webMobileSize",
//		)
//	} else if webMobile == "off" {
//		fltrString = append(fltrString,
//			"isWebNetwork",
//			"notWebMobile",
//		)
//	}
//	uniquefltrStr := UniqueStr(fltrString)
//	for u := range uniquefltrStr {
//		logrus.Info(uniquefltrStr[u])
//		resFilter = append(resFilter, allFiter[uniquefltrStr[u]])
//	}
//	//resFilter1 := []selector.FilterFunc{filter.IsWebNetwork, filter.CheckOS, filter.CheckWhiteList}
//
//	m = selector.Context{
//		RequestData: *rd,
//		Website:     website,
//		Size:        sizeNumSlice,
//		Province:    &province,
//		Campaign:    campaign,
//		App:         app,
//	}
//	filteredAds := selector.Apply(&m, selector.GetAdData(), selector.Mix(resFilter...))
//	all := make([]*mr.AdData, 0)
//	for i := range filteredAds {
//		all = append(all, filteredAds[i]...)
//	}
//	al := allDate()
//	al.Vast = vv
//	al.Data = all
//	al.Len = len(all)
//
//	buf := &bytes.Buffer{}
//	err = allAdTemplate.Execute(buf, al)
//	logrus.Info(err)
//	return c.HTML(http.StatusOK, buf.String())
//
//	return c.JSON(200, struct {
//		Count int
//		All   []*mr.AdData
//	}{
//		Count: len(all),
//		All:   all,
//	})
//}
//func allDate() AllData {
//	/*c, err := mr.NewManager().FetchCampaignAll()
//	if err != nil {
//		c = nil
//	}*/
//	p, err := mr.NewManager().FetchProvinceAll()
//	if err != nil {
//		p = nil
//	}
//	w, err := mr.NewManager().FetchWebsiteAll()
//	if err != nil {
//		w = nil
//	}
//	s := config.GetAllSize()
//	al := AllData{
//		//Campaign: c,
//		Province: p,
//		Website:  w,
//		Size:     s,
//	}
//	return al
//}
