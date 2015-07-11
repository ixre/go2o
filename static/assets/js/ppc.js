/*** PC Global Js **/
j6.event.add(window,'load',function(){
   var btns = j6.dom.getsByClass(document.body,'btn');
    j6.each(btns,function(i,e){
        var _do = e.getAttribute('do');
        if(_do && window.funcs[_do]){
           j6.event.add(e,'click',window.funcs[_do]);
        }
    });
});

window.funcs = {
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
        j6.animation.toggleHeight('top-container', f, 50);
    }
};


