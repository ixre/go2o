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
    this.call = function (funcName,arg) {
        return eval('window' +this.joinStr()+"." + funcName + '(arg)');
    };
}

window.pc = window;
var dpList= [new deviceCallProvider("android"),
    new deviceCallProvider("ios"),
    new deviceCallProvider("pc")];

function loopCall(func,arg){
    for(var i=0;i< dpList.length;i++) {
        var p = dpList[i];
        if (p.providerName && p.defined(func)) {
            p.call(func,arg);
            break;
        }
    }
}

window.cli = {
    alert:function(arg){
        loopCall('alert',arg);
    },
    login:function(arg){
        loopCall('login',arg);
    },
    close:function(arg){
        loopCall('close',arg);
    },
};