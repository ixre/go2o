function addCart(){	
//<![CDATA[
		(function(){
			// 工具类
			var Tools = {
				$			: fTool$,
				createTr	: fToolCreateTr,
				listen		: fToolsListen,
				getTarget	: fToolsGetTarget
			}

			// 获取节点
			function fTool$(sId){
				return document.getElementById(sId);
			}

			// 构造tr节点
			function fToolCreateTr(sHtml){
				var oTable = document.createElement("div");
				oTable.innerHTML = '<table><tbody><tr>' + sHtml + '</tr></tbody></table>'
				return oTable.getElementsByTagName("tr")[0];
			}
			// 事件监听
			function fToolsListen(oElement, sName, fObserver, bUseCapture){
				if (oElement.addEventListener) {
					oElement.addEventListener(sName, fObserver, bUseCapture);
				}else if(oElement.attachEvent){
					oElement.attachEvent('on' + sName, fObserver);
				}
			}
			// 获取事件源
			function fToolsGetTarget(oEvent){
				var oEvent = oEvent || window.event,
					oTarget = oEvent.target || oEvent.srcElement;
				return {
					event : oEvent,
					target : oTarget
				};
			}

			//////////////////////////////////////////////////////////////////////////////////////////////////////////

			// 业务类定义
			var MenuList = function(){}, OrderList = function(){};

			// 菜单类
			MenuList.prototype = {
				init				: fMenuListInit,
				initData			: fMenuListInitData,
				initEvent			: fMenuListInitEvent,
				handleClick			: fMenuListInitHandleClick,
				handleHover			: fMenuListInitHandleHover,
				getData				: fMenuListInitGetData,
				setStatus			: fMenuListInitSetStatus
			}

			// 方法定义
			// MenuList method
			// 对象初始化
			function fMenuListInit(oParam){
				var oThat = this;

				// 常量
				oThat.UL_ID = "box";
				oThat.LI_ID = "foodList";

				// 数据对象
				oThat.data = {};

				// 选择记录
				oThat.checkStatus = {};

				// 列表对象
				oThat.orderList = oParam.orderList;
				oThat.orderList.menuList = oThat;
				
				// 列表容器
				oThat.container = Tools.$("menu"),

				// 初始化数据
				oThat.initData();

				// 初始化事件
				oThat.initEvent();

				return this;
			}
			
			// 数据初始化
			function fMenuListInitData(){
				var oThat = this;
				var aUL = oThat.container.getElementsByTagName("ul");
				// 分类循环
				for(var i = 0, nLen1 = aUL.length, oUL, aLi, sUlId, oUlData; i < nLen1; i += 1){
					oUL = aUL[i];
					aLi = oUL.getElementsByTagName("li");

					// 按菜品类别构造关联数组
					sUlId = oUL.id.replace(oThat.UL_ID, "");
					oThat.data[sUlId] = oUlData = {
						id : sUlId,
						list : {}
					};
					// 类别下菜品循环
					for(var j = 0, nLen2 = aLi.length, oLi, sLiId; j < nLen2; j += 1){
						oLi = aLi[j];
						// 按菜品id构造关联数组
						sLiId = oLi.id.replace(oThat.LI_ID, "");
						oUlData.list[sLiId] = {
							// id
							id : sLiId,
							// 菜名
							food : oLi.getAttribute("food"),
							// 价格
							price : oLi.getAttribute("price") - 0
						}
						// 删除节点自定义属性
						try{
							oLi.setAttribute("food", null);
							oLi.setAttribute("price", null);
						}catch(e){
							oLi.food = oLi.price = null;
						}
					}
				}
			}

			// 事件初始化
			function fMenuListInitEvent(){
				var oThat = this;
				// 列表点击
				Tools.listen(oThat.container, "click", function(){oThat.handleClick.apply(oThat, arguments);});
//				Tools.listen(oThat.container, "mouseover", function(){oThat.handleHover.apply(oThat, Array.apply(null, arguments).concat([true]));});
//				Tools.listen(oThat.container, "mouseout", function(){oThat.handleHover.apply(oThat, Array.apply(null, arguments).concat([false]));});
			}

			// 点击事件句柄
			function fMenuListInitHandleClick(oEvent){
				var oThat = this;
				var oTemp = oTarget = Tools.getTarget(oEvent).target,
					sNodeName;
				if(!oTarget){
					return false;
				}
				while(oTemp && (oTemp.id !== oThat.container.id)){
					sNodeName = oTemp.nodeName.toLowerCase();
					if(sNodeName != "li"){
						oTemp = oTemp.parentNode;
						continue;
					}
					var sUlId = oTemp.parentNode.id.replace(oThat.UL_ID, ""),
						sId = oTemp.id.replace(oThat.LI_ID, "");
						oData = oThat.getData({main : sUlId, sub : sId});
					if(oData){
						var oTemps=oTemp.getElementsByTagName("span");
						if(oThat.checkStatus[sId]){
							// 去除该物品
							oThat.checkStatus[sId] = false;
							oTemps[0].style.cssText="";
							oTemps[1].style.cssText="top:5px";
							oTemp.onmouseout=function(){};
							oTemp.className="";
							try{
								oThat.orderList.remove(oData.id);
							}catch(e){}
						}
						else if(oTemp.style.background==""){
							// 加入已选列表
							oTemp.style.background="url('../images/added.jpg') no-repeat scroll 0 0 #F0F0F0";
							oTemps[0].style.background="none";
							oTemps[1].style.background="none";
							oTemp.onmouseout=function(){oThat.checkStatus[sId] = true;oTemp.className="cur";oTemp.style.cssText="";};
							try{
								oThat.orderList.add(oData);
							}catch(e){}
						}
					}
					return;
				}
				var sForId = oTarget.getAttribute("forId");
				if(sForId){
					oTarget=oTarget.className.indexOf("fr")<0?oTarget.getElementsByTagName("a")[0]:oTarget;
					var bShow = oTarget.className.indexOf("m_up") < 0, sNewClass = bShow ? "m_up" : "m_down", sOldClass = bShow ? "m_down" : "m_up";
					Tools.$(sForId).style.display = bShow ? "" : "none";
					oTarget.className = oTarget.className.replace(sOldClass, sNewClass);
				}
			}
			// 鼠标滑过事件句柄
			function fMenuListInitHandleHover(oEvent, bIn){
			//var sSign = bIn ? "--over:" : "--out:";
				var oThat = this;
				var oTarget = Tools.getTarget(oEvent).target,
					sNodeName;

				if(!oTarget){
					return false;
				}
				
				var sIndex = oThat.hoverIndex;
				while(oTarget && (oTarget.id !== oThat.container.id)){
					sNodeName = oTarget.nodeName.toLowerCase();
					if(sNodeName != "li"){
						oTarget = oTarget.parentNode;
						continue;
					}
					if(bIn){
						if(sIndex !== oTarget.id){
							var oAddLink = oTarget.getElementsByTagName("span")[0].parentNode;
							oAddLink.style.display = "";
							oThat.hoverIndex = oTarget.id;
						}
					}
					else{
						if(sIndex){
							Tools.$(sIndex).getElementsByTagName("span")[0].parentNode.style.display = "none";
							oThat.hoverIndex = null;
						}
					}
					return true;
				}
			}

			// 查询、获取数据
			function fMenuListInitGetData(oParam){
				var oThat = this;

				var oData = oThat.data;

				try{
					return oData[oParam.main].list[oParam.sub];
				}catch(e){
					return;
				}
			}
			// 设置选中状态
			function fMenuListInitSetStatus(sId, bChecked){
				var oThat = this;
				var oLi = Tools.$(oThat.LI_ID + sId);
				if(oLi){
					bChecked = !!bChecked;
					oThat.checkStatus[sId] = bChecked;
					oLi.className="";
					oLi.onmouseout=null;
					var oLis=oLi.getElementsByTagName("span");
					oLis[0].style.cssText="";
					oLis[1].style.cssText="top:5px";
				}
			}
			//////////////////////////////////////////////////////////////////////////////////////////////////////////

			// 购物车类
			OrderList.prototype = {
				init				: fOrderListInit,
				initEvent			: fOrderListInitEvent,
				handleClick			: fOrderListHandleClick,
				handleChange		: fOrderListHandleChange,
				add					: fOrderListAdd,
				remove				: fOrderListRemove,
				refresh				: fOrderListRefresh,
				getPrice			: fOrderListGetPrice
			}
			// OrderList method
			function fOrderListInit(oParam){
				var oThat = this;
				oThat.TR_ID = "cart_";
				oThat.DEL_ID = "del_";
				// 送餐费
				var cart_fare=Tools.$("cart_fare");
				var fare=cart_fare==null?"0":cart_fare.innerHTML;
				oThat.delivery = parseInt(fare.indexOf("未知")>=0?"0":fare);
				// 订餐数据
				oThat.data = {};
				// 容器节点
				oThat.container = Tools.$("rightcart");
				oThat.cartHidden=Tools.$("cartHidden");
				// 菜品表格节点
				oThat.table = Tools.$("cartTable").getElementsByTagName("tbody")[0];
				// 小计节点
				oThat.subTotleDd = Tools.$("cart_xiaoji");
				// 送餐费节点 todo
				// oThat.deliveryDd = Tools.$("");
				// 总计节点
				oThat.totleSpan = Tools.$("cart_zongjia");
				// 初始化事件
				oThat.initEvent();

				return this;
			}

			// 事件初始化
			function fOrderListInitEvent(){
				var oThat = this;
				Tools.listen(oThat.table, "click", function(){oThat.handleClick.apply(oThat, arguments);});
			}

			// 点击事件句柄
			function fOrderListHandleClick(oEvent){
				var oThat = this;
				var oTarget = Tools.getTarget(oEvent).target;

				if(oTarget && oTarget.id && oTarget.id.indexOf(oThat.DEL_ID) > -1){
					oThat.remove(oTarget.id.replace(oThat.TR_ID + oThat.DEL_ID, ""), true);
				}
			}

			// 菜品数量变化事件句柄
			function fOrderListHandleChange(){
				var oThat = this;
				oThat.refresh();
			}

			// 添加菜品
			function fOrderListAdd(oData){alert('x');
				if(oData.food==null) return;
				var oThat = this;
				oThat.container.style.display = "";
				oThat.cartHidden.style.display = "none";

				var sId = oThat.TR_ID + oData.id;
				if(oThat.data[sId]){
					var oSelect = Tools.$(sId).getElementsByTagName("select")[0];
					if(oSelect.selectedIndex < oSelect.options.length - 1){
						oSelect.selectedIndex = oSelect.selectedIndex + 1;
					}
				}
				else{
					var price=oData.price=="0"?"时价":oData.price;
					var sHtml = '\
						<td class="ttl">' + oData.food + '</td>\
						<td width="40"><select name="itemNum[]" class="cart_o_num"><option value="1" selected="true">1</option><option value="2">2</option><option value="3">3</option><option value="4">4</option><option value="5">5</option><option value="6">6</option><option value="7">7</option><option value="8">8</option></select></td>\
						<td width="30" style="font-size:12px;">' + price + '\<input type="hidden" name="itemId[]" value="'+oData.id+'"/></td>\
						<td width="30"><a id="' + oThat.TR_ID + oThat.DEL_ID + oData.id + '" class="del_btn" href="javascript:void(0);">删除</a></td>\
					',
					oTr = Tools.createTr(sHtml);
					Tools.$("noItemTips").style.display="none";
					oThat.table.appendChild(oTr);

					oTr.id = sId;

					oThat.data[sId] = oData;

					Tools.listen(oTr.getElementsByTagName("select")[0], "change", function(){oThat.handleChange.apply(oThat, arguments);});
				}

				oThat.refresh();
			}

			// 移除菜品
			function fOrderListRemove(sId, bSetMenu){
				var oThat = this;
				oThat.container.style.display = "";
				oThat.cartHidden.style.display = "none";
				var sTrId = oThat.TR_ID + sId;
				var oTr = Tools.$(sTrId);
				if(oTr){
					if(oTr.parentNode.getElementsByTagName("tr").length==1) Tools.$("noItemTips").style.display="";
					oTr.parentNode.removeChild(oTr);
					if(bSetMenu){
						try{
							oThat.menuList.setStatus(sId, false);
						}catch(e){}
					}
					if(oThat.data[sTrId]){
						delete oThat.data[sTrId];
					}
				}
				oThat.refresh();
			}

			// 刷新
			function fOrderListRefresh(){
				var oThat = this;
				var oPrice = oThat.getPrice();
				//oThat.subTotleDd.innerHTML = oPrice.subTotle + "";
				oThat.totleSpan.innerHTML = oPrice.totle + "";
				Tools.$("cartItemNum").innerHTML=oPrice.subValue;
				Tools.$("itemSum").value=oPrice.totle;
			}

			// 计算价格
			function fOrderListGetPrice(){
				var oThat = this;
				var nPrice = 0,nValue=0;
				var aTr = oThat.table.getElementsByTagName("tr");

				for(var i = 0, nLen = aTr.length, oTr; i < nLen; i += 1){
					oTr = aTr[i];
					var value=oTr.getElementsByTagName("select")[0].value-0;
					nValue+=value;
					try{
						nPrice += (oThat.data[oTr.id].price * value);
					}catch(e){}
				}

				return {
					subValue : nValue,
					subTotle : nPrice,
					totle : nPrice > 0 ? (nPrice + oThat.delivery) : nPrice
				}
			}

			//////////////////////////////////////////////////////////////////////////////////////////////////////////
			if(Tools.$("cartTable")==null){
				setTimeout(function(){addCart()},300);
				return;
			}
			// 对象初始化
			var oOrderList = (new OrderList()).init({delivery : 0});

			var oMenuList = (new MenuList()).init({orderList : oOrderList});
		})();
}
	//]]>
