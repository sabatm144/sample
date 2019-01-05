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

  $http.put('/comment/' + contentID, comment, config).then(function successCallback(response) {
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
      url: '/totalComments/' + contentID,
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

$scope.updateVote = function(contentData, Value) {

  var statusIns = {
    "status": Value
  }
  console.log(contentData, status)
  var config = {
    headers : {
      Authorization: localStorage.getItem("sample_user_token")
    }
  }

  $http.put('/likeContent/' + contentData.id, statusIns, config).then(function successCallback(response) {
      console.log("SUCCESS: ", response)
      $scope.countLikes(contentData.id)
      }, function errorCallback(response) {
          console.log("ERROR: ", response)
  });

};

$scope.countLikes = function(contentID) {
  console.log(contentID)
    $http({
      method: 'GET',
      url: '/countVotes/' + contentID,
      headers : {
        Authorization: localStorage.getItem("sample_user_token")
      }
    }).then(function successCallback(response) {
      console.log("LIKE SUCCESS: ", response)
      $scope.dLikeCount = response.data.noOfDisLikes
      $scope.likeCount = response.data.noOfLikes
      console.log("LIKE SUCCESS: ", $scope.likeCount)
      console.log("LIKE SUCCESS: ", $scope.dLikeCount)

      }, function errorCallback(response) {
          console.log("LIKE ERROR: ", response)
      });
    };
});