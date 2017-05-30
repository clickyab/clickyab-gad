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
		Funcs(template.FuncMap{"isCircleCorner": isCircleCorner}).
		Parse(l.String())
	assert.Nil(e)
	buf := &bytes.Buffer{}
	for _, ad := range ads {
		buf.WriteString(`<div class="native-grids">
		<div class="native-element  ">`)
		adTemplate.Execute(buf, ad)
		buf.WriteString(`</div></div>`)
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
		return "circle-title"
	}
	return ""
}

var layoutString = [...]string{
	`<div class="native-border {{.Corners}}">
		<a href="{{.Site}}"  data-href="{{.URL}}"   onclick="cl_na_ha_ko_blah_blah_blah(event)"><img src="{{.Image}}"></a>
	</div>
	<div class="native-content ">
		<a href="{{.Site}}" data-href="{{.URL}}"  onclick="cl_na_ha_ko_blah_blah_blah(event)"><span>{{.Title}}</span></a>
		<!--<p>{{.Lead}}</p>-->
		<a href="{{.Site}}" data-href="{{.URL}}"  onclick="cl_na_ha_ko_blah_blah_blah(event)" class="btn btn-default ">{{.More}}</a>
	</div>`,
	`<a href="{{.Site}}" onclick="cl_na_ha_ko_blah_blah_blah(event)" data-href="{{.URL}}"><span class="headline ">{{.Title}}</span></a>
	<div class="native-border {{.Corners}} ">
		<a href="{{.Site}}" onclick="cl_na_ha_ko_blah_blah_blah(event)" data-href="{{.URL}}"><img src="{{.Image}}" ></a>
	</div>
	<div class="native-content ">
		<!--<p>{{.Lead}}</p>-->
		<a href="{{.Site}}"  data-href="{{.URL}}"  onclick="cl_na_ha_ko_blah_blah_blah(event)" class="btn btn-default">{{.More}}</a>
	</div>`,
	`<div class="native-content ">
		<a href="{{.Site}}" onclick="cl_na_ha_ko_blah_blah_blah(event)" data-href="{{.URL}}"><span class="headline">{{.Title}}</span></a>
			<!--<p>{{.Lead}}</p>-->
		</div>
		<div class="native-border {{.Corners}} ">
			<a href="{{.Site}}" onclick="cl_na_ha_ko_blah_blah_blah(event)" data-href="{{.URL}}"><img src="{{.Image}}"></a>
		</div>`,
	`<a class="{{isCircleCorner .Corners}} sit-left" href="{{.Site}}" onclick="cl_na_ha_ko_blah_blah_blah(event)"  data-href="{{.URL}}"><p class="headline">{{.Title}}</p></a>
		<div class="native-border {{.Corners}} sit-right ">
			<a class="{{isCircleCorner .Corners}}" href="{{.Site}}" onclick="cl_na_ha_ko_blah_blah_blah(event)" data-href="{{.URL}}"><img src="{{.Image}}"></a>
		</div>`,
	`<a href="{{.Site}}" onclick="cl_na_ha_ko_blah_blah_blah(event)" class="{{isCircleCorner .Corners}} sit-right"  data-href="{{.URL}}"><p class="headline">{{.Title}}</p></a>
		<div class="native-border {{.Corners}} sit-left ">
			<a href="{{.Site}}"  onclick="cl_na_ha_ko_blah_blah_blah(event)" data-href="{{.URL}}"><img src="{{.Image}}"></a>
		</div>`,
}

func (t layout) String() string {
	return layoutString[t]
}

var adTemplate = `
{{define "ads"}}
{{template "head" .}}
<div class="header-line " >
	<div class="line "></div>
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
	head = `{{define "head"}}<!doctype html>
<html lang="fa">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="ie=edge">
	<style>
	{{template "style"}}
	</style>
</head>
<body class="{{.Font}}" style="font-size:{{.FontSize}};font-family:{{.Font}};">
	<div class="native-ad-wrapper {{isHorizontal  .Layout}} ">
	{{end}}`

	foot = `{{define "foot"}}</div>


		</body>
	</html>
	{{end}}`
)

const style = `{{define "style"}}
   .native-ad-wrapper {
            height: 340px;
        }

        .native-element a:hover {
            text-decoration: none;
        }

        .header-line {
            margin-bottom: 12px;
            position: relative;
        }

        .simple-header {
            position: relative;
        }

        .simple-header p {
            position: absolute;
            top: -39px;
            right: 0;
            font-size: 13px;
            font-weight: bold;
        }

        .simple-line {
            height: 2px;
            width: 100%;
            background-color: rgba(0, 0, 0, 0.2)
        }

        .header-line .line {
            height: 13px;
            background-color: rgba(0, 0, 0, 0.3);
        }

        .header-line p {
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

        .native-grids {
            width: 25%;
            max-width: 367px;
            min-width: 366px;
            float: right;
            padding: 5px;
            box-sizing: border-box;
        }

        .native-element {
            height: auto;
            width: 100%;
        }

        .native-border {
            border: 1px solid rgba(0, 0, 0, 0.3);
            width: 100%;
            height: 175px;
            padding: 8px 8px 4px 8px;
            box-sizing: border-box;
        }

        .native-border a {
            width: 100%;
            height: 100%;
        }

        a:hover, a:visited, a:link, a:active {
            text-decoration: none;
            border: none !important;
        }

        .native-ad-wrapper.horizontal {
            height: 110px;
        }

        .horizontal div.native-grids div.native-element {
            height: 94px;
        }

        .horizontal div.native-grids div.native-element div.native-border {
            height: 94px;
            padding: 3px;
            width: 44%;
        }

        .horizontal div.native-grids div.native-element > a {
            padding: 10px;
            width: 54%;
            font-size: 11px;
            box-sizing: border-box;
        }

        .sit-left {
            float: left;
        }

        .sit-right {
            float: right;
        }

        .native-border img {
            width: 100%;
            height: 100%;
        }

        .native-content {
            text-align: right;
        }

        .native-content span, .headline {
            font-size: 13px;
            color: rgba(0, 0, 0, 0.8);
            direction: rtl;
            margin: 10px 0;
            display: block;
            font-weight: bold;
            line-height: 21px;
        }

        .native-content p {
            font-size: 12px;
        }

        .round {
            border-radius: 5px;
        }

        .round img {
            border-radius: 5px;
        }

        .circle {
            width: 200px;
            height: 200px;
            border-radius: 50%;
            margin: 0 auto;
            padding: 7px;
        }

        .circle img {
            width: 185px;
            height: 185px;
            border-radius: 50%;
        }

        .horizontal .circle {
            width: 94px !important;
            height: 94px !important;
            border-radius: 50%;
            margin: 0 auto;
            padding: 4px !important;
        }

        .horizontal .circle img {
            width: 84px;
            height: 84px;
            border-radius: 50%;
        }

        .circle-title {
            width: 68% !important;
        }

        .btn {
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

        .btn:focus,
        .btn:active:focus,
        .btn.active:focus,
        .btn.focus,
        .btn:active.focus,
        .btn.active.focus {
            outline: 5px auto -webkit-focus-ring-color;
            outline-offset: -2px;
        }

        .btn:hover,
        .btn:focus,
        .btn.focus {
            color: #333;
            text-decoration: none;
        }

        .btn:active,
        .btn.active {
            background-image: none;
            outline: 0;
            -webkit-box-shadow: inset 0 3px 5px rgba(0, 0, 0, .125);
            box-shadow: inset 0 3px 5px rgba(0, 0, 0, .125);
        }

        .btn.disabled,
        .btn[disabled],
        fieldset[disabled] .btn {
            cursor: not-allowed;
            filter: alpha(opacity=65);
            -webkit-box-shadow: none;
            box-shadow: none;
            opacity: .65;
        }

        a.btn.disabled,
        fieldset[disabled] a.btn {
            pointer-events: none;
        }

        .btn-default {
            color: #333;
            background-color: #fff;
            border-color: #ccc;
            border-radius: 13px;
            font-size: 10px;
        }

        .btn-default:focus,
        .btn-default.focus {
            color: #333;
            background-color: #e6e6e6;
            border-color: #8c8c8c;
        }

        .btn-default:hover {
            color: #333;
            background-color: #e6e6e6;
            border-color: #adadad;
        }

        .btn-default:active,
        .btn-default.active,
        .open > .dropdown-toggle.btn-default {
            color: #333;
            background-color: #e6e6e6;
            border-color: #adadad;
        }

        .btn-default:active:hover,
        .btn-default.active:hover,
        .open > .dropdown-toggle.btn-default:hover,
        .btn-default:active:focus,
        .btn-default.active:focus,
        .open > .dropdown-toggle.btn-default:focus,
        .btn-default:active.focus,
        .btn-default.active.focus,
        .open > .dropdown-toggle.btn-default.focus {
            color: #333;
            background-color: #d4d4d4;
            border-color: #8c8c8c;
        }

        .btn-default:active,
        .btn-default.active,
        .open > .dropdown-toggle.btn-default {
            background-image: none;
        }

        .btn-default.disabled:hover,
        .btn-default[disabled]:hover,
        fieldset[disabled] .btn-default:hover,
        .btn-default.disabled:focus,
        .btn-default[disabled]:focus,
        fieldset[disabled] .btn-default:focus,
        .btn-default.disabled.focus,
        .btn-default[disabled].focus,
        fieldset[disabled] .btn-default.focus {
            background-color: #fff;
            border-color: #ccc;
        }

        .btn-default .badge {
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
