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
var _3=$.data(_2,"tabs").options;
if(_3.tabPosition=="left"||_3.tabPosition=="right"){
return;
}
var _4=$(_2).children("div.tabs-header");
var _5=_4.children("div.tabs-tool");
var _6=_4.children("div.tabs-scroller-left");
var _7=_4.children("div.tabs-scroller-right");
var _8=_4.children("div.tabs-wrap");
var _9=_4.outerHeight();
if(_3.plain){
_9-=_9-_4.height();
}
_5._outerHeight(_9);
var _a=0;
$("ul.tabs li",_4).each(function(){
_a+=$(this).outerWidth(true);
});
var _b=_4.width()-_5._outerWidth();
if(_a>_b){
_6.add(_7).show()._outerHeight(_9);
if(_3.toolPosition=="left"){
_5.css({left:_6.outerWidth(),right:""});
_8.css({marginLeft:_6.outerWidth()+_5._outerWidth(),marginRight:_7._outerWidth(),width:_b-_6.outerWidth()-_7.outerWidth()});
}else{
_5.css({left:"",right:_7.outerWidth()});
_8.css({marginLeft:_6.outerWidth(),marginRight:_7.outerWidth()+_5._outerWidth(),width:_b-_6.outerWidth()-_7.outerWidth()});
}
}else{
_6.add(_7).hide();
if(_3.toolPosition=="left"){
_5.css({left:0,right:""});
_8.css({marginLeft:_5._outerWidth(),marginRight:0,width:_b});
}else{
_5.css({left:"",right:0});
_8.css({marginLeft:0,marginRight:_5._outerWidth(),width:_b});
}
}
};
function _c(_d){
var _e=$.data(_d,"tabs").options;
var _f=$(_d).children("div.tabs-header");
if(_e.tools){
if(typeof _e.tools=="string"){
$(_e.tools).addClass("tabs-tool").appendTo(_f);
$(_e.tools).show();
}else{
_f.children("div.tabs-tool").remove();
var _10=$("<div class=\"tabs-tool\"><table cellspacing=\"0\" cellpadding=\"0\" style=\"height:100%\"><tr></tr></table></div>").appendTo(_f);
var tr=_10.find("tr");
for(var i=0;i<_e.tools.length;i++){
var td=$("<td></td>").appendTo(tr);
var _11=$("<a href=\"javascript:void(0);\"></a>").appendTo(td);
_11[0].onclick=eval(_e.tools[i].handler||function(){
});
_11.linkbutton($.extend({},_e.tools[i],{plain:true}));
}
}
}else{
_f.children("div.tabs-tool").remove();
}
};
function _12(_13){
var _14=$.data(_13,"tabs");
var _15=_14.options;
var cc=$(_13);
_15.fit?$.extend(_15,cc._fit()):cc._fit(false);
cc.width(_15.width).height(_15.height);
var _16=$(_13).children("div.tabs-header");
var _17=$(_13).children("div.tabs-panels");
var _18=_16.find("div.tabs-wrap");
var ul=_18.find(".tabs");
for(var i=0;i<_14.tabs.length;i++){
var _19=_14.tabs[i].panel("options");
var p_t=_19.tab.find("a.tabs-inner");
var _1a=parseInt(_19.tabWidth||_15.tabWidth)||undefined;
if(_1a){
p_t._outerWidth(_1a);
}else{
p_t.css("width","");
}
p_t._outerHeight(_15.tabHeight);
p_t.css("lineHeight",p_t.height()+"px");
}
if(_15.tabPosition=="left"||_15.tabPosition=="right"){
_16._outerWidth(_15.headerWidth);
_17._outerWidth(cc.width()-_15.headerWidth);
_16.add(_17)._outerHeight(_15.height);
_18._outerWidth(_16.width());
ul._outerWidth(_18.width()).css("height","");
}else{
_16._outerWidth(_15.width).css("height","");
ul._outerHeight(_15.tabHeight).css("width","");
_1(_13);
var _1b=_15.height;
if(!isNaN(_1b)){
_17._outerHeight(_1b-_16.outerHeight());
}else{
_17.height("auto");
}
var _1a=_15.width;
if(!isNaN(_1a)){
_17._outerWidth(_1a);
}else{
_17.width("auto");
}
}
};
function _1c(_1d){
var _1e=$.data(_1d,"tabs").options;
var tab=_1f(_1d);
if(tab){
var _20=$(_1d).children("div.tabs-panels");
var _21=_1e.width=="auto"?"auto":_20.width();
var _22=_1e.height=="auto"?"auto":_20.height();
tab.panel("resize",{width:_21,height:_22});
}
};
function _23(_24){
var _25=$.data(_24,"tabs").tabs;
var cc=$(_24);
cc.addClass("tabs-container");
var pp=$("<div class=\"tabs-panels\"></div>").insertBefore(cc);
cc.children("div").each(function(){
pp[0].appendChild(this);
});
cc[0].appendChild(pp[0]);
$("<div class=\"tabs-header\">"+"<div class=\"tabs-scroller-left\"></div>"+"<div class=\"tabs-scroller-right\"></div>"+"<div class=\"tabs-wrap\">"+"<ul class=\"tabs\"></ul>"+"</div>"+"</div>").prependTo(_24);
cc.children("div.tabs-panels").children("div").each(function(i){
var _26=$.extend({},$.parser.parseOptions(this),{selected:($(this).attr("selected")?true:undefined)});
var pp=$(this);
_25.push(pp);
_33(_24,pp,_26);
});
cc.children("div.tabs-header").find(".tabs-scroller-left, .tabs-scroller-right").hover(function(){
$(this).addClass("tabs-scroller-over");
},function(){
$(this).removeClass("tabs-scroller-over");
});
cc.bind("_resize",function(e,_27){
var _28=$.data(_24,"tabs").options;
if(_28.fit==true||_27){
_12(_24);
_1c(_24);
}
return false;
});
};
function _29(_2a){
var _2b=$.data(_2a,"tabs").options;
$(_2a).children("div.tabs-header").unbind().bind("click",function(e){
if($(e.target).hasClass("tabs-scroller-left")){
$(_2a).tabs("scrollBy",-_2b.scrollIncrement);
}else{
if($(e.target).hasClass("tabs-scroller-right")){
$(_2a).tabs("scrollBy",_2b.scrollIncrement);
}else{
var li=$(e.target).closest("li");
if(li.hasClass("tabs-disabled")){
return;
}
var a=$(e.target).closest("a.tabs-close");
if(a.length){
_49(_2a,_2c(li));
}else{
if(li.length){
_3e(_2a,_2c(li));
}
}
}
}
}).bind("contextmenu",function(e){
var li=$(e.target).closest("li");
if(li.hasClass("tabs-disabled")){
return;
}
if(li.length){
_2b.onContextMenu.call(_2a,e,li.find("span.tabs-title").html(),_2c(li));
}
});
function _2c(li){
var _2d=0;
li.parent().children("li").each(function(i){
if(li[0]==this){
_2d=i;
return false;
}
});
return _2d;
};
};
function _2e(_2f){
var _30=$.data(_2f,"tabs").options;
var _31=$(_2f).children("div.tabs-header");
var _32=$(_2f).children("div.tabs-panels");
_31.removeClass("tabs-header-top tabs-header-bottom tabs-header-left tabs-header-right");
_32.removeClass("tabs-panels-top tabs-panels-bottom tabs-panels-left tabs-panels-right");
if(_30.tabPosition=="top"){
_31.insertBefore(_32);
}else{
if(_30.tabPosition=="bottom"){
_31.insertAfter(_32);
_31.addClass("tabs-header-bottom");
_32.addClass("tabs-panels-top");
}else{
if(_30.tabPosition=="left"){
_31.addClass("tabs-header-left");
_32.addClass("tabs-panels-right");
}else{
if(_30.tabPosition=="right"){
_31.addClass("tabs-header-right");
_32.addClass("tabs-panels-left");
}
}
}
}
if(_30.plain==true){
_31.addClass("tabs-header-plain");
}else{
_31.removeClass("tabs-header-plain");
}
if(_30.border==true){
_31.removeClass("tabs-header-noborder");
_32.removeClass("tabs-panels-noborder");
}else{
_31.addClass("tabs-header-noborder");
_32.addClass("tabs-panels-noborder");
}
};
function _33(_34,pp,_35){
var _36=$.data(_34,"tabs");
_35=_35||{};
pp.panel($.extend({},_35,{border:false,noheader:true,closed:true,doSize:false,iconCls:(_35.icon?_35.icon:undefined),onLoad:function(){
if(_35.onLoad){
_35.onLoad.call(this,arguments);
}
_36.options.onLoad.call(_34,$(this));
}}));
var _37=pp.panel("options");
var _38=$(_34).children("div.tabs-header").find("ul.tabs");
_37.tab=$("<li></li>").appendTo(_38);
_37.tab.append("<a href=\"javascript:void(0)\" class=\"tabs-inner\">"+"<span class=\"tabs-title\"></span>"+"<span class=\"tabs-icon\"></span>"+"</a>");
$(_34).tabs("update",{tab:pp,options:_37});
};
function _39(_3a,_3b){
var _3c=$.data(_3a,"tabs").options;
var _3d=$.data(_3a,"tabs").tabs;
if(_3b.selected==undefined){
_3b.selected=true;
}
var pp=$("<div></div>").appendTo($(_3a).children("div.tabs-panels"));
_3d.push(pp);
_33(_3a,pp,_3b);
_3c.onAdd.call(_3a,_3b.title,_3d.length-1);
_12(_3a);
if(_3b.selected){
_3e(_3a,_3d.length-1);
}
};
function _3f(_40,_41){
var _42=$.data(_40,"tabs").selectHis;
var pp=_41.tab;
var _43=pp.panel("options").title;
pp.panel($.extend({},_41.options,{iconCls:(_41.options.icon?_41.options.icon:undefined)}));
var _44=pp.panel("options");
var tab=_44.tab;
var _45=tab.find("span.tabs-title");
var _46=tab.find("span.tabs-icon");
_45.html(_44.title);
_46.attr("class","tabs-icon");
tab.find("a.tabs-close").remove();
if(_44.closable){
_45.addClass("tabs-closable");
$("<a href=\"javascript:void(0)\" class=\"tabs-close\"></a>").appendTo(tab);
}else{
_45.removeClass("tabs-closable");
}
if(_44.iconCls){
_45.addClass("tabs-with-icon");
_46.addClass(_44.iconCls);
}else{
_45.removeClass("tabs-with-icon");
}
if(_43!=_44.title){
for(var i=0;i<_42.length;i++){
if(_42[i]==_43){
_42[i]=_44.title;
}
}
}
tab.find("span.tabs-p-tool").remove();
if(_44.tools){
var _47=$("<span class=\"tabs-p-tool\"></span>").insertAfter(tab.find("a.tabs-inner"));
if($.isArray(_44.tools)){
for(var i=0;i<_44.tools.length;i++){
var t=$("<a href=\"javascript:void(0)\"></a>").appendTo(_47);
t.addClass(_44.tools[i].iconCls);
if(_44.tools[i].handler){
t.bind("click",{handler:_44.tools[i].handler},function(e){
if($(this).parents("li").hasClass("tabs-disabled")){
return;
}
e.data.handler.call(this);
});
}
}
}else{
$(_44.tools).children().appendTo(_47);
}
var pr=_47.children().length*12;
if(_44.closable){
pr+=8;
}else{
pr-=3;
_47.css("right","5px");
}
_45.css("padding-right",pr+"px");
}
_12(_40);
$.data(_40,"tabs").options.onUpdate.call(_40,_44.title,_48(_40,pp));
};
function _49(_4a,_4b){
var _4c=$.data(_4a,"tabs").options;
var _4d=$.data(_4a,"tabs").tabs;
var _4e=$.data(_4a,"tabs").selectHis;
if(!_4f(_4a,_4b)){
return;
}
var tab=_50(_4a,_4b);
var _51=tab.panel("options").title;
var _52=_48(_4a,tab);
if(_4c.onBeforeClose.call(_4a,_51,_52)==false){
return;
}
var tab=_50(_4a,_4b,true);
tab.panel("options").tab.remove();
tab.panel("destroy");
_4c.onClose.call(_4a,_51,_52);
_12(_4a);
for(var i=0;i<_4e.length;i++){
if(_4e[i]==_51){
_4e.splice(i,1);
i--;
}
}
var _53=_4e.pop();
if(_53){
_3e(_4a,_53);
}else{
if(_4d.length){
_3e(_4a,0);
}
}
};
function _50(_54,_55,_56){
var _57=$.data(_54,"tabs").tabs;
if(typeof _55=="number"){
if(_55<0||_55>=_57.length){
return null;
}else{
var tab=_57[_55];
if(_56){
_57.splice(_55,1);
}
return tab;
}
}
for(var i=0;i<_57.length;i++){
var tab=_57[i];
if(tab.panel("options").title==_55){
if(_56){
_57.splice(i,1);
}
return tab;
}
}
return null;
};
function _48(_58,tab){
var _59=$.data(_58,"tabs").tabs;
for(var i=0;i<_59.length;i++){
if(_59[i][0]==$(tab)[0]){
return i;
}
}
return -1;
};
function _1f(_5a){
var _5b=$.data(_5a,"tabs").tabs;
for(var i=0;i<_5b.length;i++){
var tab=_5b[i];
if(tab.panel("options").closed==false){
return tab;
}
}
return null;
};
function _5c(_5d){
var _5e=$.data(_5d,"tabs").tabs;
for(var i=0;i<_5e.length;i++){
if(_5e[i].panel("options").selected){
_3e(_5d,i);
return;
}
}
if(_5e.length){
_3e(_5d,0);
}
};
function _3e(_5f,_60){
var _61=$.data(_5f,"tabs").options;
var _62=$.data(_5f,"tabs").tabs;
var _63=$.data(_5f,"tabs").selectHis;
if(_62.length==0){
return;
}
var _64=_50(_5f,_60);
if(!_64){
return;
}
var _65=_1f(_5f);
if(_65){
_65.panel("close");
_65.panel("options").tab.removeClass("tabs-selected");
}
_64.panel("open");
var _66=_64.panel("options").title;
_63.push(_66);
var tab=_64.panel("options").tab;
tab.addClass("tabs-selected");
var _67=$(_5f).find(">div.tabs-header>div.tabs-wrap");
var _68=tab.position().left;
var _69=_68+tab.outerWidth();
if(_68<0||_69>_67.width()){
var _6a=_68-(_67.width()-tab.width())/2;
$(_5f).tabs("scrollBy",_6a);
}else{
$(_5f).tabs("scrollBy",0);
}
_1c(_5f);
_61.onSelect.call(_5f,_66,_48(_5f,_64));
};
function _4f(_6b,_6c){
return _50(_6b,_6c)!=null;
};
$.fn.tabs=function(_6d,_6e){
if(typeof _6d=="string"){
return $.fn.tabs.methods[_6d](this,_6e);
}
_6d=_6d||{};
return this.each(function(){
var _6f=$.data(this,"tabs");
var _70;
if(_6f){
_70=$.extend(_6f.options,_6d);
_6f.options=_70;
}else{
$.data(this,"tabs",{options:$.extend({},$.fn.tabs.defaults,$.fn.tabs.parseOptions(this),_6d),tabs:[],selectHis:[]});
_23(this);
}
_c(this);
_2e(this);
_12(this);
_29(this);
_5c(this);
});
};
$.fn.tabs.methods={options:function(jq){
return $.data(jq[0],"tabs").options;
},tabs:function(jq){
return $.data(jq[0],"tabs").tabs;
},resize:function(jq){
return jq.each(function(){
_12(this);
_1c(this);
});
},add:function(jq,_71){
return jq.each(function(){
_39(this,_71);
});
},close:function(jq,_72){
return jq.each(function(){
_49(this,_72);
});
},getTab:function(jq,_73){
return _50(jq[0],_73);
},getTabIndex:function(jq,tab){
return _48(jq[0],tab);
},getSelected:function(jq){
return _1f(jq[0]);
},select:function(jq,_74){
return jq.each(function(){
_3e(this,_74);
});
},exists:function(jq,_75){
return _4f(jq[0],_75);
},update:function(jq,_76){
return jq.each(function(){
_3f(this,_76);
});
},enableTab:function(jq,_77){
return jq.each(function(){
$(this).tabs("getTab",_77).panel("options").tab.removeClass("tabs-disabled");
});
},disableTab:function(jq,_78){
return jq.each(function(){
$(this).tabs("getTab",_78).panel("options").tab.addClass("tabs-disabled");
});
},scrollBy:function(jq,_79){
return jq.each(function(){
var _7a=$(this).tabs("options");
var _7b=$(this).find(">div.tabs-header>div.tabs-wrap");
var pos=Math.min(_7b._scrollLeft()+_79,_7c());
_7b.animate({scrollLeft:pos},_7a.scrollDuration);
function _7c(){
var w=0;
var ul=_7b.children("ul");
ul.children("li").each(function(){
w+=$(this).outerWidth(true);
});
return w-_7b.width()+(ul.outerWidth()-ul.width());
};
});
}};
$.fn.tabs.parseOptions=function(_7d){
return $.extend({},$.parser.parseOptions(_7d,["width","height","tools","toolPosition","tabPosition",{fit:"boolean",border:"boolean",plain:"boolean",headerWidth:"number",tabWidth:"number",tabHeight:"number"}]));
};
$.fn.tabs.defaults={width:"auto",height:"auto",headerWidth:150,tabWidth:"auto",tabHeight:27,plain:false,fit:false,border:true,tools:null,toolPosition:"right",tabPosition:"top",scrollIncrement:100,scrollDuration:400,onLoad:function(_7e){
},onSelect:function(_7f,_80){
},onBeforeClose:function(_81,_82){
},onClose:function(_83,_84){
},onAdd:function(_85,_86){
},onUpdate:function(_87,_88){
},onContextMenu:function(e,_89,_8a){
}};
})(jQuery);

