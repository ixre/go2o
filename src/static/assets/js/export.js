
//
// ===========================
//  导出及查询页面通用Js
// ===========================
// 此脚本依赖于：json4html.js

//
//  配置：
//  expr.portal = '这里为导出项的类名';
//
//  expr.checkParams = function (data) {
//    //这里校验参数的准确性,data为Json格式
//    return true;
//  };
//
//
//


if (!window.$JS) { alert('请加载json4html.js文件!'); }

var expr = {
    ele: '',
    appPath: '/',
    hanlderPrefix: '/export/',
    portal: '',
    _getParams: function () {
        return encodeURIComponent($JS.json.toString(this.ele));
    },

    checkParams: function (data) {
        return true;
    },
    getDataUrl: function () {
        if (this.checkParams()) {
            var _appPath = window.appPath || this.appPath;
            return (_appPath == '/' ? '' : _appPath)
                + this.hanlderPrefix
                + 'getExportData?portal=' + this.portal
                + '&params=' + this._getParams();
        }
        return null;
    },
    showExportDialog: function (title, width, height) {
        if (!PW) {
            alert('PW对象为空！');
            return;
        }
        if (!expr.checkParams()) return;
        var _appPath = window.appPath || this.appPath;
        var url = (_appPath == '/' ? '' : _appPath)
            + this.hanlderPrefix
            + 'setup?portal=' + expr.portal
            + '&params=' + expr._getParams();
        PW.getWindow(url, title || '导出数据', width || 400, height || 300);
    },
    search: function (id) {
        $('#' + id).datagrid({ url: expr.getDataUrl() });
    },
    reload: function (id) {
        $('#' + id).datagrid('reload');
    },
    bindTotalView: function (id) {
        if (!expr.checkParams()) return;
        var _appPath = window.appPath || this.appPath;
        var url = (_appPath == '/' ? '' : _appPath)
            + this.hanlderPrefix
            + 'getTotalView?portal=' + expr.portal
            + '&params=' + expr._getParams();

        $JS.xhr.post(url, {}, function (json) {
            $JS.json.bind(id || 'totalView', json);
        });
    }
};
