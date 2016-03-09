$(document).ready(function(e) {
	$(window).resize(function(){
		$('.pxui-slide-ad').slide({srcProperty:'src2'});
	}).resize();
	$('#js-tab-style').delegate('a','click',function(){
		$(this).parent().parent().find('.pxui-shoes,.pxui-show-more').hide();
		var next = $('#js-home-tab-'+$(this).index('')).show().next();
		if(next.hasClass('pxui-show-more')){
			next.show();
			if(!next.prev().data('isone')){
				next.prev().data('isone',true);
				next.find('a').click();
			}
		}
	});
	$('#js-home-tab-0').data('isone',true);
	$('#js-show-more-btand').click(function(){
		$(this).find('a').html('查看全部品牌 <i></i>');
	});
});