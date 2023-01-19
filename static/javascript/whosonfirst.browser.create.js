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
		console.log("UPDATE", feature_collection);
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
		
		var inputs = document.getElementsByClassName("wof:property");
		var count = inputs.length;

		for (var i=0; i < count; i++){

		    var el = inputs[i];
		    var k = el.getAttribute("id");
		    var v = el.value;

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
			    console.log("Missing geometry");
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
		    
		    var uri = "/api/create/";

		    whosonfirst.browser.api.do("PUT", uri, feature).then((data) => {
			console.log("OKAY", data);
		    }).catch((err) => {
			console.log("NOT OKAY", err);
		    });			
		    
		} catch (err) {
		    console.log("SAD", err);
		}
		
		return false;
	    };
	},
    }
    
    return self;
    
})();	

