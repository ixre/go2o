var pl = $JS.$('order-confirm-panel');
var cashPl = $JS.$('cash-panel');
var couDes = $JS.$('coupon_describe');

function setHtm(id,h){
    if(h.length!=0) {
        $JS.$(id).innerHTML = h;
    }
}
function initEvents(){
    var items = $JS.getElementsByClassName(pl,'item');
    var editLinks = $JS.getElementsByClassName(pl,'edit_link');
    var confirmBtns = $JS.getElementsByClassName(pl,'confirm-button');

    $JS.each(editLinks,function(i,e){
        $JS.event.add(e.getElementsByTagName('A')[0],'click',(function(items,i,e){
            return function(){
                $JS.each(items,function(i2,e2){
                    if(i2 == i){
                        e2.className +=' active_item';
                        //e2.style.display='none';
                        //$JS.animation.toggleHeight(e2,null,15);
                    }else {
                        e2.className = e2.className.replace(' active_item', '');
                    }
                });
            };
        })(items,i,e));
    });


    $JS.each(confirmBtns,function(i,e){
        $JS.event.add(e,'click',(function(item){
            return function(){
                item.className = item.className.replace(' active_item', '');
            };
        })(items[i]));
    });


    $JS.$('cb1').onclick=function(){

    };
    $JS.$('cb2').onclick=function(){
        $JS.json.bind('ctl2',{deliver_opt:window.sctJson.deliver_opt});
    };
    $JS.$('cb2').onclick=function(){
        var data = $JS.json.toObject('ctl2');
        window.sctJson.deliver_opt = parseInt(data.deliver_opt);
        dynamicContent('deliver');
        persistData();
    };

    $JS.$('el3').onclick=function(){
        $JS.json.bind('ctl3',{pay_opt:window.sctJson.pay_opt});
    };

    $JS.$('cb3').onclick=function(){
        var data = $JS.json.toObject('ctl3');
        window.sctJson.pay_opt = parseInt(data.pay_opt);
        dynamicContent('payment');
        persistData();
    };

}

// 显示动态信息
function dynamicContent(t) {
    var showAll = t == "" || t == null;
    // 显示动态支付信息
    if(showAll || t =='payment') {
        var payOpt = parseInt(window.sctJson.pay_opt);
        var payOptEle = $JS.$('payment_opt_name');

        if (payOpt == 1) {
            payOptEle.innerHTML = '现金支付';
        } else if (payOpt == 2) {
            payOptEle.innerHTML = '网银支付(支付宝)';
        }
    }

    // 显示动态配送信息
    if(showAll || t =='deliver') {
        var dlOpt = parseInt(window.sctJson.deliver_opt);
        var dlOptEle = $JS.$('deliver_opt_name');

        if (dlOpt == 1) {
            dlOptEle.innerHTML = '智能配送';
        } else if (dlOpt == 2) {
            dlOptEle.innerHTML = '上门自提';
        }

        setHtm('deliver_rn', sctJson.deliver_person);
        setHtm('deliver_ph', sctJson.deliver_phone);
        setHtm('deliver_addr', sctJson.deliver_address);
    }
}


// 更新数据到服务器端
function persistData() {
    $JS.xhr.jsonPost('/buy/buyingPersist', window.sctJson, function (d) {
        if(d.message){
            alert(d.message);
        }
    });
}

// 选择配送地址
function selectDeliver(){
    $JS.load('deliver-panel','/buy/getDeliverAddrs?sel='+ window.sctJson.deliver_id);
}

// 从表单中恢复数据
function recoverFrom(id) {
    window.sctJson = $JS.json.toObject(id);
    if (window.sctJson.deliver_id <= 0) {
        $JS.$('item1').className += ' active_item';
        selectDeliver();
    }
}

function applyCouponCode(){
    if(this.value==''){
        $JS.validator.removeTip(this);
        couDes.innerHTML='';
        couDes.className += 'hidden';
    }else{
        var t = this;
        $JS.xhr.jsonPost('/buy/apply?type=coupon',{code:this.value},function(json){
            if(json.result == false) {
                $JS.validator.setTip(t,false,null,json.message);
                couDes.className='coupon_desc hidden';
                reloadFee();
            }else{
                $JS.validator.removeTip(t);
                if(json.couponFee) {
                    couDes.className = 'coupon_desc';
                    couDes.innerHTML = '优惠内容：' + json.couponDescribe + '<br /><em>使用该优惠券总节省：￥' +
                    json.couponFee + '元</em>';
                }
                $JS.json.bind(cashPl,{PromFee:json.discountFee,
                    OrderFee:json.payFee});
                $JS.$('final_fee').innerHTML = json.payFee;
            }
        });
    }
}

function submitOrder() {
    if ($JS.validator.validate('form_coupon')) {
        var data = window.sctJson;
        var cp = $JS.json.toObject(form_coupon);
        if (data.deliver_id <= 0) {
            var e = $JS.$('item1');
            e.className += ' active_item';
            return false;
        }
        data.coupon_code = cp.CouponCode;
        $JS.xhr.jsonPost('submit_0', data, function (j) {
            if (j.result) {
                var orderNo = j.data;
                location.replace("order_finish?order_no=" + orderNo)
            } else {
                alert(j.message);
            }
        });
    }
}

function reloadFee(){
    $JS.json.bind(cashPl,{PromFee:prom_fee,OrderFee:order_fee});
    $JS.$('final_fee').innerHTML =order_fee;
}

