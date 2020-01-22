var whosonfirst = whosonfirst || {};
whosonfirst.browser = whosonfirst.browser || {};

whosonfirst.browser.id = (function(){
    
    var map;
    
    var self = {
	
	'init': function() {
	    
	    self.init_endpoints();
	    self.init_map();
	    self.init_properties();
	    self.init_update_controls();
	    self.init_names();
	},
	
	'init_endpoints': function() {
	    
	    var body = document.body;
	    var root = body.getAttribute("data-whosonfirst-uri-endpoint");
	    
	    if (root){
		whosonfirst.uri.endpoint(root);
	    }			
	},

	'init_update_controls': function(){

	    self.init_cessate_controls();
	    self.init_deprecate_controls();	    
	},

	'init_cessate_controls': function(){

	    var els = document.getElementsByClassName("cessate");
	    var count = els.length;

	    for (var i=0; i < count; i++){

		var el = els[i];
		el.onclick = function(e){

		    var el = e.target;
		    var parent = el.parentNode;
		    
		    var wof_id = el.getAttribute("data-whosonfirst-id");
		    var wof_name = el.getAttribute("data-whosonfirst-name");
		    
		    if (! confirm("Are you sure you want to cessate " + wof_name  + " ?")){
			return false;
		    }

		    // TO DO: ASSIGN CUSTOM DATE?
				   
		    var on_success = function(rsp){
			
			// TO DO: REDRAW PROPERTIES... or not?
			// maybe just refreshing the page is the
			// simplest dumbest thing?

			// parent.removeChild(el);
			
			location.href = location.href;
		    };

		    var on_error = function(err){

			// TO DO: WHERE DO ERRORS GET REPORTED/DISPLAYED?			
			console.log("ERROR", err);
		    };

		    var parse_on_success = whosonfirst.browser.api.on_success_with_json(on_success, on_error);		    
		    whosonfirst.browser.api.cessate(wof_id, {}, parse_on_success, on_error);
		    
		    console.log("CESSATE", wof_id);
		    return false;
		};
	    }
	},

	'init_deprecate_controls': function(){

	    var els = document.getElementsByClassName("deprecate");
	    var count = els.length;

	    for (var i=0; i < count; i++){

		var el = els[i];
		el.onclick = function(e){

		    var el = e.target;
		    var parent = el.parentNode;
		    
		    var wof_id = el.getAttribute("data-whosonfirst-id");
		    var wof_name = el.getAttribute("data-whosonfirst-name");
		    
		    if (! confirm("Are you sure you want to deprecate " + wof_name  + " ?")){
			return false;
		    }

		    // TO DO: ASSIGN CUSTOM DATE?
		    
		    var on_success = function(rsp){

			// TO DO: REDRAW PROPERTIES... or not?
			// maybe just refreshing the page is the
			// simplest dumbest thing?
			
			// parent.removeChild(el);

			location.href = location.href;
		    };

		    var on_error = function(err){

			// TO DO: WHERE DO ERRORS GET REPORTED/DISPLAYED?
			console.log("ERROR", err);
		    };

		    var parse_on_success = whosonfirst.browser.api.on_success_with_json(on_success, on_error);
		    
		    whosonfirst.browser.api.deprecate(wof_id, {}, parse_on_success, on_error);
		    
		    console.log("DEPRECATE", wof_id);
		    return false;
		};
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
	    
	    var place = document.getElementById("whosonfirst-place");
	    var id = place.getAttribute("data-whosonfirst-id");
	    
	    var data_url = whosonfirst.uri.id2abspath(id)
	    console.log("FETCH", data_url);
	    
	    var on_success = function(feature){
		
		var geom = document.getElementById("whosonfirst-place");

		// remember we might not have a nextzen API key and by extension
		// no map - see above in init_map (20191219/straup)
		
		if (map){
		    
		    var lat = geom.getAttribute("data-latitude");
		    var lon = geom.getAttribute("data-longitude");
		    
		    var centroid_geom = { "type": "Point", "coordinates": [ lon, lat ] };
		    var centroid_props = {};
		    
		    var centroid = {
			"type": "Feature",
			"geometry": centroid_geom,
			"properties": centroid_props
		    };

		    var feature_style = whosonfirst.leaflet.styles.polygon();
		    var centroid_style = whosonfirst.leaflet.styles.centroid();
		    var centroid_handler = whosonfirst.leaflet.handlers.centroid(centroid_style);
		    
		    if (feature["geometry"]["type"] == "Point"){
			whosonfirst.leaflet.utils.draw_point(map, feature, feature_style, centroid_handler);
		    }
		    
		    else if (feature["geometry"]["type"] == "MultiPoint"){
			whosonfirst.leaflet.utils.draw_point(map, feature, feature_style, centroid_handler);
		    }

		    else {
			whosonfirst.leaflet.utils.draw_feature(map, feature, feature_style);
		    }
		    
		    whosonfirst.leaflet.utils.draw_point(map, centroid, centroid_style, centroid_handler);
		    
		    whosonfirst.leaflet.utils.fit_map(map, feature);
		}
		
		// sudo put all of this in a wapper function somewhere...
		// (20171224/thisisaaronland)
		
		var props = feature["properties"];
		var props_str = JSON.stringify(props, null, "\t");
		
		var props_raw = document.createElement("pre");
		props_raw.appendChild(document.createTextNode(props_str));
		
		var raw_el = document.getElementById("whosonfirst-properties-raw");
		var pretty_el = document.getElementById("whosonfirst-properties-pretty");
		
		var button_raw = document.createElement("button");
		button_raw.setAttribute("class", "raw-pretty");
		button_raw.appendChild(document.createTextNode("show pretty"));
		
		button_raw.onclick = function(){
		    raw_el.style.display = "none";
		    pretty_el.style.display = "block";					
		};
		
		var button_pretty = document.createElement("button");
		button_pretty.setAttribute("class", "raw-pretty");				
		button_pretty.appendChild(document.createTextNode("show raw"));
		
		button_pretty.onclick = function(){
		    pretty_el.style.display = "none";
		    raw_el.style.display = "block";					
		};
		
		try {
		    var props_pretty = whosonfirst.properties.render(props);
		    
		    if (props_pretty){
			raw_el.style.display = "none";
			
			raw_el.appendChild(button_raw);
			raw_el.appendChild(props_raw);
			
			pretty_el.appendChild(button_pretty);
			pretty_el.appendChild(props_pretty);
		    }
		    
		    else {
			throw "Failed to generate pretty properties";
		    }
		}
		
		catch (e) {
		    raw_el.appendChild(button_raw);				
		    raw_el.appendChild(props_raw);
		    console.log("PRETTY", "ERR", e);
		}

		var wof_el = document.getElementById("whosonfirst-wof");
		var src_el = document.getElementById("whosonfirst-src");
		
		var wof_props = {};
		var src_props = {};
		
		for (k in props){

		    var parts = k.split(":");
		    var ns = parts[0];
		    var label = parts[1];
		    
		    switch (ns){
			case "src":
			    src_props[k] = props[k];
			    break;
			case "wof":
			    wof_props[k] = props[k];
			    break;
			default:
			    continue;
		    }
		}

		var append = function(el, props){

		    var pretty = whosonfirst.properties.render(props);
		    el.appendChild(pretty);

		    var details = pretty.getElementsByTagName("details");
		    var details_count = details.length;
		
		    for (var i=0; i < details_count; i++){		    
			var d = details[i];
			d.setAttribute("open", "true");
		    }
		    
		};

		append(src_el, src_props);
		append(wof_el, wof_props);
		
		self.init_names();				
	    };
	    
	    var on_error = function(rsp){
		
		console.log("ERROR", rsp);
		
		var centroid_layer = L.geoJSON(centroid);
		centroid_layer.addTo(map);				
	    };
	    
	    whosonfirst.net.fetch(data_url, on_success, on_error);
	}
    }
    
    return self;
    
})();	

