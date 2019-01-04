
sampleApp.controller('homeCtrl', function($state, $scope, $http) { 

  $scope.token = localStorage.getItem("sample_user_token");
  $scope.userID = localStorage.getItem("user");

  console.log($scope.token)
    // List services
    $scope.listContents = function() {
    $http({
      method: 'GET',
      url: '/getContents'
    }).then(function successCallback(response) {
      $scope.contents = response.data
      console.log("SUCCESS: ", $scope.contents,  $scope.userID)
      }, function errorCallback(response) {
          console.log("ERROR: ", response)
      });
    };
    $scope.listContents();

//Delete
$scope.deleteContent = function(id) {
  $http({
    method: 'DELETE',
    url: '/deleteContent/' + id,
    headers : {
      Authorization: localStorage.getItem("sample_user_token")
    }
  }).then(function successCallback(response) {
    console.log("DELETE SUCCESS: ", response)
    alert(response.data.message)
    $scope.listContents();
    $scope.contents = response.data
    }, function errorCallback(response) {
        console.log("DELETE ERROR: ", response)
    });
  };

//log out 
$scope.logOut =  function() {
  console.log("Inside logout!")
  $scope.token = localStorage.removeItem("sample_user_token");
  $state.go("login")
}

$scope.displayContentDesc =  function(description) {
  var w = window.open();
  w.document.open();

 const markup = `<head>
  <title>Description</title>
</head><div><p>` + description + `</div></p>`
  w.document.write(markup);
  w.document.close();
}

});