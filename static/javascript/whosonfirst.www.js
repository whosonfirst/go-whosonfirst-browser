var whosonfirst = whosonfirst || {};

whosonfirst.www = (function(){

	var self = {

		'init': function(){

		},

		'abs_root_url': function(){
			var body = document.body;
			return body.getAttribute("data-abs-root-url");
		},

	};

	return self;
})();
		
