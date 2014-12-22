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
function _1(_2,_3,_4){
for(var i=0;i<_2.length;i++){
var _5=_2[i];
if(_5[_3]==_4){
return _5;
}
}
return null;
};
function _6(_7,_8){
var _9=$(_7).combo("panel");
var _a=_9.find("div.combobox-item[value=\""+_8+"\"]");
if(_a.length){
if(_a.position().top<=0){
var h=_9.scrollTop()+_a.position().top;
_9.scrollTop(h);
}else{
if(_a.position().top+_a.outerHeight()>_9.height()){
var h=_9.scrollTop()+_a.position().top+_a.outerHeight()-_9.height();
_9.scrollTop(h);
}
}
}
};
function _b(_c,_d){
var _e=$(_c).combobox("options");
var _f=$(_c).combobox("panel");
var _10=_f.children("div.combobox-item-hover");
if(!_10.length){
_10=_f.children("div.combobox-item-selected");
}
_10.removeClass("combobox-item-hover");
if(!_10.length){
_10=_f.children("div.combobox-item:visible:"+(_d=="next"?"first":"last"));
}else{
if(_d=="next"){
_10=_10.nextAll("div.combobox-item:visible:first");
if(!_10.length){
_10=_f.children("div.combobox-item:visible:first");
}
}else{
_10=_10.prevAll("div.combobox-item:visible:first");
if(!_10.length){
_10=_f.children("div.combobox-item:visible:last");
}
}
}
if(_10.length){
_10.addClass("combobox-item-hover");
_6(_c,_10.attr("value"));
if(_e.selectOnNavigation){
_11(_c,_10.attr("value"));
}
}
};
function _11(_12,_13){
var _14=$.data(_12,"combobox").options;
var _15=$.data(_12,"combobox").data;
if(_14.multiple){
var _16=$(_12).combo("getValues");
for(var i=0;i<_16.length;i++){
if(_16[i]==_13){
return;
}
}
_16.push(_13);
_17(_12,_16);
}else{
_17(_12,[_13]);
}
var _18=_1(_15,_14.valueField,_13);
if(_18){
_14.onSelect.call(_12,_18);
}
};
function _19(_1a,_1b){
var _1c=$.data(_1a,"combobox");
var _1d=_1c.options;
var _1e=$(_1a).combo("getValues");
var _1f=$.inArray(_1b+"",_1e);
if(_1f>=0){
_1e.splice(_1f,1);
_17(_1a,_1e);
}
var _20=_1(_1c.data,_1d.valueField,_1b);
if(_20){
_1d.onUnselect.call(_1a,_20);
}
};
function _17(_21,_22,_23){
var _24=$.data(_21,"combobox").options;
var _25=$.data(_21,"combobox").data;
var _26=$(_21).combo("panel");
_26.find("div.combobox-item-selected").removeClass("combobox-item-selected");
var vv=[],ss=[];
for(var i=0;i<_22.length;i++){
var v=_22[i];
var s=v;
var _27=_1(_25,_24.valueField,v);
if(_27){
s=_27[_24.textField];
}
vv.push(v);
ss.push(s);
_26.find("div.combobox-item[value=\""+v+"\"]").addClass("combobox-item-selected");
}
$(_21).combo("setValues",vv);
if(!_23){
$(_21).combo("setText",ss.join(_24.separator));
}
};
function _28(_29,_2a,_2b){
var _2c=$.data(_29,"combobox");
var _2d=_2c.options;
_2c.data=_2d.loadFilter.call(_29,_2a);
_2a=_2c.data;
var _2e=$(_29).combobox("getValues");
var dd=[];
var _2f=undefined;
for(var i=0;i<_2a.length;i++){
var _30=_2a[i];
var v=_30[_2d.valueField];
var s=_30[_2d.textField];
var g=_30[_2d.groupField];
if(g){
if(_2f!=g){
_2f=g;
dd.push("<div class=\"combobox-group\" value=\""+g+"\">");
dd.push(_2d.groupFormatter?_2d.groupFormatter.call(_29,g):g);
dd.push("</div>");
}
}else{
_2f=undefined;
}
dd.push("<div class=\"combobox-item"+(g?" combobox-gitem":"")+"\" value=\""+v+"\""+(g?" group=\""+g+"\"":"")+">");
dd.push(_2d.formatter?_2d.formatter.call(_29,_30):s);
dd.push("</div>");
if(_30["selected"]){
(function(){
for(var i=0;i<_2e.length;i++){
if(v==_2e[i]){
return;
}
}
_2e.push(v);
})();
}
}
$(_29).combo("panel").html(dd.join(""));
if(_2d.multiple){
_17(_29,_2e,_2b);
}else{
if(_2e.length){
_17(_29,[_2e[_2e.length-1]],_2b);
}else{
_17(_29,[],_2b);
}
}
_2d.onLoadSuccess.call(_29,_2a);
};
function _31(_32,url,_33,_34){
var _35=$.data(_32,"combobox").options;
if(url){
_35.url=url;
}
_33=_33||{};
if(_35.onBeforeLoad.call(_32,_33)==false){
return;
}
_35.loader.call(_32,_33,function(_36){
_28(_32,_36,_34);
},function(){
_35.onLoadError.apply(this,arguments);
});
};
function _37(_38,q){
var _39=$.data(_38,"combobox");
var _3a=_39.options;
if(_3a.multiple&&!q){
_17(_38,[],true);
}else{
_17(_38,[q],true);
}
if(_3a.mode=="remote"){
_31(_38,null,{q:q},true);
}else{
var _3b=$(_38).combo("panel");
_3b.find("div.combobox-item,div.combobox-group").hide();
var _3c=_39.data;
var _3d=undefined;
for(var i=0;i<_3c.length;i++){
var _3e=_3c[i];
if(_3a.filter.call(_38,q,_3e)){
var v=_3e[_3a.valueField];
var s=_3e[_3a.textField];
var g=_3e[_3a.groupField];
var _3e=_3b.find("div.combobox-item[value=\""+v+"\"]");
_3e.show();
if(s==q){
_17(_38,[v],true);
_3e.addClass("combobox-item-selected");
}
if(_3a.groupField&&_3d!=g){
_3b.find("div.combobox-group[value=\""+g+"\"]").show();
_3d=g;
}
}
}
}
};
function _3f(_40){
var t=$(_40);
var _41=t.combobox("panel");
var _42=t.combobox("options");
var _43=t.combobox("getData");
var _44=_41.children("div.combobox-item-hover");
if(!_44.length){
_44=_41.children("div.combobox-item-selected");
}
if(!_44.length){
return;
}
if(_42.multiple){
if(_44.hasClass("combobox-item-selected")){
t.combobox("unselect",_44.attr("value"));
}else{
t.combobox("select",_44.attr("value"));
}
}else{
t.combobox("select",_44.attr("value"));
t.combobox("hidePanel");
}
var vv=[];
var _45=t.combobox("getValues");
for(var i=0;i<_45.length;i++){
if(_1(_43,_42.valueField,_45[i])){
vv.push(_45[i]);
}
}
t.combobox("setValues",vv);
};
function _46(_47){
var _48=$.data(_47,"combobox").options;
$(_47).addClass("combobox-f");
$(_47).combo($.extend({},_48,{onShowPanel:function(){
$(_47).combo("panel").find("div.combobox-item").show();
_6(_47,$(_47).combobox("getValue"));
_48.onShowPanel.call(_47);
}}));
$(_47).combo("panel").unbind().bind("mouseover",function(e){
$(this).children("div.combobox-item-hover").removeClass("combobox-item-hover");
$(e.target).closest("div.combobox-item").addClass("combobox-item-hover");
e.stopPropagation();
}).bind("mouseout",function(e){
$(e.target).closest("div.combobox-item").removeClass("combobox-item-hover");
e.stopPropagation();
}).bind("click",function(e){
var _49=$(e.target).closest("div.combobox-item");
if(!_49.length){
return;
}
var _4a=_49.attr("value");
if(_48.multiple){
if(_49.hasClass("combobox-item-selected")){
_19(_47,_4a);
}else{
_11(_47,_4a);
}
}else{
_11(_47,_4a);
$(_47).combo("hidePanel");
}
e.stopPropagation();
});
};
$.fn.combobox=function(_4b,_4c){
if(typeof _4b=="string"){
var _4d=$.fn.combobox.methods[_4b];
if(_4d){
return _4d(this,_4c);
}else{
return this.combo(_4b,_4c);
}
}
_4b=_4b||{};
return this.each(function(){
var _4e=$.data(this,"combobox");
if(_4e){
$.extend(_4e.options,_4b);
_46(this);
}else{
_4e=$.data(this,"combobox",{options:$.extend({},$.fn.combobox.defaults,$.fn.combobox.parseOptions(this),_4b),data:[]});
_46(this);
var _4f=$.fn.combobox.parseData(this);
if(_4f.length){
_28(this,_4f);
}
}
if(_4e.options.data){
_28(this,_4e.options.data);
}
_31(this);
});
};
$.fn.combobox.methods={options:function(jq){
var _50=jq.combo("options");
return $.extend($.data(jq[0],"combobox").options,{originalValue:_50.originalValue,disabled:_50.disabled,readonly:_50.readonly});
},getData:function(jq){
return $.data(jq[0],"combobox").data;
},setValues:function(jq,_51){
return jq.each(function(){
_17(this,_51);
});
},setValue:function(jq,_52){
return jq.each(function(){
_17(this,[_52]);
});
},clear:function(jq){
return jq.each(function(){
$(this).combo("clear");
var _53=$(this).combo("panel");
_53.find("div.combobox-item-selected").removeClass("combobox-item-selected");
});
},reset:function(jq){
return jq.each(function(){
var _54=$(this).combobox("options");
if(_54.multiple){
$(this).combobox("setValues",_54.originalValue);
}else{
$(this).combobox("setValue",_54.originalValue);
}
});
},loadData:function(jq,_55){
return jq.each(function(){
_28(this,_55);
});
},reload:function(jq,url){
return jq.each(function(){
_31(this,url);
});
},select:function(jq,_56){
return jq.each(function(){
_11(this,_56);
});
},unselect:function(jq,_57){
return jq.each(function(){
_19(this,_57);
});
}};
$.fn.combobox.parseOptions=function(_58){
var t=$(_58);
return $.extend({},$.fn.combo.parseOptions(_58),$.parser.parseOptions(_58,["valueField","textField","groupField","mode","method","url"]));
};
$.fn.combobox.parseData=function(_59){
var _5a=[];
var _5b=$(_59).combobox("options");
$(_59).children().each(function(){
if(this.tagName.toLowerCase()=="optgroup"){
var _5c=$(this).attr("label");
$(this).children().each(function(){
_5d(this,_5c);
});
}else{
_5d(this);
}
});
return _5a;
function _5d(el,_5e){
var t=$(el);
var _5f={};
_5f[_5b.valueField]=t.attr("value")!=undefined?t.attr("value"):t.html();
_5f[_5b.textField]=t.html();
_5f["selected"]=t.is(":selected");
if(_5e){
_5b.groupField=_5b.groupField||"group";
_5f[_5b.groupField]=_5e;
}
_5a.push(_5f);
};
};
$.fn.combobox.defaults=$.extend({},$.fn.combo.defaults,{valueField:"value",textField:"text",groupField:null,groupFormatter:function(_60){
return _60;
},mode:"local",method:"post",url:null,data:null,keyHandler:{up:function(){
_b(this,"prev");
},down:function(){
_b(this,"next");
},enter:function(){
_3f(this);
},query:function(q){
_37(this,q);
}},filter:function(q,row){
var _61=$(this).combobox("options");
return row[_61.textField].indexOf(q)==0;
},formatter:function(row){
var _62=$(this).combobox("options");
return row[_62.textField];
},loader:function(_63,_64,_65){
var _66=$(this).combobox("options");
if(!_66.url){
return false;
}
$.ajax({type:_66.method,url:_66.url,data:_63,dataType:"json",success:function(_67){
_64(_67);
},error:function(){
_65.apply(this,arguments);
}});
},loadFilter:function(_68){
return _68;
},onBeforeLoad:function(_69){
},onLoadSuccess:function(){
},onLoadError:function(){
},onSelect:function(_6a){
},onUnselect:function(_6b){
}});
})(jQuery);

