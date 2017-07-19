package routes

import (
	"assert"
	"bytes"
	"text/template"
)

type protocol string

const (
	httpScheme  protocol = "http"
	httpsScheme          = "https"
)

type nativeContainer struct {
	Ads      []nativeAd
	Title    string
	Style    string
	FontSize string
	Position string
}

type nativeAd struct {
	Protocol protocol
	Corners  string
	Image    string
	Title    string
	More     string
	Lead     string
	URL      string
	Site     string
}

const nativeTmpl = `{{define "ads"}}<div class="cyb-holder cyb-custom-holder" style="font-size: {{.FontSize}}">
	<style>
	{{.Style}}
	</style>
    <div class="cyb-title-holder cyb-custom-title-holder">
        <div class="cyb-title-before cyb-custom-title-before"></div>
        <div class="cyb-title cyb-custom-title">{{.Title}}</div>
         <div class="cyb-title-after cyb-custom-title-after"></div>

            <div class="cyb-logo">

                <a target="_blank" href="https://www.clickyab.com/?ref=icon" class="cyb-logo-container">
                    <img class="cyb-logo-color" src="data:image/png;base64,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">
                    <img class="cyb-logo-gray" src="data:image/png;base64,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">

                </a>

            </div>

    </div>
    <div class="cyb-suggests cyb-{{.Position}} cyb-custom-suggests">
    	{{renderAds .Ads}}
    </div>
</div>
{{end}}
`

const adTmpl = `{{define "ad"}}
       <div class="cyb-suggest cyb-custom-suggest ">
                <div class="cyb-img-holder cyb-custom-img-holder">
                    <a target="_blank" href="{{.URL}}" onclick="cybOpen(event)" oncontextmenu="cybOpen(event)"
                       ondblclick="cybOpen(event)" data-href="{{.URL}}">
                        <img src="{{.Image}}" alt="{{.Title}}"
                             class="cyb-img {{isRound .Corners}} cyb-custom-img">
                    </a>
                </div>
                <div class="cyb-desc-holder cyb-custom-desc-holder">
                    <div class="cyb-desc cyb-custom-desc">
                        <a target="_blank" href="{{.URL}}" onclick="cybOpen(event)" oncontextmenu="cybOpen(event)"
                           ondblclick="cybOpen(event)" data-href="{{.URL}}">
                            {{.Title}}
                        </a>
                    </div>
                </div>
            </div>
            {{end}}
`

var addRenderer = func(ads []nativeAd) string {
	t, e := template.New("ad").Funcs(template.FuncMap{"isRound": func(s string) string {
		return "cyb-" + s
	}}).Parse(adTmpl)
	assert.Nil(e)

	b := &bytes.Buffer{}

	// remember to pack each two ad into one div like following example
	//         <div class="cyb-pack cyb-custom-pack">
	// 				<AD>
	//				<AD>
	// 			</div>
	// it's a simple hack to keep all ads (relatively) in same ratio
	p := 0
	for i, ad := range ads {
		if i != 0 && i == p {
			b.WriteString("</div>")
		}
		if i%2 == 0 {
			p += 2
			b.WriteString(`<div class="cyb-pack cyb-custom-pack">`)
		}
		e := t.Lookup("ad").Execute(b, ad)
		assert.Nil(e)

		if len(ads)-1 == i {
			b.WriteString("</div>")
		}

	}

	return b.String()
}

var native = template.New("native").Funcs(template.FuncMap{"renderAds": addRenderer})

func renderNative(imp nativeContainer) string {
	buf := &bytes.Buffer{}
	imp.Style = style
	e := native.Lookup("ads").Execute(buf, imp)
	assert.Nil(e)
	return string(buf.Bytes())
}
func init() {
	native.Parse(nativeTmpl)
	native.Parse(adTmpl)
}
