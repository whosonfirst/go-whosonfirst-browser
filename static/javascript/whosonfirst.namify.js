var whosonfirst = whosonfirst || {};
whosonfirst = whosonfirst || {};

/*
  things this expects

  - whosonfirst.uri
  - whosonfirst.brands (for brands)
  - whosonfirst.log
  - whosonfirst.net
  - whosonfirst.php
  - localforage

*/

whosonfirst.namify = (function() {

    var cache_ttl = 30000;

    var self = {
	
	'init': function(){

	},
	
	'namify_wof': function(){

	    var resolver = whosonfirst.uri.id2abspath;

	    var els = document.getElementsByClassName("whosonfirst-namify");
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

	'namify_el': function(el, resolver){

	    var wofid = el.getAttribute("data-whosonfirst-id");

	    if (! wofid){	
		return;
	    }

	    if (el.textContent != wofid){
		return;
	    }

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

	'cache_get': function(key, on_hit, on_miss){

	    if (typeof(localforage) != 'object'){
		return false;
	    }

	    var fq_key = self.cache_prep_key(key);

	    localforage.getItem(fq_key, function (err, rsp){

		if ((err) || (! rsp)){
		    on_miss();
		}

		var data = rsp['data'];

		if (! data){
		    on_miss();
		}

		var dt = new Date();
		var ts = dt.getTime();

		var then = rsp['created'];
		var diff = ts - then;

		if (diff > cache_ttl){
		    self.cache_unset(key);
		    on_miss();
		}

		on_hit(data);
	    });

	    return true;
	},

	'cache_set': function(key, value){

	    if (typeof(localforage) != 'object'){
		return false;
	    }

	    var dt = new Date();
	    var ts = dt.getTime();

	    var wrapper = {
		'data': value,
		'created': ts
	    };

	    key = self.cache_prep_key(key);

	    localforage.setItem(key, wrapper);
	    return true;
	},

	'cache_unset': function(key){

	    if (typeof(localforage) != 'object'){
		return false;
	    }

	    key = self.cache_prep_key(key);

	    localforage.removeItem(key);
	    return true;
	},

	'cache_prep_key': function(key){
	    return key + '#namify';
	}
    };

    return self;

})();
