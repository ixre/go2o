define(['jr/core'], function () {
    return {
        msg: 'hello',
        init: shopDomInit,
        parseTmpl: function (htm) {
            return htm.replace(/(\{|\})/ig, '$1$1');
        },
        getByClass: function (cls) {
            var c = null;
            this.eachByClass(cls, function (i,e) {
                c = e;
                return false;
            });
            return c;
        },
        eachByClass: function (cls, call) {
            if (call != null && call instanceof Function) {
                var arr = document.getElementsByClassName(cls);
                var len = arr.length;
                for (var i = 0; i < len; i++) {
                    if (call(i,arr[i]) == false)break;
                }
            }
        },
        closeTipBox: function () {
            var o = jr.$('g_tips_outer');
            if (o != null) {
                o.className += ' hidden';
            }
        },
        showTipBox: function (msg) {
            with (jr.$('g_tips_outer')) {
                className = className.replace(' hidden', '');
            }
            if (msg && msg.length > 0) {
                jr.$('g_tips').innerHTML = msg;
            }
        },
        cart: {
            api: '/cart_api_v1',
            key: null,
            cookieManaged: false,  //仅cookie
            xhr: function (data, call) {
                j6.xhr.jsonPost(this.api, data, function (obj) {
                    if (call)call(obj);  // 回调处理购物车项
                });
            },
            renewKey: function (newKey) {
                if (this.key == null) {
                    this.key = j6.cookie.read('_cart');
                }
                if (newKey) {
                    this.key = newKey;
                    if (this.cookieManaged) {
                        j6.cookie.write('_cart', this.key);
                    }
                }
                return this.key;
            },
            loadCart: function (call) {
                this.renewKey();
                var caller = (function (t) {
                    return function (obj) {
                        if (t.key != obj.key) {
                            t.renewKey(obj.key);
                        }
                        if (call)call(obj);  // 回调处理购物车项
                    };
                })(this);
                this.xhr({action: 'get', 'cart.key': this.key}, caller)
            },
            add: function (goodsId, skuId, num, callback) {
                this.xhr({action: 'add', 'cart.key': this.key, id: goodsId, sku: skuId, num: num},
                    function (obj) {
                        if (obj) {
                            if (obj.item == null) {
                                alert(obj.error);
                            } else {
                                if (callback)callback(obj);
                            }
                        }
                    });
            }
        }
    };
});

function shopDomInit() {
    j6.xhr.filter = null;
    var btnSearch = document.getElementById('btn-goods-search');
    if (btnSearch) {
        btnSearch.onclick = function () {
            var v = j6.dom.getsByClass(document.body, 'search-key')[0].value;
            location.href = '/search?word=' + encodeURIComponent(v);
        };
    }

    //加载购物车
    var cartNumEles = document.getElementsByClassName('top-cart-num');
    if (cartNumEles.length > 0) {
        this.cart.loadCart(function (c) {
            j6.each(cartNumEles, function (i, e) {
                e.getElementsByTagName('I')[0].innerHTML = c['total_num'];
            });
        });
    }

    //初始化按钮
    var btns = j6.dom.getsByClass(document.body,'btn');
    j6.each(btns,function(i,e){
        var _do = e.getAttribute('do');
        if(_do && window.funcs[_do]){
            jr.event.add(e,'click',window.funcs[_do]);
        }
    });
}