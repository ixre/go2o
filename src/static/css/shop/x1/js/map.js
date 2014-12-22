function _getId(_id){
  return document.getElementById(_id);
}
//初始化Geocoder，作为地址与坐标转换功能
//全局geocoder、全局相对地址
var geocoder,relativePath="",defaultCity="",map,mapType,geocoder2,indexMap="";
//全局Marker
var marker,mapLabel,promptlabel=null;	//弹出地图使用
var myIcon;
var addressHelperMap = new Array();
function initGeocoder(){
  if(null==geocoder){
    //geocoder = new google.maps.Geocoder();
  }
  if(null==geocoder2){
	geocoder2 = new BMap.Geocoder();
  }
}

function addressHelper(oInput,oDiv,count){
	this.cityId=_getId("cityId");
	this.oInput=oInput;
	this.oDiv=oDiv;
	this.oH4=oDiv.children[0];//sug标题
  this.oUl=oDiv.children[1];
	this.count=count;
  this.pointer=-1;
  this.defaultValue="";
  
  oInput.onfocus=this.focus;
  oInput.onkeyup=this.keyUp;
  oInput.onkeydown=this.keyDown;
  oInput.addressHelper=this;
  oInput.onblur=this.hideSuggest;
  if("您当前的小区、楼栋、大学名称.如：清华大学东门"==oInput.value){
	  oInput.value="";
	  oInput.style.color="#333";
  }
  else oInput.style.color="#666";

  this.oUl.onmouseover=this.mouseOver;
  this.oUl.onmousedown=this.mouseDown;
  //oDiv.addressHelper=this;
  this.oUl.addressHelper=this;
  if(_getId("searchBtn")!=null){
	  _getId("searchBtn").onclick=function(){submitForm('searchForm','','')};
	  this.showDefault();
  }
}
addressHelper.prototype.showDefault=function(){
		var hisSug=_getId("hisSug").innerHTML.replace(/(^\s*)/g, "");
		if(hisSug.length<5){
			//显示热门
			this.oH4.innerHTML="<strong>随便看看</strong>";
			this.oUl.innerHTML=_getId("hotSug").innerHTML.replace(/(^\s*)/g, "");
		}
		else{
			//显示历史搜索
			this.oUl.innerHTML=hisSug;
		}
		this.oH4.style.display="block";
		this.oDiv.style.display="block";
}
addressHelper.prototype.focus=function(oEvent)
{
  if("您当前的小区、楼栋、大学名称.如：清华大学东门"==this.addressHelper.oInput.value){
    this.addressHelper.oInput.value="";
    this.addressHelper.oInput.style.color="#333333";
  }
  this.addressHelper.showDefault();
};
addressHelper.prototype.mouseDown=function(oEvent)
{
  oEvent=window.event || oEvent;
  oSrcDiv=oEvent.target || oEvent.srcElement;
  this.addressHelper.oInput.value=oSrcDiv.innerHTML.replace(/<b>|<\/b>/ig,"");
  this.addressHelper.cityId.value=oSrcDiv.getAttribute("keyid");
  submitForm('searchForm','shopList.php','');
};

addressHelper.prototype.mouseOver=function()
{
  if(-1!=this.addressHelper.pointer)
  {
    this.addressHelper.oUl.childNodes[this.addressHelper.pointer].className="";
  }
};

addressHelper.prototype.moveDown=function()
{
	if(this.oUl.childNodes.length>0)
	{
		++this.pointer;
    if(this.pointer>this.oUl.childNodes.length-1){this.pointer=-1;}
    this.oInput.value=this.defaultValue;
		for(var i=0;i<this.oUl.childNodes.length;i++)
		{
			if(i==this.pointer)
			{
				this.oUl.childNodes[i].className="over";
        //TODO extra replace
				this.oInput.value=this.oUl.childNodes[i].innerHTML.replace(/<b>|<\/b>/ig,"");
			}
			else
			{
				this.oUl.childNodes[i].className="";
			}
		}
	}
};

addressHelper.prototype.moveUp=function()
{
	if(this.oUl.childNodes.length>0)
	{
		--this.pointer;
    if(this.pointer<-1){this.pointer=this.oUl.childNodes.length-1;}
    this.oInput.value=this.defaultValue;
		for(var i=0;i<this.oUl.childNodes.length;i++)
		{
			if(i==this.pointer)
			{
				this.oUl.childNodes[i].className="over";
        //TODO extra replace
				this.oInput.value=this.oUl.childNodes[i].innerHTML.replace(/<b>|<\/b>/ig,"");
			}
			else
			{
				this.oUl.childNodes[i].className="";
			}
		}
	}
};

addressHelper.prototype.keyDown=function(oEvent)
{
  oEvent=window.event || oEvent;
	iKeyCode=oEvent.keyCode;
	switch(iKeyCode)
	{
		case 38: //up arrow
			this.addressHelper.moveUp();
			break;
		case 40: //down arrow
			this.addressHelper.moveDown();
			break;
		case 13: //return key
			this.addressHelper.oH4.style.display="none";	//隐藏sug标题
			submitForm('searchForm','shopList.php','');
//			this.addressHelper.oInput.blur();
			//window.focus();
			break;
			/**/
	}
};

addressHelper.prototype.keyUp=function(oEvent)
{
  oEvent=oEvent || window.event;
	var iKeyCode=oEvent.keyCode;
	if(iKeyCode==8 || iKeyCode==46)
	{
		this.addressHelper.onTextChange(false); /* without autocomplete */
		this.addressHelper.oH4.style.display="none";	//隐藏sug标题
		if(this.addressHelper.oInput.value=="") this.addressHelper.showDefault();
	}
	else if (iKeyCode < 32 || (iKeyCode >= 33 && iKeyCode <= 46) || (iKeyCode >= 112 && iKeyCode <= 123)) 
	{
        //ignore
  } 
	else 
	{
		this.addressHelper.oH4.style.display="none";	//隐藏sug标题
		this.addressHelper.onTextChange(true); /* with autocomplete */
	}
};

addressHelper.prototype.onTextChange=function(bTextComplete)
{
  var txt=this.oInput.value;
  this.defaultValue=txt;
	var oThis=this;
	this.cur=-1;
	
	if(txt.length>0)
	{
    this.oUl.innerHTML = "";
    this.dynamicLoadSuggestion();
	}
	else
	{
		this.oUl.innerHTML="";
		this.oDiv.style.display="none";
	}
};

addressHelper.prototype.hideSuggest=function()
{
	if(""==this.addressHelper.oInput.value){
	    this.addressHelper.oInput.value="您当前的小区、楼栋、大学名称.如：清华大学东门";
	    this.addressHelper.oInput.style.color="#666666";
	  }
	this.addressHelper.oDiv.style.display="none";
};

addressHelper.prototype.dynamicAddSuggestions=function(aStr){
    var sList = "";
		for(i in aStr)
		{
			if(aStr[i].length>13||i>10) continue;
	var filters=addressFilter(aStr[i]).split(this.defaultValue);
	var filtered="<b>"+filters[0]+"</b>"+this.defaultValue+"<b>"+(filters[1]?filters[1]:"")+"</b>";
//      var filtered=addressFilter(aStr[i]).replace(this.defaultValue,"<>"+this.defaultValue+"<>");
	  var tmpLi="<li keyid='"+this.cityId.value+"'>"+filtered+"</li>";
      if(""!=filtered && -1==this.oUl.innerHTML.indexOf(tmpLi)){
        sList += tmpLi;
      }
		}
    this.oUl.innerHTML += sList;
	this.pointer=-1;
}

addressHelper.prototype.loading=function(bLoad){
  if(bLoad){
    this.oUl.innerHTML += "<li><img src='images/loadingapple.gif' border=0></img></li>";
  }else{
    this.oUl.innerHTML = this.oUl.innerHTML.substr(0,this.oUl.innerHTML.lastIndexOf("<li>"));
  }
}

addressHelper.prototype.failedLoadSuggestion=function(msg){
}

addressHelper.prototype.newHttpRequest=function(){
  //oReq.addressHelper=this;
  return createXMLHttp();
}

addressHelper.prototype.dynamicLoadSuggestion=function(){
	addressHelperMap.lenght=0;
  //this.loading(true);
  var ah = new Array(this);
  var div=this.oDiv;
  /*baidu map*/
  if(null==_getId("baidu_container")){
    var container=document.createElement("DIV");
    container.setAttribute("id","baidu_container");
    container.style.display="none";
    document.body.appendChild(container);
  }
  var map = new BMap.Map("baidu_container");
  var options = {
  onSearchComplete: function(results){
    if (local.getStatus() == BMAP_STATUS_SUCCESS){
      // 判断状态是否正确
      var s = [];
      for (var i = 0; i < results.getCurrentNumPois(); i ++){
        s.push(results.getPoi(i).title);
		addressHelperMap.push(
			{'name':results.getPoi(i).title,
			'lat':results.getPoi(i).point.lat,
			'lng':results.getPoi(i).point.lng
			});
      }
      //alert(ah[0]);
      ah[0].dynamicAddSuggestions(s);
      div.style.display="block";
      //document.getElementById("log").innerHTML = s.join("<br/>");
    }
    else div.style.display="none";
  }
  };
  var local = new BMap.LocalSearch(map, options);
  local.search(defaultCity+" "+this.defaultValue);
}

function createXMLHttp(){
  var oReq = null;
  if(window.ActiveXObject)
    oReq=new ActiveXObject("MSXML2.XMLHTTP");
  else if(window.createRequest)
    oReq=window.createRequest();
  else
    oReq=new XMLHttpRequest();
  
  return oReq;
}

function initIpAddressAndCityBox(oInput,oButton,oCityBox){

	defaultCity=oInput.innerHTML.replace("站","");
  oButton.onclick=function(){
	  clickOpen(oCityBox.id);
  }
  oButton.onmouseout=function(){
    displayCityBox(oCityBox,false);
  }

  oCityBox.onmouseover=function(){
    displayCityBox(oCityBox,true);
  }
  oCityBox.onmouseout=function(){
    displayCityBox(oCityBox,false);
  }
  oCityBox.onclick=function(oEvent){
    oEvent=oEvent || window.event;
    oSrc=oEvent.target || oEvent.srcElement;
    if(null!=oSrc&&null!=oSrc.parentNode&&null!=oSrc.parentNode.getAttribute("key")){
    	/*By Jerry*/
    	_getId("cityId").value=oSrc.parentNode.getAttribute("keyid");
    	/**/
      defaultCity=oSrc.innerHTML;
      oInput.innerHTML=oSrc.innerHTML;
//      saveCityName(defaultCity);//in cookie
      displayCityBox(oCityBox,false);
    }
  }
  /*By Jerry*/
//  if(null==getCityName()){
//    getIpAddress(oInput);
//  }else{
//    oInput.innerHTML=getCityName();
//    defaultCity=getCityName();
//  }
  /**/
}

function displayCityBox(oBox,_display){
  if(oBox){
    if(_display){
      oBox.display=1;
      setTimeout(function(){if(1==oBox.display) {oBox.style.display="block"} },300);
    }else{
      oBox.display=0;
      setTimeout(function(){if(0==oBox.display) {oBox.style.display="none"} },600);
    }
  }
}

function getIpAddress(oInput){
  var oHttpReq = createXMLHttp();
  var getIpUrl="iplocate.php"//ip delegrate
  oHttpReq.open("GET", getIpUrl, true);
  oHttpReq.onreadystatechange = function(){
    if(oHttpReq.readyState==4){
      var results = oHttpReq.responseText;
      try{
        eval(results);
        if(null!=remote_ip_info.city){
           defaultCity=remote_ip_info.city;
           oInput.innerHTML=remote_ip_info.city;
        }
      }catch(e){throw e;}
    }
  };
  oHttpReq.send();
}
function initAddressHelper(path){
	if(path) relativePath=path;
  initGeocoder();
  new addressHelper(_getId("txt_map"),_getId("sug"),5);
//   if(null!=getInputAddress()){
//    _getId("txt_map").value=getInputAddress();
//    _getId("txt_map").style.color="#333333";
//  }
}
function initSearchBar(){
	initIpAddressAndCityBox(_getId("selectArea"),_getId("selectArea"),_getId("city_pop_box"));
}

function addressFilter(address){

  if(0<=address.indexOf("中国")){
    address=address.substr(2);
  }
//  if(0<=address.indexOf(defaultCity)){
//    address=address.substr(defaultCity.length+address.indexOf(defaultCity));
//  }
  if("市"==address.substr(0,1)){
    address=address.substr(1);
  }
  address=address.replace(/邮政编码: [\d]{6}/,"");
  return address;
}
//获取地址坐标
function getAddress(cityId,address,funcStr){
	for(var i=0;i<addressHelperMap.length;i++){
		if(address==addressHelperMap[i].name){
			var mapx=addressHelperMap[i].lng;
			var mapy=addressHelperMap[i].lat;
			submitUrl(relativePath+'module/map/mapSug.php?funcStr='+funcStr+'&cityId='+cityId+'&address='+address+'&mapx='+mapx+"&mapy="+mapy,'','');
			return;
		}
	}
	geocoder2.getPoint(defaultCity+" "+address, function(point){
	  if (point) {
        var mapx=point.lng;
        var mapy=point.lat;
        submitUrl(relativePath+'module/map/mapSug.php?funcStr='+funcStr+'&cityId='+cityId+'&address='+address+'&mapx='+mapx+"&mapy="+mapy,'','');
      }	
	  else{
		  submitUrl(relativePath+'module/map/mapSug.php?cityId='+cityId+'&address='+address+'&mapx=0&mapy=0','','');
		  locateInfo(true);
		  if(indexMap!="") eval(indexMap);
	  }
	});
}

/*sample getPoint2('春秀路',function(x,y){alert(x);alert(y)});*/
function getPoint2(cityId,address,funcStr){
	geocoder2.getPoint(address, function(point){
		  if (point) {
	        var mapx=point.lng;
	        var mapy=point.lat;
	        submitUrl(relativePath+'module/map/mapSug.php?funcStr='+funcStr+'&cityId='+cityId+'&address='+address+'&mapx='+mapx+"&mapy="+mapy,'','');
//			return;
//	        callback(mapx,mapy);
	      }	
		  else{
			    var pointxy = window.showModalDialog("maptool.html",address,"dialogWidth=1460px;dialogHeight=600px");
		        var mapx=pointxy.split(',')[0];
		        var mapy=pointxy.split(',')[1];
		        submitUrl(relativePath+'module/map/mapSug.php?funcStr='+funcStr+'&cityId='+cityId+'&address='+address+'&mapx='+mapx+"&mapy="+mapy,'','');
		  }});	
}

//提示marker下经纬度
function __showMarkerPoint(latlng){
  var address=_getId("txt_map").value;
  var msg=address+"坐标是\n"
  +"经度"+latlng.lng()
  +"\n纬度"+latlng.lat();
  alert(msg);
  return new Array(latlng.lng(), latlng.lat());
}


function loadDirection(oDiv, p1, p2, oInfo, oInfo2){
	var dis="50米";
	    if(null==map){
		  map = new BMap.Map("s_u_map");
		}
	    map.disableDoubleClickZoom();
		var walking = new BMap.WalkingRoute(map, {onSearchComplete:function(rs){
			if(null!=oInfo){
				dis=rs.getPlan(0).getDistance(true);
				dis=dis=="0米"?"50米":dis;
				oInfo.innerHTML = "到本店的路程约：" + dis;
			}
		}, renderOptions:{map: map, autoViewport: true, selectFirstResult:true}});
		walking.search(p1, p2);
		if(null!=oInfo2){
			var diss=map.getDistance(p1,p2);
			diss=Math.ceil(isNaN(diss)?50:diss);
			dis=dis.replace("米","").replace("公里","");
			if(diss<parseInt(dis)) diss=dis;
			oInfo2.innerHTML="直线距离约：" + diss + "米";
		}		
}

function getDistance(p1, p2){
  if(null==map){
    map = new BMap.Map("s_u_map");
  }
  return isNaN(map.getDistance(p1,p2)) ? "0米" :Math.ceil(map.getDistance(p1,p2)) + "米";
}

function popShops(popId, shopList,start){
	if(null!=_getId(popId)){
		map = new BMap.Map(popId);
		var opts = {anchor: BMAP_ANCHOR_TOP_LEFT, offset: new BMap.Size(10, 10)};
		map.addControl(new BMap.NavigationControl(opts));
		map.addControl(new BMap.ScaleControl()); 
		map.addControl(new BMap.OverviewMapControl()); 
		map.enableScrollWheelZoom();//允许滚轮
		 if(null!=curPoint){      
			 var marker = new BMap.Marker(curPoint,{icon:new BMap.Icon("images/shop_map/map_me_icon.png",new BMap.Size(19,59)),offset:new BMap.Size(0,-30)});
			 marker.setTitle("您的位置");
			 marker.setTop(true);
			 map.centerAndZoom(curPoint, 16);
			 map.addOverlay(marker);		  
		 }
		 myIcon = new BMap.Icon( "images/shop_map/map_cursor.png",new BMap.Size(24, 24),{infoWindowAnchor:new BMap.Size(15,3)});
		 var points = new Array();
		for(var i=0; i<shopList.length; i++){
				if(points[shopList[i].point.lat+""+shopList[i].point.lng]){
				 shopList[i].point.lat= shopList[i].point.lat +  0.002;
				 //shopList[i].point.lng=shopList[i].point.lng + 0.002;
				}
				points[shopList[i].point.lat+""+shopList[i].point.lng] = true;
				var infoWindow = new BMap.InfoWindow(""); 
				var label=new BMap.Label(i+parseInt(start));
				label.setStyle({
					width:"29px",
					height:"21px",
					paddingTop:"2px",
					border:"none",
					color:"white",
					fontWeight:"bold",
					textAlign:"center",
					background:"url(images/shop_map/map_cursor.png)"
					});
				var markers = new BMap.Marker(shopList[i].point,{icon:myIcon,offset:new BMap.Size(11,-23)});
				markers.txt=shopList[i].content;
				markers.num=i;
				markers.setLabel(label);
				markers.setTitle(shopList[i].name);
				map.addOverlay(markers);
				markers.addEventListener("click", function(){
					var obj=document.getElementsByName("mapList");
					for(var i=0;i<obj.length;i++){
						if(i==this.num) {obj[i].className="cur";}
						else {obj[i].className="";}
					}
					infoWindow.setContent(this.txt);
					this.openInfoWindow(infoWindow);
				});
				markers.addEventListener("infowindowopen", function(){
					var label=this.getLabel();
					label.setStyle({background:"url(images/shop_map/map_cursor_hover.png)"});
					this.setTop(true);
				});
				markers.addEventListener("infowindowclose", function(){
					var label=this.getLabel();
					label.setStyle({background:"url(images/shop_map/map_cursor.png)"});
					this.setTop(false);
				});
				if(i==0){
					infoWindow.setContent(markers.txt);
					markers.openInfoWindow(infoWindow);
				}
		}
	}
}
function showMapShop(num){
	var shops=map.getOverlays();
	var infoWindow = new BMap.InfoWindow(shops[num].txt);
	shops[num].openInfoWindow(infoWindow);
	var obj=document.getElementsByName("mapList");
	for(var i=0;i<obj.length;i++){
		if(i==num-1) {obj[i].className="cur";obj[i].setAttribute('class','cur');}
		else {obj[i].className="";obj[i].setAttribute('class',"");}
	}
}
function viewShopAsMap(flag){
	var mapview = _getId('mapview');
	var mapM=_getId('mapM');
	var listM=_getId('listM');
	var shops=_getId('shops');
	//列表模式
	if(!flag){
		mapview.style.display = "none";
		setTimeout(function(){shops.style.display="block"},10);
		mapM.style.display="";
		listM.style.display="none";
	}
	//地图模式
	else{
		mapview.style.display = "block";
		shops.style.display="none";
		mapM.style.display="none";
		listM.style.display="";
		if(''==_getId('shops_map').innerHTML){
			popShops("shops_map",popshopList,_getId('startNum').innerHTML);
		}
	}
}
function expandMapView(){
	setTimeout(function(){
		var _mv = _getId('shops_map');
		if(''==_mv.style.height)
		  _mv.style.height='0px';
		
		h=parseInt(_mv.style.height);
		if(h<500){
			_mv.style.height=(h+30)+'px';
		  expandMapView();
		}
	},1);	
}

function popmap(mapx,mapy,shopx,shopy){
	if(null!=_getId("pop")){
		_getId("pop").style.display="block";
//		if(""!=_getId("pop_map").innerHTML){
			var point1 = new BMap.Point(mapx, mapy);
			var point2 = new BMap.Point(shopx, shopy);
			popDirection("pop_map", point1, point2,"pop_map_steps");
			_getId("pop_map").style.position="fixed";
//		}
	}
}

function popDirection(popId,p1,p2,stepsId){
	map = new BMap.Map(popId);
	map.addControl(new BMap.NavigationControl());
	var walking = new BMap.WalkingRoute(map, {renderOptions:{ map: map, autoViewport: true}});
	walking.search(p1, p2);
	_getId('tit_small').innerHTML="（您距离本店大约："+Math.ceil(map.getDistance(p1,p2))+"米）";
}
//地图开关
function mapSwitch(open,shade){
	var popMap=_getId("popMap");
	if(open){
		if(shade) creatShade();
		var locateInfo=_getId("locateInfo");
		if(locateInfo&&locateInfo.style.display=="block") popMap.style.top="93px";
		popMap.style.visibility="visible";
	}
	else{
		if(shade) showObj(document.getElementById("backShade"), 0, -5);
		popMap.style.visibility="hidden";
	}
}
//取marker下地址名
//function getMarkerAddress(){
//  var latlng = marker.getPosition();
//  if(null==geocoder){
//	  geocoder = new BMap.Geocoder();
//  }
//  geocoder.getLocation(latlng, function(results) {
//      if (results) {
//    	  alert(results.address);
//      } else {
//        alert("No results found");
//      }
//  });
//}
//初始化（map样式，Geocoder）
function initialize() {
	map = new BMap.Map("map_canvas");
	var opts = {anchor: BMAP_ANCHOR_TOP_RIGHT, offset: new BMap.Size(10, 10)};
	map.addControl(new BMap.NavigationControl(opts));
	map.addControl(new BMap.ScaleControl()); 
	map.addControl(new BMap.OverviewMapControl()); 
	
	map.enableScrollWheelZoom();//允许滚轮
	map.setDefaultCursor("move");
	var point = new BMap.Point(116.404, 39.915);
	
	map.clearOverlays();
	marker = new BMap.Marker(point);
	map.addOverlay(marker);
	marker.enableDragging(true); // 设置标注可拖拽
	marker.setTitle("请拖跩到您的位置");
	marker.addEventListener("dragging", function(){   
		map.removeOverlay(mapLabel);
		mapLabel=null;
	});
	marker.addEventListener("dragend", function(){   
	    addLabel(false);
	});
	map.addEventListener("zoomend", function(){
		map.setCenter(marker.getPosition());
		addLabel(true);
	});
	map.addEventListener("moveend", function(){
		var bounds=map.getBounds();
		if(!bounds.containsPoint(marker.getPosition())){
			if(promptlabel!=null) {
				map.removeOverlay(promptlabel);
				promptlabel=null;
			}
			var opts = {position:map.getCenter(),offset: new BMap.Size(-8,-175)};
			promptlabel=new BMap.Label("<a href='javascript:putMarkerInMapCenter()' style='color:#fff'>点击将标注移到中间</a>",opts);
			promptlabel.setStyle({
				padding:"3px 10px",
				border:"none",
				background: "#FAA116",
				fontSize : "14px"
				});
			map.addOverlay(promptlabel);
		}
		else if(promptlabel!=null){
			map.removeOverlay(promptlabel);
			promptlabel=null;
		}
	});
	var left=(document.documentElement.clientWidth)/2;
	_getId("popMap").style.left=left+"px";
}
//打开地图
function openMap(mapx,mapy){
	mapType="search";
	if("您当前的小区、楼栋、大学名称.如：清华大学东门"!=_getId("txt_map").value)
		_getId("address_input").value=_getId("txt_map").value.replace("周边...","");
	if(mapx==0||mapy==0){
		moveMapCenterToInputAddress();
	}
	else{
		var point = new BMap.Point(mapx,mapy);
		marker.setPosition(point);
		map.centerAndZoom(point, 15);
		addLabel();
	}
	mapSwitch(true,false);
}
function editMap(address,id,x,y){
	mapType="edit";
	_getId("address_input").value=address;
	var point = new BMap.Point(x,y);
	marker.setPosition(point);
	map.centerAndZoom(point, 15);
	mapSwitch(true,true);
}
//新增地图
function newMap(address){
	mapType="new";
	if("您当前的小区、楼栋、大学名称.如：清华大学东门"!=address)
		_getId("address_input").value=address.replace("周边...","");
	mapx=_getId("mapx").value;
	mapy=_getId("mapy").value;
	if(mapx==""||mapy==""){
		moveMapCenterToInputAddress();
	}
	else{
		var point = new BMap.Point(mapx,mapy);
		marker.setPosition(point);
		map.centerAndZoom(point, 15);
	}
	var top=document.documentElement.scrollTop==0?document.body.scrollTop:document.documentElement.scrollTop;
	var t = (document.documentElement.clientHeight) / 2+top;
	_getId("popMap").style.top=t+"px";
	mapSwitch(true,true);
}
//添加label
function addLabel(first){
	if(promptlabel!=null){
		map.removeOverlay(promptlabel);
		promptlabel=null;
	}
	if(mapLabel==null){
		var opts = {position:marker.getPosition(), offset: new BMap.Size(-278,-248)};
		var str=mapType=="search"?"搜附近餐馆":"确定";
		if(first) mapLabel=new BMap.Label("<div class='location_popbox'><div class='location_main_box'><div class='location_main'><span>请拖动图标到您的位置</span></div></div></div>",opts);
		else mapLabel=new BMap.Label("<div class='location_popbox'><div class='location_main_box'><div class='location_main'><span>使用该地址？</span><a class='green_btn' href='javascript:getMarkerPoint()'><span>"+str+"</span></a></div></div></div>",opts);
//		mapLabel=new BMap.Label("<a href='javascript:getMarkerPoint()' style='font-size:14px;text-decoration:underline;height:15px;'>"+str+"</a>",opts);
		mapLabel.setStyle({
//			width:"132px",
//			height:"26px",
//			lineHeight:"26px",
//			textAlign:"center",
//			padding:"0px 0px 8px",
			fontSize:"13px",
			cursor:"default",
			border:"none",
			background:"url('images/location_bg.png') no-repeat 0 0"
//			background: "url('images/map_btn_bg.gif') no-repeat transparent"
			});
		map.addOverlay(mapLabel);
	}
	else{
		mapLabel.setPosition(marker.getPosition());
	}
}
//将marker移动到地图中心位置
function putMarkerInMapCenter(){
marker.setPosition(map.getCenter());
addLabel(false);
}
//移动map到“address_input”输入的地址处
function moveMapCenterToInputAddress(){
var address = defaultCity+" "+__getInputAddress();
var options = {
onSearchComplete: function(results){
  if (local.getStatus() == BMAP_STATUS_SUCCESS){
    // 判断状态是否正确
	  if(null!=results.getPoi(0)){
		  marker.setPosition(results.getPoi(0).point);
		  map.centerAndZoom(results.getPoi(0).point, 15);
		  map.removeOverlay(mapLabel);
		  mapLabel=null;
		  addLabel(false);
	  }else{
	    alert("Geocode was not successful.");
	  }
    //alert(ah[0]);
    //document.getElementById("log").innerHTML = s.join("<br/>");
  }
}
};
var local = new BMap.LocalSearch(map, options);
local.search(address);
addLabel(true);
}
//取“address_input”输入内容
function __getInputAddress(){
return _getId("address_input").value;
}
//使用该地址
function getMarkerPoint(){
	var latlng = marker.getPosition();
	switch(mapType){
	case "search":
		var url="module/common/search.php?mapx="+latlng.lng+"&mapy="+latlng.lat+"&cityId=1&address="+__getInputAddress();
		self.location.href=url;
		break;
	case "new":
		_getId("mapx").value=latlng.lng;
		_getId("mapy").value=latlng.lat;
		_getId("selectAdd").innerHTML=__getInputAddress();
		_getId("selectAdd").style.display="block";
		_getId("mapAddress").value=__getInputAddress();
		mapSwitch(false,true)
		break;
	case "edit":
		var url="module/customer/myAddress.php?mapx="+latlng.lng+"&mapy="+latlng.lat+"&cityId=1&address="+__getInputAddress();
		submitUrl('module/customer/myAddress.php?');
		mapSwitch(false,true)
		break;
	default:
		var url="module/common/search.php?mapx="+latlng.lng+"&mapy="+latlng.lat+"&cityId=1&address="+__getInputAddress();
		self.location.href=url;
	}
}
//打开关闭定位错误提示
function locateInfo(open){
	var obj=_getId("locateInfo");
	if(obj==null) return;
	if(open){
//		_getId("popMap").style.top="173px";
		obj.style.display="block";
		showInfo(obj);
		setTimeout( function() {locateInfo(false);},30000);
	}
	else{
		var popMap=_getId("popMap");
		if(popMap!=null) popMap.style.top="50px";
		closeInfo(obj);
	}
	function showInfo(obj){
		if(obj.offsetHeight<34){
			obj.style.height=(obj.offsetHeight+2)+"px";
			setTimeout( function() {showInfo(obj)},10);
		}
	}
	function closeInfo(obj){
		if(obj.offsetHeight>0){
			obj.style.height=(obj.offsetHeight-2)+"px";
			if(obj.offsetHeight==0) obj.style.display="none";
			setTimeout( function() {closeInfo(obj)},5);
		}
	}
	
}
