var whosonfirst = whosonfirst || {};
whosonfirst = whosonfirst || {};

whosonfirst.cache = (function() {

	var cache_ttl = 30000;
	var local_cache = {};
	
	var self = {
		
		'init': function(){
			
		},
		
		'get': function(key, on_hit, on_miss, args){

			var fq_key = self.prep_key(key, args);
			console.log("CACHE GET", fq_key);
			
			var handle_rsp = function(rsp){
				
				if (! rsp){
					console.log("CACHE MISS", fq_key);
					on_miss();
					return;
				}
				
				var data = rsp['data'];
				
				if (! data){
					console.log("CACHE MISS", fq_key);
					on_miss();
					return;
				}
				
				var dt = new Date();
				var ts = dt.getTime();
				
				var then = rsp['created'];
				var diff = ts - then;
				
				if (diff > cache_ttl){
					console.log("CACHE EXPIRED", fq_key);					
					self.unset(key, args);
					on_miss();
					return;
				}
				
				on_hit(data);
			};
			
			if (typeof(localforage) == 'object'){
		    
				localforage.getItem(fq_key, function (err, rsp){
					
					if (err){
						on_miss();
					}

					handle_rsp(rsp);
				});
			}

			else {

				var rsp = local_cache[key];
				handle_rsp(rsp);
			}
			
			return true;
		},
		
		'set': function(key, value, args){

			var fq_key = self.prep_key(key, args);
			console.log("CACHE SET", fq_key);
			
			var dt = new Date();
			var ts = dt.getTime();
			
			var wrapper = {
				'data': value,
				'created': ts
			};		

			if (typeof(localforage) == 'object'){
				localforage.setItem(fq_key, wrapper);				
			}

			else {
				local_cache[fq_key] = wrapper;
			}
			
			return true;
		},
	    
		'unset': function(key, args){

			var fq_key = self.prep_key(key, args);
			console.log("CACHE UNSET", fq_key);
			
			if (typeof(localforage) != 'object'){
				localforage.removeItem(fq_key);				
			}

			else {
				delete(local_cache[fq_key]);
			}		
			
			return true;
		},
		
		'prep_key': function(key, args){

			if (! args){
				args = {};
			}

			var suffix = args["suffix"];

			if (! suffix){
				suffix = location.host;
			}
			
			return key + "#" + suffix;
		}
	};
	
	return self;
	
})();
