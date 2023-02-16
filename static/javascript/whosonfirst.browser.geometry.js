var whosonfirst = whosonfirst || {};
whosonfirst.browser = whosonfirst.browser || {};

whosonfirst.browser.geometry = (function(){
    
    var map;
    var events = {};
    
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

	    var edit_el = document.getElementById("edit-geometry");
	    
	    if (! edit_el){
		whosonfirst.browser.feedback.emit("Missing 'edit-geometry' element");
		return false;
	    }
	    
	    var wof_id = edit_el.getAttribute("data-whosonfirst-id");
	    
	    if (! wof_id){
		whosonfirst.browser.feedback.emit("Missing 'data-whosonfirst-id' attribute");
		return;
	    }
	    
	    var map_el = document.getElementById("map");

	    if (! map_el){
		whosonfirst.browser.feedback.emit("Missing 'map' element");
		return false
	    }

	    // To do: Determine how much of this code can be reconciled with whosonfirst.browser.create
	    
	    L.PM.setOptIn(true);
	    
	    var map_args = {
		pmIgnore: false,
	    };
	    
	    map = whosonfirst.browser.maps.getMap(map_el, map_args);

	    var geojson_pane_name = "geometry"
	    var geojson_pane = map.createPane(geojson_pane_name);
	    geojson_pane.style.zIndex = 8000;
	    
	    var data_url = whosonfirst.uri.id2abspath(wof_id)
	    
	    var on_success = function(feature){

		whosonfirst.browser.feedback.clear();
		
		var props = feature.properties;
		var name = props["wof:name"];
		var id = props["wof:id"];

		var name_el = document.getElementById("wof:name");
		var id_el = document.getElementById("wof:id");		

		if (name_el){
		    name_el.innerText = name;
		}

		if (id_el){
		    id_el.innerText = id;
		}
		
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
		    
		    var layer = L.geoJSON(feature);
		    layer.addTo(map);

		    whosonfirst.browser.feedback.emit("Missing map.pm");
		    return;
		}
		
		var on_update = function(){
		    
		    var feature_group = map.pm.getGeomanLayers(true);
		    var feature_collection = feature_group.toGeoJSON();

		    var geojson_el = document.getElementById("geojson");

		    if (! geojson_el){
			return;
		    }
		    
		    var enc_featurecollection = JSON.stringify(feature_collection, " ", 2);
		    geojson_el.innerText = enc_feature_collection;
		};

		map.pm.setGlobalOptions({
		    'panes': {
			vertexPane: geojson_pane_name,
			layerPane: geojson_pane_name,
			markerPane: geojson_pane_name,
		    }
		});
		
		map.pm.addControls({
		    position: 'topleft',
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

		const geojson_layer = L.geoJson(feature, {
		    pointToLayer: (feature, latlng) => {
			if (feature.properties.customGeometry) {
			    return new L.Circle(latlng, {
				pane: geojson_pane_name,
				radius: feature.properties.customGeometry.radius,
				pmIgnore: false,
			    });
			} else {
			    
			    return new L.Marker(latlng, {
				icon: whosonfirst.browser.leaflet.markerIcon(),
				pane: geojson_pane_name,
				pmIgnore: false,				
			    });
			}
		    },
		    pane: geojson_pane_name,
		    pmIgnore: false,
		});

		geojson_layer.addTo(map);

		if (events["load"]){

		    var count = events["load"].length;
		    
		    for (var i=0; i < count; i++){
			events["load"][i](feature);
		    }
		}
	    };

	    var on_error = function(err){
		whosonfirst.browser.feedback.emit("Failed to load record", err);		
	    };

	    whosonfirst.browser.feedback.emit("Loading record");
	    whosonfirst.net.fetch(data_url, on_success, on_error);	    
	},

	'init_controls': function(){

	    self.init_save_control();
	},

	'init_save_control': function(){

	    // Eventually this should become a Leaflet control... maybe?
	    
	    var save_button = document.getElementById("save");

	    save_button.onclick = function(){
		
		var edit_el = document.getElementById("edit-geometry");

		if (! edit_el){
		    whosonfirst.browser.feedback.emit("Missing 'edit-geometry' element");
		    return false;
		}
		
		var wof_id = edit_el.getAttribute("data-whosonfirst-id");
		
		if (! wof_id){
		    whosonfirst.browser.feedback.emit("Missing 'data-whosonfirst-id' attribute");
		    return;
		}

		try {
		    
		    var feature_group = map.pm.getGeomanLayers(true);
		    var feature_collection = feature_group.toGeoJSON();

		    // START OF reconcile with whosonfirst.browser.create
		    
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

		    // END OF reconcile with whosonfirst.browser.create
		    
		    var feature = {
			'type': 'Feature',
			'properties': {},
			'geometry': geom,
		    };

		    var edit_uri = whosonfirst.browser.uris.forLabel("edit_geometry_api") + wof_id;

		    whosonfirst.browser.api.do("POST", edit_uri, feature)
			       .then((data) => {
				   whosonfirst.browser.feedback.emit("Geometry has been successfully updated");
				   save_button.removeAttribute("disabled");			
			       })
			       .catch((err) => {
				   whosonfirst.browser.feedback.emit("Failed to update geometry", err);
				   save_button.removeAttribute("disabled");			
			       });			
		    
		    whosonfirst.browser.feedback.emit("Updating geometry...");				
		    save_button.setAttribute("disabled", "disabled");
		    
		} catch (err) {
		    whosonfirst.browser.feedback.emit("Failed to parse record before submitting update", err);
		}

		return false;
	    };
	},

	on: function(target, func){

	    var addEvent = function(target, func){

		funcs = events[target];

		if (! funcs){
		    funcs = [];
		}

		funcs.push(func)
		events[target] = funcs;
	    };
	    
	    switch (target){
		case 'load':
		    addEvent(target, func);
		    break;
		default:
		    console.log("Unsupported target", target);
		    break;
	    }
	},

    }
    
    return self;
    
})();	

