var whosonfirst = whosonfirst || {};
whosonfirst.browser = whosonfirst.browser || {};

whosonfirst.browser.maps = (function(){

    var attribution;
   
    var maps = {};

    var self = {

	'getMap': function(map_el, args){
	    return aaronland.maps.getMap(map_el, args);
	},

	'getAttribution': function(){
	    return attribution;
	},
    };

    return self;
    
})();
