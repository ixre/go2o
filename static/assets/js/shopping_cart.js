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

function ShoppingCart() {
    this.key = null;
    this.api = '/cart_api';
    this.cp = null;
    this.total_fee = 0; 					//总金额
    this.total_num = 0;                     //总件数
    this.data = '';                         //数据字符串

    this.addQua = function (goodsId) {
        //添加数量
        var e = this.getGoodsEle(goodsId).getElementsByTagName('INPUT')[0];
        e.value = parseInt(e.value) + 1;
        this.onQuaChanged();
    };

    this.subQua = function (goodsId) {
        //减少数量
        var e = this.getGoodsEle(goodsId).getElementsByTagName('INPUT')[0];
        if (parseInt(e.value) > 1) {
            e.value = parseInt(e.value) - 1;
            this.onQuaChanged();
        }
    };

    this.getGoodsEle= function(goodsId){
        return document.getElementById('cr' + goodsId);
    };

    this.getNum = function (goodsId) {
        //获得数量输入框
        return parseInt(this.getGoodsEle(goodsId).getElementsByTagName('input')[0].value);
    };

    this.removeAll = function (goodsId) {
        var ele = this.getGoodsEle(goodsId);
        this.cp.getElementsByTagName('tbody')[0].removeChild(ele);
        //触发事件
        this.totalMath();
    };
}
// 重置购物车KEY
ShoppingCart.prototype.renewKey = function (newKey) {
    if (this.key == null) {
        this.key = $JS.cookie.read('_cart');
    }
    if (newKey) {
        this.key = newKey;
        $JS.cookie.write('_cart', this.key);
    }
    return this.key;
};

ShoppingCart.prototype.xhr = function (data, call) {
    $JS.xhr.jsonPost(this.api, data, function (obj) {
        if (call)call(obj);  // 回调处理购物车项
    });
};

ShoppingCart.prototype.init = function (panel_id, usetheme) {
    this.loadCart((function (t) {
        return function (cart) {
            t.initLayout(panel_id, usetheme);
            t.retrieval(cart);
        };
    })(this));
};


ShoppingCart.prototype.notify = function (msg) {
    alert(msg);
};

ShoppingCart.prototype.loadCart = function (call) {
    var t = this;
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

ShoppingCart.prototype.initLayout = function (panel_id, usetheme) {
    this.panel = document.getElementById(panel_id);
    var pnodes = this.panel.childNodes;
    for (var i = 0; i < pnodes.length; i++) {
        if (pnodes[i].className == 'cart') {
            this.cp = pnodes[i];
            break;
        }
    }

    var css = '<style type="text/css">.cart table{background:green;border:solid 1px #f0f0f0;width:100%;}'
        + '.cart table th{background:white;}'
        + '.cart table .center{text-align:center;}'
        + '.carttable.cart td{background:white;}'
        + '.cart table .cart_q{width:15px;text-align:center;}'
        + '</style>';
    this.cp.innerHTML = (usetheme ? css : '') + '<table cellspacing="1"><tbody><tr class="cart_header"><th>名称</th><th>单价</th><th>数量</th><th>删除</th></tr></tbody></table>' +
    '<p class="center">共<span class="cart_tq">0</span>件，总价：￥<span class="cart_fee">0</span>元</p>';
};

// 添加项
ShoppingCart.prototype.addItem = function (args) {
    //如果未显示则显示
    if (this.panel.style.display != 'block') {
        this.panel.style.display = 'block';
    }

    if (document.getElementById('cr' + args.id) != null) {
        this.addQua(args.id);
    } else {

        var tr = document.createElement('tr');
        tr.id = 'cr' + args.id;
        tr.setAttribute('itemid', args.id);
        tr.setAttribute('price', args.price);
        tr.setAttribute('itemname', args.name);

        tr.apptd = function (html, className) {
            var td = document.createElement("TD");
            if (className) {
                td.className = className;
            }
            td.innerHTML = html;
            tr.appendChild(td);
            return tr;
        };

        //获取数量
        args.num = args.num || 1;

        //添加
        tr.apptd(args.name).apptd('￥' + args.price).apptd('<a class="sub_btn" href="javascript:;" onclick="return cart.remove(\''
        + args.id + '\',1)">-</a><input class="cart_q" value="' + args.num + '" type="text"/><a class="plus_btn" href="javascript:;" onclick="return cart.add(\''
        + args.id + '\',1)">+</a>', 'cart_qpanel').apptd('<a class="remove_btn" href="javascript:void(0)" onclick="cart.remove(\'' + args.id + '\')" class="cart_remove">x</a>', 'center');

        this.cp.getElementsByTagName('tbody')[0].appendChild(tr);

        //触发事件
        this.totalMath();
    }

};


ShoppingCart.prototype.onQuaChanged = function () {
    this.totalMath();
};

ShoppingCart.prototype.totalMath = function () {
    var trs = this.cp.getElementsByTagName('tr');

    this.total_fee = 0;
    this.total_num = 0;

    for (var i = 0; i < trs.length; i++) {
        var t = trs[i];
        if (t.id.indexOf('cr') != -1) {
            //计算数量
            var _q = parseInt(t.getElementsByTagName('INPUT')[0].value);
            this.total_num += _q;
            //计算金额
            var _f = parseFloat(t.getAttribute('price'));
            this.total_fee += _f * _q;
         }
    }

    var tqs = document.getElementsByClassName('cart_tq'),
        tfs = document.getElementsByClassName('cart_fee');

    for (var i = 0; i < tqs.length; i++) {
        tqs[i].innerHTML = this.total_num;
    }
    for (var i = 0; i < tfs.length; i++) {
        tfs[i].innerHTML = this.total_fee.toFixed(2);
    }
};

//恢复购物车
ShoppingCart.prototype.retrieval = function (cart) {
    if (cart == null || cart.items == null)return
    var item;
    for (var i = 0; i < cart.items.length; i++) {
        item = cart.items[i];
        this.addItem({'id': item.id, 'name': item.name, 'num': item.num, 'price': item.salePrice});
    }
};

// 购物车添加项
ShoppingCart.prototype.add = function (goodsId, num,callback) {
    this.xhr({action: 'add', 'cart.key': this.key, id: goodsId, num: num}, (function (t) {
        return function (obj) {
            if (obj) {
                if (obj.item == null) {
                    t.notify(obj.message);
                } else {
                    t.addItem(obj.item,num);
                    if(callback)callback(obj);
                }
            }
        };
    })(this));
};

// 购物车删除项
ShoppingCart.prototype.remove = function (goodsId, num) {
    var totalNum = this.getNum(goodsId);
    if(num == null){
        num = totalNum;
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

var cart = new ShoppingCart();
cart.api = '/cart_api_v1';

