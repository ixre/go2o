var memberTmpl = '';
var valletTmpl = '';

function show(e) {
    e.className = e.className.replace(' hidden', '');
}




function loadMember(m, Mustache) {
    j6.xhr.jsonPost('/json/member', '', function (data) {
        m.closeTipBox();
        var e = m.getByClass('member-profile');
        e.innerHTML = Mustache.render(memberTmpl,{member:data});
        e = m.getByClass('my-vallet');
        e.innerHTML = Mustache.render(valletTmpl,{member:data});
    });
}

require([
    'uc/main',
    'lib/mustache',
], function (m, Mustache) {
    m.init();
    j6.xhr.filter = null;
    preParseTmpl(m, Mustache);
    loadMember(m, Mustache);
    //loadVallet(m, Mustache);
});

function preParseTmpl(m, Mustache) {
    memberTmpl = m.parseTmpl(m.getByClass('template-member').innerHTML); //获取滚动模板
    valletTmpl = m.parseTmpl(m.getByClass('template-vallet').innerHTML);
    Mustache.parse(memberTmpl);
    Mustache.parse(valletTmpl);
}
