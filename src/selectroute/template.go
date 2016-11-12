package selectroute

import "html/template"

const singleAd = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta name="robots" content="nofollow">
    <meta content="text/html; charset=utf-8" http-equiv="Content-Type">
    <style>
        body{ margin: 0; padding: 0; text-align: center; }
        .o{ position:absolute; top:0; left:0; border:0; height:250px; width:300px; z-index: 99; }
        #showb{ position:absolute; top:0; left:0; border:0; line-height: 250px; height:250px; width:300px; z-index: 100; background: rgba(0, 0, 0, 0.60); text-align: center; }
        {{ if .Tiny }}
        .tiny2{ height: 18px; width: 19px; position: absolute; bottom: 0px; right: 0; z-index: 100; background: url("//static.clickyab.com/img/clickyab-tiny.png") right top no-repeat; border-top-left-radius:4px; -moz-border-radius-topleft:4px  }
        .tiny2:hover{ width: 66px;  }
        .tiny{ height: 18px; width: 19px; position: absolute; top: 0px; right: 0; z-index: 100; background: url("//static.clickyab.com/img/clickyab-tiny.png") right top no-repeat; border-bottom-left-radius:4px; -moz-border-radius-bottomleft:4px  }
        .tiny:hover{ width: 66px;  }
        .tiny3{ position: absolute; top: 0px; right: 0; z-index: 100; }
        {{ end }}
        .butl {background: #4474CB;color: #FFF;padding: 10px;text-decoration: none;border: 2px solid #FFFFFF;font-family: tahoma;font-size: 13px;}
        img.adhere {max-width:100%;height:auto;}
        video {background: #232323 none repeat scroll 0 0;}
    </style>
</head>
<body>

    {{ if .Tiny }}<a class="tiny" href="http://clickyab.com/?ref=icon" target="_blank"></a>{{ end }}
	<a href="{{ .Link }}" target="_blank"><img  src="{{ .Src }}" border="0" height="{{ .Height }}" width="{{ .Width }}"/></a>
<br style="clear: both;"/>
</body></html>`

const videoAD=`<div id="video_advertise">
    <video width="{{ .Width }}" height="{{ .Height }}" autoplay controls muted>
        <source src="{{ .Src }}" type="video/mp4">
    </video></div>
<script>
    function unwrap(wrapper) {
        // place childNodes in document fragment
        var docFrag = document.createDocumentFragment();
        while (wrapper.firstChild) {
            var child = wrapper.removeChild(wrapper.firstChild);
            docFrag.appendChild(child);
        }

        // replace wrapper with document fragment
        wrapper.parentNode.replaceChild(docFrag, wrapper);
    }
    var link = "{{ .Link }}";
    org_html = document.getElementById("video_advertise").innerHTML;
    appendHtmlLink = "<a id='a_advertise' target='_blank' href='"+ link +"'>" + org_html + "</a>";
    var FinalElementHtml = document.getElementById("video_advertise").innerHTML = appendHtmlLink;
    document.getElementById("video_advertise").addEventListener("click", function () {
        var linkElement = document.getElementById('a_advertise');
        if (typeof(linkElement) != 'undefined' && linkElement != null)
        {
            unwrap(document.getElementById('a_advertise'));
        }

    });
</script>`

var (
	singleAdTemplate = template.Must(template.New("single_ad").Parse(singleAd))
	videoAdTemplate = template.Must(template.New("video_ad").Parse(videoAD))
)
