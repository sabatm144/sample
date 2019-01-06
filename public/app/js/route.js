function routeConfig($urlRouterProvider, $stateProvider) {
    $urlRouterProvider.when("/", "/post").otherwise("/");
    var app_dir = "public/app/html/";

    $stateProvider.state('login', {
        url: '/login',
        controller: 'logInCtrl',
        templateUrl: app_dir + 'logIn.html'
    }).state('register', {
        url: '/register',
        controller: 'logInCtrl',
        templateUrl: app_dir + 'signUp.html'
    }).state('home', {
        url: "^/home",
        controller: "homeCtrl",
        templateUrl: app_dir + "home.html",
    }).state("home.createPost", {
        url: "^/create",
        controller: "createCtrl",
        templateUrl: app_dir + "create.html",
        resolve: {
            'content': function () {
                return {};
            }
        }
    }).state("home.editPost", {
        url: "^/edit/:id",
        controller: "createCtrl",
        templateUrl: app_dir + "create.html",
        resolve: {
            'content': getContent
        }
    }).state("home.list", {
        url: "^/post",
        controller: "homeCtrl",
        templateUrl: app_dir + "list.html",
    })
}

sampleApp.config(routeConfig);

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