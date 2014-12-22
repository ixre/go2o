

function closeTab(title,call) {
    $('#tabs').tabs('close', title);
    if(call)call();
}

function closeCurrent(call) {
    var tab = $('#tabs').tabs('getSelected');
    var index = $('#tabs').tabs('getTabIndex', tab);
    closeTab(index, call);
}
function closeAndReresh(call, title) {
    var tab = $('#tabs').tabs('getSelected');
    var index = $('#tabs').tabs('getTabIndex', tab);
    closeTab(index, call);

    //获取跳转的tab
    var tabLength = $('#tabs .tabs li').length;

    var selectedTab = null;
    for (var i = 0; i < tabLength; i++) {
        selectedTab = $('#tabs').tabs('getTab', i);
        if (selectedTab.panel('options').title.replace(/<[^>]+>|\s|\n/igm, '') == title) {
            break;
        } else {
            selectedTab = null;
        }
    }

    //刷新跳转的tab
    if (selectedTab) {
        var index2 = $('#tabs').tabs('getTabIndex', selectedTab);
        var panel = $('#tabs .panel')[index2];
        var iframe = panel.getElementsByTagName('IFRAME');
        if (iframe.length != 0) {
            iframe = iframe[0];
            if (iframe.contentWindow.refresh)
                iframe.contentWindow.refresh();
        }
    }
}

function addTab2(subtitle,url,icon,closeable){
	closeTab(subtitle);
	addTab(subtitle,url,icon,closeable);
}

function addTab(subtitle, url, icon, closable, id) {
    if (!$('#tabs').tabs('exists', subtitle)) {
        $('#tabs').tabs('add', {
            title: subtitle,
            content: createFrame(url, id),
            closable: closable,
            icon: icon,
            loadingMessage: '正在加载中......'
        });
    } else {
        $('#tabs').tabs('select', subtitle);
		//$('#mm-tabupdate').click();
    }
    tabClose();
}


function createFrame(url)
{
	var s = '<iframe scrolling="auto" frameborder="0"  src="'+url+'" style="width:100%;height:100%;"></iframe>';
	return s;
}

function tabClose()
{
	/*双击关闭TAB选项卡*/
	$(".tabs-inner").dblclick(function(){
		var subtitle = $(this).children(".tabs-closable").text();
		$('#tabs').tabs('close',subtitle);
	})
	/*为选项卡绑定右键*/
	$(".tabs-inner").bind('contextmenu',function(e){
		$('#mm').menu('show', {
			left: e.pageX,
			top: e.pageY
		});

		var subtitle =$(this).children(".tabs-closable").text();

		$('#mm').data("currtab",subtitle);
		$('#tabs').tabs('select',subtitle);
		return false;
	});
}
//绑定右键菜单事件
function tabCloseEven()
{
	//刷新
	$('#mm-tabupdate').click(function(){
		var currTab = $('#tabs').tabs('getSelected');
		var url = $(currTab.panel('options').content).attr('src');
		$('#tabs').tabs('update',{
			tab:currTab,
			options:{
				content:createFrame(url)
			}
		})
	})
	//关闭当前
	$('#mm-tabclose').click(function(){
		var currtab_title = $('#mm').data("currtab");
		$('#tabs').tabs('close',currtab_title);
	})
	//全部关闭
	$('#mm-tabcloseall').click(function(){
		$('.tabs-inner span').each(function(i,n){
			var t = $(n).text();
			$('#tabs').tabs('close',t);
		});
	});
	//关闭除当前之外的TAB
	$('#mm-tabcloseother').click(function(){
		$('#mm-tabcloseright').click();
		$('#mm-tabcloseleft').click();
	});
	//关闭当前右侧的TAB
	$('#mm-tabcloseright').click(function(){
		var nextall = $('.tabs-selected').nextAll();
		if(nextall.length==0){
			//msgShow('系统提示','后边没有啦~~','error');
			alert('后边没有啦~~');
			return false;
		}
		nextall.each(function(i,n){
			var t=$('a:eq(0) span',$(n)).text();
			$('#tabs').tabs('close',t);
		});
		return false;
	});
	//关闭当前左侧的TAB
	$('#mm-tabcloseleft').click(function(){
		var prevall = $('.tabs-selected').prevAll();
		if(prevall.length==0){
			alert('到头了，前边没有啦~~');
			return false;
		}
		prevall.each(function(i,n){
			var t=$('a:eq(0) span',$(n)).text();
			$('#tabs').tabs('close',t);
		});
		return false;
	});

	//退出
	$("#mm-exit").click(function(){
		$('#mm').menu('hide');
	})
}

//弹出信息窗口 title:标题 msgString:提示信息 msgType:信息类型 [error,info,question,warning]
function msgShow(title, msgString, msgType) {
	$.messager.alert(title, msgString, msgType);
}


