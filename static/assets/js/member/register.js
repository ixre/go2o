/*
 * module : 用户注册
 * author : jarryliu
 * date   : 2011/01/21
 */

var _vcode;//验证码

//对每个输入进行验证

//定义元素变量
//var uid,email,nick,pwd,repwd,realname,cardnumber;
var uname = $JS.$('username');
var pwd = $JS.$('password');
var repwd= $JS.$('repassword');
var phone = $JS.$('phone');
var valid = $JS.validator;
if(uname){
    uname.onblur=function(){
        if(this.value == undefined)return;
        if(this.value=='')valid.setTip(this,false,0,'请输入用户名!');
        //else if(!/^(?=[A-Za-z])/.test(this.value))valid.setError(this,1);     //必须字符开头
        else if(!/^[A-Za-z0-9]+$/.test(this.value))valid.setTip(this,false,'3','');
        else if(this.value.length<4||this.value.length>15)valid.setTip(this,false,'2','');
        else{
            var t=this;
             valid.setTip(t,false,null,'验证中...');
            $JS.xhr.jsonPost('/user/ValidUsr',{usr:escape(t.value)},function(json){
            	if(json.result){
            	    valid.setTip(t,true,null,'用户名可用');
            	}
            	else{
            	    valid.setTip(t,false,'4');
            	}
            });
        }
    };
}

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
 if(pwd){
    pwd.onblur=function(){
        if(this.value == undefined)return;
        if(this.value=='')valid.setTip(this,false,'0');
        else if(/^(?=_)/.test(this.value)||this.value.indexOf('_')==this.value.length-1)
            valid.setTip(this,false,'1');
       else if(!/^[A-Za-z0-9_]*$/.test(this.value)){valid.setTip(this,false,'3');}
       else if(this.value.length<6||this.value.length>12){valid.setTip(this,false,'2');}
       else{
          valid.setTip(this,true,null,'');
       }
    };
 }

 if(repwd!=null){
    repwd.onblur=function(){
        if(this.value == undefined)return;
        if(this.value!=pwd.value){
            valid.setTip(this,false,'0','两次密码输入不一致!');
        }
        else{
            valid.setTip(this,true,null,'');
        }
    };
 }

if(phone!=null){
	phone.onblur=function(){
        if(this.value == undefined)return;
        if(this.value!='' && !/^(13[0-9]|15[1|2|3|4|5|6|8|9]|18[0|6|7|8|9])(\d{8})$/.test(this.value)){
            valid.setTip(this,false,'0');
        }else{
            if(this.value!=''){
                valid.setTip(this,true,null,'填写正确');
            }else{
                valid.setTip(this,false,'2');
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
//		if(this.value==''&&location.href.indexOf('partnerid=666888')!=-1){
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
// sBirthday=sId.substr(6,4)+"-"+Number(sId.substr(10,2))+"-"+Number(sId.substr(12,2));
// var d=new Date(sBirthday.replace(/-/g,"/")) ;
// if(sBirthday!=(d.getFullYear()+"-"+ (d.getMonth()+1) + "-" + d.getDate()))return "&nbsp;";
// for(var i = 17;i>=0;i --) iSum += (Math.pow(2,i) % 11) * parseInt(sId.charAt(17 - i),11) ;
// if(iSum%11!=1) return "&nbsp;";
// return true;//aCity[parseInt(sId.substr(0,2))]+","+sBirthday+","+ (sId.substr(16,1)%2?"男":"女")
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