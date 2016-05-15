//自动完成插件

function autoCompletion(ele, url, loadCallback, selectCallback,errCallback, charMinLen) {
    var panel;
    var panelInner;
    this.charMinLen = charMinLen || 1;
    this.lastChar = '';
    this.isOnFocus = false;
    this.timer = null;
    this.url = url;

    if (!ele.nodeName) {
        ele = jr.$(ele);
    }


    var attachBox = function () {
        panel = ele.previousNode;
        if (!panel || panel.nodeName != 'DIV' || panel.className != 'ui-autocompletion-panel') {

            //为父元素设为绝对定位
            if (ele.parentNode.offsetLeft > ele.offsetLeft) {
                ele.parentNode.style.cssText += 'position:relative';
            }


            panel = document.createElement('DIV');
            panel.className = 'ui-autocompletion-panel';
            panel.style.cssText = 'curcor:default;z-index:102;position:absolute;left:' + ele.offsetLeft + 'px;top:'
                + (ele.offsetTop + ele.offsetHeight) + 'px;width:' + ele.offsetWidth + 'px;overflow:hidden;display:none';
            ele.parentNode.insertBefore(panel, ele);


            panelInner = document.createElement('DIV');
            panelInner.className = 'inner';
            panelInner.style.cssText = 'background-color:#fff;';

            panel.appendChild(panelInner);
        } else {
            panelInner = panel.getElementsByTagName('DIV')[0];
        }
    };

    attachBox();

    //筛选框
    var handler = (function (e, p, pi, lc, sc,ec, t) {
        return function (event, isOnfocus) {

            if (isOnfocus) t.isOnFocus = true;

            //排除alt和ctrl键，空格键
            var _event = window.event || event;
            if (_event.altKey || _event.keyCode == 17) return;
            //截取中间的字符,移除前后的字符
            var keyStr = e.value;

            //达不到默认长度
            if (/^\s*$/.test(keyStr) && t.charMinLen != 0) return;
            keyStr = keyStr.replace(/^(\s*)(.+?)(\s*)$/, '$2');

            if (keyStr.length < t.charMinLen) return;

            //document.getElementById('t1').innerHTML = t.lastChar + '/' + keyStr + ((t.lastChar != '' && t.lastChar == keyStr));
            //判断是否和上次一样
            if (t.lastChar != '' && t.lastChar == keyStr) {
                return;
            } else {
                t.lastChar = keyStr;
            }

            jr.xhr.request({
                uri: t.url + (t.url.indexOf('?') == -1 ? "?" : '&') + 'key=' + encodeURIComponent(keyStr),
                params: {}, method: 'GET', data: 'json'
            }, {
                success: function (json) {
                    if (lc) lc(json);
                    if (json.length != 0) {
                        p.style.display = '';

                        var reg = new RegExp(keyStr, 'i');
                        var replaceWord = '<b>' + keyStr + '</b>';
                        var html = '<ul style="margin:0;padding:0;">';
                        for (var i = 0; i < json.length; i++) {
                            html += '<li' + (i == 0 ? ' class="first"' : '')
                                + (json[i].title ? ' title="' + json[i].title + '"' : '')
                                + '>' + json[i].text.replace(reg, replaceWord) + '</li>';

                            //如果和结果相同则默认选择
                            if (json[i].text == keyStr && sc) {
                                if (e.onblur) e.onblur();
                                sc(json[i]);
                            }
                        }
                        html += '</ul>';

                        pi.innerHTML = html;

                        var lis = pi.getElementsByTagName('LI');
                        jr.each(lis, function (i, li) {
                            li.onmouseover = (function (_p, _lis) {
                                return function () {
                                    for (var j = 0; j < _lis.length; j++) {
                                        _lis[j].className = j == 0 ? 'first' : '';
                                    }
                                    this.className = this == lis[0] ? 'first selected' : 'selected';
                                };
                            })(p, lis);

                            li.onclick = (function (j, _p) {
                                return function () {
                                    e.value = j.text;
                                    _p.style.display = 'none';
                                    if (e.onblur) e.onblur();
                                    if (sc) sc(j);
                                };
                            })(json[i], p);
                        });
                    } else {
                        //隐藏输入框
                        p.style.display = 'none';
                    }

                    //temp
                    setTimeout(function () {
                        t.isOnFocus = false;
                    }, 500);
                },
                error:function() {
                    if (ec && ec instanceof Function)ec();
                }
            });
        };

    })(ele, panel, panelInner, loadCallback, selectCallback, errCallback,this);

    var closeHandler = (function (p, t) {
        return function (event) {
            if (!t.isOnFocus) {
                p.style.display = 'none';
            }
        };
    })(panel, this);


    jr.event.add(ele, 'focus', (function (t) {
        return function (event) {
            handler(event, true);
            //t.timer = setInterval(function() {
            //    handler(event);
            //}, 10);
        };
    })(this));

    jr.event.add(ele, 'keyup', handler);

    //jr.event.add(ele, 'blur', (function(t) {
    //    return function() {
    //        clearInterval(t.timer);
    //    };
    //})(this));

    //绑定事件
    jr.event.add(document, 'click', closeHandler);

    return this;
}


jr.extend({
    autoCompletion: function (ele, url, loadCallback, selectCallback,errCallback, charMinLen) {
        return new autoCompletion(ele, url, loadCallback, selectCallback, errCallback,charMinLen);
    }
});