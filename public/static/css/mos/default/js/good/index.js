(function($) {
	$.pxSizeTable = function(options) {
		var settings = {
			brand 	: "",
			part	: 5,//栏目
			main	: "",
			sex		: "",
			cid		: null,
			attrs	: {},
			callback:null
		};
		if(options) {
			$.extend(settings, options);
		}
		function success(){
			try{
				//判断当前大分类是否存在尺码对照表
				if(!shoe_size_table_json[settings.part])return false;
				//判断当前主分类是否存在尺码对照表
				if(!shoe_size_table_json[settings.part][settings.main]){
					//当前主分类不存在尺码对照表，那么当大分类为5（鞋子）时候，性别为男（2）的时候读取男鞋尺码对照表，性别为女（3）的时候读取女鞋尺码对照表，性别为中性（4）的时候读取男鞋和女鞋尺码对照表
					if(!/^[12346]$/.test(settings.main))return false;
					for(var __i=0;__i<2;__i++){
						if(__i==1){
							if(settings.part!=5){
								settings.part = 5;
							}else{
								continue;
							}
						}
						if(settings.sex=='中性'){
							if(shoe_size_table_json[settings.part]["4"]){
								shoe_size_table_json = shoe_size_table_json[settings.part]["4"];
								__i = 4;
							}else{
								var array = null;
								if(shoe_size_table_json[settings.part]['2']){
									array = shoe_size_table_json[settings.part]['2'];
								}
								if(shoe_size_table_json[settings.part]['3']){
									if(!shoe_size_table_json){
										array = shoe_size_table_json[settings.part]['3'];
									}else{
										array = array.concat(shoe_size_table_json[settings.part]['3']);
									}
								}
								if(!array){
									continue;
								}
								__i = 4;
								shoe_size_table_json=array;
								array = null;
							}
						}else{
							if(settings.sex=="男"){
								if(!shoe_size_table_json[settings.part]["2"]){
									continue;
								}else{
									shoe_size_table_json = shoe_size_table_json[settings.part]["2"];
									__i = 4;
								}
							}else if(settings.sex=="女"){
								if(!shoe_size_table_json[settings.part]["3"]){
									continue;
								}else{
									shoe_size_table_json = shoe_size_table_json[settings.part]["3"];
									__i = 4;
								}
							}
						}
					}
					if(__i==2){
						return false;
					}
				}else{
					//当前主分类存在尺码对照表，直接读取当前组分类数据
					shoe_size_table_json = shoe_size_table_json[settings.part][settings.main];
				}
				if(settings.cid){
					var _data = [];
					for(var i=0;i<shoe_size_table_json.length;){
						if(!shoe_size_table_json[i].Cid){
							shoe_size_table_json[i].Cid = "0";
						}
						if(shoe_size_table_json[i].Cid=="0"){
							_data.push(shoe_size_table_json[i]);
						}
						if(shoe_size_table_json[i].Cid!=settings.cid){
							shoe_size_table_json.splice(i,1);
						}else{
							i++;
						}
					}
					if(!shoe_size_table_json.length){
						shoe_size_table_json = _data;
					}
					_data = null;
					i = null;
				}
				//性别定位
				if(shoe_size_table_json.length){
					var _table = shoe_size_table_json[0];
					var array=[];
					if(settings.sex=="男"||settings.sex=="女"){
						$.each(shoe_size_table_json,function(ffindex,table){
							if(table.IsShow!="0"&&table.Title.indexOf(settings.sex+"童")!=0){
								if(table.Title.indexOf(settings.sex)==0){
									array.push(table);
									/*return false;*/
								}
							}
						});
					}else{
						$.each(shoe_size_table_json,function(ffindex,table){
							if(table.IsShow!="0"){
								if(table.Title.indexOf(settings.sex)==0){
									array.push(table);
									/*return false;*/
								}
							}
						});
					}
					/*shoe_size_table_json = [];*/
					if(array.length){
						/*shoe_size_table_json.push(array);*/
						shoe_size_table_json = array;
						array = null;
					}/*else{
						shoe_size_table_json.push(_table);
					}*/
				}
				/*属性筛选*/
				var array = [];
				$.each(shoe_size_table_json,function(index,table){
					if(table.IsShow!="0"){
						var isok = true;
						if(table.attrs){
							table.attrs = table.attrs.split('$');
							$.each(table.attrs,function(index,attr){
								attr = attr.split(':');
								if(attr.length==2){
									if(!settings.attrs[attr[0]]){
										isok = false;
										return false;
									}else{
										if(attr[1].indexOf('<')==0){
											if(parseFloat(settings.attrs[attr[0]])>=parseFloat(attr[1].replace('<',''))){
												isok = false;
												return false;
											}
										}else if(attr[1].indexOf('>')==0){
											if(parseFloat(settings.attrs[attr[0]])<=parseFloat(attr[1].replace('>',''))){
												isok = false;
												return false;
												
											}
										}else if(attr[1].indexOf('-')>=0){
											attr[1] = split('-');
											if(parseFloat(settings.attrs[attr[0]])<=parseFloat(attr[1][0])||parseFloat(settings.attrs[attr[0]])>=parseFloat(attr[1][1])){
												isok = false;
												return false;
												
											}
										}else{
											attr[1] = attr[1].split('%');
											var _isok = false;
											$.each(attr[1],function(index,_attr){
												if(settings.attrs[attr[0]]==_attr){
													_isok = true;
												}
											});
											if(!_isok){
												isok = false;
												return false;
											}
										}
									}
								}
							});
						}
						if(isok){
							array.push(table);
						}
					}
				});
				shoe_size_table_json = [];
				if(array.length){
					array.sort(function(a,b){
						return b.attrs.length-a.attrs.length;
					});
					shoe_size_table_json.push(array[0]);
					array = null;
				}
				//生成html代码
				var _html="";
				$.each(shoe_size_table_json,function(ffindex,table){
					if(table.IsShow!="0"){
						var list = [];
						$.each(table.SizeList,function(findex,standard){
							var _table = [];
							_table.push(standard.Title);
							$.each(standard.Sizes,function(index,size){
								for(var i = parseInt(size.Colspan);i>0;i--){
									_table.push(size.Size);
								}
							});
							list.push(_table);
						});
						_html+='<table cellpadding="0" cellspacing="0" border="0" width="100%"><thead><tr><th class="pxui-bg-blue pxui-color-white" colspan="'+list.length+'">'+table.Title+'</th></tr></thead><tbody>';
						var length = list.length;
						var bfd = parseInt(1/length*1000)/100;
						$.each(list[0],function(index){
							_html+='<tr>';
							for(var i=0,len=list.length;i<len;i++){
								if(!index){
									_html+='<th width="'+bfd+'%">'+list[i][index]+'</th>';
								}else{
									_html+='<td width="'+bfd+'%">'+list[i][index]+'</td>';
								}
							}
							_html+='</tr>';
						});
						_html+='</tbody></table>';
					}
				});
				settings.callback(_html);
			}catch (e){
				settings.callback('');
			}
		};
              //  document.write("../img-cdn2.paixie.net/brandsize/"+settings.brand+"/"+settings.brand+".2.0.json");
		$.getScript("../img-cdn2.paixie.net/brandsize/"+settings.brand+"/"+settings.brand+".2.0.json",success);
	};
})($);
$(document).ready(function(e) {
	/*尺码数量*/
	var isShowSelectSize = false;
	(function(){
		var num = 0;
		function numchange(){
			var self = this;
			$('.num select').die('change',numchange);
			num = $(self).val();
			$('.num select').val(num).change();
			$('.num select').live('change',numchange);
		};
		$('.num select').live('change',numchange);
		var $sizeSelects = $("#js-sizes-select select");
                $sizeSelects.val('');
		function setSizeChange(index,stock,obj){
			if(PX_HELP_DATA[5]){
				$('.com-footer').html($('.com-footer').html()+ index);
			}
			if(index<0){
				return;
			}
			if(!$sizeSelects.data('isremove')){
				$sizeSelects.data('isremove',true).find('option:eq(0)').remove();
				if(obj){
					index--;
				}
			}
			$('.sizes a.selected').removeClass('selected');
			$('.sizes').each(function(){
				$(this).find('a:eq('+index+')').addClass('selected');
			});
			$('.js-stock').html(stock);
			var html = '';
			for(var i=1;i<=stock&&i<=20;i++){
				html+='<option value="'+i+'">'+i+'</option>';
			}
			$('.num select').die('change',numchange);
			num = num||1;
			if(num>parseInt(stock)){num = 1;}
			$('.num select').html(html).val(num).change();
			$('.num select').live('change',numchange);
		};
		function sizechange(){
			var option = $(this).find('option:selected');
			if(option.index()<0)return false;
			setSizeChange(option.index(),option.attr('stock'),this);
		};
		$sizeSelects.bind('change',sizechange);
		
		$('.addtocart:not(:disabled)').live('click',function(){
			if($sizeSelects.val()=='0'||$sizeSelects.val()==''){
				isShowSelectSize = true;
				windowScroll();
				var msgbox = $.message({
					html 	: '<div class="good-page"><ul class="goodinfo">'+$('#js-goodinfo').html().replace(/\n/g,' ').replace(/^.*<!--size-message-->/,'').replace(/<!--size-message-end-->.*$/,'')+'</ul></div>',
					title	: '选择尺码',
					height	: 'auto',
					buttons	: [{disabled: true,light: true,text:'  加入购物车  ',class:'addtocart'}]
				});
				$(msgbox.base()).find('.content a').click(function(){
					msgbox.buttons([{light : true,text :'  加入购物车  ',class:'addtocart'}]);
				});
				$(msgbox.base()).find('h3 a').click(function(){
					isShowSelectSize = false;
					windowScroll();
				});
			}else{

                                var size_str = $sizeSelects.val();
                                var size_arr = size_str.split('_');
                                var size = size_arr[0];
                                var size_id = size_arr[1];
                                if(num < 1){
                                    alert('数量不能为零！');
                                    return false;
                                }
                                var url = '';
                                if(typeof(tuan_id) == 'undefined' || tuan_id == ''|| tuan_id == null){
                                   url =  'cart/ajax@act=addcart&brand_id='+goodInfo.BrandID+'&good_id='+goodInfo.GoodID+'&item_id='+goodInfo.ID+'&size='+size+'&size_id='+size_id+'&num='+num;
                                }else{
                                   url = 'tuan/cart@brand_id='+goodInfo.BrandID+'&good_id='+goodInfo.GoodID+'&item_id='+goodInfo.ID+'&size='+size+'&size_id='+size_id+'&num='+num+'&tuan_id='+tuan_id;
                                }
				window.location.href = url;
                                
                                
			}
		});
		$('.sizes a').live('click',function(){
			setSizeChange($(this).addClass('selected').index(),$(this).attr('stock'));
			$sizeSelects.val($(this).attr('value')).unbind('change',sizechange).change();
			$sizeSelects.bind('change',sizechange);
		});
	})();
	/*尺码数量*/
	/*加入收藏*/
	$('#js-go-favorites').click(function(){
		$.ajax({url:'member/favorites@act=add&item_id='+goodInfo.ID,success:function(data){
				try{
                                        data = $.parseJSON(data);
					if(data.IsSuccess){
						alert('收藏成功！');
					}else{
						if(data.Message=='您还没登录，是否马上登录？登录后需重新收藏！'){
							if(confirm(data.Message)){
								window.location.href = 'login/@url='+encodeURI(window.location.href);
							}
						}else{
							alert(data.Message);
						}
					}
				}catch (e){
					alert('收藏失败！');
				}
			},error:function(){
				alert('链接服务器失败，请稍后再试！');
			}
		});
	});
	/*点击查看图文详情 */
	$("#js-show-img").click(function(){
//		var r=confirm("图文详情的图片较多，可能会消耗比较多的流量。\n是否继续访问？")
//		if (r==true){
         window.location.href = 'product/imgshow@id='+goodInfo.ID;
//
//		} else {
//
//		}
	});
	/*显示尺码对照表*/
	$("#js-show-size").click(function(){
		if($(this).data('open')){
			$(this).data('open',false).html('点击查看&nbsp;&nbsp;<i class="arrow-right"></i>').prev().hide();
		}else{
			if($(this).data('load')){
				$(this).data('open',true).html('点击收起&nbsp;&nbsp;<i class="arrow-top"></i>').prev().show();
			}else{
				var self = this;
				$(this).data('open',true).data('load',true).html('点击收起&nbsp;&nbsp;<i class="arrow-top"></i>').hide().prev().show().html('正在加载尺码对照表...');
				$.pxSizeTable({
					brand	: branddir,
					part	: goodInfo.part,
					main	: goodInfo.MainID,
					sex		: goodInfo.Sex,
					cid		: goodInfo.StyleID,
					attrs	: goodInfo.attrs,
					callback:function(html){
						if(html){
							$(self).prev().html(html);
						}else{
							$(self).prev().html('无尺码对照表');
						}
					}	
				});
			}
		}
	});
	/*显示评论*/
	$("#js-show-comment").click(function(){
		if($(this).data('open')){
			$(this).data('open',false).html('点击查看&nbsp;&nbsp;<i class="arrow-right"></i>').prev().hide();
		}else{
			if($(this).data('load')){
				$(this).data('open',true).html('点击收起&nbsp;&nbsp;<i class="arrow-top"></i>').prev().show();
			}else{
				var self = this;
				$(this).data('open',true).data('load',true).html('点击收起&nbsp;&nbsp;<i class="arrow-top"></i>').prev().show();
				$('#js-commentlist').next().find('a').click();
			}
		}
	});
	$("#js-comment-list").delegate('.img60','click',function(){
		var self = this;
		var img = $(self).find('img');
		var src = img.attr('maxsrc');
		var msgbox = $.message({
			html 	: '<img style="max-width:100%;display: block;margin: auto;" src="'+src+'"/>',
			title	: '查看评论图片',
			height	: 'auto',
			buttons	: [
				{
					disabled:!$(self).prev().length,
					light	:true,
					text	:'  上一张  ',
					click	:function(){
						msgbox.close();
						$(self).prev().click();
					}
				},
				{
					light	:false,
					text	:'  关 闭  ',
					click	:function(){
						msgbox.close();
					}
				},
				{
					disabled:!$(self).next().length||$(self).next().hasClass('show-more'),
					light	:true,
					text	:'  下一张  ',
					click	:function(){
						msgbox.close();
						$(self).next().click();
					}
				}
			]
		});
	}).delegate('.show-more','click',function(){
		if($(this).data('open')){
			$(this).data('open',false).html('显示更多').parent().find('a.hide').hide();
		}else{
			$(this).data('open',true).html('收起更多').parent().find('a:hidden').css('display','inline-block').addClass('hide');
		}
	});
	/*显示更多属性*/
	$('#js-show-all-attrs').click(function(){
		if(!$(this).data('isopen')){
			$(this).data('isopen',true).html('收起完整属性').parent().parent().find('li:hidden').css('display','inline').addClass('hide');
		}else{
			$(this).data('isopen',false).html('查看完整属性').parent().parent().find('li.hide').hide();
		}
	});
	/*团购倒计时*/
	if($('#js-tuan-time').length){
		var second = parseInt($('#js-tuan-time').attr('second'));
		var _setInterval = setInterval(function(){
			if(--second<=0){
				clearInterval(_setInterval);
				window.location.reload();
			}else{
				$('#js-day').html(parseInt(second/(24*60*60)));
				$('#js-hour').html(parseInt(second/(60*60)%24));
				$('#js-minute').html(parseInt(second/60%60));
				$('#js-second').html(parseInt(second%60));
			}
		},1000);
	}
	/*幻灯片*/
	$(window).resize(function(){
		$('.pxui-slide-ad').slide({srcProperty:'src2'});
	}).resize();
	function windowScroll(){
		if(isDown)return;
		if(getScroll().t>$('#js-attrs-title').offset().top&&!isShowSelectSize){
			$('body').css('padding-bottom',$('#js-fixed-add-to-cart').show().height());
		}else{
			$('#js-fixed-add-to-cart').hide();
			$('body').css('padding-bottom',0);
		}
	};
	$(window).scroll(function(){
		windowScroll();
	});
	windowScroll();
});