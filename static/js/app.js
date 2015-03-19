(function(){
  var app = angular.module("etcdRouter", []);

  app.controller("HostsController", [ '$http', '$interval', function($http, $interval){
    var hostsCtrl = this;
    hostsCtrl.hosts = [];

    hostsCtrl.loadData = function(){
      $http.get('http://localhost:3000/api/hosts').success(function(data){
        hostsCtrl.hosts = data.hosts;
      });
    };

    hostsCtrl.loadData();

    $interval(function() {
      hostsCtrl.loadData();
    }, 3000);
  }]);

})();
