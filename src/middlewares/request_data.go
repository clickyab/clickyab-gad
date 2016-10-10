package middlewares

import (
	_ "fmt"
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mssola/user_agent"
	"mr"
	"strconv"
	"regexp"
	"banners"
)

type RequestData struct {
	CloudIP        string
	RealIP         string
	IP2Location    *mr.IP2Location
	Browser        string
	Os             string
	OsVersion      string
	BrowserVersion string
	Method         string
	Referrer       string
}

type Size []int

// RequestCollector try to collect data from request
func RequestCollector(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		e := &RequestData{}
		ua := user_agent.New(ctx.Request().UserAgent())
		name, version := ua.Browser()
		e.Browser = name
		e.BrowserVersion = version
		e.Os = ua.OS()
		e.OsVersion = ua.Platform()
		e.RealIP = ctx.Request().RealIP()
		e.Referrer = ctx.Request().Referer()
		e.Method = ctx.Request().Method()

		ctx.Set("RequestData", e)
		params := ctx.QueryParams()
		public_id, _ := strconv.Atoi(params["i"][0])
		domain := params["d"][0]

		var size =make(map[string]string)
		var sizeNumSlice []int
		reg:=regexp.MustCompile(`s\[(\d*)\]`)
		for key := range params{
			slice:=reg.FindStringSubmatch(key)
			//fmt.Println(slice,len(slice))
			if len(slice)==2{
				size[slice[1]]=params[key][0]
				//check for size
				SizeNum,_:=banners.GetSize(size[slice[1]])
				sizeNumSlice=append(sizeNumSlice,SizeNum)

			}

		}
		//set size in context
		ctx.Set("RequestSize", sizeNumSlice)


		////fetch website and set in Context
		wd, err := mr.NewManager().FetchWebsite(public_id, domain)
		if err != nil {
			logrus.Fatal(err)
		}
		ctx.Set("WebsiteData", wd)

		////fetch size and add to context
		rgd, err := mr.NewManager().FetchRegion()
		if err != nil{
			logrus.Fatal(err)
		}
		ctx.Set("RegionData",rgd)

		return next(ctx)
	}
}
