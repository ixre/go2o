function getScroll() 
{
    var t, l, w, h;
     
    if (document.documentElement && document.documentElement.scrollTop) {
        t = document.documentElement.scrollTop;
        l = document.documentElement.scrollLeft;
        w = document.documentElement.scrollWidth;
        h = document.documentElement.scrollHeight;
    } else if (document.body) {
        t = document.body.scrollTop;
        l = document.body.scrollLeft;
        w = document.body.scrollWidth;
        h = document.body.scrollHeight;
    }
    return { t: t, l: l, w: w, h: h };
}
window['pxDevicePixelRatio'] = window['devicePixelRatio'] || 1;
if (window['pxDevicePixelRatio'] > 1) {
    window['pxDevicePixelRatio'] = 2;
} else {
    window['pxDevicePixelRatio'] = 1;
}
(function($) {
    $.fn.tab = function() {
        $(this).each(function() {
            var length = $(this).children('a').length - 1;
            $(this).children().width(parseInt(($(this).width() - length * 2 - 1) / $(this).width() / (length + 1) * 1000) / 10 + "%");
        });
        return $(this);
    };
    $.fn.slide = function(options) {
        $(this).each(function() {
            var self = this;
            var settings = {
                time: 3000,
                outTime: 1000,
                srcProperty: 'truesrc',
                auto: true,
                index: 0
            };
            if (options) {
                $.extend(settings, options);
            }
            clearTimeout($(self).data('slide-timeout'));
            clearInterval($(self).data('slide-interval'));
            try {
                $(self).data('slide-loading-img').unbind('load').remove();
            } catch (e) {
            }
            var width = $(self).width();
            var p = $(self).children('p');
            p.children('a').width(width);
            var click = function() {
                if ($(this).hasClass('hover')) {
                    return;
                }
                if ($(this).find('i.arrow-left').length) {
                    goto(--settings.index);
                } else if ($(this).find('i.arrow-right').length) {
                    goto(++settings.index);
                } else {
                    goto($(this).index());
                }
            };
            try {
                $(self).unbind('swipeLeft').swipeLeft(function() {
                    goto(++settings.index);
                });
                $(self).unbind('swipeRight').swipeRight(function() {
                    goto(--settings.index);
                });
            } catch (e) {
            }
            $(self).children('div').find('a').unbind('click').bind('click', click);
            $(self).children('div').find('span a').unbind('mouseenter').bind('mouseenter', click);
            var auto = function() {
                clearInterval($(self).data('slide-interval'));
                clearTimeout($(self).data('slide-timeout'));
                try {
                    $(self).data('slide-loading-img').unbind('load').remove();
                } catch (e) {
                }
                if (settings.auto) {
                    $(self).data('slide-interval', setInterval(function() {
                        goto(++settings.index);
                    }, settings.time));
                }
            };
            var imgs = [];
            p.find('a').each(function() {
                imgs.push(this);
            });
            var goto = function(index) {
                if (index < 0) {
                    index = imgs.length - 1;
                } else if (index >= imgs.length) {
                    index = 0;
                }
                clearInterval($(self).data('slide-interval'));
                settings.index = index;
                p.animate({'left': settings.index * width * -1}, {queue: false, duration: 500});
                $(self).children('div').find('span a.hover').removeClass('hover');
                var element = $(imgs[settings.index]).find('img');
                $(self).children('div').find('span a:eq(' + settings.index + ')').addClass('hover');
                var src = element.attr(settings.srcProperty);
                if (src) {
                    (function(src, element) {
                        var img = $('<img style="display: none;"/>').appendTo('body').bind('load', function() {
                            clearTimeout($(self).data('slide-timeout'));
                            $(element).attr('src', src);
                            img.remove();
                            auto();
                            $(element).attr(settings.srcProperty, null);
                        }).attr('src', src);
                        $(self).data('slide-loading-img', img);
                        $(self).data('slide-timeout', setTimeout(function() {
                            img.trigger('load');
                        }, settings.outTime));
                    })(src, element);
                } else {
                    auto();
                }
            };
            goto(settings.index);
        });
        return $(this);
    };
    var lazyloadSize = {
        '2_8080': '160160',
        '2_120120': '220220',
        '2_160160': '320320',
        '2_180180': '320320',
        '2_220220': '320320',
        '2_320320': '320320',
        '2_480480': '480480',
        '1_8080': '8080',
        '1_120120': '120120',
        '1_160160': '160160',
        '1_180180': '180180',
        '1_220220': '220220',
        '1_320320': '320320',
        '1_480480': '480480',
        '2_8060': '160160',
        '2_12090': '220220',
        '2_160120': '320320',
        '2_180135': '320320',
        '2_220165': '320320',
        '2_320300': '320320',
        '2_480360': '480480',
        '1_8060': '8080',
        '1_12090': '120120',
        '1_160120': '160160',
        '1_180135': '180180',
        '1_220165': '220220',
        '1_320300': '320320',
        '1_480360': '480480',
    };
    $.fn.lazyload = function(options) {
        $(this).each(function() {
            var self = this;
            var settings = {
                threshold: 0, //间隔
                container: window, //
                outTime: 1000, //加载超时
                srcProperty: 'truesrc', //真实src地址
                child: 'img'		//子选择器
            };
            if (options) {
                $.extend(settings, options);
            }
            var DATA_NAME = 'touch_lazyload_event_' + settings.mode + '_' + settings.srcProperty;
            var DATA_NAME_Timeout = DATA_NAME + '_timeout';
            var CHILD_NAME = DATA_NAME + '_child';
            if (settings.act == 'run') {
                try {
                    $(self).data(DATA_NAME)();
                    return;
                } catch (e) {
                }
            }
            clearTimeout($(self).data(DATA_NAME_Timeout));
            if ($(self).data(DATA_NAME)) {
                $(settings.container).unbind('scroll', $(self).data(DATA_NAME));
            }
            $(self).data(CHILD_NAME, []);
            if (settings.act == 'stop') {
                return;
            }
            $(self).find(settings.child).each(function() {
                $(self).data(CHILD_NAME).push(this);
            });
            var length = $(self).data(CHILD_NAME).length;
            $(self).data(DATA_NAME, function() {
                clearTimeout($(self).data(DATA_NAME_Timeout));
                $(self).data(DATA_NAME_Timeout, setTimeout(function() {
                    var scrollTop = $(settings.container).scrollTop();
                    var height = $(settings.container).height();
                    var isload = false;
                    for (var i = 0; i < length; ) {
                        var element = $(self).data(CHILD_NAME)[i];
                        var top = $(element).offset().top;
                        var src = $(element).attr(settings.srcProperty);
                        if (src) {
                            if (height + scrollTop < top - settings.threshold || top + $(element).height() + settings.threshold < scrollTop) {
                                if (isload) {
                                    break;
                                }
                                i++;
                            } else {
                                isload = true;
                                $(element).attr(settings.srcProperty, null);
								clearTimeout($(this).data('setTimeout'));
								if ($(element).hasClass('shoe')) {
									var srcs = src.match(/(\d+)(\.[a-zA-Z]+)$/);
									src = src.replace(srcs[0], lazyloadSize[window['pxDevicePixelRatio'] + '_' + srcs[1]] + srcs[2]);
									srcs = null;
								}
								var tagName = "";
								try{
									tagName = $(element).attr('tagName').toLowerCase();
								}catch(e) {
									tagName = $(element)[0].tagName.toLowerCase();
								}
								if (tagName == 'img') {
									if (/(^|\s)img\d+($|\s)/.test($(element).parent().attr('class'))) {
										$(element).parent().css('background-image', 'none');
									}
									$(element).attr('src', src);
								} else {
									$(element).css({"background-image": "url(" + src + ")"});
								}
                                /*(function(src, element) {
                                    var img = $('<img style="display: none;"/>').appendTo('body').one('load', function() {
                                        clearTimeout($(this).data('setTimeout'));
                                        if ($(element).hasClass('shoe')) {
                                            var srcs = src.match(/(\d+)(\.[a-zA-Z]+)$/);
                                            src = src.replace(srcs[0], lazyloadSize[window['pxDevicePixelRatio'] + '_' + srcs[1]] + srcs[2]);
                                            srcs = null;
                                        }
                                        if ($(element).attr('tagName').toLowerCase() == 'img') {
                                            if (/(^|\s)img\d+($|\s)/.test($(element).parent().attr('class'))) {
                                                $(element).parent().css('background-image', 'none');
                                            }
                                            $(element).attr('src', src);
                                        } else {
                                            $(element).css({"background-image": "url(" + src + ")"});
                                        }
                                        img.remove();
                                    }).attr('src', src);
                                    img.data('setTimeout', setTimeout(function() {
                                        img.trigger('load');
                                    }, settings.outTime));
                                })(src, element);*/
                                $(self).data(CHILD_NAME).splice(i, 1);
                                length--;
                            }
                        } else {
                            $(self).data(CHILD_NAME).splice(i, 1);
                            length--;
                        }
                    }
                    if (!length) {
                        clearTimeout($(self).data(DATA_NAME_Timeout));
                        if ($(self).data(DATA_NAME)) {
                            $(settings.container).unbind('scroll', $(self).data(DATA_NAME));
                        }
                        $(self).data(CHILD_NAME, null);
                    }
                }, 200));
            });
            $(settings.container).bind('scroll', $(self).data(DATA_NAME));
            $(self).data(DATA_NAME)();
        });
        return $(this);
    };
    $.message = function(options) {
        var buttons = [
            {
                id: null,
                light: true,
                text: '  确 定  ',
                click: function() {
                    returnObj.close();
                }
            },
            {
                id: null,
                light: false,
                text: '  取 消  ',
                click: function() {
                    returnObj.close();
                }
            }
        ];
        var settings = {
            html: '',
            title: '',
            height: 'auto',
            buttons: buttons
        };
        if (options) {
            $.extend(settings, options);
        }
        var self = $('<div class="pxui-message"><div><h3><span></span><a><i></i></a></h3><div class="content"></div><div class="buttons"></div></div></div>').appendTo('#js-com-content-area');
        var scrollTop = getScroll().t;
        self.css('top', scrollTop + 12);
        var returnObj = {
            close: function() {
                self.remove();
            },
            title: function(title) {
                self.find('h3 span').text(title);
            },
            html: function(html) {
                self.find('div.content').html(html);
            },
            height: function(height) {
                height += '';
                height.replace(/px$/i, '');
                if (/^\d+$/.test(height)) {
                    height = parseFloat(height);
                    height -= 129;
                    if (height < 0) {
                        height = 0;
                    }
                }
                self.find('div.content').height(height);
            },
            buttons: function(_buttons) {
                var html = '';
                _buttons = _buttons || buttons;
                $.each(_buttons, function(index, btn) {
                    html += '<input ' + ((btn.id != null) ? ' id="' + btn.id + '" ' : '') + ((btn.disabled) ? ' disabled ' : '') + ((btn.light) ? ' class="pxui-light-button '+btn.class+'" ' : ' class="'+btn.class+'"') + ' type="button" value="' + btn.text + '">';
                });
                self.find('.buttons').html(html);
                var inputs = self.find('.buttons input');
                $.each(_buttons, function(index, btn) {
                    if (!btn.disabled) {
                        if (btn.click) {
                            inputs.eq(index).click(btn.click);
                        }
                    }
                });
            },
            base: function() {
                return self;
            }
        };
        self.find('h3 a').click(function() {
            returnObj.close();
        });
        returnObj.title(settings.title);
        returnObj.html(settings.html);
        returnObj.height(settings.height);
        returnObj.buttons(settings.buttons);
        return returnObj;
    };
    $.fn.getMore = function(options) {
        var settings = {
            data: {},
            url: '',
            callback: null,
            template: null,
            lastid: ''
        };
        if (options) {
            $.extend(settings, options);
        }
        if (!settings.template) {
            settings.template = $('#js-good-template').html();
        }
        settings.template = template.compile(settings.template);
        var self = this;
        if ($(self).hasClass('pxui-show-more')) {
            var lastid = settings.lastid;
            function ajax() {
                settings.lastid = lastid;

                $.ajax({url: settings.url, data: {lastid:settings.lastid}, type: "post", success: function(data) {
                        $(self).removeClass('pxui-show-more-loading');
                        try {
                            data = $.parseJSON(data);
                            var html = '';
                            if (data.list && data.list.length) {
                                $.each(data.list, function(index, value) {
                                    html += settings.template({data: value}).replace(new RegExp('&#60;',"gm"),'<').replace(new RegExp('&#62;',"gm"),'>');
                                });
                                settings.callback(html);
                            } else {
                                var nodata = '';
                                if(lastid == 1 && $(self).attr('noDataTemp')!= null && $(self).attr('noDataTemp') != ''){
                                    nodata = $(self).attr('noDataTemp');
                                    var noDataTemp = $(nodata+'').html();
                                    $(self).parent().append(noDataTemp);
                                }

                                if ($(self).find('a').attr('tourl') != null && $(self).find('a').attr('tourl') != '')
                                {
                                    window.location.href = $(self).find('a').attr('tourl');
                                }
                                $(self).remove();
                                
                            }
                            if (data.isend) {
                                $(self).remove();
                            }
                            lastid = data.lastid;
//                                                settings.lastid = lastid;


                            ;
                        } catch (e) {
                           // alert('链接服务器失败，请稍后再试！');
                        }
                    }, error: function() {
                     //   alert('链接服务器失败，请稍后再试！');
                        $(self).removeClass('pxui-show-more-loading');
                    }});
            }
            ;
            $(self).find('a').unbind('click').bind('click', function() {
                $(self).addClass('pxui-show-more-loading');
                ajax();
            });
        }
        return $(this);
    };
    $.fn.page = function(options) {
        var settings = {
            data: {},
            url: '',
            callback: null,
            template: null,
            pagesize: 20,
            count: 0,
            container: null,
            pageurl: ''
        };
        if (options) {
            $.extend(settings, options);
        }
        if (!settings.template) {
            settings.template = $('#js-good-template').html();
        }
        settings.template = template.compile(settings.template);

        var self = this;
        if ($(self).hasClass('pxui-pages')) {
            if (settings.container.css('position') != 'absolute' || settings.container.css('position') != 'relative') {
                settings.container.css('position', 'relative');
            }
            function setbutton() {
                if (settings.page == 1) {
                    $prev.hide().prev().show();
                } else {
                    $prev.show().prev().hide();
                }
                if (settings.page == $select.find('option').length) {
                    $next.hide().next().show();
                } else {
                    $next.show().next().hide();
                }
            }
            ;
            var ajaxCount = 0;
            var loading = null;
            function ajax() {
                setbutton();
                if (settings.pageurl) {
                    try {
                        history.replaceState(null, '', settings.pageurl.replace(/@page/g, settings.page));
                    } catch (e) {
                    }
                }
                try {
                    loading.remove();
                } catch (e) {
                }
                document.body.scrollTop = settings.container.offset().top;
                $select.unbind('change', $select.data('page-change'));
                $select.val(settings.page).change();
                $select.bind('change', $select.data('page-change'));
                var _ajaxCount = ++ajaxCount;
                settings.data.page = settings.page;
                settings.data.pageSize = settings.pagesize;
                settings.data.count = settings.count;
                loading = $('<img style="position: absolute;top: 5px;left: 50%;margin-left: -12px;" src="template/images/public/loading.gif" width="24" height="24"/>').appendTo(settings.container);

                $.ajax({url: settings.url, type: "POST", data: settings.data, error: function() {
                        if (_ajaxCount != ajaxCount)
                            return;
                        loading.remove();
                       // alert('链接服务器失败，请稍后再试！');
                    }, success: function(data) {
                        if (_ajaxCount != ajaxCount)
                            return;
                        try {
                            data = $.parseJSON(data);
                            settings.count = data.count;
                            if (data.list) {
                                var html = '';
                                $.each(data.list, function(index, value) {
                                    html += settings.template({data: value}).replace(new RegExp('&#60;',"gm"),'<').replace(new RegExp('&#62;',"gm"),'>');;
                                });
                                settings.container.html(html);
                                document.body.scrollTop = settings.container.offset().top;
                                if (settings.callback) {
                                    settings.callback(html);
                                }
                            }
                        } catch (e) {
                          //  alert('链接服务器失败，请稍后再试！');
                        }
                    }});
            }
            ;
            var $select = $(self).find('select');
            $(self).find('a').unbind('click');
            var $prev = $(self).find('a').eq('0');
            var $next = $(self).find('a').eq('1');
            $prev.click(function() {
                settings.page--;
                ajax();
            });
            $next.click(function() {
                settings.page++;
                ajax();
            });
            if ($select.data('page-change')) {
                $select.unbind('change', $select.data('page-change'));
            }
            settings.page = parseInt($select.val());
            setbutton();
            $select.data('page-change', function() {
                settings.page = parseInt($select.val());
                ajax();
            });
            $select.bind('change', $select.data('page-change'));
        }
        return $(this);
    };
    // 遮照层
    $.documentMask = function(options) {
        // 扩展参数
        var op = $.extend({
            opacity: 0.6,
            z: 150,
            bgcolor: '#000',
            time: 500,
            id: "jquery_addmask"
        }, options);
        // 创建一个 Mask 层，追加到 document.body
        $("#" + op.id).remove();
        $('<div id="' + op.id + '" class="jquery_addmask"> </div>').appendTo(document.body).css({
            position: 'absolute',
            top: '0px',
            left: '0px',
            'z-index': op.z,
            width: $(window).width(),
            height: $(document).height(),
            'background-color': op.bgcolor,
            opacity: 0
        }).fadeTo(op.time, op.opacity).click(function() {
            // 单击事件，Mask 被销毁
            //$(this).fadeTo('slow', 0, function(){
            //    $(this).remove();
            //});
        });
        return this;
    };
})($);
function showError(error, id) {
    $("#" + (id || "js-error-msg")).fadeIn().html(error);
    document.body.scrollTop = $("#" + (id || "js-error-msg")).offset().top;
}
;
$(document).ready(function(e) {
	$('#js-com-header-search').click(function(){
		$(this).toggleClass('selected');
		//$(this).parent().next().next().toggle();
		$('#js-com-header-area form').toggle();
	});
    /*下拉选择*/
    $('.pxui-select select').live('change', function() {
        $(this).parent().children('span').html($(this).children('option[value="' + $(this).val() + '"]').html());
    });
    /*下拉选择-end*/
    /*列表*/
    $('.pxui-list a').live('click', function() {
        try {
            if ($(this).next()[0].nodeName.toLowerCase() !== 'a') {
                if ($(this).hasClass('open')) {
                    $(this).find('i.arrow-right').removeClass('arrow-bottom');
                    $(this).removeClass('open').next().css({"display": "none"});
                } else {
                    if ($(this).parent().attr('data-model') === 'radio') {
                        $(this).siblings('a.open').click();
                    }
                    $(this).find('i.arrow-right').addClass('arrow-bottom');
                    $(this).addClass('open').next().css({"display": "block"});
                }
            }
        } catch (e) {
        }
    });
    /*列表*/
    /*标签*/
    $('.pxui-tab > a').live('click', function() {
        $(this).parent().children('.selected').removeClass('selected');
        $(this).addClass('selected');
        $(this).parent().find('input').val($(this).attr('value'));
    });
    $(window).resize(function(e) {
        $('.pxui-tab').tab();
    });
    $('.pxui-tab').tab();
    /*标签-end*/
    /*搜索*/
    $("#js-com-search-text").focus(function() {
        $("#js-com-search-recommend").fadeIn(100);
    }).blur(function() {
        $("#js-com-search-recommend").fadeOut(300);
    });
    /*搜索-end*/
    $('.pxui-show-more').each(function() {
        var self = this;
        var url = $(self).attr('url');
        var container = $(self).attr('container');
        if (url && container) {
            var srcProperty = $(self).attr('srcProperty') || 'goodsrc';

            $(self).getMore({
                template: $($(self).attr('template') || '#js-good-template').html(),
                lastid: $(self).attr('lastid'),
                url: url,
                callback: function(html) {
                    $(container).append(html).lazyload({child: '[' + srcProperty + ']', srcProperty: srcProperty});
                }
            });
            $(container).lazyload({child: '[' + srcProperty + ']', srcProperty: srcProperty});
        }
    });
    $('.pxui-pages').each(function() {
        var self = this;
        var url = $(self).attr('url');
        var container = $(self).attr('container');
        if (url && container) {
            var srcProperty = $(self).attr('srcProperty') || 'goodsrc';
            var pagesize = $(self).attr('pagesize');
            var count = $(self).attr('count');
            container = $(container);

            $(self).page({
                template: $($(self).attr('template') || '#js-good-template').html(),
                pagesize: pagesize,
                count: count,
                url: url,
                container: container,
                pageurl: $(self).attr('pageurl'),
                callback: function(html) {
                    container.lazyload({child: '[' + srcProperty + ']', srcProperty: srcProperty});
                }
            });
            container.lazyload({child: '[' + srcProperty + ']', srcProperty: srcProperty});
        }
    });
    /*延时加载示例*/
    $('body').lazyload({child: '[truesrc]'});
    $('[tip]').live('click', function() {
        if (!confirm($(this).attr('tip'))) {
            return false;
        } else {
            return true;
        }
    });
    /*弹出框示例*/
    $('#js-show-message').click(function() {
        var msgbox = $.message({
            html: '内容',
            title: '标题',
            height: 'auto',
            buttons: [
                {
                    light: true,
                    text: '  确 定1  ',
                    click: function() {
                        msgbox.html('内容2');
                        msgbox.title('标题2');
                        msgbox.height(500);
                        msgbox.buttons([
                            {
                                light: false,
                                text: '  取 消  ',
                                click: function() {
                                    msgbox.close();
                                }
                            }
                        ]);
                    }
                },
                {
                    light: false,
                    text: '  取 消  ',
                    click: function() {
                        msgbox.close();
                    }
                }
            ]
        });
    });
   // 商城公告
    $('.com-notification a.close').live('click', function() {
        $(this).parent().animate({opacity: 0}, {queue: false, duration: 200, complete: function() {
                $(this).remove();
                setCookie('touch_notice_close','1');
            }});
    });


});
// 设置cookie
function setCookie(name,value)
{
    var Days = 1;
    var exp = new Date();
    exp.setTime(exp.getTime() + Days*24*60*60*1000);
    document.cookie = name + "="+ escape (value) + ";expires=" + exp.toGMTString();

    
//    var strsec = getsec(time);
//    var exp = new Date();
//    exp.setTime(exp.getTime() + strsec*1);
//    document.cookie = name + "="+ escape (value) + ";expires=" + exp.toGMTString();
}
//读取cookies
function getCookie(name)
{
    var arr,reg=new RegExp("(^| )"+name+"=([^;]*)(;|$)");
 
    if(arr=document.cookie.match(reg))
 
        return (arr[2]);
    else
        return null;
}
