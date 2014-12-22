var loadingImg = "images/loadingapple.gif";

function CreateAjax() {
    this.xmlHTTP = false;
    this.createXMLHTTP = function () {
        if (window.ActiveXObject) {
            var MSXML = ['MSXML2.XMLHTTP.5.0', 'MSXML2.XMLHTTP.4.0', 'MSXML2.XMLHTTP.3.0', 'MSXML2.XMLHTTP', 'Microsoft.XMLHTTP'];
            for (var n = 0; n < MSXML.length; n++) {
                try {
                    this.xmlHTTP = new ActiveXObject(MSXML[n]);
                    this.xmlHTTP.setRequestHeader('Content-Type', 'text/html; charset=utf-8');
                    break
                } catch (e) {}
            }
        } else if (window.XMLHttpRequest) {
            this.xmlHTTP = new XMLHttpRequest();
            if (this.xmlHTTP.overrideMimeType) {
                this.xmlHTTP.overrideMimeType("text/xml")
            }
        }
    };
    this.createXMLHTTP();
    if (!this.xmlHTTP) alert("您的浏览器不支持XMLHttpRequest，请更换设置或浏览器");
    this.submitForm = function (form, content, pageUrl, view) {
        obj = this;
        this.xmlHTTP.onreadystatechange = function () {
            obj.getAjaxRes(pageUrl, view)
        };
        content += "&" + Math.random();
        if (form.method.toLowerCase() == "post") {
            this.xmlHTTP.open("POST", form.action, true);
            this.xmlHTTP.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
            this.xmlHTTP.send(content)
        } else {
            this.xmlHTTP.open("GET", form.action + content, true);
            this.xmlHTTP.send(null)
        }
    };
    this.submitURL = function (url, pageUrl, view) {
        obj = this;
        this.xmlHTTP.onreadystatechange = function () {
            obj.getAjaxRes(pageUrl, view)
        };
        this.xmlHTTP.open("GET", encodeURI(encodeURI(url)), true);
        this.xmlHTTP.send(null)
    };
    this.getAjaxRes = function (pageUrl, view) {
        if (this.xmlHTTP.readyState == 4) {
            if (this.xmlHTTP.status == 200) {
                var res = this.xmlHTTP.responseText;
                this.xmlHTTP.abort();
                if (res) {
                    if (res == "ok" && view != null) {
                        ajaxRefresh(pageUrl, view)
                    } else if (res == "closePop") {
                        ajaxRefresh(pageUrl, view);
                        oldForm.style.display = "none";
                        showObj(document.getElementById("backShade"), 0, -5)
                    } else if (res == "goto") {
                        if (pageUrl == "") window.location.reload();
                        else self.location.href = pageUrl
                    } else if (res.indexOf("alert") == 0) {
                        res = res.replace("alert", "");
                        alert(res)
                    } else if (res.indexOf("script ") == 0 || res.indexOf("\<script t") == 0) {
                        res = res.replace("\<script type='text/javascript'>", "");
                        res = res.replace("\</script>", "");
                        res = res.replace("script ", "");
                        eval(res)
                    } else if (view != null) {
                        setInnerHTML(view, res)
                    }
                }
            } else {
                return
            }
        }
    }
};

function ajaxRefresh(url, refreshDiv) {
    if (refreshDiv.toString() == "[object HTMLSpanElement]") setInnerHTML(refreshDiv, "<img src='" + loadingImg + "' style='margin-top:5px;' width='18px' height='18px'></img>");
    else refreshDiv.style.background = "url('" + loadingImg + "') center no-repeat";
    var ajax = new CreateAjax();
    var http_request = ajax.xmlHTTP;
    http_request.onreadystatechange = function () {
        if (http_request.readyState == 4) {
            if (http_request.status == 200) {
                ResStr = http_request.responseText;
                http_request.abort();
                if (refreshDiv != null) {
                    refreshDiv.style.background = "";
                    setInnerHTML(refreshDiv, ResStr)
                }
            } else {
                alert("error : " + http_request.status)
            }
        } else {}
    };
    var nc = "&" + Math.random();
    http_request.open("GET", url + nc, true);
    http_request.send(null)
};

function ajaxOpenForm(title, url) {
    var ajax = new CreateAjax();
    var http_request = ajax.xmlHTTP;
    http_request.onreadystatechange = function () {
        if (http_request.readyState == 4) {
            if (http_request.status == 200) {
                var ResStr = http_request.responseText;
                http_request.abort();
                if (ResStr.indexOf("alert") == 0) {
                    ResStr = ResStr.replace("alert", "");
                    alert(ResStr)
                } else {
                    openshow(ResStr, title, 300, 200, 2)
                }
            } else {
                alert("error : " + http_request.status)
            }
        } else {}
    };
    var nc = "&" + Math.random();
    http_request.open("GET", url + nc, true);
    http_request.send(null)
};
var setInnerHTML = function (el, htmlCode) {
        if (el.value != null) {
            el.value = htmlCode;
            return
        };
        var ua = navigator.userAgent.toLowerCase();
        if (ua.indexOf('msie') >= 0 && ua.indexOf('opera') < 0) {
            htmlCode = '<div style="display:none">for IE</div>' + htmlCode;
            htmlCode = htmlCode.replace(/<script([^>]*)>/gi, '<script$1 defer="true">');
            el.innerHTML = htmlCode;
            el.removeChild(el.firstChild)
        } else {
            var el_next = el.nextSibling;
            var el_parent = el.parentNode;
            el_parent.removeChild(el);
            el.innerHTML = htmlCode;
            if (el_next) {
                el_parent.insertBefore(el, el_next)
            } else {
                el_parent.appendChild(el)
            }
        }
    }
