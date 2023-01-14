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

	    // map = whosonfirst.browser.common.init_map();	    
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

	    var map_el = document.getElementById("map");
	    
	    var data_url = whosonfirst.uri.id2abspath(wof_id)

	    var on_success = function(feature){
		console.log(feature);

		var bbox = whosonfirst.geojson.derive_bbox(feature);

		var bounds = [
		    [ bbox[1], bbox[0] ],
		    [ bbox[3], bbox[2] ],
		];

		var map_args = {};
		
		map = whosonfirst.browser.maps.getMap(map_el, map_args);
		map.fitBounds(bounds);
		
		var layer = L.geoJson(feature);
		layer.addTo(map);    
	    };

	    var on_error = function(err){
		console.log("SAD", err);
	    };
	    
	    whosonfirst.net.fetch(data_url, on_success, on_error);	    
	}
    }
    
    return self;
    
})();	

