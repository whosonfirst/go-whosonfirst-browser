var whosonfirst = whosonfirst || {};
whosonfirst.browser = whosonfirst.browser || {};

whosonfirst.browser.maps = (function(){

    var attribution;
   
    var maps = {};

    var self = {

	'getMap': function(map_el, args){

	    if (! args){
		args = {};
	    }

	    var api_key = args["api_key"];
	    
	    var map_id = map_el.getAttribute("id");

	    if (! map_id){
		return;
	    }
	    
	    if (maps[map_id]){
		return maps[map_id];
	    }

	    var map = L.map("map");
	    
	    var map_provider = whosonfirst.browser.common.map_provider();
	    
	    switch (map_provider){
		case "nextzen":

		    attribution = '<a href="https://github.com/tangrams" target="_blank">Tangram</a> | <a href="http://www.openstreetmap.org/copyright" target="_blank">&copy; OpenStreetMap contributors</a> | <a href="https://www.nextzen.org/" target="_blank">Nextzen</a>';
		    
		    var tangram_opts = self.getTangramOptions(args);
		    var tangramLayer = Tangram.leafletLayer(tangram_opts);
		    
		    tangramLayer.addTo(map);
		    break

		case "protomaps":

		    var tile_url = document.body.getAttribute("data-protomaps-tile-url");

		    tile_url = "/tiles/sfo/{z}/{x}/{y}.mvt";
		    
		    var layer = protomaps.leafletLayer({url:tile_url})
		    layer.addTo(map);
		    break;
		    
		default:
		    console.log("Unsupported map provider ", map_provider);
	    }
	    
	    var attribution = self.getAttribution();
	    map.attributionControl.addAttribution(attribution);

	    return map;
	},

	'getTangramOptions': function(args){

	    if (! args){
		args = {};
	    }

	    /*
	    var sceneText = await fetch(new Request('https://somwehere.com/scene.zip', { headers: { 'Accept': 'application/zip' } })).then(r => r.text());
	    var sceneURL = URL.createObjectURL(new Blob([sceneText]));
	    scene.load(sceneURL, { base_path: 'https://somwehere.com/' });
	    */
	    
	    var api_key = args["api_key"];
	    var style_url = args["style_url"];
	    var tile_url = args["tile_url"];	    
	    
	    var tangram_opts = {
		scene: {
		    import: [
			style_url,
		    ],
		    sources: {
			mapzen: {
			    url: tile_url,
			    url_subdomains: ['a', 'b', 'c', 'd'],
			    url_params: {api_key: api_key},
			    tile_size: 512,
			    max_zoom: 18
			}
		    }
		}
	    };

	    return tangram_opts;
	},

	'getAttribution': function(){
	    return attribution;
	},
    };

    return self;
    
})();
