function getContent($http, $stateParams, $q) {
  var deferred = $q.defer();
  if (!$stateParams.id) {
      deferred.reject();
      return
  }
  var config = {
    headers : {
      Authorization: localStorage.getItem("sample_user_token")
    }
  }
  $http.get('content/'+ $stateParams.id, config).then(function successCallback(content) {
    console.log(content.data)
      deferred.resolve(content.data);
  }, function (error) {
      if (error.status === 401) deferred.reject("sessionExpired");
  });

  return deferred.promise;
}
getContent.$inject = ["$http", "$stateParams", "$q"];

sampleApp.controller('createCtrl', function($stateParams, $state, $scope, $http, content) { 

  console.log("Inside create ctrl")
    $scope.contentData = content

    console.log("Inside create ctrl")
    //Create/update content
    $scope.postContent = function(contentData) {
      
        var config = {
          headers : {
            Authorization: localStorage.getItem("sample_user_token")
          }
        }
        if (!$scope.contentData.id) {
          $http.post('/createContent', contentData, config).then(function successCallback(response) {
            console.log("SUCCESS: ", response)
            $state.go("home.list")
            }, function errorCallback(response) {
                console.log("ERROR: ", response)
                if (response.data.code === 401) {
                  alert("Sign In to post!")
                }
            });
        }
        if ($scope.contentData.id) {
          $http.put('/editContent/' + $stateParams.id, contentData, config).then(function successCallback(response) {
            console.log("SUCCESS: ", response)
            $state.go("home.list")
            }, function errorCallback(response) {
                console.log("ERROR: ", response)
            });

        }
    };

});