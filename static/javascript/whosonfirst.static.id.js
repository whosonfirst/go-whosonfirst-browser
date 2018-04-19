var whosonfirst = whosonfirst || {};
whosonfirst.static = whosonfirst.static || {};

whosonfirst.static.id = (function(){

	var map;
	
	var self = {

		'init': function() {

			self.init_endpoints();
			self.init_map();
			self.init_properties();
			self.init_names();
		},

		'init_endpoints': function() {

			var body = document.body;
			var root = body.getAttribute("data-whosonfirst-uri-endpoint");
			
			if (root){
				whosonfirst.uri.endpoint(root);
			}			
		},
		
		'init_names': function() {
			
			if (typeof(whosonfirst.namify) == 'object'){
				whosonfirst.namify.namify_whosonfirst();
			}
		},

		'init_map': function() {
			
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
			
			var body = document.body;
			var api_key = body.getAttribute("data-mapzen-api-key");
			
			/*
			L.Mapzen.apiKey = api_key;
			
			var map_opts = { tangramOptions: {
				scene: L.Mapzen.BasemapStyles.Refill
			}};
			
			map = L.Mapzen.map('map', map_opts);
			*/

			var sources = {
			    mapzen: {
				url: 'https://{s}.tile.nextzen.org/tilezen/vector/v1/512/all/{z}/{x}/{y}.topojson',
				url_subdomains: ['a', 'b', 'c', 'd'],
				url_params: {
				    api_key: api_key	// not clear this actually works... ?
				},
				tile_size: 512,
				max_zoom: 15
			    }
			};
			
			var scene = {
			    import: [
				     "/tangram/refill-style.zip",
				     "/tangram/refill-style-themes-label.zip",
				     ],
			    sources: sources,
			    global: {
				sdk_mapzen_api_key: api_key,
			    },
			};
			
			var attributions = {
			    "Tangram": "https://github.com/tangrams/",
			    "© OSM contributors": "http://www.openstreetmap.org/",
			    "Who\"s On First": "http://www.whosonfirst.org/",
			    "Nextzen": "https://nextzen.org/",
			};
			
			var attrs = [];
			
			for (var label in attributions){
			    
			    var link = attrs[label];
			    
			    if (! link){
				attrs.push(label);
				continue;
			    }
			    
			    var anchor = '<a href="' + link + '" target="_blank">' + enc_label + '</a>';
			    attrs.push(anchor);
			}
			
			var str_attributions = attrs.join(" | ");
			
			// waiting for nextzen.js to be updated to do all the things - that said it's
			// not entirely clear we need all of (map/next)zen.js and could probably get
			// away with leaflet + tangram but for now we'll just keep on as-is...
			// (20180304/thisisaaronland)
			
			L.Mapzen.apiKey = api_key;
			
			var map_opts = {
			    tangramOptions: {
				scene: scene,
				attribution: attributions,
			    }
			};
			
			map = L.Mapzen.map('map', map_opts);
			
			if ((minlat == maxlat) && (minlon == maxlon)){
			    map.setView(sw, 15);
			}
			
			else {
			    var fit_opts = {padding: [50, 50]};
			    map.fitBounds(bounds, fit_opts);
			}
			
		},

		'init_properties': function(){
			
			var place = document.getElementById("whosonfirst-place");
			var id = place.getAttribute("data-whosonfirst-id");
			
			var data_url = whosonfirst.uri.id2abspath(id)
			// console.log("FETCH", data_url);
			
			var on_success = function(feature){

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

				var feature_style = whosonfirst.leaflet.styles.polygon();
				whosonfirst.leaflet.utils.draw_feature(map, feature, feature_style);
			
				var centroid_style = whosonfirst.leaflet.styles.centroid();
				var centroid_handler = whosonfirst.leaflet.handlers.centroid(centroid_style);
				whosonfirst.leaflet.utils.draw_point(map, centroid, centroid_style, centroid_handler);

				whosonfirst.leaflet.utils.fit_map(map, feature);
				
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

				var hier_el = document.getElementById("whosonfirst-hierarchy");
				
				for (var h in props["wof:hierarchy"]){					
					var hier_pretty = whosonfirst.properties.render(props["wof:hierarchy"][h]);
					hier_el.appendChild(hier_pretty);
				}
				
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

