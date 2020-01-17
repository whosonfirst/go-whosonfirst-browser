var whosonfirst = whosonfirst || {};
whosonfirst.browser = whosonfirst.browser || {};

whosonfirst.browser.alt = (function(){
    
    var map;
    
    var self = {
	
	'init': function() {
	    
	    self.init_endpoints();
	    self.init_map();
	    self.init_properties();
	},
	
	'init_endpoints': function() {
	    
	    var body = document.body;
	    var root = body.getAttribute("data-whosonfirst-uri-endpoint");
	    
	    if (root){
		whosonfirst.uri.endpoint(root);
	    }			
	},
	
	'init_names': function() {
	    
	    if (typeof(whosonfirst.namify) == 'object'){
		whosonfirst.namify.namify_whosonfirst();
	    }
	},
	
	'init_map': function() {
	    map = whosonfirst.browser.common.init_map();
	},

	'init_properties': function(){
	},
    }
    
    return self;
    
})();	

