var slideTmpl = '';
var specialTmplA = '';
var specialTmplB = '';
var goodsTmpl = '';
var titTmplA = '';
var titTmplB = '';

function show(e) {
    e.className = e.className.replace(' hidden', '');
}
function setSpecialBlock(m, Mustache, cls, tit, tmpl, data) {
    var e = m.getByClass(cls);
    if (data == null)return;
    data.Name = tit;
    if (tmpl == specialTmplB) {
        e.innerHTML = Mustache.render(tmpl, data);
    } else if (tmpl == specialTmplA) {
        if (data.data) {
            var len = data.data.length;
            if (len > 0) {
                var lArrLen = len / 3 == 0 ? 1 : len / 3;
                var data1 = new Array(lArrLen); //左侧的大图
                var data2 = new Array(len - lArrLen); //右侧对应的2个小图

                for (var i = 0; i < lArrLen; i++) {
                    data1[i] = data.data[i * 3];
                }
                for (var i = 0, j = 0; i < len; i++) {
                    if (i % 3 != 0) {
                        data2[j++] = data.data[i];
                    }
                }
                data.data1 = data1;
                data.data2 = data2;
            }
        }
        e.innerHTML = Mustache.render(tmpl, data)
    }
    show(e)
}
function loadSpecialAd(m, Mustache, $) {
    j6.xhr.get('/json/ad?names=mobi_shop_flash|mobi_shop_g1|mobi_shop_g2', function (data) {
        data = j6.toJson(data);
        initSlideFlash(m, Mustache, $, data['mobi_shop_flash']);
        setSpecialBlock(m, Mustache, 'special-block-1',
            '推荐商品', specialTmplA, data['mobi_shop_g1']);
        setSpecialBlock(m, Mustache, 'special-block-2',
            '精挑细选', specialTmplA, data['mobi_shop_g2']);
    });
}
function loadGoods(m, Mustache) {
    j6.xhr.jsonPost('/json/simple_goods', 'params=new-goods*4|hot-sales*4*3', function (data) {
        var e = m.getByClass('goods-block-new');
        e.innerHTML = Mustache.render(titTmplB, {Title: '新品上架'}) +
            Mustache.render(goodsTmpl, {list: data['new-goods']});
        show(e);
        e = m.getByClass('goods-block-hot-sales');
        e.innerHTML = Mustache.render(titTmplB, {Title: '热销商品'}) +
            Mustache.render(goodsTmpl, {list: data['hot-sales']});
        show(e);
    });
}

function loadSaleTagGoods(m, Mustache) {
    j6.xhr.jsonPost('/json/saletag_goods', 'params=fanli*4', function (data) {
        var e = m.getByClass('goods-block-fanli');
        e.innerHTML = Mustache.render(titTmplA, {Title: '返利商品', Url: '/st/fanli'}) +
            Mustache.render(goodsTmpl, {list: data['fanli']});
        show(e);
    });
}

require([
    'shop/main',
    'lib/mustache',
    'jquery',
    'lib/lazysizes.min', //延迟加载
    'jquery.slides'
], function (m, Mustache, $,l) {
    m.init();
    j6.xhr.filter = null;
    initClickEvent(m)
    preParseTmpl(m, Mustache)
    loadSpecialAd(m, Mustache, $);
    loadGoods(m, Mustache);
    loadSaleTagGoods(m, Mustache);
});

function preParseTmpl(m, Mustache) {
    slideTmpl = m.parseTmpl(m.getByClass('template-slider').innerHTML); //获取滚动模板
    goodsTmpl = m.parseTmpl(m.getByClass('template_goods').innerHTML);
    specialTmplA = m.parseTmpl(m.getByClass('template-special-a').innerHTML);
    specialTmplB = m.parseTmpl(m.getByClass('template-special-b').innerHTML);
    titTmplA = m.parseTmpl(m.getByClass('template_tit_a').innerHTML);
    titTmplB = m.parseTmpl(m.getByClass('template_tit_b').innerHTML);
    Mustache.parse(slideTmpl);
    Mustache.parse(goodsTmpl);
    Mustache.parse(specialTmplA);
    Mustache.parse(specialTmplB);
    Mustache.parse(titTmplA);
    Mustache.parse(titTmplB);
}

function initClickEvent(m) {
    m.eachByClass('icon-search', function (i, e) {
        e.onclick = function () {
            location.href = '/search?from=main';
        };
    });
    m.eachByClass('icon-shopping-cart', function (i, e) {
        e.onclick = function () {
            location.href = '/cart?from=main';
        };
    });
    m.eachByClass('icon-user', function (i, e) {
        e.onclick = function () {
            location.href = '/user/jump_uc';
        };
    });
}

function initSlideFlash(m, Mustache, $, data) {
    var htm = Mustache.render(slideTmpl, data);
    m.eachByClass('slide-main', function (i, e) {
        e.innerHTML = htm;
        $('#slides' + data.ad.Id).slidesjs({
            width: 320,
            height: 140,
            navigation: true, play: {
                //active: true,
                effect: "slide",
                interval: 3000,
                auto: true,
                //swap: true,
                pauseOnHover: false,
                restartDelay: 2500
            }
        });
    });

}


function go(t) {
    location.href = t.getAttribute('url');
}

