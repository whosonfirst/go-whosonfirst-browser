var whosonfirst = whosonfirst || {};

whosonfirst.net = (function(){

	var default_cache_ttl = 30000; // ms

	var self = {

		'encode_query': function(query){

			enc = new Array();

			for (var k in query){
				var v = query[k];
				v = encodeURIComponent(v);
				enc.push(k + "=" + v);
			}

			return enc.join("&");
		},
		
		'fetch': function(url, on_success, on_fail, args){

		    	if (typeof(args) == "undefined") {
			    args = {};
			}

			var on_hit = function(data){

				if (on_success){
					on_success(data);
				}
			};

			var on_miss = function(){
				self.fetch_with_xhr(url, on_success, on_fail, args);
			};

			if (! self.cache_get(url, on_hit, on_miss, args)){
				self.fetch_with_xhr(url, on_success, on_fail, args);
			}
		},

		'fetch_with_xhr': function(url, on_success, on_fail, args){

			if (! args){
			    args = {};
			}

			var req = new XMLHttpRequest();

			req.onload = function(){

				try {
					var data = JSON.parse(this.responseText);
				}

				catch (e){

					if (on_fail){
						on_fail({
							url: url,
							args: args,
							xhr: req
						});
					}

					return false;
				}

				self.cache_set(url, data);

				if (on_success){
					on_success(data);
				}
			};

			try {

			    	if (args["cache-busting"]){

				    var cb = Math.floor(Math.random() * 1000000);

				    var tmp = document.createElement("a");
				    tmp.href = url;

				    if (tmp.search){
					tmp.search += "&cb=" + cb;
				    }

				    else {
					tmp.search = "?cb= " + cb;
				    }

				    url = tmp.href;
				}

				req.open("get", url, true);
				req.send();
			}

			catch(e){

				if (on_fail){
					on_fail();
				}
			}
		},

		'cache_args': function(){

			return {
				'suffix': 'whosonfirst.net',
				'cache_ttl': default_cache_ttl,
			};
		},
		
		'cache_get': function(key, on_hit, on_miss, user_args){

			if (typeof(whosonfirst.cache) != 'object'){
				console.log("CACHE GET", key, "no whosonfirst.cache");
				return false;
			}

			var args = self.cache_args();

			if (typeof(user_args) == 'object'){
				
				for (k in user_args){
					args[k] = user_args[k];
				}
			}
			
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
