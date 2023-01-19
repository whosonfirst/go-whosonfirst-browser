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

	    var map_el = document.getElementById("map");

	    if (! map_el){
		console.log("Missing 'map' element");
		return false
	    }

	    var map_args = {};
	    
	    map = whosonfirst.browser.maps.getMap(map_el, map_args);

	    if (! map.pm){
		console.log("Missing map.pm");
		return;
	    }
	    
	    var on_update = function(){
		var feature_group = map.pm.getGeomanLayers(true);
		var feature_collection = feature_group.toGeoJSON();
		console.log("UPDATE", feature_collection);
	    };
	    
	    map.pm.addControls({
		position: 'topleft',
	    });
	    
	    map.on("pm:drawend", function(e){
		console.log("draw end");
		on_update();
	    });
	    
	    map.on('pm:remove', function (e) {
		console.log("remove");
		on_update();
	    });
	    
	    // This does not appear to capture drag or edit-vertex events
	    // Not sure what's up with that...
	    
	    map.on('pm:globaleditmodetoggled', (e) => {
		console.log("remove");
		on_update();
	    });	    
	    
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

