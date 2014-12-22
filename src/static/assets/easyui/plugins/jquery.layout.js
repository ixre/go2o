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
var _1=false;
function _2(_3){
var _4=$.data(_3,"layout");
var _5=_4.options;
var _6=_4.panels;
var cc=$(_3);
if(_3.tagName=="BODY"){
cc._fit();
}else{
_5.fit?cc.css(cc._fit()):cc._fit(false);
}
function _7(pp){
var _8=pp.panel("options");
return Math.min(Math.max(_8.height,_8.minHeight),_8.maxHeight);
};
function _9(pp){
var _a=pp.panel("options");
return Math.min(Math.max(_a.width,_a.minWidth),_a.maxWidth);
};
var _b={top:0,left:0,width:cc.width(),height:cc.height()};
function _c(pp){
if(!pp.length){
return;
}
var _d=_7(pp);
pp.panel("resize",{width:cc.width(),height:_d,left:0,top:0});
_b.top+=_d;
_b.height-=_d;
};
if(_14(_6.expandNorth)){
_c(_6.expandNorth);
}else{
_c(_6.north);
}
function _e(pp){
if(!pp.length){
return;
}
var _f=_7(pp);
pp.panel("resize",{width:cc.width(),height:_f,left:0,top:cc.height()-_f});
_b.height-=_f;
};
if(_14(_6.expandSouth)){
_e(_6.expandSouth);
}else{
_e(_6.south);
}
function _10(pp){
if(!pp.length){
return;
}
var _11=_9(pp);
pp.panel("resize",{width:_11,height:_b.height,left:cc.width()-_11,top:_b.top});
_b.width-=_11;
};
if(_14(_6.expandEast)){
_10(_6.expandEast);
}else{
_10(_6.east);
}
function _12(pp){
if(!pp.length){
return;
}
var _13=_9(pp);
pp.panel("resize",{width:_13,height:_b.height,left:0,top:_b.top});
_b.left+=_13;
_b.width-=_13;
};
if(_14(_6.expandWest)){
_12(_6.expandWest);
}else{
_12(_6.west);
}
_6.center.panel("resize",_b);
};
function _15(_16){
var cc=$(_16);
cc.addClass("layout");
function _17(cc){
cc.children("div").each(function(){
var _18=$.fn.layout.parsePanelOptions(this);
if("north,south,east,west,center".indexOf(_18.region)>=0){
_1b(_16,_18,this);
}
});
};
cc.children("form").length?_17(cc.children("form")):_17(cc);
cc.append("<div class=\"layout-split-proxy-h\"></div><div class=\"layout-split-proxy-v\"></div>");
cc.bind("_resize",function(e,_19){
var _1a=$.data(_16,"layout").options;
if(_1a.fit==true||_19){
_2(_16);
}
return false;
});
};
function _1b(_1c,_1d,el){
_1d.region=_1d.region||"center";
var _1e=$.data(_1c,"layout").panels;
var cc=$(_1c);
var dir=_1d.region;
if(_1e[dir].length){
return;
}
var pp=$(el);
if(!pp.length){
pp=$("<div></div>").appendTo(cc);
}
var _1f=$.extend({},$.fn.layout.paneldefaults,{width:(pp.length?parseInt(pp[0].style.width)||pp.outerWidth():"auto"),height:(pp.length?parseInt(pp[0].style.height)||pp.outerHeight():"auto"),doSize:false,collapsible:true,cls:("layout-panel layout-panel-"+dir),bodyCls:"layout-body",onOpen:function(){
var _20=$(this).panel("header").children("div.panel-tool");
_20.children("a.panel-tool-collapse").hide();
var _21={north:"up",south:"down",east:"right",west:"left"};
if(!_21[dir]){
return;
}
var _22="layout-button-"+_21[dir];
var t=_20.children("a."+_22);
if(!t.length){
t=$("<a href=\"javascript:void(0)\"></a>").addClass(_22).appendTo(_20);
t.bind("click",{dir:dir},function(e){
_2f(_1c,e.data.dir);
return false;
});
}
$(this).panel("options").collapsible?t.show():t.hide();
}},_1d);
pp.panel(_1f);
_1e[dir]=pp;
if(pp.panel("options").split){
var _23=pp.panel("panel");
_23.addClass("layout-split-"+dir);
var _24="";
if(dir=="north"){
_24="s";
}
if(dir=="south"){
_24="n";
}
if(dir=="east"){
_24="w";
}
if(dir=="west"){
_24="e";
}
_23.resizable($.extend({},{handles:_24,onStartResize:function(e){
_1=true;
if(dir=="north"||dir=="south"){
var _25=$(">div.layout-split-proxy-v",_1c);
}else{
var _25=$(">div.layout-split-proxy-h",_1c);
}
var top=0,_26=0,_27=0,_28=0;
var pos={display:"block"};
if(dir=="north"){
pos.top=parseInt(_23.css("top"))+_23.outerHeight()-_25.height();
pos.left=parseInt(_23.css("left"));
pos.width=_23.outerWidth();
pos.height=_25.height();
}else{
if(dir=="south"){
pos.top=parseInt(_23.css("top"));
pos.left=parseInt(_23.css("left"));
pos.width=_23.outerWidth();
pos.height=_25.height();
}else{
if(dir=="east"){
pos.top=parseInt(_23.css("top"))||0;
pos.left=parseInt(_23.css("left"))||0;
pos.width=_25.width();
pos.height=_23.outerHeight();
}else{
if(dir=="west"){
pos.top=parseInt(_23.css("top"))||0;
pos.left=_23.outerWidth()-_25.width();
pos.width=_25.width();
pos.height=_23.outerHeight();
}
}
}
}
_25.css(pos);
$("<div class=\"layout-mask\"></div>").css({left:0,top:0,width:cc.width(),height:cc.height()}).appendTo(cc);
},onResize:function(e){
if(dir=="north"||dir=="south"){
var _29=$(">div.layout-split-proxy-v",_1c);
_29.css("top",e.pageY-$(_1c).offset().top-_29.height()/2);
}else{
var _29=$(">div.layout-split-proxy-h",_1c);
_29.css("left",e.pageX-$(_1c).offset().left-_29.width()/2);
}
return false;
},onStopResize:function(e){
cc.children("div.layout-split-proxy-v,div.layout-split-proxy-h").hide();
pp.panel("resize",e.data);
_2(_1c);
_1=false;
cc.find(">div.layout-mask").remove();
}},_1d));
}
};
function _2a(_2b,_2c){
var _2d=$.data(_2b,"layout").panels;
if(_2d[_2c].length){
_2d[_2c].panel("destroy");
_2d[_2c]=$();
var _2e="expand"+_2c.substring(0,1).toUpperCase()+_2c.substring(1);
if(_2d[_2e]){
_2d[_2e].panel("destroy");
_2d[_2e]=undefined;
}
}
};
function _2f(_30,_31,_32){
if(_32==undefined){
_32="normal";
}
var _33=$.data(_30,"layout").panels;
var p=_33[_31];
if(p.panel("options").onBeforeCollapse.call(p)==false){
return;
}
var _34="expand"+_31.substring(0,1).toUpperCase()+_31.substring(1);
if(!_33[_34]){
_33[_34]=_35(_31);
_33[_34].panel("panel").bind("click",function(){
var _36=_37();
p.panel("expand",false).panel("open").panel("resize",_36.collapse);
p.panel("panel").animate(_36.expand,function(){
$(this).unbind(".layout").bind("mouseleave.layout",{region:_31},function(e){
if(_1==true){
return;
}
_2f(_30,e.data.region);
});
});
return false;
});
}
var _38=_37();
if(!_14(_33[_34])){
_33.center.panel("resize",_38.resizeC);
}
p.panel("panel").animate(_38.collapse,_32,function(){
p.panel("collapse",false).panel("close");
_33[_34].panel("open").panel("resize",_38.expandP);
$(this).unbind(".layout");
});
function _35(dir){
var _39;
if(dir=="east"){
_39="layout-button-left";
}else{
if(dir=="west"){
_39="layout-button-right";
}else{
if(dir=="north"){
_39="layout-button-down";
}else{
if(dir=="south"){
_39="layout-button-up";
}
}
}
}
var _3a=$.extend({},$.fn.layout.paneldefaults,{cls:"layout-expand",title:"&nbsp;",closed:true,doSize:false,tools:[{iconCls:_39,handler:function(){
_3e(_30,_31);
return false;
}}]});
var p=$("<div></div>").appendTo(_30).panel(_3a);
p.panel("panel").hover(function(){
$(this).addClass("layout-expand-over");
},function(){
$(this).removeClass("layout-expand-over");
});
return p;
};
function _37(){
var cc=$(_30);
var _3b=_33.center.panel("options");
if(_31=="east"){
var _3c=_33["east"].panel("options");
return {resizeC:{width:_3b.width+_3c.width-28},expand:{left:cc.width()-_3c.width},expandP:{top:_3b.top,left:cc.width()-28,width:28,height:_3b.height},collapse:{left:cc.width(),top:_3b.top,height:_3b.height}};
}else{
if(_31=="west"){
var _3d=_33["west"].panel("options");
return {resizeC:{width:_3b.width+_3d.width-28,left:28},expand:{left:0},expandP:{left:0,top:_3b.top,width:28,height:_3b.height},collapse:{left:-_3d.width,top:_3b.top,height:_3b.height}};
}else{
if(_31=="north"){
var hh=cc.height()-28;
if(_14(_33.expandSouth)){
hh-=_33.expandSouth.panel("options").height;
}else{
if(_14(_33.south)){
hh-=_33.south.panel("options").height;
}
}
_33.east.panel("resize",{top:28,height:hh});
_33.west.panel("resize",{top:28,height:hh});
if(_14(_33.expandEast)){
_33.expandEast.panel("resize",{top:28,height:hh});
}
if(_14(_33.expandWest)){
_33.expandWest.panel("resize",{top:28,height:hh});
}
return {resizeC:{top:28,height:hh},expand:{top:0},expandP:{top:0,left:0,width:cc.width(),height:28},collapse:{top:-_33["north"].panel("options").height,width:cc.width()}};
}else{
if(_31=="south"){
var hh=cc.height()-28;
if(_14(_33.expandNorth)){
hh-=_33.expandNorth.panel("options").height;
}else{
if(_14(_33.north)){
hh-=_33.north.panel("options").height;
}
}
_33.east.panel("resize",{height:hh});
_33.west.panel("resize",{height:hh});
if(_14(_33.expandEast)){
_33.expandEast.panel("resize",{height:hh});
}
if(_14(_33.expandWest)){
_33.expandWest.panel("resize",{height:hh});
}
return {resizeC:{height:hh},expand:{top:cc.height()-_33["south"].panel("options").height},expandP:{top:cc.height()-28,left:0,width:cc.width(),height:28},collapse:{top:cc.height(),width:cc.width()}};
}
}
}
}
};
};
function _3e(_3f,_40){
var _41=$.data(_3f,"layout").panels;
var _42=_43();
var p=_41[_40];
if(p.panel("options").onBeforeExpand.call(p)==false){
return;
}
var _44="expand"+_40.substring(0,1).toUpperCase()+_40.substring(1);
_41[_44].panel("close");
p.panel("panel").stop(true,true);
p.panel("expand",false).panel("open").panel("resize",_42.collapse);
p.panel("panel").animate(_42.expand,function(){
_2(_3f);
});
function _43(){
var cc=$(_3f);
var _45=_41.center.panel("options");
if(_40=="east"&&_41.expandEast){
return {collapse:{left:cc.width(),top:_45.top,height:_45.height},expand:{left:cc.width()-_41["east"].panel("options").width}};
}else{
if(_40=="west"&&_41.expandWest){
return {collapse:{left:-_41["west"].panel("options").width,top:_45.top,height:_45.height},expand:{left:0}};
}else{
if(_40=="north"&&_41.expandNorth){
return {collapse:{top:-_41["north"].panel("options").height,width:cc.width()},expand:{top:0}};
}else{
if(_40=="south"&&_41.expandSouth){
return {collapse:{top:cc.height(),width:cc.width()},expand:{top:cc.height()-_41["south"].panel("options").height}};
}
}
}
}
};
};
function _14(pp){
if(!pp){
return false;
}
if(pp.length){
return pp.panel("panel").is(":visible");
}else{
return false;
}
};
function _46(_47){
var _48=$.data(_47,"layout").panels;
if(_48.east.length&&_48.east.panel("options").collapsed){
_2f(_47,"east",0);
}
if(_48.west.length&&_48.west.panel("options").collapsed){
_2f(_47,"west",0);
}
if(_48.north.length&&_48.north.panel("options").collapsed){
_2f(_47,"north",0);
}
if(_48.south.length&&_48.south.panel("options").collapsed){
_2f(_47,"south",0);
}
};
$.fn.layout=function(_49,_4a){
if(typeof _49=="string"){
return $.fn.layout.methods[_49](this,_4a);
}
_49=_49||{};
return this.each(function(){
var _4b=$.data(this,"layout");
if(_4b){
$.extend(_4b.options,_49);
}else{
var _4c=$.extend({},$.fn.layout.defaults,$.fn.layout.parseOptions(this),_49);
$.data(this,"layout",{options:_4c,panels:{center:$(),north:$(),south:$(),east:$(),west:$()}});
_15(this);
}
_2(this);
_46(this);
});
};
$.fn.layout.methods={resize:function(jq){
return jq.each(function(){
_2(this);
});
},panel:function(jq,_4d){
return $.data(jq[0],"layout").panels[_4d];
},collapse:function(jq,_4e){
return jq.each(function(){
_2f(this,_4e);
});
},expand:function(jq,_4f){
return jq.each(function(){
_3e(this,_4f);
});
},add:function(jq,_50){
return jq.each(function(){
_1b(this,_50);
_2(this);
if($(this).layout("panel",_50.region).panel("options").collapsed){
_2f(this,_50.region,0);
}
});
},remove:function(jq,_51){
return jq.each(function(){
_2a(this,_51);
_2(this);
});
}};
$.fn.layout.parseOptions=function(_52){
return $.extend({},$.parser.parseOptions(_52,[{fit:"boolean"}]));
};
$.fn.layout.defaults={fit:false};
$.fn.layout.parsePanelOptions=function(_53){
var t=$(_53);
return $.extend({},$.fn.panel.parseOptions(_53),$.parser.parseOptions(_53,["region",{split:"boolean",minWidth:"number",minHeight:"number",maxWidth:"number",maxHeight:"number"}]));
};
$.fn.layout.paneldefaults=$.extend({},$.fn.panel.defaults,{region:null,split:false,minWidth:10,minHeight:10,maxWidth:10000,maxHeight:10000});
})(jQuery);

