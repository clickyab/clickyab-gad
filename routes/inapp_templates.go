package routes

import (
	"bytes"
	"text/template"
)

type inappContext struct {
	FullScreen    bool
	ExtraStyle    string
	BodyClass     string
	Dynamic       bool
	DynamicBody   string
	FatFinger     bool
	ClickURL      string
	Src           string
	CloseClass    string
	ImpID         int
	SdkVersion    int64
	RefreshMinute int
	NoAd          bool
}

const inappTemplateText = `
<!doctype html>
<html>
<head>
    <meta charset="UTF-8">
    <meta http-equiv="content-type" content="text/html; charset=UTF-8"/>
    <meta name="viewport" content="initial-scale=1, maximum-scale=1, user-scalable=no">
    <title>Clickyab</title>
    <script type="text/javascript">function showb() {
            document.getElementById("innershow").style.display = "block";
            document.getElementById("clickyab_overlay").style.display = "block";
        }
    </script>
    <style>
    {{ if .FullScreen }}
*{padding:0;margin:0;}
html,body{width:100%;height:100%;background:#fff;}
a,img{display:block;float:left;cursor:pointer;text-decoration:none}
.portrait a,.portrait img{width:320px;height:480px;}
.landscape a,.landscape img{width:480px;height:320px;}
a.close,a.largeclose{width:24px;height:24px;line-height:24px;font-size:18px;background:rgba(62,73,90,0.92);color:#FFF;text-align:center;position:absolute;display:inline-block;font-family:verdana,helvetica,arial,sans-serif;font-weight:bold;
}
a.largeclose{width:32px;height:32px;line-height:30px;font-size:24px;}
.portrait a.close{left:296px;top:0;}
.portrait a.largeclose{left:288px;}
.landscape a.close{left:456px;top:0;}
.landscape a.largeclose{left:448px;}
    {{ else }}
*{padding: 0;margin: 0;}
html,body{width: 100%;height: 100%;background: #fff;}
a,img{display: block;float: left;width: 100%;height: 100%;cursor: pointer; text-decoration: none}
a.close,a.largeclose{
    width: 24px;
    height: 24px;
    line-height: 24px;
    font-size: 18px;
    background: rgba(62, 73, 90, 0.92);
    color: #FFF;
    text-align: center;
    position: absolute;
    left: 0;
    bottom: 0;
    display: inline-block;
    font-family: verdana,helvetica,arial,sans-serif;
    font-weight: bold;
}
a.largeclose{ width: 32px; height: 32px; line-height: 32px; font-size: 24px; }
    {{ end }}
    {{ .ExtraStyle }}
    #showb {
        position: relative;border: 0;line-height: 1;
        height: 100%;
        width: 100%;
        z-index: 100;
        text-align: center;
    }
    #showb a {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%,-50%);
        -moz-transform: translate(-50%,-50%);
        -webkit-transform: translate(-50%,-50%);
        -o-transform: translate(-50%,-50%);
        -ms-transform: translate(-50%,-50%);
        background:  #000;
        border:1px solid #eee;
        color:#fff;
        font-family: tahoma,Arial,"Times New Roman";
        display: inline-block;
        cursor: pointer;
        width: 100%;
        padding: 10px 0;
        display: none;
        z-index: 10;
        height: auto !important;
    }
    #showb #clickyab_overlay {
        display: none;
        width: 100%;
        height: 100%;
        position: absolute;
        left: 0;
        top: 0;
        z-index: 101;
        background: rgba(0,0,0,0.6);
        z-index: 1;

    }
    </style>
</head>
<body {{ .BodyClass }}>
{{ if .Dynamic }}
    {{ .DynamicBody }}
{{ else }}
    {{ if .FatFinger }}
<div id="showb" class="adhere" onclick="showb()">
    <div id="clickyab_overlay"></div>
    <a id="innershow" onClick="onClickyabClicked('{{ .ClickURL }}', null, null)">مشاهده کنید</a>
    <img src="{{ .Src }}"/>
</div>
    {{ else }}
        <a onClick="onClickyabClicked('{{ .ClickURL }}',null,null)"><img src="{{ .Src }}"></a>;
    {{ end }}
{{ end }}

<a onclick="{{ if .FullScreen }}AndroidHide(){{ else }}AndroidClose(){{ end }};" class="{{ .CloseClass }}">x</a>

<script type="text/javascript">

    function showHitted() {
        document.getElementById("hitted").style.display = 'block';
    }

    var ads = document.getElementById('ads');
    if (ads != null) document.getElementById("ads").onclick = showHitted;

    //function
    var impId  = {{ .ImpID }};
    function AndroidSetImpId(impId) {
        clickyab.setImpId(impId);
    }
    function AndroidHide() {
        clickyab.hide();
    }
    function AndroidClose() {
        clickyab.closeFullAd();
    }
    function AndroidShow() {
        clickyab.show();
    }
    function AndroidRefresh() {
        clickyab.refresh();
    }
    function AndroidRefreshIfActive() {
        clickyab.refreshIfActive();
    }
    function AndroidOpenIntent(url) {
        clickyab.openIntent(url);
    }
    function AndroidOpenIntentWithin(url,packagename) {
        clickyab.openIntentWithin(url,packagename);
    }
    function AndroidOnClick() {
        clickyab.onClicked();
    }
    function AndroidHit(url,impId) {
        clickyab.hit(url,impId);
    }
    function AndroidHasNoAds() {
        clickyab.hasNoAds();
    }
    function AndroidSetHasAds(p) {
        clickyab.setHasAds(p);
    }
    //onAdsClick
    function onClickyabClicked(url,packagename,hitUrl) {
        AndroidOnClick();
        //tell android to submit hit to server
        AndroidHit(hitUrl,impId);
        if (packagename == '' || packagename == null){
            AndroidOpenIntent(url);
        }else{
            AndroidOpenIntentWithin(url,packagename);
        }
        AndroidRefreshIfActive();
    }

    document.addEventListener('DOMContentLoaded', function() {
        AndroidSetImpId(impId);
        {{ if .NoAd }}
            {{ if eq .SdkVersion 3 }}
            setTimeout(function () {
                AndroidHasNoAds();
            }, 100);// 0.1 sec
            {{ else if eq .SdkVersion 4 }}
            setTimeout(function () {
                AndroidSetHasAds(false);
            }, 100);// 0.1 sec
            {{ else }}
            setTimeout(function () {
                {{ if .FullScreen }}AndroidHide(){{ else }}AndroidClose(){{ end }};
            }, 100);// 0.1 sec
            {{ end }}
        {{ else }}
            {{ if eq .SdkVersion 4 }}
            setTimeout(function () {
                AndroidSetHasAds(true);
            }, 100);// 0.1 sec
            {{ end }}
            setInterval(function(){
                AndroidRefreshIfActive();
            },{{ .RefreshMinute }}*60*1000);// 60 sec
        {{ end }}
    });
</script>

</body>
</html>`

var inappTemplate = template.Must(template.New("inapp-template").Parse(inappTemplateText))

func renderInApp(data inappContext) (string, error) {
	res := &bytes.Buffer{}
	err := inappTemplate.Execute(res, data)
	return res.String(), err
}
