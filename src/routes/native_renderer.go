package routes

import (
	"assert"
	"bytes"
	"text/template"
)

var nativeTemplate = template.New("native").
	Funcs(template.FuncMap{"renderAds": renderAds, "isHorizontal": isHorizontal})

func renderAds(l layout, ads []nativeAd) string {
	adTemplate, e := template.New("ad").
		Funcs(template.FuncMap{"isCircleCorner": isCircleCorner, "isCircleImage": isCircleImage}).
		Parse(l.String())
	assert.Nil(e)
	buf := &bytes.Buffer{}
	closer := 0
	for i, ad := range ads {
		if i != 0 && i == closer {
			buf.WriteString(`</div>`)

		}
		if i%4 == 0 {
			closer += 4
			buf.WriteString(`<div class="cyb-row clickyab-custom-row">`)
		}

		buf.WriteString(`<div class="cyb-native-grids">
		<div class="native-element  ">`)

		adTemplate.Execute(buf, ad)
		buf.WriteString(`</div></div>`)
		if len(ads)-1 == i {
			buf.WriteString(`</div>`)

		}
	}
	return string(buf.Bytes())
}

type protocol string

const (
	httpScheme  protocol = "http"
	httpsScheme          = "https"
)

type layout int

const (
	layoutImageFirst layout = iota
	layoutTitleFirst
	layoutImageLast
	layoutImageRight
	layoutTitleRight
)

func isHorizontal(l layout) string {
	if l == layoutImageRight || l == layoutTitleRight {
		return "horizontal"
	}
	return ""
}

func isCircleCorner(c string) string {
	if c == "circle" {
		return "cyb-circle-title"
	}
	return ""
}

func isCircleImage(c string) string {
	if c == "circle" {
		return "cyb-cl-na-image-circle"
	}
	return "cyb-cl-na-image"
}

var layoutString = [...]string{
	`<div class="cyb-native-border {{.Corners}}">
		<a target="_blank"  href="{{.Site}}"  data-href="{{.URL}}"   onclick="cybOpen(event)" oncontextmenu="cybOpen(event)" ondblclick="cybOpen(event)">
				<div class="{{isCircleImage .Corners}}" style="background-image: url('{{.Image}}');" ></div>
		</a>
	</div>
	<div class="cyb-native-content ">
		<a target="_blank" href="{{.Site}}" data-href="{{.URL}}"  onclick="cybOpen(event)" oncontextmenu="cybOpen(event)" ondblclick="cybOpen(event)"><span>{{.Title}}</span></a>
		<!--<p>{{.Lead}}</p>-->
		<a target="_blank" href="{{.Site}}" data-href="{{.URL}}"  onclick="cybOpen(event)" oncontextmenu="cybOpen(event)" ondblclick="cybOpen(event)" class="cyb-btn cyb-btn-default ">{{.More}}</a>
	</div>`,
	`<a target="_blank" href="{{.Site}}" onclick="cybOpen(event)" oncontextmenu="cybOpen(event)" ondblclick="cybOpen(event)" data-href="{{.URL}}"><span class="cyb-headline ">{{.Title}}</span></a>
	<div class="cyb-native-border {{.Corners}} ">
		<a target="_blank" href="{{.Site}}" onclick="cybOpen(event)" oncontextmenu="cybOpen(event)" ondblclick="cybOpen(event)" data-href="{{.URL}}">
		<div class="{{isCircleImage .Corners}}"  style="background-image: url('{{.Image}}');" ></div>
		</a>
	</div>
	<div class="cyb-native-content ">
		<!--<p>{{.Lead}}</p>-->
		<a target="_blank" href="{{.Site}}"  data-href="{{.URL}}"  onclick="cybOpen(event)" oncontextmenu="cybOpen(event)" ondblclick="cybOpen(event)" class="btn btn-default">{{.More}}</a>
	</div>`,
	`<div class="cyb-native-content ">
		<a target="_blank" href="{{.Site}}" onclick="cybOpen(event)" oncontextmenu="cybOpen(event)" ondblclick="cybOpen(event)" data-href="{{.URL}}"><span class="cyb-headline">{{.Title}}</span></a>
			<!--<p>{{.Lead}}</p>-->
		</div>
		<div class="cyb-native-border {{.Corners}} ">
			<a target="_blank" href="{{.Site}}" onclick="cybOpen(event)" oncontextmenu="cybOpen(event)" ondblclick="cybOpen(event)" data-href="{{.URL}}">
					<div class="{{isCircleImage .Corners}}"  style="background-image: url('{{.Image}}');" ></div>
			</a>
		</div>`,
	`<a target="_blank" class="{{isCircleCorner .Corners}} sit-left" href="{{.Site}}" onclick="cybOpen(event)" oncontextmenu="cybOpen(event)" ondblclick="cybOpen(event)"  data-href="{{.URL}}"><p class="cyb-headline">{{.Title}}</p></a>
		<div class="cyb-native-border {{.Corners}} sit-right ">
			<a target="_blank" class="{{isCircleCorner .Corners}}" href="{{.Site}}" onclick="cybOpen(event)" oncontextmenu="cybOpen(event)" ondblclick="cybOpen(event)" data-href="{{.URL}}">
					<div class="{{isCircleImage .Corners}}"  style="background-image: url('{{.Image}}');" ></div>
			</a>
		</div>`,
	`<a target="_blank" href="{{.Site}}" onclick="cybOpen(event)" oncontextmenu="cybOpen(event)" ondblclick="cybOpen(event)" class="{{isCircleCorner .Corners}} sit-right"  data-href="{{.URL}}"><p class="cyb-headline">{{.Title}}</p></a>
		<div class="cyb-native-border {{.Corners}} sit-left ">
			<a target="_blank"  href="{{.Site}}"  onclick="cybOpen(event)" oncontextmenu="cybOpen(event)" ondblclick="cybOpen(event)" data-href="{{.URL}}">
					<div class="{{isCircleImage .Corners}}"  style="background-image: url('{{.Image}}');" ></div>
			</a>
		</div>`,
}

func (t layout) String() string {
	return layoutString[t]
}

var adTemplate = `
{{define "ads"}}
{{template "head" .}}
<div class="cyb-header-line " >
	<div class="cyb-line "></div>
	<p >{{.Title}}</p>
</div>

{{renderAds .Layout .Ads}}
{{template "foot"}}
{{end}}
`

type nativeContainer struct {
	Ads      []nativeAd
	Layout   layout
	Title    string
	Style    string
	FontSize string
	Font     string
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

func renderNative(ads nativeContainer) string {
	buf := &bytes.Buffer{}
	e := nativeTemplate.Lookup("ads").Execute(buf, ads)
	assert.Nil(e)
	return string(buf.Bytes())
}

const (
	head = `{{define "head"}}

	<style>
	{{template "style"}}
	</style>
	<div class="cyb-native-ad-wrapper {{isHorizontal  .Layout}} ">
	{{end}}`

	foot = `{{define "foot"}}</div>


	{{end}}`
)

const style = `{{define "style"}}

	.cyb-cl-na-image {


	    background-repeat: no-repeat;
		background-size: cover;
	    	background-position: center;
            	padding:31%;

	}

	.cyb-cl-na-image-circle {
		background-repeat: no-repeat;
		background-size: cover;
		background-position: 50%;
		border-radius: 50%;
		width: 100%;
		padding-top: 100%;



    }
    .cyb-cyb-row {
    	    width: 100%;
    		margin: 0 -15px;
    }
   .cyb-native-ad-wrapper {
            height: 340px;
        }

        .cyb-native-element a:hover {
            text-decoration: none;
        }

        .cyb-header-line {
            margin-bottom: 12px;
            position: relative;
        }

        .cyb-simple-header {
            position: relative;
        }

        .cyb-simple-header p {
            position: absolute;
            top: -39px;
            right: 0;
            font-size: 13px;
            font-weight: bold;
        }

        .cyb-simple-line {
            height: 2px;
            width: 100%;
            background-color: rgba(0, 0, 0, 0.cyb-2)
        }

        .cyb-header-line .cyb-line {
            height: 13px;
            background-color: rgba(0, 0, 0, 0.cyb-3);
        }

        .cyb-header-line p {
            position: absolute;
            top: -13px;
            font-size: 13px;
            font-weight: bold;
            background-color: white;
            display: inline;
            right: 0;
            text-align: right;
            padding-left: 9px;
        }

        .cyb-native-grids {
            width: 25%;
            height: 100%;
            float: right;
            padding: 5px;
            box-sizing: border-box;
        }

        .cyb-native-element {
            height: 100%;
            width: 100%;
        }

        .cyb-native-border {
            box-sizing: border-box;
            overflow: hidden;
        }

        .cyb-native-border a {
            width: 100%;
            height: 100%;
        }

        a:hover, a:visited, a:link, a:active {
            text-decoration: none;
            border: none !important;
        }

        .cyb-native-ad-wrapper.cyb-horizontal {
            height: 110px;
        }

        .cyb-horizontal div.cyb-native-grids div.cyb-native-element {
            height: 94px;
        }

        .cyb-horizontal div.cyb-native-grids div.cyb-native-element div.cyb-native-border {
            height: 94px;
            padding: 3px;
            width: 44%;
        }

        .cyb-horizontal div.cyb-native-grids div.cyb-native-element > a {
            padding: 10px;
            width: 54%;
            font-size: 11px;
            box-sizing: border-box;
        }

        .cyb-sit-left {
            float: left;
        }

        .cyb-sit-right {
            float: right;
        }

        .cyb-native-border img {
            width: 100%;
            height: 100%;
        }

        .cyb-native-content {
            text-align: right;
        }

        .cyb-native-content span, .cyb-headline {
            font-size: 13px;
            color: rgba(0, 0, 0, 0.8);
            direction: rtl;
            margin: 10px 0;
            display: block;
            font-weight: bold;
            line-height: 21px;
        }

        .cyb-native-content p {
            font-size: 12px;
        }

        .cyb-round {
            border-radius: 5px;
        }

        .cyb-round img {
            border-radius: 5px;
        }

        .cyb-circle {
            width: 50%;
            margin: 0 auto;
            padding: 10px;
        }

        .cyb-circle img {
            width: 185px;
            height: 185px;
            border-radius: 50%;
        }

        .cyb-horizontal .cyb-circle {
            width: 94px !important;
            height: 94px !important;
            border-radius: 50%;
            margin: 0 auto;
            padding: 4px !important;
        }

        .cyb-horizontal .cyb-circle img {
            width: 84px;
            height: 84px;
            border-radius: 50%;
        }

        .cyb-circle-title {
            width: 68% !important;
        }

        .cyb-btn {
            display: inline-block;
            padding: 6px 12px;
            margin-bottom: 0;
            font-size: 14px;
            font-weight: normal;
            line-height: 1.42857143;
            text-align: center;
            white-space: nowrap;
            vertical-align: middle;
            -ms-touch-action: manipulation;
            touch-action: manipulation;
            cursor: pointer;
            -webkit-user-select: none;
            -moz-user-select: none;
            -ms-user-select: none;
            user-select: none;
            background-image: none;
            border: 1px solid transparent;
            border-radius: 4px;
        }

        .cyb-btn:focus,
        .cyb-btn:active:focus,
        .cyb-btn.cyb-active:focus,
        .cyb-btn.cyb-focus,
        .cyb-btn:active.cyb-focus,
        .cyb-btn.cyb-active.cyb-focus {
            outline: 5px auto -webkit-focus-ring-color;
            outline-offset: -2px;
        }

        .cyb-btn:hover,
        .cyb-btn:focus,
        .cyb-btn.cyb-focus {
            color: #333;
            text-decoration: none;
        }

        .cyb-btn:active,
        .cyb-btn.cyb-active {
            background-image: none;
            outline: 0;
            -webkit-box-shadow: inset 0 3px 5px rgba(0, 0, 0, .125);
            box-shadow: inset 0 3px 5px rgba(0, 0, 0, .125);
        }

        .cyb-btn.cyb-disabled,
        .cyb-btn[disabled],
        fieldset[disabled] .cyb-btn {
            cursor: not-allowed;
            filter: alpha(opacity=65);
            -webkit-box-shadow: none;
            box-shadow: none;
            opacity: .65;
        }

        a.cyb-btn.cyb-disabled,
        fieldset[disabled] a.cyb-btn {
            pointer-events: none;
        }

        .cyb-btn-default {
            color: #333;
            background-color: #fff;
            border-color: #ccc;
            border-radius: 13px;
            font-size: 10px;
        }

        .cyb-btn-default:focus,
        .cyb-btn-default.cyb-focus {
            color: #333;
            background-color: #e6e6e6;
            border-color: #8c8c8c;
        }

        .cyb-btn-default:hover {
            color: #333;
            background-color: #e6e6e6;
            border-color: #adadad;
        }

        .cyb-btn-default:active,
        .cyb-btn-default.cyb-active,
        .cyb-open > .cyb-dropdown-toggle.cyb-btn-default {
            color: #333;
            background-color: #e6e6e6;
            border-color: #adadad;
        }

        .cyb-btn-default:active:hover,
        .cyb-btn-default.cyb-active:hover,
        .cyb-open > .cyb-dropdown-toggle.cyb-btn-default:hover,
        .cyb-btn-default:active:focus,
        .cyb-btn-default.cyb-active:focus,
        .cyb-open > .cyb-dropdown-toggle.cyb-btn-default:focus,
        .cyb-btn-default:active.cyb-focus,
        .cyb-btn-default.cyb-active.cyb-focus,
        .cyb-open > .cyb-dropdown-toggle.cyb-btn-default.cyb-focus {
            color: #333;
            background-color: #d4d4d4;
            border-color: #8c8c8c;
        }

        .cyb-btn-default:active,
        .cyb-btn-default.cyb-active,
        .cyb-open > .cyb-dropdown-toggle.cyb-btn-default {
            background-image: none;
        }

        .cyb-btn-default.cyb-disabled:hover,
        .cyb-btn-default[disabled]:hover,
        fieldset[disabled] .cyb-btn-default:hover,
        .cyb-btn-default.cyb-disabled:focus,
        .cyb-btn-default[disabled]:focus,
        fieldset[disabled] .cyb-btn-default:focus,
        .cyb-btn-default.cyb-disabled.cyb-focus,
        .cyb-btn-default[disabled].cyb-focus,
        fieldset[disabled] .cyb-btn-default.cyb-focus {
            background-color: #fff;
            border-color: #ccc;
        }

        .cyb-btn-default .cyb-badge {
            color: #fff;
            background-color: #333;
        }

{{end}}
`

func init() {
	_, e := nativeTemplate.Parse(adTemplate)
	assert.Nil(e)
	_, e = nativeTemplate.Parse(style)
	assert.Nil(e)
	_, e = nativeTemplate.Parse(head)
	assert.Nil(e)
	_, e = nativeTemplate.Parse(foot)
	assert.Nil(e)
}
