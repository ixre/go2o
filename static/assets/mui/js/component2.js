if (!window._path) {
    window._path = 'admin';
}
window.sites = [];
window.groupname = null;

if (window.menuData == undefined) {
    //cms.xhr.get(window._path + '?module=ajax&action=appinit', function (x) {
    //    var ip, address, md, username;
    //    eval(x);
    //    window.menuData = md;
    //    j6.json.bind(document, { userName: username });
    //});
    window.menuData = [];
}

if (window.menuHandler == undefined) {
    window.menuHandler = null;
}


var FwMenu = {
    ele: null,
    menuTitles: [],
    getByCls: function (className) {
        return this.ele.getElementsByClassName ? this.ele.getElementsByClassName(className) : document.getElementsByClassName(className, this.ele);
    },
    init: function (data, menuHandler) {
        //获取菜单元素
        this.ele = document.getElementsByClassName('page-left-menu')[0];
        //第一次加载
        var md = data;

        //处理菜单数据
        if (menuHandler && menuHandler instanceof Function) {
            var hdata = menuHandler(data);
            if (hdata != undefined && hdata != null) {
                md = hdata;
            }
        }

        var menuEle = this.ele;

        menuEle.innerHTML = '';
        var title, html, linktext, url;
        for (var i1 = 0; i1 < md.length; i1++) {
            title = md[i1].text;
            html = '';
            for (var i2 = 0; i2 < md[i1].childs.length; i2++) {
                if (md[i1].childs[i2].childs.length > 0) {
                    html += '<div class="group-title" group="' + md[i1].id + '" style="cursor:pointer" title="点击展开操作菜单"><span>' + md[i1].childs[i2].text + '</span></div>';
                    html += '<div class="panel hidden"><ul id="fns_' + i2 + '">';
                    for (var i3 = 0; i3 < md[i1].childs[i2].childs.length; i3++) {
                        linktext = md[i1].childs[i2].childs[i3].text;
                        url = md[i1].childs[i2].childs[i3].uri;
                        // html += (i3 != 0 && i3 % 4 == 0 ? '<div class="clearfix"></div>' : '') +
                        html += '<li' + (i2 == 0 && i3 == 0 ? ' class="current"' : '') + '><a class="fn" style="cursor:pointer;" url="' + url + '"' +
                            //(md[i1].childs[i2].childs.length == 1 ? ' style="margin:0 ' + ((100 - linktext.length * 14) / 2) + 'px"' : '') +
                            '><span class="icon icon_' + i1 + '_' + i2 + '_' + i3 + '"></span>' + linktext + '</a></li>';
                    }
                    html += '</ul></div>';
                }
            }
            menuEle.innerHTML += html;
        }

        //获取所有的标题菜单
        this.menuTitles = this.getByCls('group-title');
        var t = this;
        j6.each(this.menuTitles, function (i, e) {
            var groupName = e.getAttribute('group');
            j6.event.add(e, 'click', (function (_t, _e) {
                return function () {
                    _t.show(_e);
                };
            })(t, e));

            j6.each(e.nextSibling.getElementsByTagName('LI'), function (i2, e2) {
                j6.event.add(e2, 'click', (function (_this, _t, g) {
                    return function () {
                        _t.set(groupName, _this);
                        var a = _this.childNodes[0];
                        if (a.url != '') {
                            FwTab.show(a.innerHTML, a.getAttribute('url'));
                        }
                    };
                })(e2, t, groupName));
            });
        });

    },
    //设置第几组显示
    change: function (id) {
        var menuTitles = this.menuTitles;
        var groupName = id;
        if (!groupName) {
            if (menuTitles.length == 0) {
                return;
            } else {
                groupName = menuTitles[0].getAttribute('group');
            }
        }
        var isFirst = true;
        var selectedLi = null;  //已经选择的功能菜单
        var firstPanel = null;
        var titleGroups = [];
        var _lis;

        j6.each(menuTitles, function (i, e) {
            if (e.getAttribute('group') != groupName) {
                e.className = 'group-title hidden';
                e.nextSibling.className = 'panel hidden';
            } else {
                titleGroups.push(e);
                e.className = 'group-title';

                //第一个panel
                if (firstPanel == null) {
                    firstPanel = e.nextSibling;
                }
            }
        });

        for (var i = 0; i < titleGroups.length; i++) {
            _lis = titleGroups[i].nextSibling.getElementsByTagName('LI');
            for (var j = 0; j < _lis.length; j++) {
                if (_lis[j].className == 'current') {
                    selectedLi = _lis[j];
                    i = titleGroups.length + 1; //使其跳出循环
                    break;
                }
            }
        }

        if (selectedLi != null) {
            selectedLi.parentNode.parentNode.className = 'panel';
        } else if (firstPanel != null) {
            firstPanel.className = 'panel';
        }

    },
    //查看菜单
    show: function (titleDiv) {
        var groupName = titleDiv.getAttribute('group');
        j6.each(this.menuTitles, function (i, e) {
            if (e.getAttribute('group') == groupName) {
                if (e != titleDiv) {
                    e.nextSibling.className = 'panel hidden';
                } else {
                    e.nextSibling.className = 'panel';
                }
            }
        });
    },
    set: function (groupName, ele) {
        j6.each(this.menuTitles, function (i, e) {
            if (e.getAttribute('group') == groupName) {
                j6.each(e.nextSibling.getElementsByTagName('LI'), function (i, e2) {
                    e2.className = ele == e2 ? 'current' : '';
                });
            }
        });
    }
};


/* Tab管理 */
var FwTab = {
    //框架集
    frames: null,
    maskEle: null,
    loadEle: null,
    tabs: null,
    initialize: function () {
        var framebox = j6.$('pageframes');
        this.tabs = j6.$('pagetabs').getElementsByTagName('UL')[0];

        var getByCls = function (cls) {
            return (framebox.getElementsByClassName ? framebox.getElementsByClassName(cls) : document.getElementsByClassName(cls, framebox))[0];
        };
        this.frames = getByCls('frames');
        this.maskEle = getByCls('mask');
        this.loadEle = getByCls('loading');

        var fx = this.frames.offsetWidth,
            fy = this.frames.offsetHeight;

        //mask位置
        if (this.maskEle) {
            this.maskEle.style.width = fx + 'px';
            this.maskEle.style.height = fy + 'px';
        }

        //加载框位置
        if (this.loadEle) {
            var fx1 = j6.screen.offsetWidth(),
                fy1 = j6.screen.offsetHeight(),
                offset = 50;

            this.loadEle.style.left = (Math.floor((fx1 - this.loadEle.offsetWidth) / 2)
                - framebox.parentNode.offsetLeft) + 'px';

            this.loadEle.style.top = (Math.floor((fy1 - this.loadEle.offsetHeight) / 2)
                - framebox.parentNode.parentNode.offsetTop - offset) + 'px';
        }
    },
    pageBeforeLoad: function () {
        this.showLoadBar();
    },
    pageLoad: function () {
        this.hiddenLoadBar();
    },
    showLoadBar: function () {
        this.loadEle.className = 'loading';
        this.maskEle.className = 'mask';
    },
    hiddenLoadBar: function () {
        this.loadEle.className = 'loading hidden';
        this.maskEle.className = 'mask hidden';
    },
    show: function (text, url, closeable) {
        var _tabs = this.tabs.getElementsByTagName('LI');
        var _indent;
        var _exits = false;
        var _cur_indents = url;
        var _li = null;

        cms.each(_tabs, function (i, obj) {
            _indent = obj.getAttribute('indent');
            if (_indent == _cur_indents) {
                _exits = true;
                obj.className = 'current';
                _li = obj;
            }
        });
        if (!_exits) {
            this.pageBeforeLoad();
            //添加框架
            var frameDiv = document.createElement('DIV');
            var frame;
            try {
                //解决ie8下有边框的问题
                frame = document.createElement('<IFRAME frameborder="0">');
            } catch (ex) {
                frame = document.createElement('IFRAME');
            }
            frame.src = url;
            frameDiv.appendChild(frame);
            this.frames.appendChild(frameDiv);

            var _loadCall = (function (t) {
                return function () {
                    t.pageLoad.apply(t);
                };
            })(this);

            frame.frameBorder = '0';
            frame.setAttribute('frameBorder', '0', 0);
            frame.setAttribute('indent', _cur_indents);
            frame.setAttribute('id', 'ifr_' + _cur_indents);
            j6.event.add(frame, 'load', _loadCall);



            //添加选项卡
            _li = document.createElement('LI');
            _li.onmouseout = (function (t) {
                return function () {
                    if (t.className != 'current') t.className = '';
                };
            })(_li);
            _li.onmouseover = (function (t) {
                return function () {
                    if (t.className != 'current') t.className = 'hover';
                };
            })(_li);
            _li.setAttribute('indent', _cur_indents);
            _li.innerHTML = '<span class="txt"><span class="tab-title" onclick="FwTab.set(this)">' + text + '</span>'
                + (closeable == false ? '' : '<span class="tab-close" title="关闭选项卡" onclick="FwTab.close(this);">x</span>')
                + '</span><span class="rgt"></span>';

            this.tabs.appendChild(_li);
        }

        //触发事件,切换IFRAME
        this.set(_li, true);
    },
    set: function (t, isOpen) {

        //如果不是刚打开的tab,则关闭加载提示 
        if (!isOpen) {
            this.hiddenLoadBar();
        }

        var li = t.nodeName != 'LI' ? t.parentNode.parentNode : t;
        var _frames = this.frames.getElementsByTagName('DIV');
        var _lis = this.tabs.getElementsByTagName('LI');
        cms.each(_lis, function (i, obj) {
            if (obj == li) {
                obj.className = 'current';
                _frames[i].className = 'current';
                _frames[i].style.height = '100%';

            } else {
                obj.className = '';
                _frames[i].className = '';
                _frames[i].style.height = '0px';
            }
        });

    },

    //关闭tab,如果不指定关闭按钮，则关闭当前页
    close: function (t) {
        var closeIndex = -1;
        var isActived = false;
        var closeLi = null;

        if (t) {
            //传递指定的tab进行关闭
            if (t.nodeName == 'SPAN') {
                var list = j6.dom.getsByClass(this.tabs, 'tab-close');
                var noCloseBtnLen = this.tabs.getElementsByTagName('LI').length - list.length;
                for (var i = 0; i < list.length; i++) {
                    if (list[i] == t) {
                        closeIndex = i + noCloseBtnLen;
                        closeLi = list[i].parentNode.parentNode;
                        break;
                    }
                }
            }
            //根据标题来关闭
            else if (typeof (t) == 'string') {
                var list = j6.dom.getsByClass(this.tabs, 'tab-title');
                for (var i = 0; i < list.length; i++) {
                    if (t == list[i].innerHTML.replace(/<[^>]+>/g, '')) {
                        closeIndex = i;
                        closeLi = list[i].parentNode.parentNode;
                        break;
                    }
                }
            }
        } else {
            //关闭当前选中的tab
            var _lis = this.tabs.getElementsByTagName('LI');
            for (var i = 0; i < _lis.length; i++) {
                if (_lis[i].className == 'current') {
                    closeIndex = i;
                    closeLi = _lis[i];
                    break;
                }
            }
        }

        //判断是否关闭当前选中的tab
        if (closeLi) {
            isActived = closeLi.className == 'current';
        }

        if (closeIndex > 0) {
            var _lis = this.tabs.getElementsByTagName('LI');
            var _ifrs = this.frames.getElementsByTagName('DIV');

            var ifr = _ifrs[closeIndex].childNodes[0];
            if (ifr.nodeName == 'IFRAME') {
                ifr.src = '';
                ifr = null;
            }

            this.tabs.removeChild(_lis[closeIndex]);
            this.frames.removeChild(_ifrs[closeIndex]);

            //如果关闭当前激活的tab,则显示其他的tab和iframe
            if (isActived) {

                this.hiddenLoadBar();  /* 避免当打开就刷新时仍然加载问题 */

                if (closeIndex >= _lis.length) {
                    closeIndex = _lis.length - 1;
                }
                _lis[closeIndex].className = 'current';
                if (_ifrs[closeIndex]) {
                    _ifrs[closeIndex].className = 'current';
                    _ifrs[closeIndex].style.height = '100%';
                }

            }
        }
    },

    //获取Tab Iframe的框架,如果不包括则返回null
    getWindow: function (t) {
        if (typeof (t) == 'string') {
            var frameIndex = -1;
            var list = j6.dom.getsByClass(this.tabs, 'tab-title');
            for (var i = 0; i < list.length; i++) {
                if (t == list[i].innerHTML.replace(/<[^>]+>/g, '')) {
                    frameIndex = i;
                    break;
                }
            }
            //没有框架或超出数量
            if (frameIndex == -1) return null;
            var frameDivs = this.frames.getElementsByTagName('DIV');
            if (frameIndex >= frameDivs.length) return null;

            //获取Iframe
            var iframes = frameDivs[frameIndex].getElementsByTagName('IFRAME');
            //不包含iframe
            if (iframes.length == 0) return null;
            return iframes[0].contentWindow;
        }
        return null;
    }
};



//加载app
function loadApps() {
    var ele;
    j6.each(document.getElementsByTagName('H2'), function (i, e) {
        if (e.innerHTML == 'APPS') {
            ele = e.parentNode.getElementsByTagName('DIV')[0];
        }
    });
    if (ele) {
        ele.id = 'ribbon-apps';
        j6.load(ele, window._path + '?module=plugin&action=miniapps&ajax=1');
    }
}


window.M = {
    dialog: function (id, title, url, isAjax, width, height, closeCall) {
        newDialog(id, title, url, isAjax, width, height, closeCall);
    },
    alert: function (html, func) {
        cms.tipbox.show(html, false, 100, 2000, 'up');
        if (func) {
            setTimeout(func, 1000);
        }
    },
    msgtip: function (arg, func) {
        cms.tipbox.show(arg.html, false, 100, arg.autoClose ? 2000 : -1, 'up');
        if (func) {
            setTimeout(func, 1000);
        }
    },
    tip: function (msg, func) {
        this.msgtip({ html: msg, autoClose: true }, func);
    },
    loadCatTree: function () {
        _loadCategoryTree();
    },
    clearCache: function (t) {
        window.M.msgtip({ html: '清除中....' });
        cms.xhr.post(window._path, 'module=ajax&action=clearcache', function (x) {
            window.M.msgtip({ html: '缓存清除完成!', autoClose: true });
            cms.xhr.get('/');
        }, function (x) { });
    },
    addFavorite: function () {
        var url = location.href;
        var title = document.title;
        try {
            window.external.addFavorite(url, title);
        }
        catch (e) {
            try {
                window.sidebar.addPanel(title, url, "");
            }
            catch (e) {
                alert("浏览器不支持,请手动添加！");
            }
        }
    },
    setFullScreen: function (event) {
        //var leftWidth = $(e_SD).offsetWidth;
        //if (leftWidth >= window.M.epix.leftWidth) {
        if (!$(e_SD).parentNode.style || $(e_SD).parentNode.style.display != 'none') {
            //全屏
            $(e_HD).style.height = '0px';
            $(e_SD).style.width = '0px';
            $(e_FT).style.height = '0px';
            $(e_HD).style.overflow = 'hidden';
            $(e_SD).parentNode.style.cssText += 'display:none';
        } else {
            //取消全屏
            $(e_HD).style.overflow = '';
            $(e_HD).style.height = (window.M.epix.topHeight - 5) + 'px';
            $(e_SD).style.width = (window.M.epix.leftWidth - 1) + 'px';
            $(e_FT).style.height = (window.M.epix.footHeight - 1) + 'px';
            $(e_SD).parentNode.style.display = '';
        }
        window.onresize();
    }
};



var mainDiv = document.getElementsByClassName('page-main')[0];

function getDivByCls(cls, ele) {
    var e = ele || mainDiv;
    return (e.getElementsByClassName ?
        e.getElementsByClassName(cls) :
        document.getElementsByClassName(cls, e))[0];
}

//左栏div
var leftDiv = getDivByCls('page-main-left');
//右栏div
var rightDiv = getDivByCls('page-main-right');
//框架div
var frameDiv = getDivByCls('page-frames');
//分割div
var splitDiv = getDivByCls('page-main-split');
//框架遮盖层
var frameShadowDiv = getDivByCls('page-frame-shadow');
//图标操作栏
var iconCtrlDiv = getDivByCls('icon-ctrl', document.body);
//用户操作栏
var userDiv = getDivByCls('page-user', document.body);

//重置窗口尺寸
function _resizeWin() {
    var height = document.documentElement.clientHeight;
    var width = document.documentElement.clientWidth;

    mainDiv.style.height = (height - mainDiv.offsetTop) + 'px';
    frameDiv.style.height = (mainDiv.offsetHeight - frameDiv.offsetTop) + 'px';

    //设置右栏的宽度
    rightDiv.style.width = (width - leftDiv.offsetWidth - splitDiv.offsetWidth + 1) + 'px';

    //设置图标操作栏居中显示
    //iconCtrlDiv.className = 'icon-ctrl';
    var iconCtrlLeft = Math.floor((width - iconCtrlDiv.offsetWidth - userDiv.offsetWidth) / 2);
    iconCtrlDiv.style.left = (iconCtrlDiv < 180 ? 180 : iconCtrlLeft) + 'px';
}

j6.event.add(window, 'resize', _resizeWin);

//设置按键
window.onload = function () {
    document.onkeydown = function (event) {
        var e = window.event || event;
        //按键ALT+F11,启用全屏
        if (e.altKey && e.keyCode == 122) {
            window.M.setFullScreen();
            e.returnvalue = false;
            return false;
        } else if (e.keyCode == 122) {
            window.M.setFullScreen();
            e.returnvalue = false;
            return false;
        } else if (!e.ctrlKey && e.keyCode == 116) {
            var ifr = null;
            var ifrs = document.getElementsByTagName('IFRAME');
            for (var i = 0; i < ifrs.length; i++) {
                if (ifrs[i].className == 'current') {
                    ifr = ifrs[i];
                    break;
                }
            }
            if (ifr != null) {
                var src = ifr.src;
                ifr.src = '';
                ifr.src = src;
            }
            e.returnvalue = false;
            return false;
        }
    };

    //FwMenu.init(window.menuData, window.menuHandler);
    //FwMenu.change();

    _resizeWin();
    FwTab.initialize();

    //添加左右栏改变大小功能
    new drag(splitDiv, window).custom(null, 'w-resize', (function (ld, rd, sd, minWidth, maxWidth) {
            return function (event) {
                //显示遮罩层以支持drag
                frameShadowDiv.className = frameShadowDiv.className.replace(' hidden', '');

                var e = event || window.event;
                window.getSelection ? window.getSelection().removeAllRanges() : document.selection.empty();
                if (e.preventDefault) e.preventDefault();                       //这两句便是解决firefox拖动问题的.
                var mx = e.clientX;
                if (mx > minWidth && mx < maxWidth) {
                    sd.style.left = mx + 'px';
                    ld.style.width = mx + 'px';
                    ld.style.marginRight = -mx + 'px';
                    rd.style.marginLeft = (mx + 5) + 'px';
                    _resizeWin();
                }
            };
        })(leftDiv,
            rightDiv,
            splitDiv,
            splitDiv.getAttribute('min'),
            splitDiv.getAttribute('max')),
        function () {
            frameShadowDiv.className += ' hidden';
        });

};

