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
var _3=$.data(_2,"accordion");
var _4=_3.options;
var _5=_3.panels;
var cc=$(_2);
_4.fit?$.extend(_4,cc._fit()):cc._fit(false);
if(_4.width>0){
cc._outerWidth(_4.width);
}
var _6="auto";
if(_4.height>0){
cc._outerHeight(_4.height);
var _7=_5.length?_5[0].panel("header").css("height","")._outerHeight():"auto";
var _6=cc.height()-(_5.length-1)*_7;
}
for(var i=0;i<_5.length;i++){
var _8=_5[i];
_8.panel("header")._outerHeight(_7);
_8.panel("resize",{width:cc.width(),height:_6});
}
};
function _9(_a){
var _b=$.data(_a,"accordion").panels;
for(var i=0;i<_b.length;i++){
var _c=_b[i];
if(_c.panel("options").collapsed==false){
return _c;
}
}
return null;
};
function _d(_e,_f){
var _10=$.data(_e,"accordion").panels;
for(var i=0;i<_10.length;i++){
if(_10[i][0]==$(_f)[0]){
return i;
}
}
return -1;
};
function _11(_12,_13,_14){
var _15=$.data(_12,"accordion").panels;
if(typeof _13=="number"){
if(_13<0||_13>=_15.length){
return null;
}else{
var _16=_15[_13];
if(_14){
_15.splice(_13,1);
}
return _16;
}
}
for(var i=0;i<_15.length;i++){
var _16=_15[i];
if(_16.panel("options").title==_13){
if(_14){
_15.splice(i,1);
}
return _16;
}
}
return null;
};
function _17(_18){
var _19=$.data(_18,"accordion").options;
var cc=$(_18);
if(_19.border){
cc.removeClass("accordion-noborder");
}else{
cc.addClass("accordion-noborder");
}
};
function _1a(_1b){
var cc=$(_1b);
cc.addClass("accordion");
var _1c=[];
cc.children("div").each(function(){
var _1d=$.extend({},$.parser.parseOptions(this),{selected:($(this).attr("selected")?true:undefined)});
var pp=$(this);
_1c.push(pp);
_20(_1b,pp,_1d);
});
cc.bind("_resize",function(e,_1e){
var _1f=$.data(_1b,"accordion").options;
if(_1f.fit==true||_1e){
_1(_1b);
}
return false;
});
return {accordion:cc,panels:_1c};
};
function _20(_21,pp,_22){
pp.panel($.extend({},_22,{collapsible:false,minimizable:false,maximizable:false,closable:false,doSize:false,collapsed:true,headerCls:"accordion-header",bodyCls:"accordion-body",onBeforeExpand:function(){
if(_22.onBeforeExpand){
if(_22.onBeforeExpand.call(this)==false){
return false;
}
}
var _23=_9(_21);
if(_23){
var _24=$(_23).panel("header");
_24.removeClass("accordion-header-selected");
_24.find(".accordion-collapse").triggerHandler("click");
}
var _24=pp.panel("header");
_24.addClass("accordion-header-selected");
_24.find(".accordion-collapse").removeClass("accordion-expand");
},onExpand:function(){
if(_22.onExpand){
_22.onExpand.call(this);
}
var _25=$.data(_21,"accordion").options;
_25.onSelect.call(_21,pp.panel("options").title,_d(_21,this));
},onBeforeCollapse:function(){
if(_22.onBeforeCollapse){
if(_22.onBeforeCollapse.call(this)==false){
return false;
}
}
var _26=pp.panel("header");
_26.removeClass("accordion-header-selected");
_26.find(".accordion-collapse").addClass("accordion-expand");
}}));
var _27=pp.panel("header");
var t=$("<a class=\"accordion-collapse accordion-expand\" href=\"javascript:void(0)\"></a>").appendTo(_27.children("div.panel-tool"));
t.bind("click",function(e){
var _28=$.data(_21,"accordion").options.animate;
_35(_21);
if(pp.panel("options").collapsed){
pp.panel("expand",_28);
}else{
pp.panel("collapse",_28);
}
return false;
});
_27.click(function(){
$(this).find(".accordion-collapse").triggerHandler("click");
return false;
});
};
function _29(_2a,_2b){
var _2c=_11(_2a,_2b);
if(!_2c){
return;
}
var _2d=_9(_2a);
if(_2d&&_2d[0]==_2c[0]){
return;
}
_2c.panel("header").triggerHandler("click");
};
function _2e(_2f){
var _30=$.data(_2f,"accordion").panels;
for(var i=0;i<_30.length;i++){
if(_30[i].panel("options").selected){
_31(i);
return;
}
}
if(_30.length){
_31(0);
}
function _31(_32){
var _33=$.data(_2f,"accordion").options;
var _34=_33.animate;
_33.animate=false;
_29(_2f,_32);
_33.animate=_34;
};
};
function _35(_36){
var _37=$.data(_36,"accordion").panels;
for(var i=0;i<_37.length;i++){
_37[i].stop(true,true);
}
};
function add(_38,_39){
var _3a=$.data(_38,"accordion");
var _3b=_3a.options;
var _3c=_3a.panels;
if(_39.selected==undefined){
_39.selected=true;
}
_35(_38);
var pp=$("<div></div>").appendTo(_38);
_3c.push(pp);
_20(_38,pp,_39);
_1(_38);
_3b.onAdd.call(_38,_39.title,_3c.length-1);
if(_39.selected){
_29(_38,_3c.length-1);
}
};
function _3d(_3e,_3f){
var _40=$.data(_3e,"accordion");
var _41=_40.options;
var _42=_40.panels;
_35(_3e);
var _43=_11(_3e,_3f);
var _44=_43.panel("options").title;
var _45=_d(_3e,_43);
if(_41.onBeforeRemove.call(_3e,_44,_45)==false){
return;
}
var _43=_11(_3e,_3f,true);
if(_43){
_43.panel("destroy");
if(_42.length){
_1(_3e);
var _46=_9(_3e);
if(!_46){
_29(_3e,0);
}
}
}
_41.onRemove.call(_3e,_44,_45);
};
$.fn.accordion=function(_47,_48){
if(typeof _47=="string"){
return $.fn.accordion.methods[_47](this,_48);
}
_47=_47||{};
return this.each(function(){
var _49=$.data(this,"accordion");
var _4a;
if(_49){
_4a=$.extend(_49.options,_47);
_49.opts=_4a;
}else{
_4a=$.extend({},$.fn.accordion.defaults,$.fn.accordion.parseOptions(this),_47);
var r=_1a(this);
$.data(this,"accordion",{options:_4a,accordion:r.accordion,panels:r.panels});
}
_17(this);
_1(this);
_2e(this);
});
};
$.fn.accordion.methods={options:function(jq){
return $.data(jq[0],"accordion").options;
},panels:function(jq){
return $.data(jq[0],"accordion").panels;
},resize:function(jq){
return jq.each(function(){
_1(this);
});
},getSelected:function(jq){
return _9(jq[0]);
},getPanel:function(jq,_4b){
return _11(jq[0],_4b);
},getPanelIndex:function(jq,_4c){
return _d(jq[0],_4c);
},select:function(jq,_4d){
return jq.each(function(){
_29(this,_4d);
});
},add:function(jq,_4e){
return jq.each(function(){
add(this,_4e);
});
},remove:function(jq,_4f){
return jq.each(function(){
_3d(this,_4f);
});
}};
$.fn.accordion.parseOptions=function(_50){
var t=$(_50);
return $.extend({},$.parser.parseOptions(_50,["width","height",{fit:"boolean",border:"boolean",animate:"boolean"}]));
};
$.fn.accordion.defaults={width:"auto",height:"auto",fit:false,border:true,animate:true,onSelect:function(_51,_52){
},onAdd:function(_53,_54){
},onBeforeRemove:function(_55,_56){
},onRemove:function(_57,_58){
}};
})(jQuery);

