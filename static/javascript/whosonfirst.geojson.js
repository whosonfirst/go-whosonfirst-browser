var whosonfirst = whosonfirst || {};

whosonfirst.geojson = (function(){

	var self = {

		'derive_bbox': function(feature){

			if (feature['bbox']){
				return feature['bbox'];
			}

			if (feature['type'] == 'FeatureCollection'){

				var features = feature['features'];
				var count = features.length;

				var swlat = undefined;
				var swlon = undefined;
				var nelat = undefined;
				var nelon = undefined;

				for (var i=0; i < count; i++){

					var bbox = self.derive_bbox(features[i]);

					var _swlat = bbox[0];
					var _swlon = bbox[1];
					var _nelat = bbox[2];
					var _nelon = bbox[3];

					if ((! swlat) || (_swlat < swlat)){
						swlat = _swlat;
					}

					if ((! swlon) || (_swlon < swlon)){
						swlon = _swlon;
					}

					if ((! nelat) || (_nelat > nelat)){
						nelat = _nelat;
					}

					if ((! nelon) || (_nelon > nelon)){
						nelon = _nelon;
					}
				}

				return [ swlat, swlon, nelat, nelon ];
			}

			else if (feature['type'] == 'Feature'){

				// Adapted from http://gis.stackexchange.com/a/172561
				var geom = feature['geometry'];
				var coords = geom.coordinates;

				var lats = [],
				    lngs = [];

				for (var i = 0; i < coords[0].length; i++) {
					lats.push(coords[0][i][1]);
					lngs.push(coords[0][i][0]);
				}

				var minlat = Math.min.apply(null, lats),
				    maxlat = Math.max.apply(null, lats);
				var minlng = Math.min.apply(null, lngs),
				    maxlng = Math.max.apply(null, lngs);

				return [ minlat, minlng,
					 maxlat, maxlng ];
			}

			else {}
		}

	};

	return self;

})();
