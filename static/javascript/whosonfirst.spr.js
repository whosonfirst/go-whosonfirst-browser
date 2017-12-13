window.addEventListener("load", function load(event){

	var geom = document.getElementById("whosonfirst-geom");

	if (! geom){
		return;
	}

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

	var body = document.body;
	var api_key = body.getAttribute("data-mapzen-apikey");

	L.Mapzen.apiKey = api_key;

	var map_opts = { tangramOptions: {
		scene: L.Mapzen.BasemapStyles.Refill
	}};
	
	var map = L.Mapzen.map('map', map_opts);

	var fit_opts = {padding: [50, 50]};
	map.fitBounds(bounds, fit_opts);

	var place = document.getElementById("whosonfirst-place");
	var path = place.getAttribute("data-whosonfirst-path");

	var data_url = "https://data.whosonfirst.org/" + path;
	console.log("FETCH", data_url);
	
	var on_success = function(feature){
		
		console.log("SUCCESS", feature);

		var feature_layer = L.geoJSON(feature);
		feature_layer.addTo(map);

		var centroid_layer = L.geoJSON(centroid);
		centroid_layer.addTo(map);

		var props = feature["properties"];
		var props_str = JSON.stringify(props, null, "\t");
		
		var props_pre = document.createElement("pre");
		props_pre.appendChild(document.createTextNode(props_str));
		
		var props_el = document.getElementById("whosonfirst-properties-raw");
		props_el.appendChild(props_pre);
	};

	var on_error = function(rsp){
		
		console.log("ERROR", rsp);

		var centroid_layer = L.geoJSON(centroid);
		centroid_layer.addTo(map);				
	};

	var req = new XMLHttpRequest();

	req.onload = function(){

		try {
			var data = JSON.parse(this.responseText);
			on_success(data);
		}
		
		catch (e){
			console.log("FAIL", this);
			console.log("ERROR", e);
			on_error();
		}
	};

	req.open("GET", data_url, true);
	req.send();
	
});
