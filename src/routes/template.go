package routes

import (
	"config"
	"html/template"
)

const singleAd = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta name="robots" content="nofollow">
    <meta content="text/html; charset=utf-8" http-equiv="Content-Type">
    <style>
html, body, div, span, applet, object, iframe,
h1, h2, h3, h4, h5, h6, p, blockquote, pre,
a, abbr, acronym, address, big, cite, code,
del, dfn, em, font, img, ins, kbd, q, s, samp,
small, strike, strong, sub, sup, tt, var,
dl, dt, dd, ol, ul, li,
fieldset, form, label, legend,
table, caption, tbody, tfoot, thead, tr, th, td {
  margin: 0;
  padding: 0;
  border: 0;
  outline: 0;
  font-weight: inherit;
  font-style: inherit;
  font-size: 100%;
  font-family: inherit;
  vertical-align: baseline;
}
/* remember to define focus styles! */
:focus {
  outline: 0;
}
body {
  line-height: 1;
  color: black;
  background: white;
}
ol, ul {
  list-style: none;
}
/* tables still need 'cellspacing="0"' in the markup */
table {
  border-collapse: separate;
  border-spacing: 0;
}
caption, th, td {
  text-align: left;
  font-weight: normal;
}
blockquote:before, blockquote:after,
q:before, q:after {
  content: "";
}
blockquote, q {
  quotes: "" "";
}
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

const videoAD = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta name="robots" content="nofollow">
    <meta content="text/html; charset=utf-8" http-equiv="Content-Type">
    <script type="text/javascript">function showb(){ document.getElementById("showb").style.display = "block"; }</script>
    <style>
html, body, div, span, applet, object, iframe,
h1, h2, h3, h4, h5, h6, p, blockquote, pre,
a, abbr, acronym, address, big, cite, code,
del, dfn, em, font, img, ins, kbd, q, s, samp,
small, strike, strong, sub, sup, tt, var,
dl, dt, dd, ol, ul, li,
fieldset, form, label, legend,
table, caption, tbody, tfoot, thead, tr, th, td {
  margin: 0;
  padding: 0;
  border: 0;
  outline: 0;
  font-weight: inherit;
  font-style: inherit;
  font-size: 100%;
  font-family: inherit;
  vertical-align: baseline;
}
/* remember to define focus styles! */
:focus {
  outline: 0;
}
body {
  line-height: 1;
  color: black;
  background: white;
}
ol, ul {
  list-style: none;
}
/* tables still need 'cellspacing="0"' in the markup */
table {
  border-collapse: separate;
  border-spacing: 0;
}
caption, th, td {
  text-align: left;
  font-weight: normal;
}
blockquote:before, blockquote:after,
q:before, q:after {
  content: "";
}
blockquote, q {
  quotes: "" "";
}
        body{ margin: 0; padding: 0; text-align: center; }
        .o{ position:absolute; top:0; left:0; border:0; height:{{ .Height }}px; width:{{ .Width }}px; z-index: 99; }
        #showb{ position:absolute; top:0; left:0; border:0; line-height:{{ .Height }}px; height:{{ .Height }}px; width:{{ .Width }}px; z-index: 100; background: rgba(0, 0, 0, 0.60); text-align: center; }
        .tiny2{ height: 18px; width: 19px; position: absolute; bottom: 0px; right: 0; z-index: 100; background: url("//static.clickyab.com/img/clickyab-tiny.png") right top no-repeat; border-top-left-radius:4px; -moz-border-radius-topleft:4px  }
        .tiny2:hover{ width: 66px;  }
        .tiny{ height: 18px; width: 19px; position: absolute; top: 0px; right: 0; z-index: 100; background: url("//static.clickyab.com/img/clickyab-tiny.png") right top no-repeat; border-bottom-left-radius:4px; -moz-border-radius-bottomleft:4px  }
        .tiny:hover{ width: 66px;  }
        .tiny3{ position: absolute; top: 0px; right: 0; z-index: 100; }
        .butl {background: #4474CB;color: #FFF;padding: 10px;text-decoration: none;border: 2px solid #FFFFFF;font-family: tahoma;font-size: 13px;}
        img.adhere {max-width:100%;height:auto;}
        video {background: #232323 none repeat scroll 0 0;}
    </style>
</head>
<body>
	<div id="video_advertise_{{ .Rand }}">
	    <video id="e_video_{{ .Rand }}" width="{{ .Width }}" height="{{ .Height }}" autoplay muted loop>
		<source src="{{ .Src }}" type="video/mp4">
	    </video>
	    <div class="call-to-action-holder" style="position: absolute;height: 62px;background: rgba(0, 0, 0,0.6);bottom: 0;width: {{ .Width }}px;">
        <a href="{{ .Link }}" target="_blank" style="background-color: rgb(226, 62, 67);width: 100px;text-decoration: none;display: block;text-align: center;padding: 5px;margin: 16px auto 0;color: #fff;font-size: 12px;font-family: tahoma;">مشاهده آگهی</a>
    </div>
	</div>
	<script>

	var clickyab_video = document.getElementById('e_video_{{ .Rand }}');

   document.getElementById("video_advertise_{{ .Rand }}").addEventListener("mouseover", function(t) {
        t.target.muted = false;
    });
    document.getElementById("video_advertise_{{ .Rand }}").addEventListener("mouseout", function(t) {
        t.target.muted = true;
    });

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
	    org_html = document.getElementById('video_advertise_{{ .Rand }}').innerHTML;
	    appendHtmlLink = "<a id='a_advertise' target='_blank' href='"+ link +"'>" + org_html + "</a>";
	    var FinalElementHtml = document.getElementById('video_advertise_{{ .Rand }}').innerHTML = appendHtmlLink;
	    document.getElementById('video_advertise_{{ .Rand }}').addEventListener("click", function () {
		var linkElement = document.getElementById('a_advertise');
		if (typeof(linkElement) != 'undefined' && linkElement != null)
		{
		    unwrap(document.getElementById('a_advertise'));
		}

	    });
	</script>
	</div>
	</body></html>`

const allAds = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <!-- The above 3 meta tags *must* come first in the head; any other head content must come *after* these tags -->
    <title>View AD</title>

    <!-- Bootstrap -->
 <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
<style type="text/css">
	html, body {
		height: 100%;
	}
	input[type="checkbox"] {
		margin-top: 0;
	}
	[class*="col-"] {
		height: 100%;
	}
	.inputbar {
		border-right: 1px solid azure;
	}
	img{
		max-height:200px;
	}
</style>
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
      <script src="https://oss.maxcdn.com/html5shiv/3.7.3/html5shiv.min.js"></script>
      <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
    <![endif]-->
  </head>
  <body>
  	<div class="container">
  		<div class="col-md-12 col-sm-12">
			<div class="row" style="margin-top:20px;">
  				<form method="GET">
					<div class="checkbox">
						<label>
							<input type="checkbox" name="v" /> vast
						</label>
					</div>
					<label>
							province
					</label>

					<div class="form-group">
						<select class="form-control" name="p">
						{{ range  $i :=  .Province}}
							<option value="{{$i.ID}}">{{$i.Name}}</option>
						{{end}}
						</select>
					</div>
					<label>
							campaign
					</label>
					<div class="form-group">
						<select class="form-control" name="cam">
						{{ range  $a :=  .Data}}
							<option value="{{$a.CampaignID}}">{{$a.CampaignName.String}}</option>
						{{end}}
						</select>
					</div>
					<label>
							Size
					</label>
					<div class="form-group">
						<select class="form-control" name="s">
						{{ range  $k,$s :=  .Size}}
							<option value="{{$s}}">{{$k}}</option>
						{{end}}
						</select>
					</div>
					<label>
						 website
					</label>
					<div class="form-group">
						<select class="form-control" name="w">
						{{ range  $b :=  .Website}}
							<option value="{{$b.WDomain.String}}">{{$b.WDomain.String}}</option>
						{{end}}
						</select>
					</div>
					<label>
						all active campaign {{.Len}}
					</label>
					<button type="submit" class="btn btn-primary btn-block">Submit</button>
  				</form>
			</div>


			<div class="row">
  			{{range $kk,$d := .Data}}
			{{if div $kk }}
				</div>
				<div class="row">
			{{end}}
					<div class="col-sm-3 col-md-3 ">
						<a href="{{$d.AdURL.String}}">
							<div class="thumbnail">
							<img src="{{$d.AdImg.String}}">
								<div class="caption">
									<h3>campaign : {{$d.CampaignName.String}}</h3>
									<p>size : {{siz $d.AdSize }}</p>
								</div>
							</div>
						</a>
					</div>
  			{{end}}
  			</div>

  		</div>
  	</div>
  </body>
</html>`

var (
	singleAdTemplate = template.Must(template.New("single_ad").Parse(singleAd))
	videoAdTemplate  = template.Must(template.New("video_ad").Parse(videoAD))
	allAdTemplate    = template.Must(template.New("all_ad").Funcs(funcMap).Parse(allAds))
)
var funcMap = template.FuncMap{
	"siz": siz,
	"div": div,
}

func div(g int) bool {
	return g%4 == 0 && g != 0
}

func siz(g int) string {
	return config.GetSizeByNumString(g)
}
