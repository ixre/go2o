/*
 * module : 用户注册
 * author : jarryliu
 * date   : 2011/01/21
 */

var _vcode;//验证码

//对每个输入进行验证

//定义元素变量

j6.$('usr').onblur = function () {
    if (this.value == undefined)return;
    if (this.value == '')j6.validator.setTip(this, false, 0, '请输入用户名!');
    //else if(!/^(?=[A-Za-z])/.test(this.value))valid.setError(this,1);     //必须字符开头
    else if (!/^[A-Za-z0-9]+$/.test(this.value))j6.validator.setTip(this, false,'3','');
    else {
        var t = this;
        j6.validator.setTip(t, false, null, '验证中...');
        j6.xhr.jsonPost('/user/ValidUsr', {usr: escape(t.value)}, function (json) {
            if (json.result) {
                j6.validator.setTip(t, true, null, '用户名可用');
            }
            else {
                j6.validator.setTip(t, false,null,json.message);
            }
        });
    }
};

/*
 nick.onblur=function(){
 if(this.value=='')valid.setError(this,0);
 //else if(!/^(?![_\d]+)/.test(this.value))valid.setError(this,1);  //必须字符开头
 else if(!/^[a-zA-Z0-9\u4e00-\u9fa5]+$/.test(this.value))valid.setError(this,2);
 else{
 var leng=this.value.length;
 var match=this.value.match(/[^\\\\\\\\\\\\\\\\x00-\\\\\\\\\\\\\\\\x80]/ig) ;
 if(match!=null)leng+=match.length;
 if(leng>10||leng<4)valid.setError(this,3);
 else{
 //检测昵称
 var t=this;
 j.ajax.get('app.axd?task=register,validfield,nickname,'+escape(this.value),
 function(x){if(x=="")valid.displayRight(t);else valid.displayError(t,x);});
 }
 }
 }

 */
j6.$('pwd').onblur = function () {
    if (/^(?=_)/.test(this.value) || this.value.indexOf('_') == this.value.length - 1)
        j6.validator.setTip(this, false, '1');
    else if (!/^[A-Za-z0-9_]*$/.test(this.value)) {
        j6.validator.setTip(this, false, '3');
    }
    else if (this.value.length < 6 || this.value.length > 12) {
        j6.validator.setTip(this, false, '2');
    }
    else {
        j6.validator.setTip(this, true, null, '');
    }
};

j6.$('rePwd').onblur = function () {
    if (this.value != j6.$('pwd').value) {
        j6.validator.setTip(this, false, null, "两次密码输入不一致")
    } else {
        j6.validator.removeTip(this);
    }
};

var phone = j6.$('phone');
if (phone != null) {
    phone.onblur = function () {
        if (this.value == undefined)return;
        if (this.value != '' && !/^(13[0-9]|15[0|1|2|3|4|5|6|8|9]|18[0|1|2|3|5|6|7|8|9]|17[0|6])(\d{8})$/.test(this.value)) {
            j6.validator.setTip(this, false, '0');
        } else {
            if (this.value != '') {
                j6.validator.removeTip(this);
            } else {
                j6.validator.setTip(this, false, '2');
            }
            /*
             if (this.value != ''){
             var t=this;
             j.ajax.get("app.axd?task=register,validfield,accountname,"+escape(this.value),
             function(x){
             if(x=="")valid.displayRight(t);
             else valid.setError(t,1);
             });


             }*/
        }
    };
}

var inviCode = j6.$('invi_code');
if (inviCode != null) {
    inviCode.onblur = function () {
        var val = this.value;
        if(val.length == 0) {
            j6.validator.removeTip(this);
        }else{
            var t = this;
            j6.validator.setTip(this, false, null, '验证中...');
            j6.xhr.jsonPost('/user/valid_invitation', {invi_code: val}, function (json) {
                if (json.result) {
                    //valid.setTip(t, true, null, '邀请人为:'+json.data.Name);
                    j6.validator.removeTip(t);
                }
                else {
                    j6.validator.setTip(t, false, null, json.message);
                }
            });
        }
    };
}

//if(email!=null){
//	email.onblur=function(){
//		if(!/\w+([-+.']\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*/.test(this.value)){
//			valid.setTip(this,false,0);
//		}else{
//            valid.setTip(this,true);
//		}
//	}
//}


//if(tguser!=null){
//	tguser.onblur=function(){
//		//valid.displayRight(this);
//		//return true;
//		if(this.value==''&&location.href.indexOf('partnerid=101')!=-1){
//			valid.setError(this,1);
//		}else{
//
//            var t=this;
//            j.ajax.get('/valid_tguser?partnerid=&user='+escape(t.value),
//            function(x){
//            	if(x=="")valid.displayRight(t);
//            	else valid.setError(t,0);
//            });
//		}
//	};
//}

//var aCity={11:"北京",12:"天津",13:"河北",14:"山西",15:"内蒙古",21:"辽宁",22:"吉林",23:"黑龙江",31:"上海",32:"江苏",33:"浙江",34:"安徽",35:"福建",36:"江西",37:"山东",41:"河南",42:"湖北",43:"湖南",44:"广东",45:"广西",46:"海南",50:"重庆",51:"四川",52:"贵州",53:"云南",54:"西藏",61:"陕西",62:"甘肃",63:"青海",64:"宁夏",65:"新疆",71:"台湾",81:"香港",82:"澳门",91:"国外"}
//function isCardID(sId){
// var iSum=0 ;
// var info="" ;
// if(!/^\d{17}(\d|x)$/i.test(sId)) return "&nbsp;";
// sId=sId.replace(/x$/i,"a");
// if(aCity[parseInt(sId.substr(0,2))]==null) return "&nbsp;";
// sBirthDay=sId.substr(6,4)+"-"+Number(sId.substr(10,2))+"-"+Number(sId.substr(12,2));
// var d=new Date(sBirthDay.replace(/-/g,"/")) ;
// if(sBirthDay!=(d.getFullYear()+"-"+ (d.getMonth()+1) + "-" + d.getDate()))return "&nbsp;";
// for(var i = 17;i>=0;i --) iSum += (Math.pow(2,i) % 11) * parseInt(sId.charAt(17 - i),11) ;
// if(iSum%11!=1) return "&nbsp;";
// return true;//aCity[parseInt(sId.substr(0,2))]+","+sBirthDay+","+ (sId.substr(16,1)%2?"男":"女")
//}


//if(cardnumber!=null){
//    cardnumber.onblur=function(){
//        /*var t=this;
//        j.ajax.get("app.axd?task=register,valididnumber,"+escape(this.value),
//        function(x){if(x=="True")valid.displayRight(t);else valid.setError(t,0);});
//        */
//        if(isCardID(this.value)==true){
//            valid.displayRight(this);
//        }else{
//            valid.setError(this,0);
//        }
//    };
//}

var btnRegister = j6.$('btn_register');
btnRegister.disabled='';
btnRegister.onclick = function () {
    var t = this;
    if (j6.validator.validate('reg_panel')) {
        var d = j6.json.toObject('reg_panel');
        if (d.remember != 'on') {
            alert('请同意注册条款')
        } else {
            var tip = j6.$('tip');
            t.disabled=true;
            j6.xhr.jsonPost('/user/postRegisterInfo', d, function (json) {
                if (json.result) {
                    var returnUrl = j6.request('return_url');
                    tip.className='tip-panel';
                    tip.innerHTML = '<span style="color:#0a0">注册成功，请等待页面跳转</span>';

                    setTimeout(function(){
                        if (returnUrl != '') {
                            location.replace(returnUrl);
                        } else {
                            location.replace('/user/login?return_url=');
                        }
                    },3000);
                } else {
                    tip.className='tip-panel';
                    tip.innerHTML = json.message;
                    t.disabled=false;
                }
                if(window.registerCallback){
                    window.registerCallback(json);
                }
            }, function () {
                tip.className='tip-panel';
                tip.innerHTML = '注册失败，请重试';
                t.disabled=false;
            });
        }
    }
}

function showRegTip(msg) {
    var win = window.parent || window;

    win.xhrCt1 = document.createElement("DIV");
    win.xhrCt1.className = 'xhr1-container hidden';
    win.xhrCt1.style.cssText = 'position:absolute;left:0;top:0;';
    win.xhrCt1.innerHTML = j6.template('<div class="gate" style="width:100%;height:{h}px;opacity:0.7;filter:alpha(opacity=70);background:#FFF;position:fixed;"' +
        'id="xhr1_gate_layout"></div><div class="msg" id="xhr1_msg_layout"></div>', {
        h: j6.screen.height()
    });
    win.document.body.appendChild(win.xhrCt1);
    win.xhrGate1 = win.document.getElementById('xhr1_gate_layout');
    win.xhrMsg1 = win.document.getElementById('xhr1_msg_layout');
    win.xhrMsg1.style.cssText += "position:fixed;top:200px;border:solid 1px #FC0;background:#FFE;border-radius:5px;color:#F00;font-weight:bold;z-index:200;width:40%;margin:0 auto;left:30%;text-align:center;padding:0.5em 0;line-height:1.5em;";

    win.xhrMsg1.innerHTL = msg;
}
