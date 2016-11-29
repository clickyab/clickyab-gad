function getCookie(name) {
    var re = new RegExp(name + "=([^;]+)"); var value = re.exec(document.cookie);
    return (value != null) ? unescape(value[1]) : null;
}
function getQueryStrings() {
    var assoc  = {};
    var decode = function (s) { return decodeURIComponent(s.replace(/\+/g, " ")); };
    var queryString = location.search.substring(1);
    var keyValues = queryString.split('&');
    for(var i in keyValues) {
        var key = keyValues[i].split('=');
        if (key.length > 1) {
            assoc[decode(key[0])] = decode(key[1]);
        }
    }
    return assoc;
}
var cy_click , cy_imp ,  img_hit , getWholeQuery;
getWholeQuery = getQueryStrings();
cy_click = getWholeQuery["cy_click"];
cy_imp = getWholeQuery["cy_imp"];
if(cy_click !== undefined && cy_imp !== undefined) {
    document.cookie = "cy_click="+ cy_click+"; expires=Fri, 31 Dec 2020 23:59:59 GMT";
    document.cookie = "cy_imp="+ cy_imp +"; expires=Fri, 31 Dec 2020 23:59:59 GMT";
} else {
    var getCookieClick = parseInt(getCookie("cy_click"));
    var getCookieImp = parseInt(getCookie("cy_imp"));
}
function clickyab_callback(action_id) {
    if(action_id === undefined) {
        action_id = ""
    }
    if (cy_click === undefined && cy_imp === undefined) {

        if(!isNaN(getCookieClick) ) {
            img_hit = document.createElement("img");
            img_hit.setAttribute("src", "https://a.clickyab.com/conversion/?click_id=" + getCookieClick + "&imp_id=" + getCookieImp + "&action_id=" + action_id);
            document.body.appendChild(img_hit);
        }
    } else {
        if(!isNaN(cy_click) ) {
            img_hit = document.createElement("img");
            img_hit.setAttribute("src", "https://a.clickyab.com/conversion/?click_id=" + cy_click + "&imp_id=" + cy_imp + "&action_id=" + action_id);
            document.body.appendChild(img_hit);
        }
    }
}