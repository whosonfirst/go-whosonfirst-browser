var whosonfirst = whosonfirst || {};
whosonfirst.browser = whosonfirst.browser || {};

whosonfirst.browser.geometry = (function(){
    
    var map;
    
    var self = {
	
	'init': function() {
	    
	    self.init_endpoints();
	    self.init_map();
	    self.init_geometry();
	},
	
	'init_endpoints': function() {
	    
	    var body = document.body;
	    var root = body.getAttribute("data-whosonfirst-uri-endpoint");
	    
	    if (root){
		whosonfirst.uri.endpoint(root);
	    }			
	},
		
	'init_map': function() {

	    map = whosonfirst.browser.common.init_map();	    
	},
	
	'init_geometry': function(){

	    var pl = document.getElementById("whosonfirst-place");

	    if (! pl){
		console.log("Missing 'whosonfirst-place' element");
		return false;
	    }

	    var wof_id = pl.getAttribute("data-whosonfirst-id");

	    if (! wof_id){
		console.log("Missing 'data-whosonfirst-id' attribute");
		return;
	    }

	    var data_url = whosonfirst.uri.id2abspath(id)
	    console.log("FETCH", data_url);
	    
	}
    }
    
    return self;
    
})();	

