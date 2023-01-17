var whosonfirst = whosonfirst || {};
whosonfirst.browser = whosonfirst.browser || {};

whosonfirst.browser.api = (function(){

    var self = {

	'do': function(http_method, rel_url, body){

	    return new Promise((resolve, reject) => {
		
		var abs_url = self.abs_url(rel_url);
	    
		var req = new XMLHttpRequest();
		
		req.onload = function(){
		    
		    var rsp;
		    
		    try {
			rsp = JSON.parse(this.responseText);
            	    }
		    
		    catch (e){
			console.log("ERR", abs_url, e);
			reject(e);
			return false;
		    }
		    
		    resolve(rsp);
       		};
		
		req.open(http_method, abs_url, true);
				
		var enc_body = JSON.stringify(body);
		req.send(enc_body);	    
	    });
	},

	'abs_url': function(rel_url){
	    return rel_url;
	},
    };
    
    return self;
    
})();
