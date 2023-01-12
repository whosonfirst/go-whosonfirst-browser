var whosonfirst = whosonfirst || {};
whosonfirst.browser = whosonfirst.browser || {};

whosonfirst.browser.common = (function(){
    
    var map;
    
    var self = {

	'init_map': function(){

	    var map_svg = document.getElementById("map-svg");

	    if (map_svg){
		
		var parent = map_svg.parentNode;

		if (parent){
		    parent.removeChild(map_svg);
		}
	    }
	    
	    var geom = document.getElementById("whosonfirst-place");
	    
	    var lat = geom.getAttribute("data-latitude");
	    var lon = geom.getAttribute("data-longitude");
	    
	    var centroid_geom = { "type": "Point", "coordinates": [ lon, lat ] };
	    var centroid_props = {};
	    
	    var centroid = {
		"type": "Feature",
		"geometry": centroid_geom,
		"properties": centroid_props
	    };
	    
	    var minlat = geom.getAttribute("data-min-latitude");
	    var minlon = geom.getAttribute("data-min-longitude");
	    
	    var maxlat = geom.getAttribute("data-max-latitude");
	    var maxlon = geom.getAttribute("data-max-longitude");	
	    
	    minlat = parseFloat(minlat);
	    minlon = parseFloat(minlon);
	    maxlat = parseFloat(maxlat);
	    maxlon = parseFloat(maxlon);	
	    
	    var sw = [ minlat, minlon ];
	    var ne = [ maxlat, maxlon ];
	    
	    var bounds = [ sw, ne ];
	    
	    var map_el = document.getElementById("map");
	    
	    if (! map_el){
		console.log("Missing map element");	
		return;
	    }
	    
	    var map_args = {};
	    
	    map = whosonfirst.browser.maps.getMap(map_el, map_args);

	    if (! map){
		console.log("Failed to get map");
		return false;
	    }

	    if ((minlat == maxlat) && (minlon == maxlon)){
		map.setView(sw, 15);
	    }
	    
	    else {
		var fit_opts = {padding: [50, 50]};
		map.fitBounds(bounds, fit_opts);
	    }

	    return map;
	},

	'init_properties': function(){
	    
	    var place = document.getElementById("whosonfirst-place");
	    var id = place.getAttribute("data-whosonfirst-id");
	    
	    var is_alt = place.getAttribute("data-whosonfirst-is-alternate");	    
	    var src = place.getAttribute("data-whosonfirst-alt-source");
	    var func = place.getAttribute("data-whosonfirst-alt-function");	    

	    var uri_args = {
		'alt': is_alt,
		'source': src,
		'function': func,
		'extras': [],
	    };

	    var data_url = whosonfirst.uri.id2abspath(id, uri_args)
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
		var wof_props = {};
		
		for (k in props){

		    var parts = k.split(":");

		    if (parts[0] != "wof"){
			continue;
		    }

		    var label = parts[1];
		    wof_props[k] = props[k];
		}
		
		var pretty = whosonfirst.properties.render(wof_props);

		var details = pretty.getElementsByTagName("details");
		var details_count = details.length;
		
		for (var i=0; i < details_count; i++){		    
		    var d = details[i];
		    d.setAttribute("open", "true");
		}
		
		wof_el.appendChild(pretty);
		
		self.init_names();				
	    };
	    
	    var on_error = function(rsp){
		
		console.log("ERROR", rsp);
		
		var centroid_layer = L.geoJSON(centroid);
		centroid_layer.addTo(map);				
	    };
	    
	    whosonfirst.net.fetch(data_url, on_success, on_error);
	}
	
    };

    return self;

})();    
