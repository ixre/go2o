function getCartNums()
{
   var url = 'ajax@act=cartnum';
   $.ajax({url:url,type:'post',
	  error:function(a){
                    alert('链接服务器失败！');
	  },success:function(data){
		   $('#header-cart-num').html(data.toString());
	  }
   });         
}

$(document).ready(function(e) {
	function total(){
		var num = 0;
		var price = 0;
		$('.goodlist li').each(function(){
			if($(this).find('input:checked').length){
                                var j = $(this).find('select').val();
                                j = parseInt(j);
				num += j;
				price+=parseInt($(this).find('select').attr('price'))*$(this).find('select').val();
			}
		});
		$('#js-total').text(price.toFixed(2));
		$('#js-num').text(num);
	};
	$('h2 input').change(function(){
		$(this).parent().next().find('input[type="checkbox"]').prop('checked',$(this).prop('checked')).filter(':disabled').prop('checked',false);
		total();
	});
	$('ul input').change(function(){
		if($(this).prop('checked')){
			var checkbox = $(this).parent().parent().find('input[type="checkbox"]');
			if(checkbox.length==checkbox.filter(':checked').length+checkbox.filter(':disabled').length){
				$(this).parent().parent().prev().find('input[type="checkbox"]').prop('checked',true);
			}
		}else{
			$(this).parent().parent().prev().find('input[type="checkbox"]').prop('checked',false);
		}
		total();
	});
	$('select').change(function(){
		$(this).parents('li').find('input[type="checkbox"]').prop('checked',true);
                var cart_num = parseInt($(this).val());
                var _select = $(this);
         	$.ajax({url:'ajax',type:"POST",data:{
			brand_id:_select.attr('brandid'),
			item_id :_select.attr('itemid'),
			num     :cart_num,
			act     :'goodnum'
		}});
		total();
                getCartNums();
	});
	function remove(){
		if($(this).siblings().length){
			$(this).remove();
		}else{
			var p = $(this).parent();
			if(p.siblings('ul').length){
				p.prev().remove();
				p.remove();
			}else{
				p.parent().html('<div style="text-align:center;padding: 50px 0;font-size: 16px;">您当前购物车空荡荡的，赶快去添加吧！<br /> <a href="../default.htm">返回首页</a></div>');
			}
		}
		var _select = $(this).find('select');
                //document.write('ajax@act=delcart&brand_id='+_select.attr('brandid')+'&item_id='+_select.attr('itemid'));
		$.ajax({url:'ajax',type:"POST",data:{
			brand_id:_select.attr('brandid'),
			item_id :_select.attr('itemid'),
			num:0,
			act:'delcart'
		}});
		total();
	};
	$('a.del').click(function(){
		var p = $(this).parents('li');
		if(p.hasClass('end')){
			return remove.call(p);
		}
		if(confirm('你确定要移除商品吗？')){
                    
			remove.call(p);
                        getCartNums();
                        
		}
	});
	total();
	$('#js-form').submit(function(){
                var boo = false;
                $('.goodlist li').each(function(){
			if($(this).find('input:checked').length){
                                boo = true;
                        }
		});
                if(!boo)
                {
                    alert('请先勾选商品再提交结算！');
                    return false;
                }
		$('.goodlist li').each(function(){
			if(!$(this).find('input:checked').length){
				$(this).find('input,select').prop('disabled',true);
			}
		});

		$('#js-form').submit(function(){return false;});
	});
        // 全部选择
        $('#js-all').click(function(){
            $('.goodlist li').each(function(){
                $(this).find('input[type="checkbox"]').prop('checked',true);
            });
            total();
        });
});