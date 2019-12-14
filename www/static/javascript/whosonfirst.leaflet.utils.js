var whosonfirst = whosonfirst || {};
whosonfirst.leaflet = whosonfirst.leaflet || {};

whosonfirst.leaflet.utils = (function(){

	var self = {

		'draw_feature': function(map, feature, style, handler){

			if (feature["geometry"]["type"] == "Point"){
				return self.draw_point(map, feature, style, handler);
			}

			self.draw_polygon(map, feature, style);
		},
		
		'draw_point': function(map, feature, style, handler){

			if (! handler){
				return function(feature, latlon){
					return L.circle(latlon, style);
				}
			}
			
			var layer = L.geoJson(feature, {
				'style': style,
				'pointToLayer': handler,
			});

			layer.addTo(map);
			return layer;
		},

		'draw_polygon': function(map, feature, style){

			var layer = L.geoJson(feature, {
				'style': style
			});

			layer.addTo(map);
			return layer;
		},

		'draw_bbox': function(map, feature, style){

			if (typeof(whosonfirst.geojson) != 'object'){
				console.log("MISSING whosonfirst.geojson library");
				return null;
			}
			
			var bbox = whosonfirst.geojson.derive_bbox(feature);

			if (! bbox){
				console.log("no bounding box");
				return false;
			}

			var bbox = feature['bbox'];
			var swlat = bbox[1];
			var swlon = bbox[0];
			var nelat = bbox[3];
			var nelon = bbox[2];

			var geom = {
				'type': 'Polygon',
				'coordinates': [[
					[swlon, swlat],
					[swlon, nelat],
					[nelon, nelat],
					[nelon, swlat],
					[swlon, swlat],
				]]
			};

			var bbox_feature = {
				'type': 'Feature',
				'geometry': geom
			}

			return self.draw_polygon(map, bbox_feature, style);
		},

		'fit_map': function(map, feature, force){

			if (typeof(whosonfirst.geojson) != 'object'){
				console.log("MISSING whosonfirst.geojson library");
				return null;
			}

			var bbox = whosonfirst.geojson.derive_bbox(feature);
			
			if (! bbox){
				console.log("no bounding box");
				return false;
			}

			if ((bbox[1] == bbox[3]) && (bbox[2] == bbox[4])){
				map.setView([bbox[1], bbox[0]], 14);
				return;
			}

			var sw = [bbox[1], bbox[0]];
			var ne = [bbox[3], bbox[2]];

			var bounds = new L.LatLngBounds([sw, ne]);
			var current = map.getBounds();

			var redraw = true;

			if (! force){

				var redraw = false;

				/*
				  console.log("south bbox: " + bounds.getSouth() + " current: " + current.getSouth().toFixed(6));
				  console.log("west bbox: " + bounds.getWest() + " current: " + current.getWest().toFixed(6));
				  console.log("north bbox: " + bounds.getNorth() + " current: " + current.getNorth().toFixed(6));
				  console.log("east bbox: " + bounds.getEast() + " current: " + current.getEast().toFixed(6));
				*/

				if (bounds.getSouth() <= current.getSouth().toFixed(6)){
					redraw = true;
				}

				else if (bounds.getWest() <= current.getWest().toFixed(6)){
					redraw = true;
				}

				else if (bounds.getNorth() >= current.getNorth().toFixed(6)){
					redraw = true;
				}

				else if (bounds.getEast() >= current.getEast().toFixed(6)){
					redraw = true;
				}

				else {}
			}

			if (redraw){
			    var opts = { 'padding': [ 50, 50 ] };
			    map.fitBounds(bounds, opts);
			}
		}
	};

	return self;
	
})();
