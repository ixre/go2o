$(document).ready(function(e) {
	var issubmit = false;
	$('#js-login-form').submit(function(){
		if(issubmit)return false;
		var username = $('#js-username').val($.trim($('#js-username').val())).val();
		var password = $('#js-password').val($.trim($('#js-password').val())).val();
		if(username==''){
			showError('请输入您的账号。');
			return false;
		}
		if(password==''){
			showError('请输入您的密码。');
			return false;
		}
                try{
		var error = PXVerify.Login(username,password,false,true,function(isok,error){
			if(isok){
				window.location.href = returnurl||$('.com-header-logo').attr('href');
			}else{
				issubmit = false;
				$('#js-login').val('   登  录   ');
				showError(error);
			}
		});
                }catch(e){
                    alert(e.getMessage());
                }
		if(error){
			showError(error);
			return false;
		}else{
			issubmit = true;
			$('#js-login').val('正在登录...');
		}
		return false;
	});
});