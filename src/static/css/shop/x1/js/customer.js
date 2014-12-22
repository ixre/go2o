function submitScore(shopId, score) {
    var url = "module/common/inShop.php?shopId=" + shopId + "&score=" + score;
    submitUrl(url, '', 'scoreInfo')
};

function copyCode() {
    var txt = "";
    str = "\n欢迎来到【香哈网】 \n链接:" + window.location.href;
    if (window.clipboardData) {
        txt = document.selection.createRange().htmlText;
        window.clipboardData.clearData();
        window.clipboardData.setData("Text", txt + str)
    } else if (navigator.userAgent.indexOf("Opera") != -1) {
        txt = window.location;
        window.location = txt + str
    };
    return false
};

function getShopIllustrate(requestStr) {
    var timer = null;
    if (document.getElementById("theSecond") != null) {
        eval("addCart();tagFloat('left_menu_pop','height');");
    } else timer = setTimeout(function () {
        getShopIllustrate(requestStr)
    }, 500);
    setTimeout(function () {
        if (timer != null) clearTimeout(timer)
    }, 4000)
};

function getShopInfo(requestStr) {
    var timer = null;
    if (document.getElementById("theSecond") != null) {
        if (document.getElementById("secondInfo") != null) {} else {
            eval("addCart();tagFloat('left_menu_pop','height');");
            submitUrl('module/leftShopInfo.php?requestStr=' + requestStr, '', 'leftInfo');
            timer = setTimeout(function () {
                getShopInfo(requestStr)
            }, 500)
        }
    } else timer = setTimeout(function () {
        getShopInfo(requestStr)
    }, 500);
    setTimeout(function () {
        if (timer != null) clearTimeout(timer)
    }, 8000)
};

function showItemSort(obj, num) {
    var sort = document.getElementById("sort");
    var oldSort = sort.getElementsByTagName("span")[0];
    var obj1 = document.createElement("a");
    obj1.setAttribute("title", oldSort.title);
    obj1.setAttribute("href", "javascript:;");
    obj1.setAttribute("name", oldSort.attributes["name"].value);
    obj1.onclick = function () {
        eval("showItemSort(this," + oldSort.attributes["name"].value + ")")
    };
    obj1.innerHTML = oldSort.innerHTML;
    sort.replaceChild(obj1, oldSort);
    var obj2 = document.createElement("span");
    obj2.setAttribute("title", obj.title);
    obj2.setAttribute("name", num);
    obj2.innerHTML = obj.innerHTML;
    sort.replaceChild(obj2, obj);
    var menu = document.getElementById("menu").childNodes;
    var j = 0;
    for (i = 0; i < menu.length; i++) {
        if (menu[i].nodeType == 1) {
            if (j == num || num == -1) {
                menu[i].style.display = "block";
                var menuSort = menu[i].childNodes;
                if (menuSort[3] != null) {
                    menuSort[3].style.display = "block"
                } else menuSort[1].style.display = "block"
            } else {
                menu[i].style.display = "none"
            };
            j++
        }
    }
}
