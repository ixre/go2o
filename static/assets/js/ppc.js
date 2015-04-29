/*** PC Global Js **/
$JS.event.add(window,'load',function(){
   var btns = $JS.getElementsByClassName(document.body,'btn');
    $JS.each(btns,function(i,e){
        var _do = e.getAttribute('do');
        if(_do && window.funcs[_do]){
           $JS.event.add(e,'click',window.funcs[_do]);
        }
    });
});

window.funcs = {
    register:function(){
        if($JS.validator.result('regpanel')){
           var d = $JS.json.toObject('regpanel');
           $JS.xhr.jsonPost('/postRegistInfo',d,function(json){
               if(json.result){
                   var returnUrl = $JS.request('return_url');
                   if(returnUrl!=''){
                       location.replace(returnUrl);
                   }else{
                       location.replace('/');
                   }
               }else{
                   //注册其他
                   alert('注册失败!');
               }
           },function(){
                   alert('注册失败!');
           });
       }
    },
    toggleTop:function() {
        var f = (function (t) {
            return function () {
                if (t.className.indexOf('up') == -1) {
                    t.className = 'btn up';
                    t.innerHTML = "显示导航";
                } else {
                    t.className = 'btn';
                    t.innerHTML = "隐藏导航";
                }
            };
        })(this);
        $JS.animation.toggleHeight('top-container', f, 50);
    }
};


