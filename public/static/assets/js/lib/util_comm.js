
function tabCard(e,opt) {
    opt = opt || {};
    opt.event = opt.event || 'mouseover';
    if(opt.call != null && opt.call instanceof Function)opt.call = null;
    opt.tabClass = opt.tabClass || 'item';
    opt.frameClass = opt.frameClass || 'frame';
    opt.frames = opt.frames || e;
  
    var tabItems = jr.dom.getsByClass(e,opt.tabClass);
    var tabFrames = jr.dom.getsByClass(opt.frames,opt.frameClass);
    var len = tabItems.length;
    for (var i = 0; i < len; i++) {
        var func = function () {
            for (var j = 0; j < len; j++) {
                var isCurr = tabItems[j] == this;
                tabItems[j].className = opt.tabClass + (isCurr ? ' current':'');
                with (tabFrames[j]) {
                    if (isCurr) {
                        className = className.replace(' hidden', '');
                        if (opt.call != null)opt.call(tabItems[j], tabFrames[j]);
                    } else if (className.indexOf(' hidden') == -1) {
                        className = className + ' hidden';
                    }
                }
            }
        };
        if (opt.event == 'click') {
            tabItems[i].onclick = func;
        } else {
            tabItems[i].onmouseover = func;
        }
    }
}

function fmtAmount(amount){
    return amount.toFixed(2).replace(/([^\.]+)(\.|(\.[1-9]))0*$/ig,'$1$3');
}