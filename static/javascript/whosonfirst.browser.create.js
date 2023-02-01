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
		whosonfirst.browser.feedback.emit("Missing 'map' element");
		return false
	    }

	    // To do: Determine how much of this code can be reconciled with whosonfirst.browser.geometry
	    
	    L.PM.setOptIn(true);
	    
	    var map_args = {
		pmIgnore: false,
	    };
	    
	    map = whosonfirst.browser.maps.getMap(map_el, map_args);

	    var geojson_pane_name = "geometry"
	    var geojson_pane = map.createPane(geojson_pane_name);
	    geojson_pane.style.zIndex = 8000;
	    
	    var on_update = function(){
		var feature_group = map.pm.getGeomanLayers(true);
		var feature_collection = feature_group.toGeoJSON();
		// console.log("UPDATE", feature_collection);
	    };
	    
	    map.pm.setGlobalOptions({
		'panes': {
		    vertexPane: geojson_pane_name,
		    layerPane: geojson_pane_name,
		    markerPane: geojson_pane_name,
		}
	    });

	    map.on('pm:create', (e) => {
		e.layer.options.pmIgnore = false;
		L.PM.reInitLayer(e.layer);
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

	    map.pm.addControls({
		position: 'topleft',
	    });
	    
	},
	
	'init_controls': function(){

	    self.init_save_control();
	},
	
	'init_save_control': function(){

	    // Eventually this should become a Leaflet control... maybe?
	    
	    var save_button = document.getElementById("save");

	    save_button.onclick = function(){
		
		var props = {};
		
		var inputs = document.getElementsByClassName("wof-property");
		var count = inputs.length;

		for (var i=0; i < count; i++){

		    var el = inputs[i];
		    var k = el.getAttribute("id");
		    var v = el.value;

		    if (v == ""){
			continue;
		    }

		    // Account for arbitrary JSON in textarea elements
		    
		    if ((v.startsWith("[")) || (v.startsWith("{"))){
			
			try {
			    v = JSON.parse(v);
			} catch (err) {
			    whosonfirst.browser.feedback.emit("Failed to parse '" + k + "=" + v + "' property, " + err);
			    return false;
			}
		    }
		    
		    props[k] = v;
		}

		try {
		    
		    var feature_group = map.pm.getGeomanLayers(true);
		    var feature_collection = feature_group.toGeoJSON();

		    // START OF reconcile with whosonfirst.browser.geometry
		    
		    var count = feature_collection.features.length;
		    var geom;

		    switch (count){
			case 0:
			    whosonfirst.browser.feedback.emit("Missing geometry");
			    return false;
			    break;
			case 1:
			    geom = feature_collection.features[0].geometry;
			    break;
			default:

			    _geoms = [];

			    for (var i=0; i < count; i++){
				_geoms.push(feature_collection.features[i].geometry);
			    }

			    geom = {
				'type':'GeometryCollection',
				'geometries': _geoms,
			    };

			    break;
		    }

		    // END OF reconcile with whosonfirst.browser.geometry
		    
		    var feature = {
			'type': 'Feature',
			'properties': props,
			'geometry': geom,
		    };
		    
		    var create_uri = whosonfirst.browser.uris.forLabel("create_feature_api");

		    var str_f = JSON.stringify(feature);

		    // https://github.com/whosonfirst/go-whosonfirst-validate-wasm
		    // https://github.com/whosonfirst/go-whosonfirst-validate
		    
		    validate_feature(str_f).then(rsp => {

			whosonfirst.browser.api.do("PUT", create_uri, feature)
				   .then((data) => {
				       var props = data["properties"];
				       var id = props["wof:id"];
				       whosonfirst.browser.feedback.emit("New feature created with ID " + id);
				       save_button.removeAttribute("disabled");			
				   })
				   .catch((err) => {
				       whosonfirst.browser.feedback.emit("Failed to create new feature", err);
				       save_button.removeAttribute("disabled");						
				   });			
			
			whosonfirst.browser.feedback.emit("Creating new feature...");
			save_button.setAttribute("disabled", "disabled");
			
		    }).catch(err => {

			whosonfirst.browser.feedback.emit("Document failed to validate:", err);			
		    });
		   		    
		} catch (err) {
		    whosonfirst.browser.feedback.emit("Unable to prepare data to create new feature", err);
		}
		
		return false;
	    };
	},
    }
    
    return self;
    
})();	

