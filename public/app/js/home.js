sampleApp.controller('homeCtrl', function($state, $scope, $http) { 

  // console.log(content)
  $scope.showComment = false
  $scope.showNComment = {
    "index": "",
    "show": false
  }

  $scope.comment = {
    "text" : "",
    "id": "",
    "childID": ""
  }

  $scope.token = localStorage.getItem("sample_user_token");
  $scope.userID = localStorage.getItem("user");

  console.log($scope.token)
  //Pagination
  $scope.totalItems = 0;
  $scope.currentPage = 1;
  $scope.limit = 5;
 
  // List services
  $scope.listContents = function () {
    $http({
      method: 'GET',
      url: '/getContents?page=' + $scope.currentPage + '&limit=' + $scope.limit
    }).then(function successCallback(response) {
      console.log(response)
      $scope.contents = response.data.contents || []
      $scope.totalItems = response.data.total || 0;
      $scope.currentPage = response.data.currentPage || 1;
      $scope.limit = response.data.limit || 5;
      console.log("SUCCESS: ", $scope.contents, $scope.userID)
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
    alert(response.data)
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

//
$scope.openComment = function(contentID) {
  $scope.showComment = !$scope.showComment
  $scope.countComments(contentID)
}

$scope.openNComment = function(commentID, index, childID) {

  $scope.nComment =  {
    "text" : "",
    "id": commentID,
    "childID": ""
  }

  if (childID) {
    console.log("Child present!")
    $scope.nComment.childID = childID
  }

  console.log($scope.nComment, childID)
  $scope.showNComment.index = index
  $scope.showNComment.show = true
}

//
$scope.postComment = function(contentID, comment) {
  console.log(contentID, comment)
  var config = {
    headers : {
      Authorization: localStorage.getItem("sample_user_token")
    }
  }

  $http.put('content/' + contentID + '/comment' + contentID, comment, config).then(function successCallback(response) {
    console.log("COMMENT SUCCESS: ", response)
    $scope.openComment(contentID)
    }, function errorCallback(response) {
        console.log("COMMENT ERROR: ", response)
    });
  };


$scope.countComments = function(contentID) {
  console.log(contentID)
    $http({
      method: 'GET',
      url: 'content/' + contentID + '/comments',
      headers : {
        Authorization: localStorage.getItem("sample_user_token")
      }
    }).then(function successCallback(response) {
      console.log("COUNT SUCCESS: ", response)
      $scope.comments = response.data.comments
      $scope.commentList = response.data.commentList

      }, function errorCallback(response) {
          console.log("COUNT ERROR: ", response)
      });
    };

$scope.vote = function(contentData, Value) {

  var statusIns = {
    "status": Value
  }
  console.log(contentData, status)
  var config = {
    headers : {
      Authorization: localStorage.getItem("sample_user_token")
    }
  }

  $http.put('content/' + contentData.id + '/vote', statusIns, config).then(function successCallback(response) {
      console.log("SUCCESS: ", response)
      $scope.listContents();
      }, function errorCallback(response) {
          console.log("ERROR: ", response)
  });

};

$scope.open = function (id) {
  console.log(id)
  $state.go("home.showPost", {
    id: id
  })
}

});