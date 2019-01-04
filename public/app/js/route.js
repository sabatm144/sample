function routeConfig($urlRouterProvider, $stateProvider) {
    $urlRouterProvider.when("/", "/list").otherwise("/");
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
        templateUrl: app_dir + "home.html"
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
        url: "^/list",
        controller: "homeCtrl",
        templateUrl: app_dir + "list.html"
    })
}

sampleApp.config(routeConfig);