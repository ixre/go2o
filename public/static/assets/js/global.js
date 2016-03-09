function setGeo(id){
 j6.xhr.get('/main/geoLocation',function(r){
        var j=j6.toJson(r);
        j6.$(id).innerHTML='<span style="color:green">'+ j.addr+'</span>';
    });
}