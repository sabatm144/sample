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