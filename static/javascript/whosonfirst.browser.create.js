var whosonfirst = whosonfirst || {};
whosonfirst.browser = whosonfirst.browser || {};

whosonfirst.browser.create = (function(){
    
    var map;
    
    var self = {
	
	'init': function() {
	    
	    self.init_endpoints();
	    self.init_map();

	    self.init_controls();	    
	},
	
	'init_endpoints': function() {
	    
	    var body = document.body;
	    var root = body.getAttribute("data-whosonfirst-uri-endpoint");
	    
	    if (root){
		whosonfirst.uri.endpoint(root);
	    }			
	},
		
	'init_map': function() {

	    // map = whosonfirst.browser.common.init_map();	    
	},
	
	'init_controls': function(){

	    self.init_save_control();
	},

	'init_save_control': function(){

	    // Eventually this should become a Leaflet control... maybe?
	    
	    var save_button = document.getElementById("save");

	    save_button.onclick = function(){

		var pl = document.getElementById("whosonfirst-place");
		
		if (! pl){
		    console.log("Missing 'whosonfirst-place' element");
		    return false;
		}
		
		try {
		    var feature_group = map.pm.getGeomanLayers(true);
		    var feature_collection = feature_group.toGeoJSON();
		    
		    var first = feature_collection.features[0];

		    var uri = "/api/create/";

		    whosonfirst.browser.api.do("PUT", uri, first).then((data) => {
			console.log("OKAY", data);
		    }).catch((err) => {
			console.log("NOT OKAY", err);
		    });			
		    
		    console.log("SAVE", wof_id, first);
		    
		} catch (err) {
		    console.log("SAD", err);
		}
		
		return false;
	    };
	},
    }
    
    return self;
    
})();	

