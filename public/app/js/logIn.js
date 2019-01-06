sampleApp.controller('logInCtrl', function($state, $scope, $http) { 
   
    $scope.error = false

    $scope.LogIn = function(contentData) {
        $http.post('/authenticateUser', contentData).then(function successCallback(response) {
                console.log("USER LOGIN SUCCESS: ", response)
                localStorage.setItem("sample_user_token", response.data.token);
                localStorage.setItem("user", response.data.customer.id);
                $scope.logInError = response.data.message
                $state.go("home.list")
                }, function errorCallback(response) {
                  $scope.error = true
                  $scope.lError = response.data.message
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
            $scope.error = true
           $scope.rError = response.data.message
            console.log("USER REGISTRATION ERROR: ", response)
        });;
    }   

    $scope.cancel = function() {
        $state.go("home.list")
    }
});