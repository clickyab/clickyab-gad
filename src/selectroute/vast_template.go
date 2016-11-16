package selectroute

import "html/template"

const linearTemplate = `<VAST version="3.0">
        <Ad id="{{.RND}}">
            <InLine>
                <AdSystem version="0.1.0">clickyab-vast</AdSystem>
                <AdTitle>
                    {{.Title}}
                </AdTitle>
                <Description>
                    {{.Description}}
                </Description>
                <Creatives>
                    <Creative sequence="1" id="{{.RND2}}">
                        <Linear skipOffset="{{.SkipOffset}}">
                            <Duration>{{.Duration}}</Duration>
                            <VideoClicks>
                                <ClickThrough>
                                   {{.Link}}
                                </ClickThrough>
                            </VideoClicks>
    							{{ if .Video }}
                            <TrackingEvents>
                                <Tracking event="complete">
                                    {{ .Link }}
                                </Tracking>
                            </TrackingEvents>
        						{{end}}
                            <MediaFiles>
                                <MediaFile delivery="progressive" bitrate="24" width="{{.Width}}" height="{{.Height}}" type="{{if .Video}}video/mp4{{else}}image{{end}}">
                                    {{.Src}}
                                </MediaFile>
                            </MediaFiles>
                        </Linear>
                    </Creative>
                </Creatives>
            </InLine>
        </Ad>
    </VAST>`
const nonLinearTemplate = `<VAST version="3">
        <Ad id="{{.RND}}">
            <InLine>
                <AdSystem version="0.1.0">clickyab-vast</AdSystem>
                <Creatives>
                    <Creative AdID="{{.RND}}">
                        <NonLinearAds>
                            <TrackingEvents>
                            </TrackingEvents>
                            <NonLinear height="{{.Height}}" width="{{.Width}}" minSuggestedDuration="{{.Duration}}">
                                <StaticResource creativeType="{{if .Video}}video/mp4{{else}}image{{end}}">
                                    {{.Src}}
                                </StaticResource>
                                <NonLinearClickThrough>
                                    {{.Link}}
                                </NonLinearClickThrough>
                            </NonLinear>
                        </NonLinearAds>
                    </Creative>
                </Creatives>
            </InLine>
        </Ad>
    </VAST>
`

var (
	linear    = template.Must(template.New("vast_linear_ad").Parse(linearTemplate))
	nonlinear = template.Must(template.New("vast_nonlinear_ad").Parse(nonLinearTemplate))
)
