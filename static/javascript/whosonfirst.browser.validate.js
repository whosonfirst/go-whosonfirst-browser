// This is a local clone of whosonfirst.validate.feature.js which hooks for
// deriving the URI of the wasm binary from whosonfirst.browser.uris
var whosonfirst = whosonfirst || {};
whosonfirst.browser = whosonfirst.browser || {};

whosonfirst.browser.validate = (function(){

    var self = {

	init: function(cb){

	    var wasm_uri = whosonfirst.browser.uris.forCustomLabel("validate_wasm");
	    var custom_wasm_uri = whosonfirst.browser.uris.forCustomLabel("custom_validate_wasm");

	    return new Promise((resolve, reject) => {

		whosonfirst.browser.wasm.fetch(wasm_uri).then(rsp => {		
		    
		    if (! custom_wasm_uri){
			resolve();
			return;
		    }

		    whosonfirst.browser.wasm.fetch(custom_wasm_uri).then(rsp => {				    
			resolve();
			return;
		    }).catch(err => {
			reject(err);
		    });
		    
		}).catch(err => {
		    reject(err);
		})

	    });
	},
	},	// init
	
    };

    return self;
    
})();
