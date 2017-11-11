package routes

import (
	"encoding/base64"
	"clickyab.com/gad/middlewares"
	"text/template"

	"github.com/clickyab/services/assert"
	"net/url"

	"math/rand"

	"clickyab.com/gad/mr"

	"time"

	"gopkg.in/labstack/echo.v3"
)

const showTemplate = `;(function (name, context, definition) {
  if (typeof module !== 'undefined' && module.exports) { module.exports = definition(); }
  else if (typeof define === 'function' && define.amd) { define(definition); }
  else { context[name] = definition(); }
})('Fingerprint', this, function () {
  'use strict';

  var Fingerprint = function (options) {
    var nativeForEach, nativeMap;
    nativeForEach = Array.prototype.forEach;
    nativeMap = Array.prototype.map;

    this.each = function (obj, iterator, context) {
      if (obj === null) {
        return;
      }
      if (nativeForEach && obj.forEach === nativeForEach) {
        obj.forEach(iterator, context);
      } else if (obj.length === +obj.length) {
        for (var i = 0, l = obj.length; i < l; i++) {
          if (iterator.call(context, obj[i], i, obj) === {}) return;
        }
      } else {
        for (var key in obj) {
          if (obj.hasOwnProperty(key)) {
            if (iterator.call(context, obj[key], key, obj) === {}) return;
          }
        }
      }
    };

    this.map = function(obj, iterator, context) {
      var results = [];
      if (obj == null) return results;
      if (nativeMap && obj.map === nativeMap) return obj.map(iterator, context);
      this.each(obj, function(value, index, list) {
        results[results.length] = iterator.call(context, value, index, list);
      });
      return results;
    };

    if (typeof options == 'object'){
      this.hasher = options.hasher;
      this.screen_resolution = options.screen_resolution;
      this.screen_orientation = options.screen_orientation;
      this.canvas = options.canvas;
      this.ie_activex = options.ie_activex;
    } else if(typeof options == 'function'){
      this.hasher = options;
    }
  };

  Fingerprint.prototype = {
    get: function(){
      var keys = [];
      keys.push(navigator.userAgent);
      keys.push(navigator.language);
      keys.push(screen.colorDepth);
      if (this.screen_resolution) {
        var resolution = this.getScreenResolution();
        if (typeof resolution !== 'undefined'){ // headless browsers, such as phantomjs
          keys.push(resolution.join('x'));
        }
      }
      keys.push(new Date().getTimezoneOffset());
      keys.push(this.hasSessionStorage());
      keys.push(this.hasLocalStorage());
      keys.push(!!window.indexedDB);
      //body might not be defined at this point or removed programmatically
      if(document.body){
        keys.push(typeof(document.body.addBehavior));
      } else {
        keys.push(typeof undefined);
      }
      keys.push(typeof(window.openDatabase));
      keys.push(navigator.cpuClass);
      keys.push(navigator.platform);
      keys.push(navigator.doNotTrack);
      keys.push(this.getPluginsString());
      if(this.canvas && this.isCanvasSupported()){
        keys.push(this.getCanvasFingerprint());
      }
      if(this.hasher){
        return this.hasher(keys.join('###'), 31);
      } else {
        return this.murmurhash3_32_gc(keys.join('###'), 31);
      }
    },

    murmurhash3_32_gc: function(key, seed) {
      var remainder, bytes, h1, h1b, c1, c2, k1, i;

      remainder = key.length & 3; // key.length % 4
      bytes = key.length - remainder;
      h1 = seed;
      c1 = 0xcc9e2d51;
      c2 = 0x1b873593;
      i = 0;

      while (i < bytes) {
          k1 =
            ((key.charCodeAt(i) & 0xff)) |
            ((key.charCodeAt(++i) & 0xff) << 8) |
            ((key.charCodeAt(++i) & 0xff) << 16) |
            ((key.charCodeAt(++i) & 0xff) << 24);
        ++i;

        k1 = ((((k1 & 0xffff) * c1) + ((((k1 >>> 16) * c1) & 0xffff) << 16))) & 0xffffffff;
        k1 = (k1 << 15) | (k1 >>> 17);
        k1 = ((((k1 & 0xffff) * c2) + ((((k1 >>> 16) * c2) & 0xffff) << 16))) & 0xffffffff;

        h1 ^= k1;
            h1 = (h1 << 13) | (h1 >>> 19);
        h1b = ((((h1 & 0xffff) * 5) + ((((h1 >>> 16) * 5) & 0xffff) << 16))) & 0xffffffff;
        h1 = (((h1b & 0xffff) + 0x6b64) + ((((h1b >>> 16) + 0xe654) & 0xffff) << 16));
      }

      k1 = 0;

      switch (remainder) {
        case 3: k1 ^= (key.charCodeAt(i + 2) & 0xff) << 16;
        case 2: k1 ^= (key.charCodeAt(i + 1) & 0xff) << 8;
        case 1: k1 ^= (key.charCodeAt(i) & 0xff);

        k1 = (((k1 & 0xffff) * c1) + ((((k1 >>> 16) * c1) & 0xffff) << 16)) & 0xffffffff;
        k1 = (k1 << 15) | (k1 >>> 17);
        k1 = (((k1 & 0xffff) * c2) + ((((k1 >>> 16) * c2) & 0xffff) << 16)) & 0xffffffff;
        h1 ^= k1;
      }

      h1 ^= key.length;

      h1 ^= h1 >>> 16;
      h1 = (((h1 & 0xffff) * 0x85ebca6b) + ((((h1 >>> 16) * 0x85ebca6b) & 0xffff) << 16)) & 0xffffffff;
      h1 ^= h1 >>> 13;
      h1 = ((((h1 & 0xffff) * 0xc2b2ae35) + ((((h1 >>> 16) * 0xc2b2ae35) & 0xffff) << 16))) & 0xffffffff;
      h1 ^= h1 >>> 16;

      return h1 >>> 0;
    },

    // https://bugzilla.mozilla.org/show_bug.cgi?id=781447
    hasLocalStorage: function () {
      try{
        return !!window.localStorage;
      } catch(e) {
        return true; // SecurityError when referencing it means it exists
      }
    },

    hasSessionStorage: function () {
      try{
        return !!window.sessionStorage;
      } catch(e) {
        return true; // SecurityError when referencing it means it exists
      }
    },

    isCanvasSupported: function () {
      var elem = document.createElement('canvas');
      return !!(elem.getContext && elem.getContext('2d'));
    },

    isIE: function () {
      if(navigator.appName === 'Microsoft Internet Explorer') {
        return true;
      } else if(navigator.appName === 'Netscape' && /Trident/.test(navigator.userAgent)){// IE 11
        return true;
      }
      return false;
    },

    getPluginsString: function () {
      if(this.isIE() && this.ie_activex){
        return this.getIEPluginsString();
      } else {
        return this.getRegularPluginsString();
      }
    },

    getRegularPluginsString: function () {
      return this.map(navigator.plugins, function (p) {
        var mimeTypes = this.map(p, function(mt){
          return [mt.type, mt.suffixes].join('~');
        }).join(',');
        return [p.name, p.description, mimeTypes].join('::');
      }, this).join(';');
    },

    getIEPluginsString: function () {
      if(window.ActiveXObject){
        var names = ['ShockwaveFlash.ShockwaveFlash',//flash plugin
          'AcroPDF.PDF', // Adobe PDF reader 7+
          'PDF.PdfCtrl', // Adobe PDF reader 6 and earlier, brrr
          'QuickTime.QuickTime', // QuickTime
          // 5 versions of real players
          'rmocx.RealPlayer G2 Control',
          'rmocx.RealPlayer G2 Control.1',
          'RealPlayer.RealPlayer(tm) ActiveX Control (32-bit)',
          'RealVideo.RealVideo(tm) ActiveX Control (32-bit)',
          'RealPlayer',
          'SWCtl.SWCtl', // ShockWave player
          'WMPlayer.OCX', // Windows media player
          'AgControl.AgControl', // Silverlight
          'Skype.Detection'];

        // starting to detect plugins in IE
        return this.map(names, function(name){
          try{
            new ActiveXObject(name);
            return name;
          } catch(e){
            return null;
          }
        }).join(';');
      } else {
        return ""; // behavior prior version 0.5.0, not breaking backwards compat.
      }
    },

    getScreenResolution: function () {
      var resolution;
       if(this.screen_orientation){
         resolution = (screen.height > screen.width) ? [screen.height, screen.width] : [screen.width, screen.height];
       }else{
         resolution = [screen.height, screen.width];
       }
       return resolution;
    },

    getCanvasFingerprint: function () {
      var canvas = document.createElement('canvas');
      var ctx = canvas.getContext('2d');
      var txt = 'clickyab';
      ctx.textBaseline = "top";
      ctx.font = "14px 'Arial'";
      ctx.textBaseline = "alphabetic";
      ctx.fillStyle = "#f60";
      ctx.fillRect(125,1,62,20);
      ctx.fillStyle = "#069";
      ctx.fillText(txt, 2, 15);
      ctx.fillStyle = "rgba(102, 204, 0, 0.7)";
      ctx.fillText(txt, 4, 17);
      return canvas.toDataURL();
    }
  };


  return Fingerprint;

});
function effect() {
    var style = document.createElement("style");
    style.innerHTML("" +
        "    #inScreen-close {" +
        "            background: #fff url('//static.clickyab.com/img/close.png') no-repeat center center;" +
        "            cursor: pointer;" +
        "            display: block;" +
        "            filter: alpha(opacity=60);" +
        "            -moz-opacity: 0.6;" +
        "            -webkit-opacity: 0.6;" +
        "            -ms-filter: alpha(opacity=60);" +
        "            opacity: 0.6;" +
        "            float: right;" +
        "            height: 20px;" +
        "            width: 20px;" +
        "            cursor: pointer;" +
        "        }" +
        "    #inScreen-close:hover {" +
        "            filter: alpha(opacity=100);" +
        "            -moz-opacity: 1.0;" +
        "            -webkit-opacity: 1.0;" +
        "            -ms-filter: alpha(opacity=100);" +
        "            opacity: 1.0;" +
        "        }" +
        "    #inScreen-container {" +
        "            left: 50%;" +
        "            position: fixed;" +
        "            transform: translate(-50% , 100%);" +
        "            -webkit-transform: translate(-50% , 100%);" +
        "            -moz-transform: translate(-50% , 100%);" +
        "            -ms-transform: translate(-50% , 100%);" +
        "            -o-transform: translate(-50% , 100%);" +
        "            transition: transform 0.6s ease 0s;" +
        "            bottom: 0;" +
        "        }" +
        "    #inScreen-container {" +
        "            display: block;" +
        "            margin: 0;" +
        "            transition: transform 0.6s ease 0.1s;" +
        "            -webkit-transition: transform 0.6s ease 0.1s;" +
        "            -moz-transition: transform 0.6s ease 0.1s;" +
        "            -ms-transition: transform 0.6s ease 0.1s;" +
        "            -o-transition: transform 0.6s ease 0.1s;" +
        "            width: 100%;" +
        "        }" +
        "    #inScreen-container.close {" +
        "            transform: translate(-50%,100%) !important;" +
        "            -webkit-transform: translate(-50%,100%) !important;" +
        "            -moz-transform: translate(-50%,100%) !important;" +
        "            -ms-transform: translate(-50%,100%) !important;" +
        "            -o-transform: translate(-50%,100%) !important;" +
        "        }" +
        "");

    var header = document.getElementsByTagName("header");
    header.append(style);

    var InScreen = (function () {
        "use strict";
        /*global document: false */
        /*global window: false */
        // create object method
        var method = {},
            settings = {},

            inScreenContainer = document.createElement('div'),
            inScreenHeader = document.createElement('div'),
            inScreenContent = document.createElement('div'),
            inScreenClose = document.createElement('div'),

            centerModal,

            closeModalEvent,

            defaultSettings = {
                width: 'auto',
                height: 'auto',
                lock: false,
                hideClose: false,
                closeAfter: 0,
                openCallback: false,
                closeCallback: false,
                hideOverlay: false
            };

        function insertAfter(referenceNode, newNode) {
            referenceNode.parentNode.insertBefore(newNode, referenceNode.nextSibling);
        }

        // Open the modal
        method.open = function (parameters) {
            settings.width = parameters.width || defaultSettings.width;
            settings.height = parameters.height || defaultSettings.height;
            settings.lock = parameters.lock || defaultSettings.lock;
            settings.hideClose = parameters.hideClose || defaultSettings.hideClose;
            settings.closeAfter = parameters.closeAfter || defaultSettings.closeAfter;
            settings.closeCallback = parameters.closeCallback || defaultSettings.closeCallback;
            settings.openCallback = parameters.openCallback || defaultSettings.openCallback;
            settings.hideOverlay = parameters.hideOverlay || defaultSettings.hideOverlay;

            centerModal = function () {
                method.show({});
            };

            if (parameters.content) {
                var showAdsElement = document.getElementById(xadhtml);
                showAdsElement.innerHTML = inner_loop[adcount];
                inScreenContent.appendChild(showAdsElement);
            } else {
                inScreenContent.innerHTML = '';
            }


            inScreenContainer.style.width = settings.width;
            inScreenContainer.style.height = settings.height;

            method.show({});

            if (settings.lock || settings.hideClose) {
                inScreenClose.style.visibility = 'hidden';
            }

            inScreenContainer.style.visibility = 'visible';

            document.onkeypress = function (e) {
                if (e.keyCode === 27 && settings.lock !== true) {
                    method.close();
                }
            };

            inScreenClose.onclick = function () {
                if (!settings.hideClose) {
                    method.close();
                } else {
                    return false;
                }
            };


            if (window.addEventListener) {
                window.addEventListener('resize', centerModal, false);
            } else if (window.attachEvent) {
                window.attachEvent('onresize', centerModal);
            }

            inScreenHeader.onmousedown = function () {
                return false;
            };
            if (settings.closeAfter > 0) {
                closeModalEvent = window.setTimeout(function () {
                    method.close();
                }, settings.closeAfter * 1000);
            }
            if (settings.openCallback) {
                settings.openCallback();
            }
        };

        method.close = function () {
            inScreenContainer.classList.add("close");

            if (closeModalEvent) {
                window.clearTimeout(closeModalEvent);
            }
            if (settings.closeCallback) {
                settings.closeCallback();
            }

            if (window.removeEventListener) {
                window.removeEventListener('resize', centerModal, false);
            } else if (window.detachEvent) {
                window.detachEvent('onresize', centerModal);
            }
        };
        // Center the modal in the viewport
        method.show = function (parameters) {
            var inScreenContainerVar = getComputedStyle(inScreenContainer).getPropertyValue("transform");
            var showAdEvent = window.setTimeout(function () {
                inScreenContainer.style.transform = "translate(-50%,0)";
                inScreenContainer.style.transform = "msTransform(-50%,0)";
                inScreenContainer.style.transform = "WebkitTransform(-50%,0)";
            }, 2000);

        };
        // Set the id's, append the nested elements, and append the complete modal to the document body
        inScreenContainer.setAttribute('id', 'inScreen-container');
        inScreenContainer.setAttribute('id', 'inScreen-container');
        inScreenHeader.setAttribute('id', 'inScreen-header');
        inScreenContent.setAttribute('id', 'inScreen-content');
        inScreenClose.setAttribute('id', 'inScreen-close');
        inScreenHeader.appendChild(inScreenClose);
        inScreenContainer.appendChild(inScreenHeader);
        inScreenContainer.appendChild(inScreenContent);

        inScreenContainer.style.visibility = 'hidden';

        if (window.addEventListener) {
            window.addEventListener('load', function () {
                document.body.insertBefore(inScreenContainer, document.body.firstChild);
            }, false);
        } else if (window.attachEvent) {
            window.attachEvent('onload', function () {
                document.body.insertBefore(inScreenContainer, document.body.firstChild);
            });
        }

        return method;
    }());

    InScreen.open({
        content: document.getElementById(xadhtml).innerHTML = inner_loop[adcount],
    });

}
var mobad = '{{.Mobad}}';
if(typeof cy_event_page === 'undefined') var cy_event_page = '{{.Rand}}';
var hostofpage = '{{.Host}}';
{{if .NotMobile}}var nativead = 1;{{else}}var nativead = 0;{{end}}
{{if .Mobile}}var ismob = 1;{{else}}var ismob = 0;{{end}}
{{if .Alexa}}{{if .Random}}window.setTimeout(function(){ location.href = "//{{.RURL}}" }, 2000000 );{{end}}{{end}}
var activenative = 0;
function addtoq(pram, val) {
    if (val) {
        clickyab_ad['ad_url'] += "&" + pram + "=" + val;
    }
}
function addtoq2(pram, val) {
    if (val) {
        clickyab_ad['ad_url_m'] += "&" + pram + "=" + val;
    }
}

function encodeuri(b) {
    if (typeof encodeURIComponent == "function") {
        return encodeURIComponent(b);
    } else {
        return escape(b);
    }
}

var a = document, effect;
if (typeof adcount === 'undefined') var adcount = 0;
if (typeof inner_loop === 'undefined') var inner_loop = [];
var fixmob = 0;

adcount++;

function setCookie(cname, cvalue, exdays) {
    var d = new Date();
    d.setTime(d.getTime() + (exdays * 24 * 60 * 60 * 1000));
    var expires = "expires=" + d.toUTCString();
    document.cookie = cname + "=" + cvalue + "; " + expires;
}

function getCookie(cname) {
    var name = cname + "=";
    var ca = document.cookie.split(';');
    for (var i = 0; i < ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}

if (clickyab_ad['native']) {
    clickyab_ad['native'] = true
} else {
    clickyab_ad['native'] = false;
}
if (typeof clickyab_ad['responsive'] === 'undefined') clickyab_ad['responsive'] = 'true';
{{if .Mobile}}
	if(clickyab_ad['width'] > 468 && clickyab_ad['width'] > 60 && clickyab_ad['responsive'] != 'false'){
        clickyab_ad['width'] = 468;
        clickyab_ad['height'] = 60;
    }
{{end}}
if(adcount <= 30){

    a.write('<style> .adhere iframe {  max-width:100%; display: block;margin: 0 auto; }</style><div class="adhere" id="spot_'+adcount+'"></div>');
    clickyab_ad['ad_url'] = "//{{.CHost}}/ads/?a="+clickyab_ad['id'];
    addtoq("width",clickyab_ad['width']);
    addtoq("height",clickyab_ad['height']);
    addtoq("slot",clickyab_ad['slot']);
    addtoq("adtype",clickyab_ad['ad_type']);
    addtoq("domainname",clickyab_ad['domain']);
    addtoq("logo",clickyab_ad['logo']);
    addtoq("tracking",clickyab_ad['tracking']);
    addtoq("eventpage",cy_event_page);
    addtoq("loc",encodeuri(a.location));
    addtoq("ref",encodeuri(a.referrer));
    addtoq("adcount",adcount);
    var fp = new Fingerprint();

    var tid = fp.get();
    if(tid > 0) addtoq("tid",tid);
    var xadhtml = 'spot_' + adcount;
    var ignoreAdBecauseCookie = false;
    var effectString = "";
    effectString = clickyab_ad['effect'];
    clickyab_ad['effect'] = '';
    if (effectString == "interstitial" && getCookie("cy_interstitial")) {
        ignoreAdBecauseCookie = true;
    } else {
        inner_loop[adcount] = '<iframe name="clickyab_frame" width=' + clickyab_ad['width'] + ' height=' + clickyab_ad['height'] + ' frameborder=0 src="' + clickyab_ad['ad_url'] + '" marginwidth="0" marginheight="0" vspace="0" hspace="0" allowtransparency="true" scrolling="no"></iframe>';
    }
    if (effectString == "interstitial" && !getCookie("cy_interstitial")) {
        setCookie("cy_interstitial", true, 0.5);
    }
    if (ignoreAdBecauseCookie != true) {
        document.getElementById(xadhtml).innerHTML = inner_loop[adcount];
        if (typeof effectString !== "undefined" && effectString !== ""  ) {
            if(effectString == "inscreen") {
                effect = "inscreen";
            }
            else if (effectString == "slidedown") {
                effect = "slidedown";
            }
            else if (effectString == "interstitial") {
                effect = "interstitial";
            } else if (effectString == "fade") {
                effect = "fade";
            } else {
                effect = "";
            }
            if (effectString != "") {
				effect();
            }
            clickyab_ad['effect'] = "";
        }
    }
};
	{{if .Mobile}}
	{{if .Mobad}}
		if(adcount <= 1 && window.top == window.self && fixmob == 0){
    clickyab_ad['ad_url_m'] = "//{{.CHost}}/ads/?a="+clickyab_ad['id'];
    addtoq2("width",320);
    addtoq2("height",50);
    addtoq2("slot",clickyab_ad['slot']+"1");
    addtoq2("adtype",clickyab_ad['ad_type']);
    addtoq2("domainname",clickyab_ad['domain']);
    addtoq2("logo",clickyab_ad['logo']);
    addtoq2("tracking",clickyab_ad['tracking']);
    addtoq2("loc",encodeuri(a.location));
    addtoq2("ref",encodeuri(a.referrer));
    a.write('<style>  .adhere iframe {  max-width:100%; display: block;margin: 0 auto; }</style><div class="adhere" style="position: fixed; width: 100%; z-index:99999999; left: 0; bottom: 0px; margin: 0; padding: 0; text-align: center;" class="adhere"><iframe name="clickyab_ads_frame_m" width=320 height=50 frameborder=0 src="' + clickyab_ad['ad_url_m'] +'" marginwidth="0" marginheight="0" vspace="0" hspace="0" allowtransparency="true" scrolling="no">');
            a.write('</iframe></div>');
    };
	{{end}}
	{{end}}
`

type data struct {
	Mobile    bool
	Mobad     int
	Scheme    string
	Host      string
	Alexa     bool
	Rand      int
	Random    bool
	NotMobile bool
	RURL      string
	CHost     string
}

func (tc *selectController) showjs(c echo.Context) error {
	rd := middlewares.MustGetRequestData(c)
	u, err := url.Parse(rd.Referrer)
	assert.Nil(err)
	var wmobad int
	var domain string
	website, err := mr.NewManager().FetchWebsiteByDomain(u.Host, "clickyab")
	if err == nil {
		wmobad = website.WMobad
		domain = website.WDomain.String
	}
	t, err := template.New("show").Parse(showTemplate)
	assert.Nil(err)

	g := data{
		Mobile:    rd.Mobile,
		Mobad:     wmobad,
		Scheme:    rd.Scheme + "://",
		CHost:     c.Request().Host,
		Host:      domain,
		Alexa:     rd.Alexa,
		Rand:      random(999999, 999999999),
		NotMobile: !rd.Mobile,
		Random:    random(1, 2) == 1,
		RURL:      c.Request().Host + "/datacollection?b=" + base64.StdEncoding.EncodeToString([]byte(rd.Referrer)),
	}
	c.Response().Header().Set("Content-Type", "application/javascript")
	c.Response().Header().Set("Pragma", "no-cache")
	c.Response().Header().Set("Cache-Control", "max-age=0, private, no-cache, no-store, must-revalidate")
	return t.Execute(c.Response(), g)
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
