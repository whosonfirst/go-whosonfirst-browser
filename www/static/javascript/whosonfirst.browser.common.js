var whosonfirst = whosonfirst || {};
whosonfirst.browser = whosonfirst.browser || {};

whosonfirst.browser.common = (function(){
    
    var map;
    
    var self = {

	'init_map': function(){
	    
	    var api_key = document.body.getAttribute("data-nextzen-api-key");
	    var style_url = document.body.getAttribute("data-nextzen-style-url");
	    var tile_url = document.body.getAttribute("data-nextzen-tile-url");    
	    
	    if (! api_key){
		console.log("Missing API key");
		return;
	    }
	    
	    if (! style_url){
		console.log("Missing style URL");
		return;
	    }
	    
	    if (! tile_url){
		console.log("Missing tile URL");
		return;
	    }

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
	    
	    var map_args = {
		"api_key": api_key,
		"style_url": style_url,
		"tile_url": tile_url,
	    };
	    
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
    };

    return self;

})();    
