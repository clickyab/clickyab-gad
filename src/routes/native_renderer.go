package routes

import (
	"assert"
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

type protocol string

const (
	httpScheme  protocol = "http"
	httpsScheme          = "https"
)

type nativeContainer struct {
	Ads        []nativeAd
	Title      string
	Style      string
	FontSize   string
	FontFamily string
	Position   string
	MinSize    string
	IsVertical bool
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
	Extra    string
}

const nativeTmpl = `{{define "ads"}}<div class="cyb-holder {{if .IsVertical}}cyb-side {{end}} cyb-custom-holder" style="font-size: {{.FontSize}};font-family:{{.FontFamily}}">
	<style>
	@font-face {
  font-family: 'sahel';
  src: url('https://static.clickyab.com/font/Sahel-FD.eot');
  src: url('https://static.clickyab.com/font/Sahel-FD.eot?#iefix') format('embedded-opentype'),
       url('https://static.clickyab.com/font/Sahel-FD.woff') format('woff'),
       url('https://static.clickyab.com/font/Sahel-FD.ttf')  format('truetype');
}

@font-face {
  font-family: 'samim';
  src: url('https://static.clickyab.com/font/Samim-FD.eot');
  src: url('https://static.clickyab.com/font/Samim-FD.eot?#iefix') format('embedded-opentype'),
       url('https://static.clickyab.com/font/Samim-FD.woff') format('woff'),
       url('https://static.clickyab.com/font/Samim-FD.ttf')  format('truetype');
}

@font-face {
  font-family: 'vazir';
  src: url('https://static.clickyab.com/font/Vazir-FD.eot');
  src: url('https://static.clickyab.com/font/Vazir-FD.eot?#iefix') format('embedded-opentype'),
       url('https://static.clickyab.com/font/Vazir-FD.woff2') format('woff2'),
       url('https://static.clickyab.com/font/Vazir-FD.woff') format('woff'),
       url('https://static.clickyab.com/font/Vazir-FD.ttf')  format('truetype');
}

@font-face {
  font-family: 'behdad';
  src: url('https://static.clickyab.com/font/Behdad-Regular.ttf')  format('truetype'),
       url('https://static.clickyab.com/font/Behdad-Regular.woff2') format('woff2'),
       url('https://static.clickyab.com/font/Behdad-Regular.woff') format('woff'),
       url('https://static.clickyab.com/font/Behdad-Regular.ttf')  format('truetype'),
       url('https://static.clickyab.com/font/Behdad-Regular.otf') format('opentype');
}

@font-face {
	font-family: "nazanin";
	src: url('https://static.clickyab.com/font/nazanin.ttf') format('truetype');
}
	{{.Style}}
	</style>
    <div class="cyb-title-holder  cyb-custom-title-holder">
        <div class="cyb-title-before cyb-custom-title-before"></div>
        <div class="cyb-title cyb-custom-title">{{.Title}}</div>
         <div class="cyb-title-after cyb-custom-title-after"></div>

            <div class="cyb-logo">

                <a rel="nofollow" target="_blank" href="https://www.clickyab.com/?ref=icon" class="cyb-logo-container">
                <img src="data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48c3ZnIHdpZHRoPSIxNXB4IiBoZWlnaHQ9IjE4cHgiIHZpZXdCb3g9IjAgMCAxNSAxOCIgdmVyc2lvbj0iMS4xIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHhtbG5zOnhsaW5rPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hsaW5rIj4gICAgICAgIDx0aXRsZT5Hcm91cDwvdGl0bGU+ICAgIDxkZXNjPkNyZWF0ZWQgd2l0aCBTa2V0Y2guPC9kZXNjPiAgICA8ZGVmcz48L2RlZnM+ICAgIDxnIGlkPSJQYWdlLTIiIHN0cm9rZT0ibm9uZSIgc3Ryb2tlLXdpZHRoPSIxIiBmaWxsPSJub25lIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiPiAgICAgICAgPGcgaWQ9IkRlc2t0b3AtSEQiIHRyYW5zZm9ybT0idHJhbnNsYXRlKC0xMTcuMDAwMDAwLCAtMjA2LjAwMDAwMCkiIGZpbGwtcnVsZT0ibm9uemVybyI+ICAgICAgICAgICAgPGcgaWQ9Ikdyb3VwLTIiIHRyYW5zZm9ybT0idHJhbnNsYXRlKDExNy4wMDAwMDAsIDIwMi4wMDAwMDApIj4gICAgICAgICAgICAgICAgPGcgaWQ9ImNsaWNreWFiLWVuIiB0cmFuc2Zvcm09InRyYW5zbGF0ZSgwLjAwMDAwMCwgMi4wMDAwMDApIj4gICAgICAgICAgICAgICAgICAgIDxnIGlkPSJHcm91cCIgdHJhbnNmb3JtPSJ0cmFuc2xhdGUoMC4wMDAwMDAsIDIuMDAwMDAwKSI+ICAgICAgICAgICAgICAgICAgICAgICAgPGcgaWQ9IlNoYXBlIiBmaWxsPSIjNDFCNkU2Ij4gICAgICAgICAgICAgICAgICAgICAgICAgICAgPHBhdGggaWQ9InVwTG9nbyIgZD0iTTEzLjM0NzE0NTUsMy41MTgxODE4MiBDMTEuOTk5Nzc2MSwxLjI3OTMzODg0IDkuNjE3Njg2NTcsMC4wMzcxOTAwODI2IDcuMTgzNDg4ODEsMC4wMzcxOTAwODI2IEw3LjE4MzQ4ODgxLDAuMDUyMDY2MTE1NyBMNy4xODM0ODg4MSwwLjA1MjA2NjExNTcgTDcuMTgzNDg4ODEsMC4wMzcxOTAwODI2IEw3LjE4MzQ4ODgxLDAuMDM3MTkwMDgyNiBDNS45MjU0NDc3NiwwLjAzNzE5MDA4MjYgNC42NDUwNzQ2MywwLjM3MTkwMDgyNiAzLjQ5MTI1LDEuMDYzNjM2MzYgQzAuMDg5MzI4MzU4MiwzLjEwOTA5MDkxIC0xLjAxMjM4ODA2LDcuNTEyMzk2NjkgMS4wMzQ3MjAxNSwxMC45MTE1NzAyIEMyLjM4MjA4OTU1LDEzLjE1MDQxMzIgNC43NjQxNzkxLDE0LjM5MjU2MiA3LjE5ODM3Njg3LDE0LjM5MjU2MiBDOC40NTY0MTc5MSwxNC4zOTI1NjIgOS43MzY3OTEwNCwxNC4wNTc4NTEyIDEwLjg5MDYxNTcsMTMuMzY2MTE1NyBDMTQuMjg1MDkzMywxMS4zMjgwOTkyIDE1LjM4NjgwOTcsNi45MTczNTUzNyAxMy4zNDcxNDU1LDMuNTE4MTgxODIgWiI+PC9wYXRoPiAgICAgICAgICAgICAgICAgICAgICAgIDwvZz4gICAgICAgICAgICAgICAgICAgICAgICA8cGF0aCBkPSJNMTAuNzQxNzM1MSw2LjEyODkyNTYyIEMxMC42ODk2MjY5LDYuMDQ3MTA3NDQgMTAuNjMwMDc0Niw1Ljk3MjcyNzI3IDEwLjU1NTYzNDMsNS45MDU3ODUxMiBDMTAuNTEwOTcwMSw1Ljg2ODU5NTA0IDEwLjQ1ODg2MTksNS44MzE0MDQ5NiAxMC4zOTkzMDk3LDUuODAxNjUyODkgTDEwLjM5MTg2NTcsNS44MTY1Mjg5MyBMMTAuMzkxODY1Nyw1LjgxNjUyODkzIEwxMC4zOTkzMDk3LDUuODAxNjUyODkgTDEwLjM4NDQyMTYsNS43OTQyMTQ4OCBMMTAuMzg0NDIxNiw1Ljc5NDIxNDg4IEw1Ljc2OTEyMzEzLDMuMzM5NjY5NDIgTDUuNzYxNjc5MSwzLjMzMjIzMTQgTDUuNzYxNjc5MSwzLjMzMjIzMTQgQzUuNzMxOTAyOTksMy4zMTczNTUzNyA1LjcwMjEyNjg3LDMuMzAyNDc5MzQgNS42NzIzNTA3NSwzLjI5NTA0MTMyIEM1LjYzNTEzMDYsMy4yODAxNjUyOSA1LjU5NzkxMDQ1LDMuMjY1Mjg5MjYgNS41NTMyNDYyNywzLjI1Nzg1MTI0IEM1LjMxNTAzNzMxLDMuMTk4MzQ3MTEgNS4wNjkzODQzMywzLjIzNTUzNzE5IDQuODYwOTUxNDksMy4zNjE5ODM0NyBDNC42NTI1MTg2NiwzLjQ4ODQyOTc1IDQuNTAzNjM4MDYsMy42ODkyNTYyIDQuNDQ0MDg1ODIsMy45MjcyNzI3MyBDNC40MzY2NDE3OSwzLjk3MTkwMDgzIDQuNDI5MTk3NzYsNC4wMDkwOTA5MSA0LjQyMTc1MzczLDQuMDUzNzE5MDEgQzQuNDIxNzUzNzMsNC4wODM0NzEwNyA0LjQxNDMwOTcsNC4xMTMyMjMxNCA0LjQxNDMwOTcsNC4xNDI5NzUyMSBMNC40MTQzMDk3LDQuMTQyOTc1MjEgTDQuNDE0MzA5Nyw5LjM3OTMzODg0IEw0LjQyMTc1MzczLDkuMzc5MzM4ODQgTDQuNDIxNzUzNzMsOS4zOTQyMTQ4OCBDNC40MjE3NTM3Myw5LjQ1MzcxOTAxIDQuNDI5MTk3NzYsOS41MjA2NjExNiA0LjQ0NDA4NTgyLDkuNTgwMTY1MjkgQzQuNDY2NDE3OTEsOS42NzY4NTk1IDQuNTAzNjM4MDYsOS43NjYxMTU3IDQuNTU1NzQ2MjcsOS44NTUzNzE5IEM0LjY4MjI5NDc4LDEwLjA2MzYzNjQgNC44ODMyODM1OCwxMC4yMTIzOTY3IDUuMTIxNDkyNTQsMTAuMjcxOTAwOCBDNS4zNTk3MDE0OSwxMC4zMzE0MDUgNS42MTI3OTg1MSwxMC4yOTQyMTQ5IDUuODIxMjMxMzQsMTAuMTY3NzY4NiBDNS45MjU0NDc3NiwxMC4xMDA4MjY0IDYuMDIyMjIwMTUsMTAuMDE5MDA4MyA2LjA4OTIxNjQyLDkuOTIyMzE0MDUgQzYuMTExNTQ4NTEsOS44ODUxMjM5NyA2LjE0MTMyNDYzLDkuODQ3OTMzODggNi4xNTYyMTI2OSw5LjgxMDc0MzggTDYuMTYzNjU2NzIsOS43OTU4Njc3NyBMNi4xNjM2NTY3Miw5Ljc5NTg2Nzc3IEM2LjE3ODU0NDc4LDkuNzczNTUzNzIgNi4xODU5ODg4MSw5Ljc0MzgwMTY1IDYuMjAwODc2ODcsOS43MTQwNDk1OSBMNi4yMDA4NzY4Nyw5LjcxNDA0OTU5IEw2Ljc1OTE3OTEsOC4zMTU3MDI0OCBMOC4yMzMwOTcwMSwxMC43NzAyNDc5IEw4LjI1NTQyOTEsMTAuOCBMOC4yNTU0MjkxLDEwLjggQzguNDM0MDg1ODIsMTEuMDY3NzY4NiA4LjczMTg0NzAxLDExLjIxNjUyODkgOS4wMjk2MDgyMSwxMS4yMTY1Mjg5IEM5LjE5MzM3Njg3LDExLjIxNjUyODkgOS4zNTcxNDU1MiwxMS4xNzE5MDA4IDkuNDk4NTgyMDksMTEuMDkwMDgyNiBDOS45MjI4OTE3OSwxMC44MzcxOTAxIDEwLjA2NDMyODQsMTAuMjg2Nzc2OSA5LjgzMzU2MzQzLDkuODU1MzcxOSBMOS44MzM1NjM0Myw5Ljg1NTM3MTkgTDguMzM3MzEzNDMsNy4zNzEwNzQzOCBMOS44NDEwMDc0Niw3LjUzNDcxMDc0IEw5Ljg0MTAwNzQ2LDcuNTM0NzEwNzQgQzkuODcwNzgzNTgsNy41NDIxNDg3NiA5LjkwODAwMzczLDcuNTQyMTQ4NzYgOS45MzAzMzU4Miw3LjU0MjE0ODc2IEw5Ljk0NTIyMzg4LDcuNTQyMTQ4NzYgQzkuOTg5ODg4MDYsNy41NDIxNDg3NiAxMC4wMjcxMDgyLDcuNTQyMTQ4NzYgMTAuMDcxNzcyNCw3LjUzNDcxMDc0IEMxMC4xOTgzMjA5LDcuNTE5ODM0NzEgMTAuMzA5OTgxMyw3LjQ3NTIwNjYxIDEwLjQyMTY0MTgsNy40MDgyNjQ0NiBDMTAuNjMwMDc0Niw3LjI4MTgxODE4IDEwLjc3ODk1NTIsNy4wODA5OTE3NCAxMC44Mzg1MDc1LDYuODQyOTc1MjEgQzEwLjkwNTUwMzcsNi41ODI2NDQ2MyAxMC44NjgyODM2LDYuMzM3MTkwMDggMTAuNzQxNzM1MSw2LjEyODkyNTYyIFoiIGlkPSJTaGFwZSIgZmlsbD0iI0ZGRkZGRiI+PC9wYXRoPiAgICAgICAgICAgICAgICAgICAgICAgIDxwYXRoIGlkPSJkb3duTG9nbyIgZD0iTTEzLjk0MjY2NzksMTYuMzA0MTMyMiBMMTIuODMzNTA3NSwxNC40NTk1MDQxIEMxMi42ODQ2MjY5LDE0LjIwNjYxMTYgMTIuNDM4OTczOSwxNC4wMjgwOTkyIDEyLjE1NjEwMDcsMTMuOTUzNzE5IEMxMS44NjU3ODM2LDEzLjg3OTMzODggMTEuNTc1NDY2NCwxMy45MjM5NjY5IDExLjMyMjM2OTQsMTQuMDgwMTY1MyBDMTAuODAxMjg3MywxNC4zOTI1NjIgMTAuNjMwMDc0NiwxNS4wNzY4NTk1IDEwLjk0MjcyMzksMTUuNTk3NTIwNyBMMTIuMDUxODg0MywxNy40NDIxNDg4IEMxMi4yNTI4NzMxLDE3Ljc2OTQyMTUgMTIuNjE3NjMwNiwxNy45Nzc2ODYgMTMuMDA0NzIwMSwxNy45Nzc2ODYgTDEzLjAwNDcyMDEsMTcuOTc3Njg2IEMxMy4yMDU3MDksMTcuOTc3Njg2IDEzLjM5OTI1MzcsMTcuOTI1NjE5OCAxMy41NzA0NjY0LDE3LjgyMTQ4NzYgQzEzLjgyMzU2MzQsMTcuNjcyNzI3MyAxNC4wMDIyMjAxLDE3LjQyNzI3MjcgMTQuMDc2NjYwNCwxNy4xNDQ2MjgxIEMxNC4xMzYyMTI3LDE2Ljg1NDU0NTUgMTQuMDkxNTQ4NSwxNi41NTcwMjQ4IDEzLjk0MjY2NzksMTYuMzA0MTMyMiBaIiBmaWxsPSIjNDFCNkU2Ij48L3BhdGg+ICAgICAgICAgICAgICAgICAgICA8L2c+ICAgICAgICAgICAgICAgIDwvZz4gICAgICAgICAgICA8L2c+ICAgICAgICA8L2c+ICAgIDwvZz48L3N2Zz4="/>
                  پیشنهاد شده توسط کلیک‌یاب
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
                <div class="cyb-img-holder  cyb-custom-img-holder">
                    <a rel="nofollow" target="_blank" href="{{.URL}}" onclick="cybOpen(event)" oncontextmenu="cybOpen(event)"
                       ondblclick="cybOpen(event)" data-href="{{.URL}}">
                        <img src="{{.Image}}" alt="{{.Title}}"
                             class="cyb-img {{isRound .Corners}} cyb-custom-img">
                    </a>
                </div>
                <div class="cyb-desc-holder cyb-custom-desc-holder">
                    <div class="cyb-desc cyb-custom-desc">
                        <a rel="nofollow" target="_blank" href="{{.URL}}" onclick="cybOpen(event)" oncontextmenu="cybOpen(event)"
                           ondblclick="cybOpen(event)" data-href="{{.URL}}">
                            {{.Title}}
                        </a>
                        <!--{{ .Extra }}-->
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
			b.WriteString(`<div class="cyb-pack cyb-min-size cyb-custom-pack">`)
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
	var resStyle = style
	if imp.MinSize != "" {
		num, err := strconv.ParseInt(imp.MinSize, 10, 64)
		halfNum := num / 2
		if err == nil {
			rep := strings.NewReplacer("cyb-minSize", fmt.Sprintf("%dpx", halfNum), "cyb-doubleMinSize", fmt.Sprintf("%dpx", num))
			resStyle = rep.Replace(style)
		}
	}
	imp.Style = resStyle
	e := native.Lookup("ads").Execute(buf, imp)
	assert.Nil(e)
	return string(buf.Bytes())
}
func init() {
	native.Parse(nativeTmpl)
	native.Parse(adTmpl)
}
