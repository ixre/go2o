var loadNum = 0; //已经加载数量
var isOver = false; //是否已经加载完毕
var isLoad = false; //是否正在加载
var loadSize = 10;  //每次加载数量
var categoryId = 0;
var noRecordTxt = "未找到相关商品";
var loadEle,loadMsgEle = null;
require([
    'shop/main',"lib/mustache",
    'lib/scroll_load',
    'lib/parabola'],function(m,Mustache){
    m.init();
    jr.xhr.filter = null;
    var itemPanel = m.getByClass('item-list');
    var itemTmpl = m.parseTmpl(m.getByClass('template_goods').innerHTML);
    Mustache.parse(itemTmpl);
    loadEle = m.getByClass('load-ele');
    loadMsgEle = m.getByClass('load-msg');
    var arr = location.pathname.substring(1,location.pathname.indexOf(".")).split('-');
    categoryId = parseInt(arr[arr.length-1]);
    ajaxLoad(itemPanel,m,Mustache,itemTmpl); //加载商品数据

    window.onscroll = function () {
        //监听事件内容
        if (getScrollHeight() == getWindowHeight() + getDocumentTop()) {
            ajaxLoad(itemPanel,m,Mustache,itemTmpl);
        }
    };
});

function ajaxLoad(ele,m,Mustache,itemTmpl) {
    if(isOver || isLoad)return;
    m.showTipBox();
   //loadEle.className = 'load-ele';
    //loadMsgEle.innerHTML = '加载中...';
    isLoad = true;
    jr.xhr.jsonPost('/list/GetGoodsJsonBySaleTag',{
            code:_code||'',
            size:loadSize,begin:loadNum,
            sort:jr.request('sort')}, function(d){
        isLoad = false;
        if(d.total == 0){
            isOver = true;
            loadMsgEle.innerHTML='<span style="color:#F00">'+ noRecordTxt +'</span>';
        }else if(d.total == loadNum){
            isOver = true;
            loadEle.className = 'load-ele';
            loadMsgEle.innerHTML='没有了';
            window.onscroll = null;
        }else{
            loadEle.className += ' hidden';
            appendRecord(ele,m,Mustache,itemTmpl,d.rows);
        }

        loadNum += d.rows.length; //累计数量
        m.closeTipBox(); //关闭加载框
    },function(){
        isLoad = false;
        loadEle.className = 'load-ele';
        loadMsgEle.innerHTML='<span style="color:#F00">加载失败</span>';
        m.closeTipBox();
    });
}

function appendRecord(ele,m,Mustache,itemTmpl,rows) {
   ele.innerHTML += Mustache.render(itemTmpl,{list:rows});
    jr.each(ele.getElementsByClassName('join-cart'),function(i,e){
        e.onclick=function() {
            var goodsId = parseInt(e.getAttribute('goodsId'));
            var skuId = parseInt(e.getAttribute('skuId'));
            m.cart.add(goodsId, skuId, 1, function (d) {
                //加载购物车
                var cartNumEles = document.getElementsByClassName('top-cart-num');
                //创建图标
                var parent = e.parentNode;
                var img = parent.getElementsByClassName('goods-img')[0];
                var moveImg = document.createElement("IMG");
                moveImg.style.cssText='width:40px;height:40px;position:absolute;left:10px;top:10px;z-index:200;';
                moveImg.src = img.src;
                parent.appendChild(moveImg);
                funParabola(moveImg,cartNumEles[0],{
                    complete:function(){
                        //删除运动对象
                        parent.removeChild(moveImg);
                        //更新数量
                        if (cartNumEles.length > 0) {
                            jr.each(cartNumEles, function (i, e) {
                                var e1 = e.getElementsByTagName('I')[0];
                                e1.innerHTML = parseInt(e1.innerHTML) + 1;
                            });
                        }
                    }
                }).mark().init();

            });
        };
    });
}

function go(t){
    location.href=t.getAttribute('url');
}