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
		},
		markerStyle: {
		    draggle: true,
		    icon: whosonfirst.browser.leaflet.markerIcon(),
		    pane: geojson_pane_name,
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

	    self.init_property_controls();	    	    
	    self.init_save_controls();
	},

	'init_property_controls': function(){

	    var optional_els = document.getElementsByClassName("wof-property-optional");
	    var count_optional = optional_els.length;

	    for (var i=0; i < count_optional; i++){
		self.add_remove_button(optional_els[i]);
	    }

	    var add_button = document.getElementById("add-property");

	    if (! add_button){
		console.log("Missing 'add-property' element.");
		return false;
	    };
	    
	    add_button.onclick = function(){

		var form = document.getElementById("wof-properties-form");

		if (! form){
		    whosonfirst.browser.feedback.emit("Failed to locate 'wof-properties-form' element.");
		    return false;
		}
		
		var id = prompt("Property name");

		if (! id){
		    return false;
		}

		var existing_el = document.getElementById(id);

		if (existing_el){
		    alert("Property with this name already exists");
		    return false;
		}
		
		var property_type = "string";	// To do: Make this an option (once we've figured out what those options are/need to be)

		var new_el = self.create_property_el(id, property_type);

		form.appendChild(new_el);
		return false;
	    };
	},

	'create_property_el': function(id, property_type){

	    var ta = document.createElement("textarea");
	    ta.setAttribute("class", "form-control-plaintext wof-property wof-property-" + property_type);
	    ta.setAttribute("id", id);
	    
	    var l = document.createElement("label");
	    l.setAttribute("for", id);
	    l.appendChild(document.createTextNode(id));
	    
	    var el = document.createElement("div");
	    el.setAttribute("class", "col-auto wof-property-block");
	    el.setAttribute("id", id + "-wrapper");

	    el.appendChild(l);
	    el.appendChild(ta);

	    self.add_remove_button(el);
	    return el;
	},

	'add_remove_button': function(el){

	    // This is all a bit brittle but will have to do for now...
	    
	    var labels = el.getElementsByTagName("label");
	    var count = labels.length;

	    if (count != 1){
		console.log("Invalid count for labels", count);
	    }

	    var l = labels[0];

	    var id = el.getAttribute("id");
	    id = id.replace("-wrapper", "");
	    
	    var remove = document.createElement("span");
	    remove.setAttribute("class", "remove-property");
	    remove.setAttribute("data-property", id);
	    remove.appendChild(document.createTextNode("[x]"));

	    remove.onclick = function(e){
		
		var this_el = e.target;
		var prop = this_el.getAttribute("data-property");

		var wrapper = document.getElementById(prop + "-wrapper");

		if (! wrapper){
		    return false;
		}

		if (confirm("Are you sure you want to remove this property?")){
		    wrapper.remove();
		}
		
		return false;
	    };

	    l.appendChild(remove);
	},
	
	'init_save_controls': function(){

	    // Eventually this should become a Leaflet control... maybe?
	    
	    var save_button = document.getElementById("save");

	    save_button.onclick = function(){
		
		var props = {};
		
		var inputs = document.getElementsByClassName("wof-property");
		var count = inputs.length;

		for (var i=0; i < count; i++){

		    var el = inputs[i];
		    var el_class = el.getAttribute("class");

		    var k = el.getAttribute("id");
		    var v = el.value;

		    // console.log("Debug", k, v);
		    
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

		    var do_create = function(){

			whosonfirst.browser.api.do("PUT", create_uri, feature).then((data) => {
			    var props = data["properties"];
			    var id = props["wof:id"];
			    whosonfirst.browser.feedback.emit("New feature created with ID " + id);
			    save_button.removeAttribute("disabled");			
			}).catch((err) => {
			    whosonfirst.browser.feedback.emit("Failed to create new feature", err);
			    save_button.removeAttribute("disabled");						
			});			
			
			whosonfirst.browser.feedback.emit("Creating new feature...");
			save_button.setAttribute("disabled", "disabled");
			
		    };
		    
		    validate_feature(str_f).then(rsp => {

			/*

			   See this? As of this writing it is expected that:
			   
			   1) in whosonfirst.browser.validate.js
			   2) if whosonfirst.browser.uris.forCustomLabel("custom_validate_wasm") is non-nil
			   3) there is a wasm binary at the end of that URI
			   4) that wasm binary exports a function named "whosonfirst_validate_feature_custom"

			   Is this the best way to do this? It doesn't feel like. It doesn't feel elegant
			   at least. But, for now, we'll live it...
			   
			 */
			
			if (typeof(whosonfirst_validate_feature_custom) != "function"){
			    do_create();
			    return;
			}
			
			whosonfirst_validate_feature_custom(str_f).then(rsp => {
			    do_create();
			}).catch(err => {
			    whosonfirst.browser.feedback.emit("Document failed custom validation:", err);
			})
			
		    }).catch(err => {
			whosonfirst.browser.feedback.emit("Document failed to validate:", err);
			console.log("Document failed to validate", err, feature);
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

