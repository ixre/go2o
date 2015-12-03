//j6.__WORKPATH__ = '/assets/js/plugin/';

//
//j6.extend({
//    tab: {
//        check:function(){
//            if(window.parent.FwTab){
//                return true;
//            }else{
//                alert('不支持此功能');
//                return false;
//            }
//        },
//        open: function (tabTitle, url, closeable) {
//            if (this.check()) {
//                window.parent.FwTab.show(tabTitle, url,closeable);
//            }
//        },
//        open2: function (tabTitle, url, icon, closeable) {
//            if (this.check()) {
//                window.parent.FwTab.show(tabTitle, url,closeable);
//            }
//        },
//        close: function (title, call) {
//            if (this.check()) {
//                window.parent.FwTab.close(title, call);
//            }
//        },
//        closeCurrent:function(call) {
//            window.parent.FwTab.close();
//        },
//        closeAndRefresh:function(title) {
//            if (this.check()) {
//                var win = window.parent.FwTab.getWindow(title);
//                if(win && win.refresh){
//                    win.refresh();
//                }
//                window.parent.FwTab.close();
//            }
//        }
//    }
//});

/*
 j6.extend({
 repeater: function (ele, url, data, format, loaded) {
 j6.lazyRun(function() {
 var dataLoader = j6.dataLoader(ele, {
 url: url,
 data: data,
 loaded: function (json) {
 ele.innerHTML = '';
 var html = '<ul>';
 for (var i in json.rows) {
 html += '<li>' + j6.template(format, json.rows[i]) + '</li>';
 }
 html += '</ul>';
 ele.innerHTML = html;
 if (loaded) {
 loaded(json);
 }
 }
 });
 });
 },
 completion: function (ele, url, loadCallback, selectCallback,minLen) {
 j6.lazyRun(function () {
 j6.autoCompletion(ele, url, loadCallback, selectCallback, minLen);
 });
 }
 });

 */


j6.extend({
    float:function(val){
        return parseFloat(val);
    }
})
