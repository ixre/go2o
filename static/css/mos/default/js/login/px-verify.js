//朝首登录注册相关验证
var PXVerify = {
	_trim:function(str){
		return $.trim(str);
	},
	_callback:function(callback,isok,error,data){
		if(callback!=null){
			callback(isok,error,data);
		}
		return error;
	},
	//功能：验证邮箱
	//参数：邮箱地址(string)、是否验证已存在(bool)、验证成功后回调函数(function)
	//		验证成功后回调函数(function)参数：是否成功、相关信息(错误提示信息)、返回原始数据
	//返回：错误信息(string)，正确返回null
	Email:function(email,exist,callback){
		var _this = this;
		email = _this._trim(email);
		if(email==''){
			return '邮箱地址不能为空！';
		};
		if (!/^\w+((-\w+)|(\.\w+))*\@[A-Za-z0-9]+((\.|-)[A-Za-z0-9]+)*\.[A-Za-z0-9]+$/.test(email)) {
			return '邮箱地址格式错误！';
		};
		if(exist){
			$.ajax({
				url:DOMIN.MAIN+'/register?jsoncallback=?',
				type:'GET',
				data:{
					act		:'verify_email',
					phone	:email
				},
				dataType:'json',
				error:function(data){
					return _this._callback(callback,false,'连接服务器失败，请稍后再试！');
				}
				,success:function(data){
					return _this._callback(callback,data.IsSuccess,data.Message,data);
				}
			});
		}
		return null;
	},
	//功能：发送邮箱验证码
	//参数：手机号码(string)、是否验证已绑定(bool)、验证成功后回调函数(function)
	//		验证成功后回调函数(function)参数：是否成功、相关信息(错误提示信息)、返回原始数据
	//返回：错误信息(string)，正确返回null
	SendEmailCode:function(email,exist,callback,url){
		var _this = this;
		email = _this._trim(email);
		if(email==''){
			return '邮箱地址不能为空！';
		};
		if (!/^\w+((-\w+)|(\.\w+))*\@[A-Za-z0-9]+((\.|-)[A-Za-z0-9]+)*\.[A-Za-z0-9]+$/.test(email)) {
			return '邮箱地址格式错误！';
		};
		url = url||DOMIN.MAIN+'/register?jsoncallback=?';
		$.ajax({
			url:url,
			type:'GET',
			data:{
				act		:'get_email_code',
				email	:email,
				forgot	:((exist==true)?1:0)
			},
			dataType:'json',
			error:function(){
				return _this._callback(callback,false,'连接服务器失败，请稍后再试！');
			}
			,success:function(data){
				return _this._callback(callback,data.IsSuccess,data.Message,data);
			}
		});
		return null;
	},
	//功能：验证邮箱验证码
	//参数：邮箱地址(string)、验证码(string)、是否ajax验证(bool)、验证成功后回调函数(function)
	//		验证成功后回调函数(function)参数：是否成功、相关信息(错误提示信息)、返回原始数据
	//返回：错误信息(string)，正确返回null
	EmailCode:function(email,code,ajax,callback){
		var _this = this;
		email = _this._trim(email);
		code = _this._trim(code);
		if(email==''){
			return '邮箱地址不能为空！';
		};
		if (!/^\w+((-\w+)|(\.\w+))*\@[A-Za-z0-9]+((\.|-)[A-Za-z0-9]+)*\.[A-Za-z0-9]+$/.test(email)) {
			return '邮箱地址格式错误！';
		};
		if(code==''){
			return '请输入您的手机验证码。';
		};
		if (!/^[0-9a-zA-Z]{6}$/.test(code)) {
			return '您输入的验证码错误，请重新输入或重新获取验证码。';
		};
		if(ajax){
			$.ajax({
				url:DOMIN.MAIN+'/register?jsoncallback=?',
				type:'GET',
				data:{
					act		:'verify_email_code',
					email	:email,
					code	:code
				},
				dataType:'json',
				error:function(){
					return _this._callback(callback,false,'连接服务器失败，请稍后再试！');
				}
				,success:function(data){
					return _this._callback(callback,data.IsSuccess,data.Message,data);
				}
			});
		}
		return null;
	},
	//功能：绑定邮箱
	//参数：邮箱地址(string)、验证码(string)、绑定成功后回调函数(function)
	//		绑定成功后回调函数(function)参数：是否成功、相关信息(错误提示信息)、返回原始数据
	//返回：错误信息(string)，正确返回null
	BindEmail:function(email,code,callback){
		var _this = this;
		email = _this._trim(email);
		code = _this._trim(code);
		if(email==''){
			return '邮箱地址不能为空！';
		};
		if (!/^\w+((-\w+)|(\.\w+))*\@[A-Za-z0-9]+((\.|-)[A-Za-z0-9]+)*\.[A-Za-z0-9]+$/.test(email)) {
			return '邮箱地址格式错误！';
		};
		if(code==''){
			return '请输入您的手机验证码。';
		};
		if (!/^[0-9a-zA-Z]{6}$/.test(code)) {
			return '您输入的验证码错误，请重新输入或重新获取验证码。';
		};
		$.ajax({
			url:DOMIN.MAIN+'/register?jsoncallback=?',
			type:'GET',
			data:{
				act		:'bind_email',
				email	:email,
				code	:code
			},
			dataType:'json',
			error:function(data){
				return _this._callback(callback,false,'连接服务器失败，请稍后再试！');
			}
			,success:function(data){
				return _this._callback(callback,data.IsSuccess,data.Message,data);
			}
		});
		return null;
	},
	//功能：验证手机
	//参数：手机号码(string)、是否验证已绑定(bool)、验证成功后回调函数(function)
	//		验证成功后回调函数(function)参数：是否成功、相关信息(错误提示信息)、返回原始数据
	//返回：错误信息(string)，正确返回null
	Phone:function(phone,exist,callback){
		var _this = this;
		phone = _this._trim(phone);
		if(phone==''){
			return '手机号码不能为空！';
		};
		if (!/^\d{10,11}$/.test(phone)) {
			return '手机号码格式错误！';
		};
		if(exist){
			$.ajax({
				url:DOMIN.MAIN+'/register?jsoncallback=?',
				type:'GET',
				data:{
					act		:'verify_phone',
					phone	:phone
				},
				dataType:'json',
				error:function(data){
					return _this._callback(callback,false,'连接服务器失败，请稍后再试！');
				}
				,success:function(data){
					return _this._callback(callback,data.IsSuccess,data.Message,data);
				}
			});
		}
		return null;
	},
	ShowImageVerify : function(phone,forgot,callback){
		var box = $('<div class="send-phone-verify"><h3>请输入验证码<a></a></h3><p><span>验证码：</span><input type="text"/><img title="点击刷新" src="'+DOMIN.MAIN+'/register/vercode""></p><div><span></span><a></a></div></div>').appendTo('body').css({"z-index":"5001",top:$(window).scrollTop()+$(window).height()/2-80});
		$.documentMask({z:5000,id:"js-send-phone-verify-bg"});
		var input = box.find("input");
		var error = box.find("div span");
		box.find("img").click(function(){
			$(this).attr('src',$(this).attr('src').split('?')[0]+"?"+(new Date()).valueOf());
		});
		input.focus(function(){error.html('');});
		var ajaxindex = 0;
		box.find("h3 a").click(function(){
			$("#js-send-phone-verify-bg").fadeOut(function(){$(this).remove();});
			box.remove();
			callback({IsSuccess:false,Message:"验证码发送失败，请稍后再试！"});
		});
		box.find("div a").click(function(){
			var code = input.val($.trim(input.val())).val();
			if(code==""){
				error.html('请输入您的手机验证码。');
				return;
			}
			if(!/^.{4}$/.test(code)){
				error.html('您输入的验证码错误，请重新输入或重新获取验证码。');
				return false;
			}
			var _ajaxindex = ++ajaxindex;
			$.ajax({
				url:DOMIN.MAIN+"/register?jsoncallback=?",
				dataType:"json",
				data:{code:code,act:"verifyphoneimg",phone:phone,forgot:forgot},
				error:function(){
					error.html('链接服务器失败，请稍后再试！');
				},
				success:function(data){
					if(data.IsSuccess){
						$("#js-send-phone-verify-bg").fadeOut(function(){$(this).remove();});
						box.remove();
						callback({IsSuccess:true,Message:""});
					}else{
						error.html(data.Message);
					}
				}
			});
		});
	},
	//功能：发送手机验证码
	//参数：手机号码(string)、是否验证已绑定(bool)、验证成功后回调函数(function)
	//		验证成功后回调函数(function)参数：是否成功、相关信息(错误提示信息)、返回原始数据
	//返回：错误信息(string)，正确返回null
	SendPhoneCode:function(phone,exist,callback,url){
		var _this = this;
		phone = _this._trim(phone);
		if(phone==''){
			return '手机号码不能为空！';
		};
		if (!/^\d{10,11}$/.test(phone)) {
			return '手机号码格式错误！';
		};
		$.ajax({url:DOMIN.MAIN+'/register?jsoncallback=?',data:{act:'get_rand_code'},dataType:"json",success:function(data){
			url = url||DOMIN.MAIN+'/register?jsoncallback=?';
			$.ajax({
				url:url,
				type:'GET',
				data:{
					rand	:data.code,
					act		:'get_code',
					phone	:phone,
					forgot	:((exist==true)?1:0)
				},
				dataType:'json',
				error:function(data){
					return _this._callback(callback,false,'连接服务器失败，请稍后再试！');
				}
				,success:function(data){
					if(!data.IsSuccess&&data.Message=="显示图灵验证"){
						PXVerify.ShowImageVerify(phone,((exist==true)?1:0),function(data){
							return _this._callback(callback,data.IsSuccess,data.Message,data);
						});
					}else{
						return _this._callback(callback,data.IsSuccess,data.Message,data);
					}
				}
			});
		},error:function(){
			return _this._callback(callback,false,'连接服务器失败，请稍后再试！');
		}});
		return null;
	},
	//功能：验证手机验证码
	//参数：手机号码(string)、验证码(string)、是否ajax验证(bool)、验证成功后回调函数(function)
	//		验证成功后回调函数(function)参数：是否成功、相关信息(错误提示信息)、返回原始数据
	//返回：错误信息(string)，正确返回null
	PhoneCode:function(phone,code,ajax,callback){
		var _this = this;
		phone = _this._trim(phone);
		code = _this._trim(code);
		if(phone==''){
			return '手机号码不能为空！';
		};
		if (!/^\d{10,11}$/.test(phone)) {
			return '手机号码格式错误！';
		};
		if(code==''){
			return '请输入您的手机验证码。';
		};
		if (!/^\d{6}$/.test(code)) {
			return '您输入的验证码错误，请重新输入或重新获取验证码。';
		};
		if(ajax){
			$.ajax({
				url:DOMIN.MAIN+'/register?jsoncallback=?',
				type:'GET',
				data:{
					act		:'verify_phone_code',
					phone	:phone,
					code	:code
				},
				dataType:'json',
				error:function(data){
					return _this._callback(callback,false,'连接服务器失败，请稍后再试！');
				}
				,success:function(data){
					return _this._callback(callback,data.IsSuccess,data.Message,data);
				}
			});
		}
		return null;
	},
	//功能：绑定手机
	//参数：手机号码(string)、验证码(string)、绑定成功后回调函数(function)
	//		绑定成功后回调函数(function)参数：是否成功、相关信息(错误提示信息)、返回原始数据
	//返回：错误信息(string)，正确返回null
	BindPhone:function(phone,code,callback){
		var _this = this;
		phone = _this._trim(phone);
		code = _this._trim(code);
		if(phone==''){
			return '手机号码不能为空！';
		};
		if (!/^\d{10,11}$/.test(phone)) {
			return '手机号码格式错误！';
		};
		if(code==''){
			return '请输入您的手机验证码。';
		};
		if (!/^\d{6}$/.test(code)) {
			return '您输入的验证码错误，请重新输入或重新获取验证码。';
		};
		$.ajax({
			url:DOMIN.MAIN+'/register?jsoncallback=?',
			type:'GET',
			data:{
				act		:'bind_phone',
				phone	:phone,
				code	:code
			},
			dataType:'json',
			error:function(data){
				return _this._callback(callback,false,'连接服务器失败，请稍后再试！');
			}
			,success:function(data){
				return _this._callback(callback,data.IsSuccess,data.Message,data);
			}
		});
		return null;
	},
	//功能：验证密码
	//参数：密码(string)、验证成功后回调函数(function)
	//		验证成功后回调函数(function)参数：是否成功、相关信息(错误提示信息)、返回原始数据
	//返回：错误信息(string)，正确返回null
	Password:function(password,callback){
		var _this = this;
		password = _this._trim(password);
		if(password==''){
			return '请输入您的密码。';
		};
		if (!/^[A-Za-z0-9]{6,15}$/.test(password)) {
			return '您输入的密码长度应为6-16个字符，请重新输入。';
		};
		return null;
	},
	//功能：验证二次密码
	//参数：密码(string)、二次密码(string)、验证成功后回调函数(function)
	//		验证成功后回调函数(function)参数：是否成功、相关信息(错误提示信息)、返回原始数据
	//返回：错误信息(string)，正确返回null
	Password2:function(password,password2,callback){
		var _this = this;
		password = _this._trim(password);
		password2 = _this._trim(password2);
		if(password2==''){
			return '请输入您的密码。';
		};
		if (!/^[A-Za-z0-9]{6,15}$/.test(password2)) {
			return '您输入的密码长度应为6-16个字符，请重新输入。';
		};
		if(password2!=password){
			return '两次密码不相同！';
		}
		return null;
	},
	//功能：验证图片验证码
	//参数：验证码(string)、是否ajax验证(bool)、验证成功后回调函数(function)
	//		验证成功后回调函数(function)参数：是否成功、相关信息(错误提示信息)、返回原始数据
	//返回：错误信息(string)，正确返回null
	ImageCode:function(code,ajax,callback){
		var _this = this;
		code = _this._trim(code);
		if(code==''){
			return '请输入您的手机验证码。';
		};
		if (!/^[0-9a-zA-Z]{4}$/.test(code)) {
			return '您输入的验证码错误，请重新输入或重新获取验证码。';
		};
		if(ajax){
			
		}
		return null;
	},
	//功能：登录
	//参数：账号名(string)、密码(string)、下次自动登录(bool)、是否ajax验证(bool)、验证成功后回调函数(function)
	//		验证成功后回调函数(function)参数：是否成功、相关信息(错误提示信息)、返回原始数据
	//返回：错误信息(string)，正确返回null
	Login:function(user,password,remember,ajax,callback){
		var _this = this;
		user = _this._trim(user);
		password = _this._trim(password);
		if(user==''){
			return '账号名不能为空！';
		};
		if(password==''){
			return '密码名不能为空！';
		};
		if(ajax){
			$.ajax({
				url:DOMIN.MAIN+'/login?act=ajaxlogin&jsoncallback=?',
				type: 'GET',
				dataType: 'json',
				cache:false,
				data: '&username='+encodeURIComponent(user)+'&password='+password+'&ckusername='+((remember)?'on':'off'), 
				error:function(data){
					return _this._callback(callback,false,'连接服务器失败，请稍后再试！');
				},success: function(data){
					return _this._callback(callback,data.IsSuccess,data.Message,data);
				}
			});
		}
		return null;
	},
	//功能：注册
	//参数：手机/邮箱(string)、密码(string)、验证码(string)、学校ID(string)、模式(string[Phone|Email])、是否ajax验证(bool)、验证成功后回调函数(function)
	//		验证成功后回调函数(function)参数：是否成功、相关信息(错误提示信息)、返回原始数据
	//返回：错误信息(string)，正确返回null
	Register:function(phone,password,code,campus_id,mode,ajax,callback,url){
		var _this = this;
		phone = _this._trim(phone);
		campus_id = _this._trim(campus_id);
		mode = _this._trim(mode);
		password = _this._trim(password);
		code = _this._trim(code);
		var data = {
			phone 		: phone,
			password 	:password,
			campus_id	:campus_id
		};
		var msg = null;
		if(mode=='Phone'){
			data['act'] = 'register_phone';
			data['code'] = code;
			msg = _this.Phone(phone,false);
			if(msg){
				return msg;
			}
		}else{
			data['act'] = 'register_email';
			data['imgcode'] = code;
			msg = _this.Email(phone,false);
			if(msg){
				return msg;
			}
		}
		msg = _this.Password(password,false);
		if(msg){
			return msg;
		}
		if(mode=='Phone'){
			msg = _this.PhoneCode(phone,code,false);
			if(msg){
				return msg;
			}
		}else{
			msg = _this.ImageCode(code,false);
			if(msg){
				return msg;
			}
		}
		if(ajax){
			$.ajax({
				url:url||(DOMIN.MAIN+'/login?act=ajaxlogin&jsoncallback=?'),
				type: 'GET',
				dataType: 'json',
				cache:false,
				data:data, 
				error:function(data){
					return _this._callback(callback,false,'连接服务器失败，请稍后再试！');
				},success: function(data){
					return _this._callback(callback,data.IsSuccess,data.Message,data);
				}
			});
		}
		return null;
	}
};