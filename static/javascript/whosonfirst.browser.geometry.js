var whosonfirst = whosonfirst || {};
whosonfirst.browser = whosonfirst.browser || {};

whosonfirst.browser.geometry = (function(){
    
    var map;
    
    var self = {
	
	'init': function() {
	    
	    self.init_endpoints();
	    self.init_geometry();
	    self.init_controls();	    
	},
	
	'init_endpoints': function() {
	    
	    var body = document.body;
	    var root = body.getAttribute("data-whosonfirst-uri-endpoint");
	    
	    if (root){
		whosonfirst.uri.endpoint(root);
	    }			
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

	    if (! map_el){
		console.log("Missing 'map' element");
		return false
	    }

	    var map_args = {};
	    
	    map = whosonfirst.browser.maps.getMap(map_el, map_args);

	    var geojson_pane_name = "geometry"
	    var geojson_pane = map.createPane(geojson_pane_name);
	    geojson_pane.style.zIndex = 8000;
	    
	    var data_url = whosonfirst.uri.id2abspath(wof_id)
	    
	    var on_success = function(feature){

		var bbox = whosonfirst.geojson.derive_bbox(feature);
		
		var bounds = [
		    [ bbox[1], bbox[0] ],
		    [ bbox[3], bbox[2] ],
		];
		
		switch (map_el.getAttribute("data-map-provider")) {
			
		    case "protomaps":
			break; 	// Y U NO WORK PROTOMAPS?
		    default:
			map.fitBounds(bounds);
			break;
		}

		if (! map.pm){
		    
		    var layer = L.geoJson(feature);
		    layer.addTo(map);

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

		const geojson_layer = L.geoJson(feature, {
		    pointToLayer: (feature, latlng) => {
			if (feature.properties.customGeometry) {
			    return new L.Circle(latlng, {"pane": geojson_pane_name, "radius": feature.properties.customGeometry.radius});
			} else {
			    return new L.Marker(latlng, {"pane": geojson_pane_name});
			}
		    },
		    pane: geojson_pane_name,
		});
		
		geojson_layer.addTo(map);
	    };

	    var on_error = function(err){
		console.log("SAD", err);
	    };
	    
	    whosonfirst.net.fetch(data_url, on_success, on_error);	    
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
		
		var wof_id = pl.getAttribute("data-whosonfirst-id");
		
		if (! wof_id){
		    console.log("Missing 'data-whosonfirst-id' attribute");
		    return;
		}

		try {
		    var feature_group = map.pm.getGeomanLayers(true);
		    var feature_collection = feature_group.toGeoJSON();
		    
		    var first = feature_collection.features[0];

		    var uri = "/api/geometry/" + wof_id;

		    whosonfirst.browser.api.do("POST", uri, first).then((data) => {
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

