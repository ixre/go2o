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
    this.cp = null;
    this.totalFee = 0; 					    // 总金额
    this.totalNum = 0;                      // 总件数
    this.data = '';                         // 数据字符串
    this.cookieManaged = false;             // 自动管理cookie

    // this.cartData = new cartData();


    this.addQua = function (goodsId) {
        //添加数量
        var goodsEle = this.getGoodsEle(goodsId);
        var e = goodsEle.getElementsByTagName('INPUT')[0];
        var num = parseInt(e.value) + 1;
        e.value = num;
        this.onQuaChanged(goodsId, num,goodsEle);
    };

    this.subQua = function (goodsId) {
        //减少数量
        var goodsEle = this.getGoodsEle(goodsId);
        var e = goodsEle.getElementsByTagName('INPUT')[0];
        if (parseInt(e.value) > 1) {
            var num = parseInt(e.value) - 1;
            e.value = num;
            this.onQuaChanged(goodsId, num,goodsEle);
        }
    };

    this.getGoodsEle= function(goodsId){
        return document.getElementById('cr' + goodsId);
    };

    this.getNum = function (goodsId) {
        //获得数量输入框
        return parseInt(this.getGoodsEle(goodsId).getElementsByTagName('INPUT')[0].value);
    };

    this.removeAll = function (goodsId) {
        var ele = this.getGoodsEle(goodsId);
        this.cp.getElementsByTagName('TBODY')[0].removeChild(ele);
        //触发事件
        this.totalMath();
    };
}
// 重置购物车KEY
shoppingCart.prototype.renewKey = function (newKey) {
    if (this.key == null) {
        this.key = j6.cookie.read('_cart');
    }
    if (newKey) {
        this.key = newKey;
        if(this.cookieManaged) {
            j6.cookie.write('_cart', this.key);
        }
    }
    return this.key;
};

shoppingCart.prototype.xhr = function (data, call) {
    j6.xhr.jsonPost(this.api, data, function (obj) {
        if (call)call(obj);  // 回调处理购物车项
    });
};

shoppingCart.prototype.init = function (panel_id,callback) {
    this.loadCart((function (t) {
        return function (cart) {
            t.initLayout(panel_id, false);
            t.retrieval(cart,callback);
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

shoppingCart.prototype.initLayout = function (panel_id, usetheme) {
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
    this.cp.innerHTML = (usetheme ? css : '') + '<table cellspacing="1"><tbody><tr class="cart_header"><th>名称</th><th>单价</th><th>总价</th><th>删除</th></tr></tbody></table>' +
        '<p class="total">共<span class="cart_tq">0</span>件，总价：￥<span class="cart_fee">0</span>元</p>';
};

// 添加项
shoppingCart.prototype.addItem = function (args) {
    //如果未显示则显示
    if (this.panel.style.display != 'block') {
        this.panel.style.display = 'block';
    }

    if (document.getElementById('cr' + args.id) != null) {
        this.addQua(args.id);
    } else {

        var tr = document.createElement('tr');
        tr.id = 'cr' + args.id;
        tr.setAttribute('item-id', args.id);
        tr.setAttribute('item-price', args.price);
        tr.setAttribute('item-name', args.name);

        tr.appTd = function (html, className) {
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
        var nameTdHtml = '<p class="name">' + args.name+'</p>';
        if(args.image){
            nameTdHtml = '<img src="'+args.image+'" class="image"/>'+ nameTdHtml
        }
        nameTdHtml += '<p class="qp"><a class="sub_btn" href="javascript:;" onclick="return cart.remove(\''
            + args.id + '\',1)">-</a><input class="cart_q" value="' + args.num + '" type="text"/><a class="plus_btn" href="javascript:;" onclick="return cart.add(\''
            + args.id + '\',1)">+</a></p>';

        tr.appTd(nameTdHtml,"goods-title").appTd('￥' + args.price,'goods-price center')
            .appTd('￥<span>' + args.price * args.num+'<span>','goods-fee center')
            .appTd('<a class="remove_btn" href="javascript:void(0)" onclick="cart.remove(\''
            + args.id + '\')" class="cart_remove">x</a>', 'goods-del center');

        this.cp.getElementsByTagName('TBODY')[0].appendChild(tr);

        //触发事件
        this.totalMath();
    }

};


shoppingCart.prototype.onQuaChanged = function (goodsId,num,goodsEle) {
    this.totalMath();
    if(goodsEle){
        var price = parseFloat(goodsEle.getAttribute("item-price"));
        goodsEle.getElementsByTagName('TD')[2].getElementsByTagName('SPAN')[0].innerHTML = price * num;
    }
};

shoppingCart.prototype.totalMath = function () {
    var trs = this.cp.getElementsByTagName('tr');

    this.totalFee = 0;
    this.totalNum = 0;

    for (var i = 0; i < trs.length; i++) {
        var t = trs[i];
        if (t.id.indexOf('cr') != -1) {
            //计算数量
            var _q = parseInt(t.getElementsByTagName('INPUT')[0].value);
            this.totalNum += _q;
            //计算金额
            var _f = parseFloat(t.getAttribute('item-price'));
            this.totalFee += _f * _q;
        }
    }

    var tqs = document.getElementsByClassName('cart_tq'),
        tfs = document.getElementsByClassName('cart_fee');

    for (var i = 0; i < tqs.length; i++) {
        tqs[i].innerHTML = this.totalNum;
    }
    for (var i = 0; i < tfs.length; i++) {
        tfs[i].innerHTML = this.totalFee.toFixed(2);
    }
};

//恢复购物车
shoppingCart.prototype.retrieval = function (cart,callback){
    if (cart == null || cart.items == null)return;
    var item;
    for (var i = 0; i < cart.items.length; i++) {
        item = cart.items[i];
        this.addItem({'id': item.id, 'name': item.name, 'image':item.image,'num': item.num, 'price': item.sale_price});
    }
    this.totalFee = cart.fee;
    this.totalNum = cart.total_num;
    if(callback)callback(cart);
};

// 购物车添加项
shoppingCart.prototype.add = function (goodsId, num,callback,notify) {
    this.xhr({action: 'add', 'cart.key': this.key, id: goodsId, num: num}, (function (t) {
        return function (obj) {
            if (obj) {
                if (obj.item == null) {
                    if(notify){
                        notify(obj.message || obj.error);
                    }else {
                        t.notify(obj.message || obj.error);
                    }
                } else {
                    t.addItem(obj.item,num);
                    if(callback)callback(obj);
                }
            }
        };
    })(this));
};

// 购物车删除项
shoppingCart.prototype.remove = function (goodsId, num) {
    var totalNum = this.getNum(goodsId);
    if(num == null){
        num = totalNum;
    }

    if(totalNum - num == 0) {
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

var cart = new shoppingCart();
cart.api = '/cart_api_v1';

