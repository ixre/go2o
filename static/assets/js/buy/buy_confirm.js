var pl = j6.$('order-confirm-panel');
var cashPl = j6.$('cash-panel');
var couDes = j6.$('coupon_describe');
var lbBalance = j6.$('lb_balance');
var ckBalance = j6.$('ck_useBalance');
var tbCouponCode = j6.$('coupon_code');
var coupon_fee = 0;

function setHtm(id, h) {
    if (h.length != 0) {
        j6.$(id).innerHTML = h;
    }
}

function initEvents() {
    var items = j6.dom.getsByClass(pl, 'item');
    var editLinks = j6.dom.getsByClass(pl, 'edit_link');
    var confirmBtns = j6.dom.getsByClass(pl, 'confirm-button');

    j6.each(editLinks,
        function(i, e) {
            j6.event.add(e.getElementsByTagName('A')[0], 'click', (function(items, i, e) {
                return function() {
                    j6.each(items,
                        function(i2, e2) {
                            if (i2 == i) {
                                e2.className += ' active_item';
                                //e2.style.display='none';
                                //j6.animation.toggleHeight(e2,null,15);
                            } else {
                                e2.className = e2.className.replace(' active_item', '');
                            }
                        });
                };
            })(items, i, e));
        });

    j6.each(confirmBtns,
        function(i, e) {
            j6.event.add(e, 'click', (function(item) {
                return function() {
                    item.className = item.className.replace(' active_item', '');
                };
            })(items[i]));
        });

    j6.$('cb1').onclick = function() {};
    j6.$('el2').onclick = function() {
        j6.json.bind('ctl2', {
            deliver_opt: window.sctJson.deliver_opt
        });
    };
    j6.$('cb2').onclick = function() {
        var data = j6.json.toObject('ctl2');
        window.sctJson.deliver_opt = parseInt(data.deliver_opt);
        dynamicContent('deliver');
        persistData();
    };

    j6.$('el3').onclick = function() {
        j6.json.bind('ctl3', {
            pay_opt: window.sctJson.pay_opt
        });
    };

    j6.$('cb3').onclick = function() {
        var data = j6.json.toObject('ctl3');
        window.sctJson.pay_opt = parseInt(data.pay_opt);
        dynamicContent('payment');
        persistData();
    };

}

// 页面初始化
function pageInit(){
    initEvents();

    if(tbCouponCode){
        j6.event.add(tbCouponCode,'blur',applyCouponCode);
        if(tbCouponCode.value.replace(/\\s+/ig,'').length > 0 ){
            applyCouponCode.apply(tbCouponCode);
        }
    }

    if(lbBalance){
        lbBalance.innerHTML = acc_balance;
    }
    j6.event.add(ckBalance,'click',reloadFee);

    if(acc_balance>0){
        ckBalance.setAttribute('checked','checked');
        ckBalance.removeAttribute('disabled');
    }else{
        ckBalance.removeAttribute('checked');
        ckBalance.setAttribute('disabled','disabled');
    }
}

// 显示动态信息
function dynamicContent(t) {
    var showAll = t == "" || t == null;
    // 显示动态支付信息
    if (showAll || t == 'payment') {
        var payOpt = parseInt(window.sctJson.pay_opt);
        var payOptEle = j6.$('payment_opt_name');

        if (payOpt == 1) {
            payOptEle.innerHTML = '在线支付';
        } else if (payOpt == 2) {
            payOptEle.innerHTML = '货到付款';
        }else if(payOpt == 3){
            payOptEle.innerHTML = '转账汇款';
        }
    }

    // 显示动态配送信息
    if (showAll || t == 'deliver') {
        var dlOpt = parseInt(window.sctJson.deliver_opt);
        var dlOptEle = j6.$('deliver_opt_name');

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
    j6.xhr.jsonPost('/buy/buyingPersist', window.sctJson, function (d) {
        if (d.message) {
            window.cli.alert(d.message);
        }
    });
}

// 选择配送地址
function selectDeliver() {
    j6.load('deliver-panel', '/buy/getDeliverAddress?sel=' + window.sctJson.deliver_id);
}

// 从表单中恢复数据
function recoverFrom(id) {
    window.sctJson = j6.json.toObject(id);
    if (window.sctJson.deliver_id <= 0) {
        j6.$('item1').className += ' active_item';
        selectDeliver();
    }
}

function applyCouponCode() {
    if (this.value == '') {
        j6.validator.removeTip(this);
        couDes.innerHTML = '';
        if (couDes.indexOf(' hidden') == -1) {
            couDes.className += ' hidden';
        }
        coupon_fee = 0;
    } else {
        var t = this;
        j6.xhr.jsonPost('/buy/apply?type=coupon', {
                code: this.value
            },
            function(json) {
                if (json.result == false) {
                    j6.validator.setTip(t, false, null, json.message);
                    couDes.className = 'coupon_desc hidden';
                    coupon_fee = 0;
                    reloadFee();
                } else {
                    j6.validator.removeTip(t);
                    if (json.couponFee) {
                        couDes.className = 'coupon_desc';
                        couDes.innerHTML = '优惠内容：' + json.couponDescribe +
                            '<br /><em>使用该优惠券总节省：￥' + json.couponFee + '元</em>';
                    }
                    coupon_fee = json.couponFee;
                    reloadFee();
                }
            });
    }
}

function submitOrder() {
    var ele = this;

    var unDis =function(){
        ele.removeAttribute('disabled');
        ele.className +=' btn-orange';
    }

    if (j6.validator.validate('form_coupon')) {

        var data = window.sctJson;
        var cp = j6.json.toObject(form_coupon);
        if (data.deliver_id <= 0) {
            var e = j6.$('item1');
            e.className += ' active_item';
            window.cli.alert("请选择配送地址");
            return false;
        }

        ele.disabled = 'disabled';
        ele.className = ele.className.replace(' btn-orange','')

        data.coupon_code = cp.CouponCode;

        data.balance_discount = ckBalance.checked?'1':'0';

        j6.xhr.jsonPost('submit_0', data, function (j) {
            if (j.result) {
                var orderNo = j.data;
                location.replace("payment?order_no=" + orderNo)
            } else {
                unDis();
                window.cli.alert(j.message);
            }
        }, function (x) {
            unDis();
            window.cli.alert('订单提交失败');
        });
    }
}

function reloadFee(promFee,payFee) {
    var balancePay = 0;
    var _payFee = payFee || order_fee - coupon_fee;
    var _promFee = promFee || prom_fee + coupon_fee;

    if(ckBalance.checked && acc_balance > 0 ){
        balancePay = acc_balance > _payFee ?_payFee:acc_balance;
        _payFee -= balancePay;
    }
    j6.json.bind(cashPl, {
        PromFee: _promFee,
        BalanceFee :balancePay,
        OrderFee: _payFee,
    });
    j6.$('final_fee').innerHTML = _payFee;
}