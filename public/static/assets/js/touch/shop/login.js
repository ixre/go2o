require([
    'shop/main',
    'lib/device_adapter',
    'jr/form_require',
    'lib/sha1'],
    function(m) {
        m.init();
        jr.validator.init();
        window.pc.login = null;
        window.cli.login(jr.request('return_url'));
        if(window.android || window.ios) {
            window.history.go(-1);
        }else{
            document.body.className='';
        }
        var tip = jr.$('tip');
        document.body.onkeydown = function (e) {
            var event = window.event || e;
            if (event.keyCode == 13) {
                subLogin();
            }
        };
        jr.$('btn_login').onclick=subLogin;
    }
);

function subLogin() {
    var data = jr.json.toObject('form1');
    if (jr.validator.validate('form1')){
        data.pwd = sha1(data.pwd);
        jr.json.bind('form1',data);

        jr.xhr.jsonPost(location.href, data, function (json) {
            if (json.result) {
                window.parent.location.replace(decodeURIComponent(jr.request('return_url')||'/'));
            }else {
                if (json.message.indexOf('验证码') != -1) {
                    refreshImage();
                }
                tip.className= 'tip-panel';
                tip.innerHTML = '<span style="color:red">' + json.message + '</span>';
            }
        }, function (x) {
            tip.className= 'tip-panel';
            tip.innerHTML = '<span style="color:red">登陆服务器失败请重试!</span>';
        });
    }
}
