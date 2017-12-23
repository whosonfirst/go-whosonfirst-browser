var whosonfirst = whosonfirst || {};
whosonfirst = whosonfirst || {};

whosonfirst.namify = (function() {

	var cache_ttl = 30000;
	
	var self = {
		
		'init': function(){
			
		},
		
		'namify_whosonfirst': function(){
			
			var resolver = whosonfirst.uri.id2abspath;
			
			var els = document.getElementsByClassName("whosonfirst-namify");
			var count = els.length;

			for (var i=0; i < count; i++){
				self.namify_el(els[i], resolver);
			}
		},

		/*
		'namify_iso_countries': function(){
			
			var resolver = whosonfirst.uri.id2abspath;
			
			var els = document.getElementsByClassName("whosonfirst-namify-country");
			var count = els.length;
			
			for (var i=0; i < count; i++){
				self.namify_el(els[i], resolver);
			}
		},
		
		'namify_brands': function(){
			
			var resolver = whosonfirst.brands.id2abspath;
			
			var els = document.getElementsByClassName("whosonfirst-namify-brand");
			var count = els.length;
			
			for (var i=0; i < count; i++){
				
				self.namify_el(els[i], resolver);
			}
		},
		*/
		
		'namify_el': function(el, resolver){
			
			var wofid = el.getAttribute("data-whosonfirst-id");
			
			if (! wofid){	
				return;
			}

			/*
			if (el.textContent != wofid){
				return;
			}
			*/
			
			var url = resolver(wofid);
			
			var on_hit = function(feature){
				self.apply_namification(el, feature);
			};
			
			var on_miss = function(){
				self.namify_el_from_source(url, el);
			};
			
			if (! self.cache_get(url, on_hit, on_miss)){
				self.namify_el_from_source(url, el);
			}
			
		},
		
		'namify_el_from_source': function(url, el){
			
			var on_fetch = function(feature){
				
				self.apply_namification(el, feature);
				self.cache_set(url, feature);
			};
			
			var on_fail = function(rsp){
				// console.log("sad face");
			};
			
			whosonfirst.net.fetch(url, on_fetch, on_fail);
		},
		
		'apply_namification': function(el, feature){
			
			var props = feature['properties'];
			
			// to account for whosonfirst-brands which needs to be updated
			// to grow a 'properties' hash... (20160319/thisisaaronland)
			
			if (! props){
				props = feature;
			}
			
			// console.log(props);
			
			var label = props['wof:label'];
			
			if ((! label) || (label == '')){
				
				var possible = [
					'wof:name',
					'wof:brand_name'
				];
				
				var count = possible.length;
				
				for (var i = 0; i < count; i++) {
					
					var k = possible[i];
					
					if (props[k]){
						label = props[k];
						break;
					}
				}
			}
			
			// var enc_label = whosonfirst.php.htmlspecialchars(label);
			el.innerHTML = label;
		},
		
		'cache_args': function(){
			
			return {
				'suffix': 'whosonfirst.namify',
				'ttl': cache_ttl,
			};
		},
		
		'cache_get': function(key, on_hit, on_miss){
			
			if (typeof(whosonfirst.cache) != 'object'){
				console.log("CACHE GET", key, "no whosonfirst.cache");
				return false;
			}
			
			var args = self.cache_args();
			return whosonfirst.cache.get(key, on_hit, on_miss, args);
		},

		'cache_set': function(key, value){

			if (typeof(whosonfirst.cache) != 'object'){
				console.log("CACHE SET", key, "no whosonfirst.cache");				
				return false;
			}

			var args = self.cache_args();			
			return whosonfirst.cache.set(key, value, args);
		},

		'cache_unset': function(key){

			if (typeof(whosonfirst.cache) != 'object'){
				console.log("CACHE UNSET", key, "no whosonfirst.cache");								
				return false;
			}

			var args = self.cache_args();
			return whosonfirst.cache.unset(key, args);
		}	    
    };

    return self;

})();
