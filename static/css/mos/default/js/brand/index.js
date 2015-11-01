$(document).ready(function(e) {
	$('.pxui-tab a').click(function(){
		var index = $(this).index();
		$('.tab-box').hide().eq(index).show();
		if(index==0){
			$('#js-letter').lazyload({child:'[lettersrc]',srcProperty:'lettersrc'});
			$('#js-style').lazyload({child:'[stylesrc]',srcProperty:'stylesrc',act:'stop'});
		}else{
			$('#js-letter').lazyload({child:'[lettersrc]',srcProperty:'lettersrc',act:'stop'});
			$('#js-style').lazyload({child:'[stylesrc]',srcProperty:'stylesrc'});
		}
	});
	$('#js-letter').lazyload({child:'[lettersrc]',srcProperty:'lettersrc'});
});