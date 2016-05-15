//
// JS CART 
// CREATE ONE ONLINE JS CARTBASKET
// COPY 2012 OPSOFT&WLY
// =================================
//  2012-04-19  添加恢复数据功能
// =================================
// HTML Panel:
//
// <div id="plcart">
//	<h3>我的购物车</h3>
//	<div class="cart"></div>
// <div>
// =================================
//

/*
 function cartItem(){
 this.id = 0;
 this.salePrice = 0;
 this.price = 0;
 this.num = 0;
 this.title = '';
 this.image = '';
 }

 function cartData(){
 this.total = 0;
 this.fee = 0;
 this.totalNum = 0;
 this.isBought = 0;
 this.items = new Array();
 }*/

function shoppingCart() {
    this.key = null;
    this.api = '/cart_api';
    this.cp = null;  //cart panel
    this.tp = null;  //total panel
    this.totalFee = 0; 					    // 总金额
    this.totalNum = 0;                      // 总件数
    this.data = '';                         // 数据字符串
    this.cookieManaged = false;             // 自动管理cookie
    this.defaultCartHtml = '';             // 默认的购物车HTML

    // this.cartData = new cartData();

    this.getByClass = function (ele, cls) {
        if (ele.getElementsByClassName) {
            return ele.getElementsByClassName(cls);
        }
        return jr.dom.getsByClass(ele, cls);
    };

    this.addQua = function (goodsId, num) {
        //添加数量
        var goodsEle = this.getGoodsEle(goodsId);
        var e = goodsEle.getElementsByTagName('INPUT')[0];
        num = parseInt(e.value) + (num || 1);
        e.value = num;
        this.onQuaChanged(goodsId, num, goodsEle);
    };

    this.subQua = function (goodsId) {
        //减少数量
        var goodsEle = this.getGoodsEle(goodsId);
        var e = goodsEle.getElementsByTagName('INPUT')[0];
        if (parseInt(e.value) > 1) {
            var num = parseInt(e.value) - 1;
            e.value = num;
            this.onQuaChanged(goodsId, num, goodsEle);
        }
    };

    this.getGoodsEle = function (goodsId) {
        return document.getElementById('cr' + goodsId);
    };

    this.getNum = function (goodsId) {
        //获得数量输入框
        return parseInt(this.getGoodsEle(goodsId).getElementsByTagName('INPUT')[0].value);
    };

    this.removeAll = function (goodsId) {
        var ele = this.getGoodsEle(goodsId);
        this.cp.removeChild(ele);
        if (this.cp.innerHTML.trim() == '') {//恢复默认的内容
            this.cp.innerHTML = this.defaultCartHtml;
            this.tp.style.display = 'none';
        }
        //触发事件
        this.totalMath();
    };
}
// 重置购物车KEY
shoppingCart.prototype.renewKey = function (newKey) {
    if (this.key == null) {
        this.key = jr.cookie.read('_cart');
    }
    if (newKey) {
        this.key = newKey;
        if (this.cookieManaged) {
            jr.cookie.write('_cart', this.key);
        }
    }
    return this.key;
};

shoppingCart.prototype.xhr = function (data, call) {
    jr.xhr.jsonPost(this.api, data, function (obj) {
        if (call)call(obj);  // 回调处理购物车项
    });
};

shoppingCart.prototype.init = function (panel_id, callback) {
    this.loadCart((function (t) {
        return function (cart) {
            t.initLayout(panel_id, false);
            t.retrieval(cart, callback);
        };
    })(this));
};


shoppingCart.prototype.notify = function (msg) {
    alert(msg);
};

shoppingCart.prototype.loadCart = function (call) {
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
};

shoppingCart.prototype.initLayout = function (ele) {
    var ns = this.getByClass(ele, 'cart');
    if (ns.length > 0) {
        this.cp = ns[0];
        this.defaultCartHtml = this.cp.innerHTML;  //获取购物车区间的HTML
    } else {
        this.cp = document.createElement('DIV');
        this.cp.className = 'cart';
        ele.appendChild(this.cp);
    }

    ns = this.getByClass(ele, 'total');
    if (ns.length > 0) {
        this.tp = ns[0];
    } else {
        this.tp = document.createElement('DIV');
        this.tp.className = 'total';
        ele.appendChild(this.tp);
        this.tp.innerHTML = '共<span class="cart_tq">0</span>件，总价：￥<span class="cart_fee">0</span>元';
    }
};

// 添加项
shoppingCart.prototype.addItem = function (args) {
    //获取数量
    args.num = args.num || 1;
    if (document.getElementById('cr' + args.id) != null) {
        this.addQua(args.id, args.num);
    } else {
        var item = document.createElement('DIV');
        item.className = 'cart-item';
        item.id = 'cr' + args.id;
        item.setAttribute('item-id', args.id);
        item.setAttribute('item-price', args.price);
        item.setAttribute('item-name', args.name);

        item.appendNode = function (html, className) {
            var td = document.createElement("SPAN");
            if (className) {
                td.className = className;
            }
            td.innerHTML = html;
            item.appendChild(td);
            return item;
        };


        //添加
        var nameHtml = '<p class="name">' + args.name + '</p>';
        if (args.image) {
            nameHtml = '<img src="' + args.image + '" class="image"/>' + nameHtml
        }
        nameHtml += jr.template('<span class="qp"><a class="sub_btn" href="javascript:;" onclick="return' +
            ' cart.remove({goodsId},1)">-</a>' +
            '<input class="cart_q" value="{num}" type="text" onlur="return cart.setNum({goodsId});"/>' +
            '<a class="plus_btn" href="javascript:;" onclick="return cart.add({goodsId},1)">+</a></span>',
            {
                goodsId: args.id,
                num: args.num,
            });
        var goodsFee = this.fmtAmount(args.price * args.num);
        item.appendNode(nameHtml, "goods-title").appendNode('￥' + args.price, 'goods-price center')
            .appendNode('￥<span>' + goodsFee + '</span>', 'goods-fee center')
            .appendNode('<a class="remove_btn" href="javascript:void(0)" onclick="cart.remove(\''
                + args.id + '\')" class="cart_remove">x</a>', 'goods-del center');

        this.cp.appendChild(item);

        //触发事件
        this.totalMath();
    }

};

shoppingCart.prototype.fmtAmount= function(amount){
    return amount.toFixed(2).replace(/([^\.]+)(\.|(\.[1-9]))0*$/ig,'$1$3');
};

shoppingCart.prototype.onQuaChanged = function (goodsId, num, goodsEle) {
    this.totalMath();
    if (goodsEle) {
        var price = parseFloat(goodsEle.getAttribute("item-price"));
        this.getByClass(goodsEle, 'goods-fee')[0].innerHTML = this.fmtAmount(price * num)
    }
};

shoppingCart.prototype.totalMath = function () {
    var items = this.getByClass(this.cp, 'cart-item');
    this.totalFee = 0;
    this.totalNum = 0;

    for (var i = 0; i < items.length; i++) {
        //计算数量
        var _q = parseInt(items[i].getElementsByTagName('INPUT')[0].value);
        this.totalNum += _q;
        //计算金额
        var _f = parseFloat(items[i].getAttribute('item-price'));
        this.totalFee += _f * _q;
    }

    var tqs = this.getByClass(this.tp,'cart_tq'),
        tfs = this.getByClass(this.tp,'cart_fee');

    for (var i = 0; i < tqs.length; i++) {
        tqs[i].innerHTML = this.totalNum;
    }
    for (var i = 0; i < tfs.length; i++) {
        tfs[i].innerHTML = this.fmtAmount(this.totalFee);
    }
};

//恢复购物车
shoppingCart.prototype.retrieval = function (cart, callback) {
    if (this.cp.style.display != 'block') {
        this.cp.style.display = 'block';
    }
    if (cart && cart.items && cart.items.length > 0) {
        var len = cart.items.length;
        if (len > 0) {
            this.cp.innerHTML = '';
            if (this.tp.style.display != 'block') {
                this.tp.style.display = 'block';
            }
        }
        for (var i = 0; i < len; i++) {
            var item = cart.items[i];
            this.addItem({
                'id': item.id,
                'name': item.name,
                'image': item.image,
                'num': item.num,
                'price': item.sale_price
            });
        }
        this.totalFee = cart.fee;
        this.totalNum = cart.total_num;
    }
    if (callback)callback(cart);
};

// 购物车添加项
shoppingCart.prototype.add = function (goodsId, num, callback, notify) {
    this.xhr({action: 'add', 'cart.key': this.key, id: goodsId, num: num}, (function (t) {
        return function (obj) {
            if (obj) {
                if (obj.item == null) {
                    if (notify) {
                        notify(obj.message || obj.error);
                    } else {
                        t.notify(obj.message || obj.error);
                    }
                } else {
                    t.addItem(obj.item, num);
                    if (callback)callback(obj);
                }
            }
        };
    })(this));
};

// 购物车删除项
shoppingCart.prototype.remove = function (goodsId, num) {
    var totalNum = this.getNum(goodsId);
    if (num == null) {
        num = totalNum;
    }

    if (this.confirmBeforeClean && totalNum - num == 0) {
        var goodsEle = this.getGoodsEle(goodsId);
        if (!confirm('确定要从购物车中删除商品:' + goodsEle.getAttribute('item-name') + '吗？')) {
            return;
        }
    }

    this.xhr({action: 'remove', 'cart.key': this.key, id: goodsId, num: num}, (function (t) {
        return function (obj) {
            if (obj) {
                if (obj.message) {
                    t.notify(obj.message);
                } else {
                    if (num >= totalNum) {
                        t.removeAll(goodsId);
                    } else {
                        t.subQua(goodsId, num);
                    }
                }
            }
        };
    })(this));
};

shoppingCart.prototype.setNum = function (goodsId, num) {
    var totalNum = this.getNum(goodsId);
    if (num == null) {
        num = totalNum;
    }

    this.xhr({action: 'set', 'cart.key': this.key, id: goodsId, num: num}, (function (t) {
        return function (obj) {
            if (obj) {
                if (obj.message) {
                    t.notify(obj.message);
                } else {
                    if (num >= totalNum) {
                        t.removeAll(goodsId);
                    }
                }
            }
        };
    })(this));
};


var cart = new shoppingCart();
cart.api = '/cart_api_v1';

