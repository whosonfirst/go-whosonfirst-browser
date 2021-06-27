var whosonfirst = whosonfirst || {};
whosonfirst.leaflet = whosonfirst.leaflet || {};

whosonfirst.leaflet.handlers = (function(){

	var self = {

		'centroid': function(style){

			return function(feature, latlon){

				var m = L.circleMarker(latlon, style);
				
				try {
					var props = feature['properties'];
					m.bindTooltip(props["wof:name"]).openTooltip();
				}
				
				catch (e){
					console.log("failed to bind label because " + e);
				}
				
				return m;
			};
		},
	};
	
	return self;
})();
