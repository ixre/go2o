require([
    'shop/main',
    "lib/mustache",
    "lib/shopping_cart"
], function (m, Mustache) {
    m.init();
    j6.xhr.filter = null;
    window.cart.init('cart-panel', function (c) {
        m.closeTipBox();
    });
    jr.$('btn_settle').onclick = function () {
        location.href = '/buy/confirm';
    }
    jr.$('btn_back').onclick = function () {
        if (document.referrer) {
            window.history.go(-1);
        } else {
            location.href = '/';
        }
    }
});