sampleApp.controller('logInCtrl', function($state, $scope, $http) { 
   
    $scope.LogIn = function(contentData) {
        $http.post('/authenticateUser', contentData).then(function successCallback(response) {
                console.log("USER LOGIN SUCCESS: ", response)
                localStorage.setItem("sample_user_token", response.data.token);
                localStorage.setItem("user", response.data.customer.id);
                $scope.logInError = false
                $state.go("home.list")
                }, function errorCallback(response) {
                  $scope.logInError = true
                    console.log("USER LOGIN ERROR: ", response)
        });;
    }   
    
    $scope.SignUp = function(contentData) {
        $http.post('/registerUser', contentData).then(function successCallback(response) {
              console.log("USER REGISTRATION SUCCESS: ", response)
              localStorage.setItem("sample_user_token", response.data.token);
              $scope.registrationError = false
              $state.go("home")
        }, function errorCallback(response) {
           $scope.registrationError = true
            console.log("USER REGISTRATION ERROR: ", response)
        });;
    }   

    $scope.cancel = function() {
        $state.go("home.list")
    }
});