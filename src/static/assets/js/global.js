function setGeo(id){
 $JS.xhr.get('/comm/geoLocation',function(r){
        var j=$JS.toJson(r);
        $JS.$(id).innerHTML='<span style="color:green">'+ j.addr+'</span>';
    });
}