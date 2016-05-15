var cartNum;
function initTabCard(m) {
    m.eachByClass('goods-tab', function (i, e) {
        tabCard(e, {event: 'click', frames: m.getByClass('goods-frames')});
    });
}

require([
    'shop/main',
    'lib/util_comm'
    //"lib/mustache",
    //'lib/scroll_load',
    //'lib/parabola'
], function (m) {
    m.init();
    jr.xhr.filter = null;
    cartNum = jr.$('cart-num');
    m.cart.loadCart(function (c) {
        cartNum.innerHTML = c.total_num;
    });

    initTabCard(m);

    jr.$('btn_add_cart').onclick = function () {
        addToCart(m, this, !true);
    };
    jr.$('btn_view_buy').onclick = function () {
        addToCart(m, this, false, function () {
            location.href = '/buy/confirm?t=' + new Date().getSeconds();
        });
        this.onblur();
    };
});

function addToCart(m, t, tip, call) {
    var num = parseInt(jr.$('buyNum').value);
    var goodsId = parseInt(t.getAttribute("goods-id"));
    var skuId = 0;
    m.cart.add(goodsId, skuId, num, function (d) {
        if (tip) {
            showMsg('加入购物车成功', 500);
        }
        cartNum.innerHTML = parseInt(cartNum.innerHTML) + num;
        if (call != null) {
            call(d);
        }
    });
}