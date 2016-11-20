(function($){
  $(function(){

    $('.button-collapse').sideNav();
		$('.tooltipped').tooltip({delay: 50});

  }); // end of document ready
})(jQuery); // end of jQuery name space

Vue.directive('pooler', {
  deep: true,
	update: function (x) {

    var divstr = '<h5>Stack is currently: ' + this.vm.response.data[x].state + '</h5> <p>Click on blue and the following pool members will be enabled</p> <ul class="collection with-header bluepool">';

		$.each(this.vm.response.data[x].pools, function(pkey, pvalue) {
      divstr += '<li class="collection-header"><h6>pool: ' + pkey + '<h6></li>';
			$.each(pvalue["blue"], function(mkey, mvalue) {
							divstr += '<li class="collection-item">' + mvalue + '</li>';
      });
    });
		divstr += '</ul>';

    divstr += '</p><p>Click on green and the following pool members will be enabled</p> <ul class="collection with-header greenpool">';

		$.each(this.vm.response.data[x].pools, function(pkey, pvalue) {
      divstr += '<li class="collection-header"><h6>pool: ' + pkey + '<h6></li>';
			$.each(pvalue["green"], function(mkey, mvalue) {
							divstr += '<li class="collection-item">' + mvalue + '</li>';
      });
    });
		divstr += '</ul>';

		this.el.innerHTML = divstr;

	}
});

var vm = new Vue({
	el: '#app',
		data: {
						response: ''
		},
		methods: {
			updateGroup: function(name, state) {

				var pdata = {"name": name, "state": state};
        this.$http.put('/api/group', pdata, function (resp, status, request) {
          // set data on vm
          this.$set('response', resp)
		      console.log(JSON.stringify(resp));
        }).error(function (resp, status, request) {
          // handle error
					console.log(status);
        });

			}
		},
    ready: function () {

      // GET request
      this.$http.get('/api/group', function (resp, status, request) {
          // set data on vm
          this.$set('response', resp)
          this.$set('user', resp.user)
      }).error(function (resp, status, request) {
          // handle error
					console.log(status)
      });
		$('.modal-trigger').leanModal();
		$('.tooltipped').tooltip({delay: 50});
		}
});


