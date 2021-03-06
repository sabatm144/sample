sampleApp.controller("showCtrl", function ($scope, content, $http) {

    console.log("Inside show ctrl", content)
    $scope.content = content
    $scope.comment = {
      "text": ""
    }
   
    $scope.getComments = function (contentID) {
      $http({
        method: 'GET',
        url: '/content/' + contentID + '/comments',
        headers: {
          Authorization: localStorage.getItem("sample_user_token")
        }
      }).then(function successCallback(response) {
        $scope.comments = response.data.comments || []
      }, function errorCallback(response) {
        console.log("COUNT ERROR: ", response)
      });
    }
    $scope.getComments(content.id)
   
    $scope.postComment = function (contentID, comment) {
      console.log(contentID, comment)
      var config = {
        headers: {
            Authorization: localStorage.getItem("sample_user_token")
        }
      }
   
      $http.put('/content/' + contentID + '/comment', comment, config).then(function successCallback(response) {
        console.log("COMMENT SUCCESS: ", response)
        $scope.getComments(content.id)
        alert(response.data)
        $scope.comment.text = ""
      }, function errorCallback(response) {
        console.log("COMMENT ERROR: ", response)
      });
    };
   
    $scope.replyComment = function (commentID, comment) {
      console.log(commentID, comment)
      var config = {
        headers: {
          Authorization: localStorage.getItem("sample_user_token")
        }
      }
   
      $http.put('/comment/' + commentID + '/reply', comment, config).then(function successCallback(response) {
        console.log("COMMENT SUCCESS: ", response)
        $scope.getComments(content.id)
      }, function errorCallback(response) {
        console.log("COMMENT ERROR: ", response)
      });
    }
   
   });