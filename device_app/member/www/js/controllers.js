angular.module('starter.controllers', [])

.controller('HomeCtrl', function($scope) {
	//$scope.merchant_name = APP.config.partnerName;
})

.controller('OrdingCtrl', function($scope) {
  //$scope.friends = Friends.all();

  //load categories
  $JS.send('category',{},function(r){
  	var html = "<ul>";
  	for(var i=0;i<r.length;i++){
  		html += '<li index="'+i+'">'+ r[i].Name +'</li>';
  	}

  	html += '</ul>';
  	$JS.$('category-panel').innerHTML=html;
  });
})

.controller('DashCtrl', function($scope) {
})

.controller('FriendsCtrl', function($scope, Friends) {
  $scope.friends = Friends.all();
})

.controller('FriendDetailCtrl', function($scope, $stateParams, Friends) {
  $scope.friend = Friends.get($stateParams.friendId);
})

.controller('AccountCtrl', function($scope) {
});
