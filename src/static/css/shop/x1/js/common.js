function submitForm(e, f, g) {
    var h = document.getElementById(e);
    var i = h.getAttribute("valid") ? validForm(h) : true;
    if (i) {
        var j = document.getElementById(g);
        var k = getFormContent(h);
        var l = new CreateAjax();
        l.submitForm(h, k, f, j)
    }
};

function normalSubmitForm(e) {
    var f = document.getElementById(e);
    var g = f.getAttribute("valid") ? validForm(f) : true;
    if (g) {
        f.submit()
    }
};

function clickOpen(e, f, g) {
    var h = document.getElementById(e);
    if (h == null) return;
    if (f == null) {
        if (h.style.display == "none") {
            if (g) creatShade(h);
            h.style.display = ""
        } else {
            if (g) showObj(document.getElementById("backShade"), 0, -5);
            h.style.display = "none"
        }
    } else if (f == true) {
        if (g) creatShade(h);
        h.style.display = ""
    } else {
        if (g) showObj(document.getElementById("backShade"), 0, -5);
        h.style.display = "none"
    }
};

function backToTop(e) {
    var f = document.getElementById("backToTop");
    window.a = set;

    function set() {
        var g = document.documentElement.scrollTop == 0 ? document.body.scrollTop : document.documentElement.scrollTop;
        f.style.display = g > e ? 'block' : "none";
        if (!window.XMLHttpRequest) {
            var h = (document.documentElement.clientHeight) - 58 + g;
            f.style.top = h + "px";
            f.style.right = "10px";
            f.style.position = "absolute"
        } else {
            f.style.bottom = "100px"
        }
    }
};

function displayBox(e, _display) {
    oBox = document.getElementById(e);
    if (oBox) {
        oBox.onmouseover = function () {
            displayBox(e, true)
        };
        oBox.onmouseout = function () {
            displayBox(e, false)
        };
        if (_display) {
            oBox.display = 1;
            setTimeout(function () {
                if (1 == oBox.display) {
                    oBox.style.display = ""
                }
            }, 100)
        } else {
            oBox.display = 0;
            setTimeout(function () {
                if (0 == oBox.display) {
                    oBox.style.display = "none"
                }
            }, 300)
        }
    }
};
var b = null;
var c = "";

function inputFocus(e) {
    if (document.all) {
        if (b != null) b.setAttribute("style", c);
        b = e;
        c = e.getAttribute("style");
        e.setAttribute("style", c + ";padding:5px 4px; border:1px solid #F2BC5F;")
    }
};

function clearDefault(e) {
    var f = "",
        g = "";
    if (e.getAttribute("notD") != 'true') {
        f = e.value;
        g = e.getAttribute("style");
        e.value = "";
        e.setAttribute("style", "");
        e.onblur = function () {
            if (e.value == "") {
                e.value = f;
                e.setAttribute("style", g);
                e.setAttribute("notD", false)
            } else {
                e.setAttribute("notD", true);
                return
            }
        }
    }
};

function clickCheck(e) {
    var f = document.getElementById(e);
    f.checked = !f.checked
};

function selectAll(e, f) {
    var g = document.getElementById(e);
    var h = document.getElementsByName(f);
    for (var i = 0; i < h.length; i++) {
        h[i].checked = g.checked
    }
};

function findChecked(e) {
    var f = document.getElementById(e);
    if (f == null) return;
    var g = f.getElementsByTagName("input");
    var h = "";
    for (var i = 0; i < g.length; i++) {
        if (g[i].checked && g[i].value != "") {
            h += g[i].value + "-"
        }
    };
    return h.substring(0, h.length - 1)
};

function submitUrl(e, f, g) {
    var h = document.getElementById(g);
    var i = new CreateAjax();
    i.submitURL(e, f, h)
};

function deleteThings(e, f, g, h) {
    if (h == null) {
        if (confirm("确认删除此项？")) {
            submitUrl(e, f, g)
        }
    } else {
        if (confirm("确认删除多项？")) {
            submitUrl(e + findChecked(h), f, g)
        }
    }
};
var d = null;

function openForm(e, f) {
    var g = document.getElementById(f);
    ajaxRefresh(e, g);
    var h = g.parentNode.parentNode;
    if (d != null && d != h) d.style.display = "none";
    h.style.display = "";
    creatShade(g);
    var i = (document.documentElement.clientWidth) / 2;
    h.style.left = i + "px";
    if (!window.XMLHttpRequest) {
        var j = document.documentElement.scrollTop == 0 ? document.body.scrollTop : document.documentElement.scrollTop;
        var k = (document.documentElement.clientHeight) / 2 + j;
        h.style.top = k + "px";
        h.style.position = "absolute"
    } else h.style.top = "50%";
    d = h
};

function alertForm(e, f) {
    var g = document.createElement("div");
    g.setAttribute("class", "p_s_tips " + e + "_pop");
    g.className = "p_s_tips " + e + "_pop";
    g.setAttribute("id", "pop_alert");
    var h = document.documentElement.scrollTop == 0 ? document.body.scrollTop : document.documentElement.scrollTop;
    var i = (document.documentElement.clientHeight) / 2 + h;
    g.style.top = i + "px";
    g.style.left = "40%";
    g.innerHTML = f;
    var j = document.getElementsByTagName("body")[0];
    setTimeout(function () {
        closeAlert()
    }, 2000);
    var k = document.getElementById("pop_alert");
    if (k != null) j.removeChild(k);
    j.appendChild(g)
};

function closeAlert() {
    var e = document.getElementById("pop_alert");
    var f = document.getElementsByTagName("body")[0];
    f.removeChild(e)
};

function creatShade() {
    var e = (document.all) ? true : false;
    var f = document.createElement("div");
    var g = parseInt(document.documentElement.scrollWidth);
    var h = document.body.clientHeight;
    var i = "display:block;top:0px;left:0px;position:absolute;z-index:198;background:#aaa;width:" + g + "px;height:" + h + "px;";
    f.id = "backShade";
    i += (e) ? "filter:alpha(opacity=0);" : "opacity:0;";
    f.style.cssText = i;
    document.body.appendChild(f);
    showObj(f, 35, 4)
};

function showObj(e, f, g) {
    var h = (document.all) ? true : false;
    if (e == null) return;
    if (h) {
        if (f == 0 && e.filters.alpha.opacity <= 0) {
            if (e.id == "backShade") e.parentNode.removeChild(e);
            else e.style.display = "none";
            return
        };
        e.filters.alpha.opacity += g;
        if (e.filters.alpha.opacity < f || f == 0) {
            setTimeout(function () {
                showObj(e, f, g)
            }, 1)
        }
    } else {
        al = parseFloat(e.style.opacity);
        al += g / 100;
        if (f == 0 && al <= 0) {
            if (e.id == "backShade") e.parentNode.removeChild(e);
            else e.style.display = "none";
            return
        };
        e.style.opacity = al;
        if (al < (f / 100) || f == 0) {
            setTimeout(function () {
                showObj(e, f, g)
            }, 1)
        }
    }
};

function getId(e) {
    return document.getElementById(e)
};

function twoTagChange(e, f, g) {
    var h = getId(e);
    h.setAttribute("class", g ? "mrt cur" : "mrt");
    h.className = g ? "mrt cur" : "mrt";
    clickOpen(f, g)
};

function goBack() {
    var e = document.referrer;
    if (e == "" || e.indexOf("html") > 0 || e.indexOf("www.xiangha.com/validEmail") > 0) self.location.href = "index.php";
    else self.location.href = e
};

function clickTag(e, f, g) {
    var h = e.innerHTML;
    var i = e.onclick;
    var j = e.parentNode;
    var k = j.title;
    j.innerHTML = h;
    j.className = "cur";
    e = j.parentNode;
    var l = e.getElementsByTagName("li");
    for (var m = 0; m < l.length; m++) {
        if (l[m].className.indexOf("cur") >= 0 && l[m] != j) {
            var n = document.createElement("a");
            n.href = "javascript:;";
            n.onclick = i;
            n.innerHTML = l[m].innerHTML;
            l[m].className = "";
            l[m].innerHTML = "";
            l[m].appendChild(n);
            break
        }
    };
    submitUrl('module/' + f + '&' + g + '=' + k, '', g)
};

function tagFloat(e, f, g) {
    var h = document.getElementById(e);
    var i = window.onscroll ||
    function () {};
    var j = parseInt(h.getAttribute("top") ? h.getAttribute("top") : 0);
    var k = k ? ofHight : 0;
    var l = null;
    var m = h.getAttribute("start") ? h.getAttribute("start") : 0;
    if (window.addEventListener) window.addEventListener("scroll", func, false);
    else if (window.attachEvent) window.attachEvent("onscroll", func);
    else window["onscroll"] = func;
    func();

    function func() {
        var n = document.documentElement.scrollTop == 0 ? document.body.scrollTop : document.documentElement.scrollTop;
        var o = parseInt(h.style.top.replace("px", ""));
        if (f == "bottom") j = document.documentElement.clientHeight;
        else if (f == "center") j = (document.documentElement.clientHeight) / 2;
        else if (f == "height") {
            h.style.position = "absolute";
            if (l != null) clearTimeout(l);
            if (n < m) {
                h.style.top = h.getAttribute("theTop") + "px";
                eval(i);
                return
            };
            var p = j + n;

            function addPosition(q) {
                if (Math.abs(p - q) > 10) {
                    q += (p - q) / 10;
                    h.style.top = q + "px";
                    l = setTimeout(function () {
                        addPosition(q)
                    }, 5)
                }
            };
            addPosition(o)
        } else if (f == "fixed") {
            if (n < m) {
                h.style.position = "";
                eval(i);
                return
            };
            if (window.XMLHttpRequest) {
                h.style.position = "fixed";
                h.style.top = j + "px"
            } else {
                h.style.position = "absolute";
                h.style.top = j + n - k + "px"
            }
        };
        if (f == "bottom" || f == "center") {
            h.style.position = "absolute";
            h.style.top = j + n - k + "px"
        };
        eval(i)
    }
};

function getFormContent(j) {
    var k = j.elements;
    var l;
    var m = "";
    for (var n = 0; n < k.length; n++) {
        var o = k[n];
        if (o.type == "text" || o.type == "textarea" || o.type == "hidden") {
            m += encodeURIComponent(o.name) + "=" + encodeURIComponent(o.value) + "&"
        } else if (o.type == "select-one" || o.type == "select-multiple") {
            var p = o.options,
                q, r;
            for (q = 0; q < p.length; ++q) {
                r = p[q];
                if (r.selected) {
                    m += encodeURIComponent(o.name) + "=" + encodeURIComponent(r.value) + "&"
                }
            }
        } else if (o.type == "checkbox" || o.type == "radio") {
            if (o.checked) {
                m += encodeURIComponent(o.name) + "=" + encodeURIComponent(o.value) + "&"
            }
        } else if (o.type == "file") {
            if (o.value != "") {
                m += encodeURIComponent(o.name) + "=" + encodeURIComponent(o.value) + "&"
            }
        } else {
            m += encodeURIComponent(o.name) + "=" + encodeURIComponent(o.value) + "&"
        }
    };
    return m
}
