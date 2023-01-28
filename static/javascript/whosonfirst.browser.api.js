var whosonfirst = whosonfirst || {};
whosonfirst.browser = whosonfirst.browser || {};

whosonfirst.browser.api = (function(){

    var self = {

	'do': function(http_method, rel_url, body){

	    return new Promise((resolve, reject) => {
		
		var abs_url = self.abs_url(rel_url);

		var onload = function(){

		    if (this.status >= 300){
			reject(this.responseText);
			return;
		    }
		    
		    var rsp;
		    
		    try {
			rsp = JSON.parse(this.responseText);
            	    } catch (err){
			console.log("ERR", abs_url, err);
			reject(err);
			return false;
		    }
		    
		    resolve(rsp);
       		};
		
		var onprogress = function(rsp){
		    // console.log("progress");
		};
		
		var onfailed = function(rsp){
		    reject("Connection failed " + rsp);
		};
		
		var onabort = function(rsp){
		    reject("Connection aborted " + rsp);
		};
		
		var req = new XMLHttpRequest();
		
		// https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequest/withCredentials		    
		req.withCredentials = true;
		
		req.addEventListener("load", onload);
		req.addEventListener("progress", onprogress);
		req.addEventListener("error", onfailed);
		req.addEventListener("abort", onabort);
		
		req.open(http_method, abs_url, true);
		
		var enc_body = JSON.stringify(body);
		req.send(enc_body);	    
	    });
	},

	'abs_url': function(rel_url){

	    // To do: read root URL from ... ?
	    
	    return rel_url;
	},
    };
    
    return self;
    
})();
