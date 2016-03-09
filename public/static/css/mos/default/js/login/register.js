$(document).ready(function(e) {
	$('#js-show-psw').click(function(){
		if($(this).attr('checked')){
			$('#js-password').attr('type','text');
			/*$('#js-password-2').attr('type','text');*/
		}else{
			$('#js-password').attr('type','password');
			/*$('#js-password-2').attr('type','password');*/
		}
	});
	var phoneBox = $('.phone-box');
	var mode = 'Email';
	var isloaded = false;
	$('#js-username').change(function(){
		var username = $(this).val($.trim($(this).val())).val();
		if(/\d{10,11}/.test(username)){
			PXVerify.Phone(username,true,function(isok,error){
				if(isok){
					phoneBox.show();
				}else{
					if(error=='该手机已存在，请登录！'){
						showError('对不起！账号：'+username+'已被注册！您可以使用该账号<a class="a" href="../login/@username='+username+'">直接登录</a>，或者使用其他账号注册。');
					}
				}
			});
			mode = 'Phone';
		}else{
			var error = PXVerify.Email(username,true,function(isok,error){
				if(!isok){
					if(error=='该邮箱已存在，请登录！'){
						showError('对不起！账号：'+username+'已被注册！您可以使用该账号<a class="a" href="../login/@username='+username+'">直接登录</a>，或者使用其他账号注册。');
					}else{
						showError(error);
					}
				}
			});
			if(error){
				if(error=='邮箱地址格式错误！'){
					showError('您输入的账号格式错误，请重新输入。');
				}else if(error=='邮箱地址不能为空！'&&!isloaded){
					
				}else{
					showError(error);
				}
			}
			phoneBox.hide();
			mode = 'Email';
		}
	}).change();
	isloaded = true;
	$('#js-get-phone').click(function(){
		var self = this;
		if($(self).data('issending'))return;
		$(self).data('issending',true);
		$(self).data('issend',false);
		$(self).prop('disabled',true).val('正在发送验证码...');
		$('#js-phone-code-tip').show();
		var error = PXVerify.SendPhoneCode($.trim($('#js-username').val()),true,function(isok,error){
			if(isok){
				$(self).data('issend',true);
				var time = 60;
				$(self).val(time+'后重可新获取');
				$(self).data('sendcode-setInterval',setInterval(function(){
					if(time==0){
						clearInterval($(self).data('sendcode-setInterval'));
						$(self).prop('disabled',false).val('发送验证码').data('issending',false);
						return;
					}
					$(self).val((--time)+'后重可新获取');
				},1000));
			}else{
				alert(error);
				$(self).prop('disabled',false).val('发送验证码').data('issending',false);
			}
		});
		if(error){
			alert(error);
			$(self).prop('disabled',false).val('发送验证码').data('issending',false);
		}
	});
	
	var issubmit = false;
	$('#js-register').click(function(){
		if(issubmit)return false;
		var phone = $('#js-username').val($.trim($('#js-username').val())).val();
		var password = $('#js-password').val($.trim($('#js-password').val())).val();
		/*var password2 = $('#js-password-2').val($.trim($('#js-password-2').val())).val();*/
		var code = $('#js-code').val($.trim($('#js-code').val())).val();
                if(typeof(phone) == 'undefined' || phone == ''){
                    showError('请输入您的账号。');
                    return false;
                }
		if(mode=='Email'){
			code = '1111';
			var error = PXVerify.Email(phone);
			if(error){
				if(error=='邮箱地址格式错误！'){
					showError('您输入的账号格式错误，请重新输入。');
				}else{
					showError(error);
				}
				return false;
			}
		}
		var error = PXVerify.Password(password);
		if(error){
			if(error=='密码长度应为6-16个字符！'){
				showError('您输入的密码长度应为6-16个字符，请重新输入。');
			}else{
				showError(error);
			}
			return false;
		}
		/*error = PXVerify.Password2(password,password2);
		if(error){
			showError(error);
			return false;
		}*/
		if(mode=='Phone'){
			error = PXVerify.PhoneCode(phone,code);
			if(error){
				showError(error);
				return false;
			}
		}
               
		var error = PXVerify.Register(phone, password, code, '', mode, true, function(isok,error){
                        if(isok){
                                 url = '../msg/@msg='+error+'&url='+(returnurl||$('.com-header-logo').attr('href'));
                                 window.location.href  = url;
				//window.location.href = returnurl||$('.com-header-logo').attr('href');
			}else{
				if(error=='该邮箱已存在，请登录！'){
					showError('对不起！账号：'+phone+'已被注册！您可以使用该账号<a class="a" href="../login/@username='+phone+'">直接登录</a>，或者使用其他账号注册。');
				}else{
					showError(error);
				}
				issubmit = false;
				$('#js-register').val('   注 册   ');
			}
		}, DOMIN.MAIN+'/register?jsoncallback=?');
                
		if(error){
			showError(error);
			return false;
		}else{
			issubmit = true;
			$('#js-register').val('正在注册...');
		}
	});
});