/**
 * Created by sonven on 15/7/4.
 */

function deviceCallProvider(providerName) {
    this.providerName = providerName;
    this.joinStr = function(){
        return !this.providerName || this.providerName.length == 0 ?'':"."+this.providerName;
    };

    this.defined = function(funcName){
        return eval('window'+this.joinStr()+' && window'+this.joinStr()+'.'+funcName+'!= undefined');
    }
    this.call = function (funcName,arg1,arg2,arg3,arg4,arg5,arg6) {
        return eval('window' +this.joinStr()+"." + funcName + '(arg1,arg2,arg3,arg4,arg5,arg6)');
    };
}

window.cli = {
    pds :[new deviceCallProvider("android"),
        new deviceCallProvider("ios"),
        new deviceCallProvider(null)],
    alert:function(arg1,arg2,arg3,arg4,arg5,arg6){
        var arr = new Array(arguments.length+1);
        var i = 0;
        for(var i=0;i< arguments.length;i++){
            arr[i+1] = arguments[i];
        }

        for(var i=0;i< this.pds.length;i++) {
            var p = this.pds[i];
            if (p.providerName && !p.defined('alert')) {
                continue;
            }
            p.call("alert",arg1,arg2,arg3,arg4,arg5,arg6);
            break;
        }
    }
};