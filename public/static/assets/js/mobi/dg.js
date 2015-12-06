//
//文件：数据表格
//版本: 1.0
//时间：2014-04-01
//


function datagrid(ele, config) {
    this.panel = ele.nodeName ? ele : j6.$(ele);
    this.columns = config.columns;
    //Id域
    this.idField = config.idField || "id";
    this.data_url = config.url;
    this.data = config.data;

    //加载完成后触发
    this.onLoaded = config.loaded;

    //列的长度
    //this.columns_width = [];

    this.loadbox = null;
    this.gridView = null;



    this.loading = function () {
        //初始化高度
        if (this.gridView.offsetHeight == 0) {
            var header_height = this.gridView.previousSibling.offsetHeight;
            var gridview_height = this.panel.offsetHeight - this.gridView.previousSibling.offsetHeight;

            this.gridView.style.cssText = this.gridView.style.cssText
                .replace(/(\s*)height:[^;]+;/ig, ' height:' + (gridview_height > header_height ? gridview_height + 'px;' : 'auto'));


            var ldLft = Math.ceil((this.gridView.clientWidth - this.loadbox.offsetWidth) / 2);
            var ldTop = Math.ceil((this.gridView.clientHeight - this.loadbox.offsetHeight) / 2);

            this.loadbox.style.cssText = this.loadbox.style.cssText
                .replace(/(;\s*)*left:[^;]+;([\s\S]*)(\s)top:([^;]+)/ig,
                '$1left:' + ldLft + 'px;$2 top:'
                + (ldTop < 0 ? -ldTop : ldTop) + 'px');

        }

        this.loadbox.style.display = '';
    };

    this._initLayout = function () {
        var html = '';
        if (this.columns && this.columns.length != 0) {

            //添加头部
            html += '<div class="ui-datagrid-header"><table width="100%" cellspacing="0" cellpadding="0"><tr>';
            for (var i in this.columns) {

                // this.columns_width.push(this.columns[i].width);

                html += '<td'
                    + (i == 0 ? ' class="first"' : '')
                    + (this.columns[i].align ? ' align="' + this.columns[i].align + '"' : '')
                    + (this.columns[i].width ? ' width="' + this.columns[i].width + '"' : '')
                    + '><span class="ui-datagrid-header-title">'
                    + this.columns[i].title
                    + '</span></td>';
            }

            html += '</tr></table></div>';

            //添加内容页
            html += '<div class="ui-datagrid-view" style="position:relative;overflow:auto;height:0;">'
                + '<div class="loading" style="position: absolute; display: inline-block; left:0; top:0;">加载中...</div>'
                + '<div class="view"></div>'
                + '</div>';

        }
        this.panel.innerHTML = html;

        this.gridView = (this.panel.getElementsByClassName
            ? this.panel.getElementsByClassName('ui-datagrid-view')
            : j6.dom.getsByClass(this.panel, 'ui-datagrid-view'))[0];

        this.loadbox = this.gridView.getElementsByTagName('DIV')[0];
    };


    this._fill_data = function (data) {
        if (!data) return;

        var item;
        var col;
        var val;
        var html = '';
        var rows = data['rows'] || data;

        html += '<table width="100%" cellspacing="0" cellpadding="0">';

        for (var i = 0; i < rows.length; i++) {
            item = rows[i];
            html += '<tr'
                + (item[this.idField] ? ' data-indent="' + item[this.idField] + '"' : '')
                + '>';

            for (var j in this.columns) {
                col = this.columns[j];
                val = item[col.field];

                html += '<td'
                    + (j == 0 ? ' class="first"' : '')
                    + (i == 0 && col.width ? ' width="' + col.width + '"' : '')
                    + (col.align?' align="'+col.align+'"':'')
                    + '>'
                    + (col.formatter && col.formatter instanceof Function ? col.formatter(val, item, i) : val)
                    + '</td>';

            }
            html += '</tr>';
        }

        html += '</table><div style="clear:both"></div>';

        //gridview的第1个div
        var gv = this.gridView.getElementsByTagName('DIV')[1];
        gv.innerHTML = html;

        //this._fixPosition();

        gv.srcollTop = 0;

        this.loadbox.style.display = 'none';

        if (this.onLoaded && this.onLoaded instanceof Function)
            this.onLoaded(data);
    };


    this._fixPosition= function(){
    };

    this._load_data = function (func) {
        if (!this.data_url) return;
        var t = this;

        if (func) {
            if (!(func instanceof Function)) {
                func = null;
            }
        }

        j6.xhr.request({
            uri: this.data_url,
            data: 'json',
            params: this.data,
            method: 'POST'
        }, {
            success: function (json) {
                t._fill_data(json);
            }, error: function () {
                //alert('加载失败!');
            }
        });

    };

    /* 为兼容IE6 */
    //var resizeFunc = (function (t) {
    //    return function () {
    //        t.resize.apply(t);
    //    };
    //})(this);
    //j6.event.add(window, 'load', resizeFunc);
    //window.attachEvent('resize', resizeFunc);
    //j6.event.add(window, 'resize', this.resize.apply(this));

    this._initLayout();

    //重置尺寸
    //this._resize();

    //加载数据
    this.load();
}

datagrid.prototype.resize = function () {
    this._fixPosition();
};

datagrid.prototype.load = function (data) {
    //显示加载框
    this.loading();
    if (data && data instanceof Object) {
        this._fill_data(data);
        return;
    }
    this._load_data();
};

/* 重新加载 */
datagrid.prototype.reload = function (params, data) {
    if (params) {
        this.data = params;
    }
    this.load(data);
};

j6.extend({
    grid: function (ele, config) {
        return new datagrid(ele, config);
    },
    datagrid: function (ele, config) {
        return new datagrid(ele, config);
    }
});


