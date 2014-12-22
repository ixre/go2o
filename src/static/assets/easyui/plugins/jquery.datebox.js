/**
 * jQuery EasyUI 1.3.4
 * 
 * Copyright (c) 2009-2013 www.jeasyui.com. All rights reserved.
 *
 * Licensed under the GPL or commercial licenses
 * To use it on other terms please contact us: info@jeasyui.com
 * http://www.gnu.org/licenses/gpl.txt
 * http://www.jeasyui.com/license_commercial.php
 *
 */
(function($){
function _1(_2){
var _3=$.data(_2,"datebox");
var _4=_3.options;
$(_2).addClass("datebox-f").combo($.extend({},_4,{onShowPanel:function(){
_5();
_4.onShowPanel.call(_2);
}}));
$(_2).combo("textbox").parent().addClass("datebox");
if(!_3.calendar){
_6();
}
function _6(){
var _7=$(_2).combo("panel");
_3.calendar=$("<div></div>").appendTo(_7).wrap("<div class=\"datebox-calendar-inner\"></div>");
_3.calendar.calendar({fit:true,border:false,onSelect:function(_8){
var _9=_4.formatter(_8);
_11(_2,_9);
$(_2).combo("hidePanel");
_4.onSelect.call(_2,_8);
}});
_11(_2,_4.value);
var _a=$("<div class=\"datebox-button\"></div>").appendTo(_7);
var _b=$("<a href=\"javascript:void(0)\" class=\"datebox-current\"></a>").html(_4.currentText).appendTo(_a);
var _c=$("<a href=\"javascript:void(0)\" class=\"datebox-close\"></a>").html(_4.closeText).appendTo(_a);
_b.click(function(){
_3.calendar.calendar({year:new Date().getFullYear(),month:new Date().getMonth()+1,current:new Date()});
});
_c.click(function(){
$(_2).combo("hidePanel");
});
};
function _5(){
if(_4.panelHeight!="auto"){
var _d=$(_2).combo("panel");
var ci=_d.children("div.datebox-calendar-inner");
var _e=_d.height();
_d.children().not(ci).each(function(){
_e-=$(this).outerHeight();
});
ci._outerHeight(_e);
}
_3.calendar.calendar("resize");
};
};
function _f(_10,q){
_11(_10,q);
};
function _12(_13){
var _14=$.data(_13,"datebox");
var _15=_14.options;
var c=_14.calendar;
var _16=_15.formatter(c.calendar("options").current);
_11(_13,_16);
$(_13).combo("hidePanel");
};
function _11(_17,_18){
var _19=$.data(_17,"datebox");
var _1a=_19.options;
$(_17).combo("setValue",_18).combo("setText",_18);
_19.calendar.calendar("moveTo",_1a.parser(_18));
};
$.fn.datebox=function(_1b,_1c){
if(typeof _1b=="string"){
var _1d=$.fn.datebox.methods[_1b];
if(_1d){
return _1d(this,_1c);
}else{
return this.combo(_1b,_1c);
}
}
_1b=_1b||{};
return this.each(function(){
var _1e=$.data(this,"datebox");
if(_1e){
$.extend(_1e.options,_1b);
}else{
$.data(this,"datebox",{options:$.extend({},$.fn.datebox.defaults,$.fn.datebox.parseOptions(this),_1b)});
}
_1(this);
});
};
$.fn.datebox.methods={options:function(jq){
var _1f=jq.combo("options");
return $.extend($.data(jq[0],"datebox").options,{originalValue:_1f.originalValue,disabled:_1f.disabled,readonly:_1f.readonly});
},calendar:function(jq){
return $.data(jq[0],"datebox").calendar;
},setValue:function(jq,_20){
return jq.each(function(){
_11(this,_20);
});
},reset:function(jq){
return jq.each(function(){
var _21=$(this).datebox("options");
$(this).datebox("setValue",_21.originalValue);
});
}};
$.fn.datebox.parseOptions=function(_22){
var t=$(_22);
return $.extend({},$.fn.combo.parseOptions(_22),{});
};
$.fn.datebox.defaults=$.extend({},$.fn.combo.defaults,{panelWidth:180,panelHeight:"auto",keyHandler:{up:function(){
},down:function(){
},enter:function(){
_12(this);
},query:function(q){
_f(this,q);
}},currentText:"Today",closeText:"Close",okText:"Ok",formatter:function(_23){
var y=_23.getFullYear();
var m=_23.getMonth()+1;
var d=_23.getDate();
return m+"/"+d+"/"+y;
},parser:function(s){
var t=Date.parse(s);
if(!isNaN(t)){
return new Date(t);
}else{
return new Date();
}
},onSelect:function(_24){
}});
})(jQuery);

