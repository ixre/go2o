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

var cart = {
    cp: null, 						//cart panel
    total_fee: 0, 					//总金额
    total_num: 0, 					//总件数
    data: '', 						//数据字符串

    init: function (panel_id, usetheme) {
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

        this.retrieval();

    },
    retrieval: function () {
        //恢复购物车
        var _cart = document.cookie.match(/cart=([^;]+)/);
        if (_cart) {
            this.data = _cart[0];
            var items = this.data.match(/([^\*]+)\*([^\*]+)\*([^|]+)\*([^|]+)/g);
            var arr;
            for (var i = 0; i < items.length; i++) {
                arr = unescape(items[i]).replace(/\||cart=/, '').split('*');
                this.add({ 'id': arr[1], 'name': arr[0], 'num': arr[2], 'price': arr[3] });
            }
        }
    },
    add: function (args) {
        //如果未显示则显示
        if (this.panel.style.display != 'block') {
            this.panel.style.display = 'block';
        }

        if (document.getElementById('cr' + args.id) != null) {
            this.addqua(args.id);
        } else {

            var tr = document.createElement('tr');
            tr.id = 'cr' + args.id;
            tr.setAttribute('itemid', args.id);
            tr.setAttribute('price', args.price);
            tr.setAttribute('itemname', args.name);

            tr.appTd = function (html, className) {
                var td = document.createElement("TD");
                if (className) { td.className = className; }
                td.innerHTML = html;
                tr.appendChild(td);
                return tr;
            };

            //获取数量
            args.num = args.num || 1;

            //添加
            tr.appTd(args.name).appTd('￥' + args.price).appTd('<a href="javascript:;" onclick="return cart.subqua(\''
            + args.id + '\')">-</a><input class="cart_q" value="' + args.num + '" type="text"/><a href="javascript:;" onclick="return cart.addqua(\''
            + args.id + '\')">+</a>', 'cart_qpanel').appTd('<a href="javascript:cart.remove(\'' + args.id + '\')" class="cart_remove">x</a>', 'center');
      
            this.cp.getElementsByTagName('tbody')[0].appendChild(tr);

            //触发事件
            this.totalMath();
        }

    },
    remove: function (id) {
        this.cp.getElementsByTagName('tbody')[0].removeChild(document.getElementById('cr' + id));
        //触发事件
        this.totalMath();
    },

    getQIPT: function (id) {
        //获得数量输入框
        return document.getElementById('cr' + id).getElementsByTagName('input')[0];
    },
    addqua: function (id) {
        //添加数量
        var e = this.getQIPT(id);
        e.value = parseInt(e.value) + 1;
        this.onQuaChanged();
    },
    subqua: function (id) {
        //减少数量
        var e = this.getQIPT(id);
        if (parseInt(e.value) > 1) {
            e.value = parseInt(e.value) - 1;
            this.onQuaChanged();
        }
    },

    onQuaChanged: function () {
        this.totalMath();
    },

    totalMath: function () {
        var trs = this.cp.getElementsByTagName('tr');

        this.total_fee = 0;
        this.totalNum = 0;
        this.data = '';

        for (var i = 0; i < trs.length; i++) {
            var t = trs[i];
            if (t.id.indexOf('cr') != -1) {
                //计算数量
                var _q = parseInt(t.getElementsByTagName('INPUT')[0].value);
                this.totalNum += _q;
                //计算金额
                var _f = parseFloat(t.getAttribute('price'));
                this.total_fee += _f * _q;
                //保存数据
                this.data += escape(t.getAttribute('itemname')) + '*' + t.getAttribute('itemid') + '*' + _q + '*' + _f + '|';
            }
        }

        var tqs = document.getElementsByClassName('cart_tq'),
			tfs = document.getElementsByClassName('cart_fee');

        for (var i = 0; i < tqs.length; i++) {
            tqs[i].innerHTML = this.totalNum;
        } for (var i = 0; i < tfs.length; i++) {
            tfs[i].innerHTML = this.total_fee.toFixed(2);
        }

        this.data += this.total_fee;
        this.async(); 										//同步数据
    },
    async: function () {
        var exp = new Date();
        exp.setTime(exp.getTime() + 3600000000);
        document.cookie = 'cart=' + this.data + ';expire=' + exp.toGMTString() + ';path=/';
    }
};

function checkCart(){
	 if(document.getElementsByClassName('cart_tq')[0].innerHTML.indexOf('空')!=-1){
	 	alert('您的参盒还是空的,请先订餐!');
	 }else{
	 	location.href='/order';
	 }
}


