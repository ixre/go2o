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
var _3=$(_2);
_3.addClass("tree");
return _3;
};
function _4(_5){
var _6=[];
_7(_6,$(_5));
function _7(aa,_8){
_8.children("li").each(function(){
var _9=$(this);
var _a=$.extend({},$.parser.parseOptions(this,["id","iconCls","state"]),{checked:(_9.attr("checked")?true:undefined)});
_a.text=_9.children("span").html();
if(!_a.text){
_a.text=_9.html();
}
var _b=_9.children("ul");
if(_b.length){
_a.children=[];
_7(_a.children,_b);
}
aa.push(_a);
});
};
return _6;
};
function _c(_d){
var _e=$.data(_d,"tree").options;
$(_d).unbind().bind("mouseover",function(e){
var tt=$(e.target);
var _f=tt.closest("div.tree-node");
if(!_f.length){
return;
}
_f.addClass("tree-node-hover");
if(tt.hasClass("tree-hit")){
if(tt.hasClass("tree-expanded")){
tt.addClass("tree-expanded-hover");
}else{
tt.addClass("tree-collapsed-hover");
}
}
e.stopPropagation();
}).bind("mouseout",function(e){
var tt=$(e.target);
var _10=tt.closest("div.tree-node");
if(!_10.length){
return;
}
_10.removeClass("tree-node-hover");
if(tt.hasClass("tree-hit")){
if(tt.hasClass("tree-expanded")){
tt.removeClass("tree-expanded-hover");
}else{
tt.removeClass("tree-collapsed-hover");
}
}
e.stopPropagation();
}).bind("click",function(e){
var tt=$(e.target);
var _11=tt.closest("div.tree-node");
if(!_11.length){
return;
}
if(tt.hasClass("tree-hit")){
_86(_d,_11[0]);
return false;
}else{
if(tt.hasClass("tree-checkbox")){
_39(_d,_11[0],!tt.hasClass("tree-checkbox1"));
return false;
}else{
_d8(_d,_11[0]);
_e.onClick.call(_d,_14(_d,_11[0]));
}
}
e.stopPropagation();
}).bind("dblclick",function(e){
var _12=$(e.target).closest("div.tree-node");
if(!_12.length){
return;
}
_d8(_d,_12[0]);
_e.onDblClick.call(_d,_14(_d,_12[0]));
e.stopPropagation();
}).bind("contextmenu",function(e){
var _13=$(e.target).closest("div.tree-node");
if(!_13.length){
return;
}
_e.onContextMenu.call(_d,e,_14(_d,_13[0]));
e.stopPropagation();
});
};
function _15(_16){
var _17=$(_16).find("div.tree-node");
_17.draggable("disable");
_17.css("cursor","pointer");
};
function _18(_19){
var _1a=$.data(_19,"tree");
var _1b=_1a.options;
var _1c=_1a.tree;
_1a.disabledNodes=[];
_1c.find("div.tree-node").draggable({disabled:false,revert:true,cursor:"pointer",proxy:function(_1d){
var p=$("<div class=\"tree-node-proxy\"></div>").appendTo("body");
p.html("<span class=\"tree-dnd-icon tree-dnd-no\">&nbsp;</span>"+$(_1d).find(".tree-title").html());
p.hide();
return p;
},deltaX:15,deltaY:15,onBeforeDrag:function(e){
if(_1b.onBeforeDrag.call(_19,_14(_19,this))==false){
return false;
}
if($(e.target).hasClass("tree-hit")||$(e.target).hasClass("tree-checkbox")){
return false;
}
if(e.which!=1){
return false;
}
$(this).next("ul").find("div.tree-node").droppable({accept:"no-accept"});
var _1e=$(this).find("span.tree-indent");
if(_1e.length){
e.data.offsetWidth-=_1e.length*_1e.width();
}
},onStartDrag:function(){
$(this).draggable("proxy").css({left:-10000,top:-10000});
_1b.onStartDrag.call(_19,_14(_19,this));
var _1f=_14(_19,this);
if(_1f.id==undefined){
_1f.id="easyui_tree_node_id_temp";
_cb(_19,_1f);
}
_1a.draggingNodeId=_1f.id;
},onDrag:function(e){
var x1=e.pageX,y1=e.pageY,x2=e.data.startX,y2=e.data.startY;
var d=Math.sqrt((x1-x2)*(x1-x2)+(y1-y2)*(y1-y2));
if(d>3){
$(this).draggable("proxy").show();
}
this.pageY=e.pageY;
},onStopDrag:function(){
$(this).next("ul").find("div.tree-node").droppable({accept:"div.tree-node"});
for(var i=0;i<_1a.disabledNodes.length;i++){
$(_1a.disabledNodes[i]).droppable("enable");
}
_1a.disabledNodes=[];
var _20=_d5(_19,_1a.draggingNodeId);
if(_20&&_20.id=="easyui_tree_node_id_temp"){
_20.id="";
_cb(_19,_20);
}
_1b.onStopDrag.call(_19,_20);
}}).droppable({accept:"div.tree-node",onDragEnter:function(e,_21){
if(_1b.onDragEnter.call(_19,this,_14(_19,_21))==false){
_22(_21,false);
$(this).removeClass("tree-node-append tree-node-top tree-node-bottom");
$(this).droppable("disable");
_1a.disabledNodes.push(this);
}
},onDragOver:function(e,_23){
if($(this).droppable("options").disabled){
return;
}
var _24=_23.pageY;
var top=$(this).offset().top;
var _25=top+$(this).outerHeight();
_22(_23,true);
$(this).removeClass("tree-node-append tree-node-top tree-node-bottom");
if(_24>top+(_25-top)/2){
if(_25-_24<5){
$(this).addClass("tree-node-bottom");
}else{
$(this).addClass("tree-node-append");
}
}else{
if(_24-top<5){
$(this).addClass("tree-node-top");
}else{
$(this).addClass("tree-node-append");
}
}
if(_1b.onDragOver.call(_19,this,_14(_19,_23))==false){
_22(_23,false);
$(this).removeClass("tree-node-append tree-node-top tree-node-bottom");
$(this).droppable("disable");
_1a.disabledNodes.push(this);
}
},onDragLeave:function(e,_26){
_22(_26,false);
$(this).removeClass("tree-node-append tree-node-top tree-node-bottom");
_1b.onDragLeave.call(_19,this,_14(_19,_26));
},onDrop:function(e,_27){
var _28=this;
var _29,_2a;
if($(this).hasClass("tree-node-append")){
_29=_2b;
_2a="append";
}else{
_29=_2c;
_2a=$(this).hasClass("tree-node-top")?"top":"bottom";
}
if(_1b.onBeforeDrop.call(_19,_28,_c4(_19,_27),_2a)==false){
$(this).removeClass("tree-node-append tree-node-top tree-node-bottom");
return;
}
_29(_27,_28,_2a);
$(this).removeClass("tree-node-append tree-node-top tree-node-bottom");
}});
function _22(_2d,_2e){
var _2f=$(_2d).draggable("proxy").find("span.tree-dnd-icon");
_2f.removeClass("tree-dnd-yes tree-dnd-no").addClass(_2e?"tree-dnd-yes":"tree-dnd-no");
};
function _2b(_30,_31){
if(_14(_19,_31).state=="closed"){
_7a(_19,_31,function(){
_32();
});
}else{
_32();
}
function _32(){
var _33=$(_19).tree("pop",_30);
$(_19).tree("append",{parent:_31,data:[_33]});
_1b.onDrop.call(_19,_31,_33,"append");
};
};
function _2c(_34,_35,_36){
var _37={};
if(_36=="top"){
_37.before=_35;
}else{
_37.after=_35;
}
var _38=$(_19).tree("pop",_34);
_37.data=_38;
$(_19).tree("insert",_37);
_1b.onDrop.call(_19,_35,_38,_36);
};
};
function _39(_3a,_3b,_3c){
var _3d=$.data(_3a,"tree").options;
if(!_3d.checkbox){
return;
}
var _3e=_14(_3a,_3b);
if(_3d.onBeforeCheck.call(_3a,_3e,_3c)==false){
return;
}
var _3f=$(_3b);
var ck=_3f.find(".tree-checkbox");
ck.removeClass("tree-checkbox0 tree-checkbox1 tree-checkbox2");
if(_3c){
ck.addClass("tree-checkbox1");
}else{
ck.addClass("tree-checkbox0");
}
if(_3d.cascadeCheck){
_40(_3f);
_41(_3f);
}
_3d.onCheck.call(_3a,_3e,_3c);
function _41(_42){
var _43=_42.next().find(".tree-checkbox");
_43.removeClass("tree-checkbox0 tree-checkbox1 tree-checkbox2");
if(_42.find(".tree-checkbox").hasClass("tree-checkbox1")){
_43.addClass("tree-checkbox1");
}else{
_43.addClass("tree-checkbox0");
}
};
function _40(_44){
var _45=_91(_3a,_44[0]);
if(_45){
var ck=$(_45.target).find(".tree-checkbox");
ck.removeClass("tree-checkbox0 tree-checkbox1 tree-checkbox2");
if(_46(_44)){
ck.addClass("tree-checkbox1");
}else{
if(_47(_44)){
ck.addClass("tree-checkbox0");
}else{
ck.addClass("tree-checkbox2");
}
}
_40($(_45.target));
}
function _46(n){
var ck=n.find(".tree-checkbox");
if(ck.hasClass("tree-checkbox0")||ck.hasClass("tree-checkbox2")){
return false;
}
var b=true;
n.parent().siblings().each(function(){
if(!$(this).children("div.tree-node").children(".tree-checkbox").hasClass("tree-checkbox1")){
b=false;
}
});
return b;
};
function _47(n){
var ck=n.find(".tree-checkbox");
if(ck.hasClass("tree-checkbox1")||ck.hasClass("tree-checkbox2")){
return false;
}
var b=true;
n.parent().siblings().each(function(){
if(!$(this).children("div.tree-node").children(".tree-checkbox").hasClass("tree-checkbox0")){
b=false;
}
});
return b;
};
};
};
function _48(_49,_4a){
var _4b=$.data(_49,"tree").options;
var _4c=$(_4a);
if(_4d(_49,_4a)){
var ck=_4c.find(".tree-checkbox");
if(ck.length){
if(ck.hasClass("tree-checkbox1")){
_39(_49,_4a,true);
}else{
_39(_49,_4a,false);
}
}else{
if(_4b.onlyLeafCheck){
$("<span class=\"tree-checkbox tree-checkbox0\"></span>").insertBefore(_4c.find(".tree-title"));
}
}
}else{
var ck=_4c.find(".tree-checkbox");
if(_4b.onlyLeafCheck){
ck.remove();
}else{
if(ck.hasClass("tree-checkbox1")){
_39(_49,_4a,true);
}else{
if(ck.hasClass("tree-checkbox2")){
var _4e=true;
var _4f=true;
var _50=_51(_49,_4a);
for(var i=0;i<_50.length;i++){
if(_50[i].checked){
_4f=false;
}else{
_4e=false;
}
}
if(_4e){
_39(_49,_4a,true);
}
if(_4f){
_39(_49,_4a,false);
}
}
}
}
}
};
function _52(_53,ul,_54,_55){
var _56=$.data(_53,"tree").options;
_54=_56.loadFilter.call(_53,_54,$(ul).prev("div.tree-node")[0]);
if(!_55){
$(ul).empty();
}
var _57=[];
var _58=[];
var _59=$(ul).prev("div.tree-node").find("span.tree-indent, span.tree-hit").length;
_5a(ul,_54,_59);
if(_56.dnd){
_18(_53);
}else{
_15(_53);
}
if(_57.length){
_39(_53,_57[0],false);
}
for(var i=0;i<_58.length;i++){
_39(_53,_58[i],true);
}
setTimeout(function(){
_62(_53,_53);
},0);
var _5b=null;
if(_53!=ul){
var _5c=$(ul).prev();
_5b=_14(_53,_5c[0]);
}
_56.onLoadSuccess.call(_53,_5b,_54);
function _5a(ul,_5d,_5e){
for(var i=0;i<_5d.length;i++){
var li=$("<li></li>").appendTo(ul);
var _5f=_5d[i];
if(_5f.state!="open"&&_5f.state!="closed"){
_5f.state="open";
}
var _60=$("<div class=\"tree-node\"></div>").appendTo(li);
_60.attr("node-id",_5f.id);
$.data(_60[0],"tree-node",{id:_5f.id,text:_5f.text,iconCls:_5f.iconCls,attributes:_5f.attributes});
$("<span class=\"tree-title\"></span>").html(_56.formatter.call(_53,_5f)).appendTo(_60);
if(_56.checkbox){
if(_56.onlyLeafCheck){
if(_5f.state=="open"&&(!_5f.children||!_5f.children.length)){
if(_5f.checked){
$("<span class=\"tree-checkbox tree-checkbox1\"></span>").prependTo(_60);
}else{
$("<span class=\"tree-checkbox tree-checkbox0\"></span>").prependTo(_60);
}
}
}else{
if(_5f.checked){
$("<span class=\"tree-checkbox tree-checkbox1\"></span>").prependTo(_60);
_58.push(_60[0]);
}else{
$("<span class=\"tree-checkbox tree-checkbox0\"></span>").prependTo(_60);
if(_5d==_54){
_57.push(_60[0]);
}
}
}
}
if(_5f.children&&_5f.children.length){
var _61=$("<ul></ul>").appendTo(li);
if(_5f.state=="open"){
$("<span class=\"tree-icon tree-folder tree-folder-open\"></span>").addClass(_5f.iconCls).prependTo(_60);
$("<span class=\"tree-hit tree-expanded\"></span>").prependTo(_60);
}else{
$("<span class=\"tree-icon tree-folder\"></span>").addClass(_5f.iconCls).prependTo(_60);
$("<span class=\"tree-hit tree-collapsed\"></span>").prependTo(_60);
_61.css("display","none");
}
_5a(_61,_5f.children,_5e+1);
}else{
if(_5f.state=="closed"){
$("<span class=\"tree-icon tree-folder\"></span>").addClass(_5f.iconCls).prependTo(_60);
$("<span class=\"tree-hit tree-collapsed\"></span>").prependTo(_60);
}else{
$("<span class=\"tree-icon tree-file\"></span>").addClass(_5f.iconCls).prependTo(_60);
$("<span class=\"tree-indent\"></span>").prependTo(_60);
}
}
for(var j=0;j<_5e;j++){
$("<span class=\"tree-indent\"></span>").prependTo(_60);
}
}
};
};
function _62(_63,ul,_64){
var _65=$.data(_63,"tree").options;
if(!_65.lines){
return;
}
if(!_64){
_64=true;
$(_63).find("span.tree-indent").removeClass("tree-line tree-join tree-joinbottom");
$(_63).find("div.tree-node").removeClass("tree-node-last tree-root-first tree-root-one");
var _66=$(_63).tree("getRoots");
if(_66.length>1){
$(_66[0].target).addClass("tree-root-first");
}else{
if(_66.length==1){
$(_66[0].target).addClass("tree-root-one");
}
}
}
$(ul).children("li").each(function(){
var _67=$(this).children("div.tree-node");
var ul=_67.next("ul");
if(ul.length){
if($(this).next().length){
_68(_67);
}
_62(_63,ul,_64);
}else{
_69(_67);
}
});
var _6a=$(ul).children("li:last").children("div.tree-node").addClass("tree-node-last");
_6a.children("span.tree-join").removeClass("tree-join").addClass("tree-joinbottom");
function _69(_6b,_6c){
var _6d=_6b.find("span.tree-icon");
_6d.prev("span.tree-indent").addClass("tree-join");
};
function _68(_6e){
var _6f=_6e.find("span.tree-indent, span.tree-hit").length;
_6e.next().find("div.tree-node").each(function(){
$(this).children("span:eq("+(_6f-1)+")").addClass("tree-line");
});
};
};
function _70(_71,ul,_72,_73){
var _74=$.data(_71,"tree").options;
_72=_72||{};
var _75=null;
if(_71!=ul){
var _76=$(ul).prev();
_75=_14(_71,_76[0]);
}
if(_74.onBeforeLoad.call(_71,_75,_72)==false){
return;
}
var _77=$(ul).prev().children("span.tree-folder");
_77.addClass("tree-loading");
var _78=_74.loader.call(_71,_72,function(_79){
_77.removeClass("tree-loading");
_52(_71,ul,_79);
if(_73){
_73();
}
},function(){
_77.removeClass("tree-loading");
_74.onLoadError.apply(_71,arguments);
if(_73){
_73();
}
});
if(_78==false){
_77.removeClass("tree-loading");
}
};
function _7a(_7b,_7c,_7d){
var _7e=$.data(_7b,"tree").options;
var hit=$(_7c).children("span.tree-hit");
if(hit.length==0){
return;
}
if(hit.hasClass("tree-expanded")){
return;
}
var _7f=_14(_7b,_7c);
if(_7e.onBeforeExpand.call(_7b,_7f)==false){
return;
}
hit.removeClass("tree-collapsed tree-collapsed-hover").addClass("tree-expanded");
hit.next().addClass("tree-folder-open");
var ul=$(_7c).next();
if(ul.length){
if(_7e.animate){
ul.slideDown("normal",function(){
_7e.onExpand.call(_7b,_7f);
if(_7d){
_7d();
}
});
}else{
ul.css("display","block");
_7e.onExpand.call(_7b,_7f);
if(_7d){
_7d();
}
}
}else{
var _80=$("<ul style=\"display:none\"></ul>").insertAfter(_7c);
_70(_7b,_80[0],{id:_7f.id},function(){
if(_80.is(":empty")){
_80.remove();
}
if(_7e.animate){
_80.slideDown("normal",function(){
_7e.onExpand.call(_7b,_7f);
if(_7d){
_7d();
}
});
}else{
_80.css("display","block");
_7e.onExpand.call(_7b,_7f);
if(_7d){
_7d();
}
}
});
}
};
function _81(_82,_83){
var _84=$.data(_82,"tree").options;
var hit=$(_83).children("span.tree-hit");
if(hit.length==0){
return;
}
if(hit.hasClass("tree-collapsed")){
return;
}
var _85=_14(_82,_83);
if(_84.onBeforeCollapse.call(_82,_85)==false){
return;
}
hit.removeClass("tree-expanded tree-expanded-hover").addClass("tree-collapsed");
hit.next().removeClass("tree-folder-open");
var ul=$(_83).next();
if(_84.animate){
ul.slideUp("normal",function(){
_84.onCollapse.call(_82,_85);
});
}else{
ul.css("display","none");
_84.onCollapse.call(_82,_85);
}
};
function _86(_87,_88){
var hit=$(_88).children("span.tree-hit");
if(hit.length==0){
return;
}
if(hit.hasClass("tree-expanded")){
_81(_87,_88);
}else{
_7a(_87,_88);
}
};
function _89(_8a,_8b){
var _8c=_51(_8a,_8b);
if(_8b){
_8c.unshift(_14(_8a,_8b));
}
for(var i=0;i<_8c.length;i++){
_7a(_8a,_8c[i].target);
}
};
function _8d(_8e,_8f){
var _90=[];
var p=_91(_8e,_8f);
while(p){
_90.unshift(p);
p=_91(_8e,p.target);
}
for(var i=0;i<_90.length;i++){
_7a(_8e,_90[i].target);
}
};
function _92(_93,_94){
var c=$(_93).parent();
while(c[0].tagName!="BODY"&&c.css("overflow-y")!="auto"){
c=c.parent();
}
var n=$(_94);
var _95=n.offset().top;
if(c[0].tagName!="BODY"){
var _96=c.offset().top;
if(_95<_96){
c.scrollTop(c.scrollTop()+_95-_96);
}else{
if(_95+n.outerHeight()>_96+c.outerHeight()-18){
c.scrollTop(c.scrollTop()+_95+n.outerHeight()-_96-c.outerHeight()+18);
}
}
}else{
c.scrollTop(_95);
}
};
function _97(_98,_99){
var _9a=_51(_98,_99);
if(_99){
_9a.unshift(_14(_98,_99));
}
for(var i=0;i<_9a.length;i++){
_81(_98,_9a[i].target);
}
};
function _9b(_9c){
var _9d=_9e(_9c);
if(_9d.length){
return _9d[0];
}else{
return null;
}
};
function _9e(_9f){
var _a0=[];
$(_9f).children("li").each(function(){
var _a1=$(this).children("div.tree-node");
_a0.push(_14(_9f,_a1[0]));
});
return _a0;
};
function _51(_a2,_a3){
var _a4=[];
if(_a3){
_a5($(_a3));
}else{
var _a6=_9e(_a2);
for(var i=0;i<_a6.length;i++){
_a4.push(_a6[i]);
_a5($(_a6[i].target));
}
}
function _a5(_a7){
_a7.next().find("div.tree-node").each(function(){
_a4.push(_14(_a2,this));
});
};
return _a4;
};
function _91(_a8,_a9){
var ul=$(_a9).parent().parent();
if(ul[0]==_a8){
return null;
}else{
return _14(_a8,ul.prev()[0]);
}
};
function _aa(_ab,_ac){
_ac=_ac||"checked";
if(!$.isArray(_ac)){
_ac=[_ac];
}
var _ad=[];
for(var i=0;i<_ac.length;i++){
var s=_ac[i];
if(s=="checked"){
_ad.push("span.tree-checkbox1");
}else{
if(s=="unchecked"){
_ad.push("span.tree-checkbox0");
}else{
if(s=="indeterminate"){
_ad.push("span.tree-checkbox2");
}
}
}
}
var _ae=[];
$(_ab).find(_ad.join(",")).each(function(){
var _af=$(this).parent();
_ae.push(_14(_ab,_af[0]));
});
return _ae;
};
function _b0(_b1){
var _b2=$(_b1).find("div.tree-node-selected");
if(_b2.length){
return _14(_b1,_b2[0]);
}else{
return null;
}
};
function _b3(_b4,_b5){
var _b6=$(_b5.parent);
var _b7=_b5.data;
if(!_b7){
return;
}
_b7=$.isArray(_b7)?_b7:[_b7];
if(!_b7.length){
return;
}
var ul;
if(_b6.length==0){
ul=$(_b4);
}else{
if(_4d(_b4,_b6[0])){
var _b8=_b6.find("span.tree-icon");
_b8.removeClass("tree-file").addClass("tree-folder tree-folder-open");
var hit=$("<span class=\"tree-hit tree-expanded\"></span>").insertBefore(_b8);
if(hit.prev().length){
hit.prev().remove();
}
}
ul=_b6.next();
if(!ul.length){
ul=$("<ul></ul>").insertAfter(_b6);
}
}
_52(_b4,ul[0],_b7,true);
_48(_b4,ul.prev());
};
function _b9(_ba,_bb){
var ref=_bb.before||_bb.after;
var _bc=_91(_ba,ref);
var _bd=_bb.data;
if(!_bd){
return;
}
_bd=$.isArray(_bd)?_bd:[_bd];
if(!_bd.length){
return;
}
_b3(_ba,{parent:(_bc?_bc.target:null),data:_bd});
var li=$();
var _be=_bc?$(_bc.target).next().children("li:last"):$(_ba).children("li:last");
for(var i=0;i<_bd.length;i++){
li=_be.add(li);
_be=_be.prev();
}
if(_bb.before){
li.insertBefore($(ref).parent());
}else{
li.insertAfter($(ref).parent());
}
};
function _bf(_c0,_c1){
var _c2=_91(_c0,_c1);
var _c3=$(_c1);
var li=_c3.parent();
var ul=li.parent();
li.remove();
if(ul.children("li").length==0){
var _c3=ul.prev();
_c3.find(".tree-icon").removeClass("tree-folder").addClass("tree-file");
_c3.find(".tree-hit").remove();
$("<span class=\"tree-indent\"></span>").prependTo(_c3);
if(ul[0]!=_c0){
ul.remove();
}
}
if(_c2){
_48(_c0,_c2.target);
}
_62(_c0,_c0);
};
function _c4(_c5,_c6){
function _c7(aa,ul){
ul.children("li").each(function(){
var _c8=$(this).children("div.tree-node");
var _c9=_14(_c5,_c8[0]);
var sub=$(this).children("ul");
if(sub.length){
_c9.children=[];
_c7(_c9.children,sub);
}
aa.push(_c9);
});
};
if(_c6){
var _ca=_14(_c5,_c6);
_ca.children=[];
_c7(_ca.children,$(_c6).next());
return _ca;
}else{
return null;
}
};
function _cb(_cc,_cd){
var _ce=$.data(_cc,"tree").options;
var _cf=$(_cd.target);
var _d0=_14(_cc,_cd.target);
if(_d0.iconCls){
_cf.find(".tree-icon").removeClass(_d0.iconCls);
}
var _d1=$.extend({},_d0,_cd);
$.data(_cd.target,"tree-node",_d1);
_cf.attr("node-id",_d1.id);
_cf.find(".tree-title").html(_ce.formatter.call(_cc,_d1));
if(_d1.iconCls){
_cf.find(".tree-icon").addClass(_d1.iconCls);
}
if(_d0.checked!=_d1.checked){
_39(_cc,_cd.target,_d1.checked);
}
};
function _14(_d2,_d3){
var _d4=$.extend({},$.data(_d3,"tree-node"),{target:_d3,checked:$(_d3).find(".tree-checkbox").hasClass("tree-checkbox1")});
if(!_4d(_d2,_d3)){
_d4.state=$(_d3).find(".tree-hit").hasClass("tree-expanded")?"open":"closed";
}
return _d4;
};
function _d5(_d6,id){
var _d7=$(_d6).find("div.tree-node[node-id=\""+id+"\"]");
if(_d7.length){
return _14(_d6,_d7[0]);
}else{
return null;
}
};
function _d8(_d9,_da){
var _db=$.data(_d9,"tree").options;
var _dc=_14(_d9,_da);
if(_db.onBeforeSelect.call(_d9,_dc)==false){
return;
}
$("div.tree-node-selected",_d9).removeClass("tree-node-selected");
$(_da).addClass("tree-node-selected");
_db.onSelect.call(_d9,_dc);
};
function _4d(_dd,_de){
var _df=$(_de);
var hit=_df.children("span.tree-hit");
return hit.length==0;
};
function _e0(_e1,_e2){
var _e3=$.data(_e1,"tree").options;
var _e4=_14(_e1,_e2);
if(_e3.onBeforeEdit.call(_e1,_e4)==false){
return;
}
$(_e2).css("position","relative");
var nt=$(_e2).find(".tree-title");
var _e5=nt.outerWidth();
nt.empty();
var _e6=$("<input class=\"tree-editor\">").appendTo(nt);
_e6.val(_e4.text).focus();
_e6.width(_e5+20);
_e6.height(document.compatMode=="CSS1Compat"?(18-(_e6.outerHeight()-_e6.height())):18);
_e6.bind("click",function(e){
return false;
}).bind("mousedown",function(e){
e.stopPropagation();
}).bind("mousemove",function(e){
e.stopPropagation();
}).bind("keydown",function(e){
if(e.keyCode==13){
_e7(_e1,_e2);
return false;
}else{
if(e.keyCode==27){
_ed(_e1,_e2);
return false;
}
}
}).bind("blur",function(e){
e.stopPropagation();
_e7(_e1,_e2);
});
};
function _e7(_e8,_e9){
var _ea=$.data(_e8,"tree").options;
$(_e9).css("position","");
var _eb=$(_e9).find("input.tree-editor");
var val=_eb.val();
_eb.remove();
var _ec=_14(_e8,_e9);
_ec.text=val;
_cb(_e8,_ec);
_ea.onAfterEdit.call(_e8,_ec);
};
function _ed(_ee,_ef){
var _f0=$.data(_ee,"tree").options;
$(_ef).css("position","");
$(_ef).find("input.tree-editor").remove();
var _f1=_14(_ee,_ef);
_cb(_ee,_f1);
_f0.onCancelEdit.call(_ee,_f1);
};
$.fn.tree=function(_f2,_f3){
if(typeof _f2=="string"){
return $.fn.tree.methods[_f2](this,_f3);
}
var _f2=_f2||{};
return this.each(function(){
var _f4=$.data(this,"tree");
var _f5;
if(_f4){
_f5=$.extend(_f4.options,_f2);
_f4.options=_f5;
}else{
_f5=$.extend({},$.fn.tree.defaults,$.fn.tree.parseOptions(this),_f2);
$.data(this,"tree",{options:_f5,tree:_1(this)});
var _f6=_4(this);
if(_f6.length&&!_f5.data){
_f5.data=_f6;
}
}
_c(this);
if(_f5.lines){
$(this).addClass("tree-lines");
}
if(_f5.data){
_52(this,this,_f5.data);
}else{
if(_f5.dnd){
_18(this);
}else{
_15(this);
}
}
_70(this,this);
});
};
$.fn.tree.methods={options:function(jq){
return $.data(jq[0],"tree").options;
},loadData:function(jq,_f7){
return jq.each(function(){
_52(this,this,_f7);
});
},getNode:function(jq,_f8){
return _14(jq[0],_f8);
},getData:function(jq,_f9){
return _c4(jq[0],_f9);
},reload:function(jq,_fa){
return jq.each(function(){
if(_fa){
var _fb=$(_fa);
var hit=_fb.children("span.tree-hit");
hit.removeClass("tree-expanded tree-expanded-hover").addClass("tree-collapsed");
_fb.next().remove();
_7a(this,_fa);
}else{
$(this).empty();
_70(this,this);
}
});
},getRoot:function(jq){
return _9b(jq[0]);
},getRoots:function(jq){
return _9e(jq[0]);
},getParent:function(jq,_fc){
return _91(jq[0],_fc);
},getChildren:function(jq,_fd){
return _51(jq[0],_fd);
},getChecked:function(jq,_fe){
return _aa(jq[0],_fe);
},getSelected:function(jq){
return _b0(jq[0]);
},isLeaf:function(jq,_ff){
return _4d(jq[0],_ff);
},find:function(jq,id){
return _d5(jq[0],id);
},select:function(jq,_100){
return jq.each(function(){
_d8(this,_100);
});
},check:function(jq,_101){
return jq.each(function(){
_39(this,_101,true);
});
},uncheck:function(jq,_102){
return jq.each(function(){
_39(this,_102,false);
});
},collapse:function(jq,_103){
return jq.each(function(){
_81(this,_103);
});
},expand:function(jq,_104){
return jq.each(function(){
_7a(this,_104);
});
},collapseAll:function(jq,_105){
return jq.each(function(){
_97(this,_105);
});
},expandAll:function(jq,_106){
return jq.each(function(){
_89(this,_106);
});
},expandTo:function(jq,_107){
return jq.each(function(){
_8d(this,_107);
});
},scrollTo:function(jq,_108){
return jq.each(function(){
_92(this,_108);
});
},toggle:function(jq,_109){
return jq.each(function(){
_86(this,_109);
});
},append:function(jq,_10a){
return jq.each(function(){
_b3(this,_10a);
});
},insert:function(jq,_10b){
return jq.each(function(){
_b9(this,_10b);
});
},remove:function(jq,_10c){
return jq.each(function(){
_bf(this,_10c);
});
},pop:function(jq,_10d){
var node=jq.tree("getData",_10d);
jq.tree("remove",_10d);
return node;
},update:function(jq,_10e){
return jq.each(function(){
_cb(this,_10e);
});
},enableDnd:function(jq){
return jq.each(function(){
_18(this);
});
},disableDnd:function(jq){
return jq.each(function(){
_15(this);
});
},beginEdit:function(jq,_10f){
return jq.each(function(){
_e0(this,_10f);
});
},endEdit:function(jq,_110){
return jq.each(function(){
_e7(this,_110);
});
},cancelEdit:function(jq,_111){
return jq.each(function(){
_ed(this,_111);
});
}};
$.fn.tree.parseOptions=function(_112){
var t=$(_112);
return $.extend({},$.parser.parseOptions(_112,["url","method",{checkbox:"boolean",cascadeCheck:"boolean",onlyLeafCheck:"boolean"},{animate:"boolean",lines:"boolean",dnd:"boolean"}]));
};
$.fn.tree.defaults={url:null,method:"post",animate:false,checkbox:false,cascadeCheck:true,onlyLeafCheck:false,lines:false,dnd:false,data:null,formatter:function(node){
return node.text;
},loader:function(_113,_114,_115){
var opts=$(this).tree("options");
if(!opts.url){
return false;
}
$.ajax({type:opts.method,url:opts.url,data:_113,dataType:"json",success:function(data){
_114(data);
},error:function(){
_115.apply(this,arguments);
}});
},loadFilter:function(data,_116){
return data;
},onBeforeLoad:function(node,_117){
},onLoadSuccess:function(node,data){
},onLoadError:function(){
},onClick:function(node){
},onDblClick:function(node){
},onBeforeExpand:function(node){
},onExpand:function(node){
},onBeforeCollapse:function(node){
},onCollapse:function(node){
},onBeforeCheck:function(node,_118){
},onCheck:function(node,_119){
},onBeforeSelect:function(node){
},onSelect:function(node){
},onContextMenu:function(e,node){
},onBeforeDrag:function(node){
},onStartDrag:function(node){
},onStopDrag:function(node){
},onDragEnter:function(_11a,_11b){
},onDragOver:function(_11c,_11d){
},onDragLeave:function(_11e,_11f){
},onBeforeDrop:function(_120,_121,_122){
},onDrop:function(_123,_124,_125){
},onBeforeEdit:function(node){
},onAfterEdit:function(node){
},onCancelEdit:function(node){
}};
})(jQuery);

