(function () {
    var clone = {};
    for (var key in clickyab_ad) {
        if (clickyab_ad.hasOwnProperty(key)) {
            clone[key] = clickyab_ad[key];
        }
    }


    if (typeof window._clickyab_stack == 'undefined') {
        window._clickyab_stack = [];
        window.string = [];
    }
    function takeAllClickYabShowJS() {
        var count = 0;
        var scripts = [];
        for (var i = 0; i < document.scripts.length; i++) {
            if (!!document.scripts[i].src.match('show.js')) {
                scripts.push(document.scripts[i]);
                count = count + 1;
            }
        }

        return {
            scripts: scripts,
            count: count
        };
    }

    var dataString = window.string;
    function ArrayToURL(dataString) {
        var pairs = [];
        for (var key in dataString) {
            if (dataString.hasOwnProperty(key)) {
                pairs.push('s[' + dataString[key].slot + ']=' + dataString[key].width + 'x' + dataString[key].height);
            }

        }
        // one time add into url D=[domain name] - I=[id]
        pairs.push('d=' + dataString[0].domain);
        pairs.push('i=' + dataString[0].id);
        return pairs.join('&');
    }
    function renderFarm(objectParameter, config) {
        for (var key in objectParameter) {
            console.log(key);
            console.log(objectParameter);
            var element = document.getElementById("clickyab_iframe_" +key);
            element.setAttribute("src", objectParameter[key]);
        }

    }

    function insertDivAfterClickYabJS(data, showJS) {
        if (typeof pairs == 'undefined') {
            var pairs = [];
        }

        var divElement = document.createElement("Div");
        divElement.id = data.slot;
        var iframeElement = document.createElement("iframe");
        iframeElement.id = "clickyab_iframe_" +data.slot;
        iframeElement.height = data.height;
        iframeElement.width = data.width;
        divElement.appendChild(iframeElement);
        showJS.parentNode.insertBefore(divElement, showJS.nextSibling);

        return divElement;

    }

    function pushToClickYabStack(data, callback) {
        var showJsDomResult = takeAllClickYabShowJS();
        window._clickyab_stack.push(data);
        window.string.push(data);

        if (window._clickyab_stack.length === showJsDomResult.count) {
            window.url= ArrayToURL(dataString);
            callback(window._clickyab_stack, showJsDomResult.scripts);
            var scriptFile = document.createElement('script');
            scriptFile.setAttribute("src","http://192.168.88.207/select?" + window.url);
            document.body.appendChild(scriptFile);

            renderFarm(
                {
                    "2635768282": "hp30download.com",
                    "2635768282": "hp30download.com",
                    "2635768282": "hp30download.com",

                }
            );
        }
    }
    if (typeof window.onLoadCallBack == 'undefined') {
        window.onLoadCallBack = [];
    }

    window.onLoadCallBack.push(function () {
        pushToClickYabStack(clone, function (stack, scripts) {
            var divisions = [];
            for (var s = 0; s < scripts.length; s++) {
                divisions.push(insertDivAfterClickYabJS(stack[s], scripts[s]));
                // cartDiv.id

            }
            // console.log(divisions);
            // do the job with the scripts
            // do the job with the stack
        });
    });

    window.onload = function () {
        for (var i = 0; i < window.onLoadCallBack.length; i++) {
            window.onLoadCallBack[i]();
        }
    };
})();
