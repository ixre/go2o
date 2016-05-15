var t1,t2;
var nDiv,mDiv;
var cacheObject = {};
require([
    'shop/main',
    "lib/mustache"],function(m,Mustache){
    m.init()
    j6.xhr.filter = null;
    nDiv = m.getByClass('l');
    mDiv = m.getByClass('m');
    t1 = m.parseTmpl(m.getByClass('template1').innerHTML);
    t2 = m.parseTmpl(m.getByClass('template2').innerHTML);
    Mustache.parse(t1);
    Mustache.parse(t2);
    j6.xhr.jsonPost('CategoryJson',{parent_id:0},function(json){
       nDiv.innerHTML = Mustache.render(t1,{list:json});
        var iLs = nDiv.getElementsByClassName('i');
        bindEvents(Mustache,iLs,nDiv);
        if(iLs.length > 0){
            iLs[0].onclick();
        }
        m.closeTipBox(); //关闭提示框
    });
});

function bindEvents(Mustache,iLs,nDiv) {
    j6.each(iLs, function (i, e) {
        e.onclick = (function (list) {
            return function () {
                for (var i = 0; i < list.length; i++) {
                    list[i].className = list[i] == this ? 'i curr' : 'i';
                }
                var data = this.getAttribute('data');
                fillChildData(Mustache,data,mDiv,t2);
            };
        })(iLs);
    });
}

function fillChildData(Mustache,data,m,t){
    var bd = function(d){
        m.innerHTML = Mustache.render(t,{list:d});
    };

    if(cacheObject[data] != null){
        bd(cacheObject[data]);
    }else {
        j6.xhr.jsonPost('CategoryJson', {parent_id: data}, function (json) {
            cacheObject[data] = json;
            bd(json);
        });
    }
}

