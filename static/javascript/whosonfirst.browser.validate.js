// This is a local clone of whosonfirst.validate.feature.js which hooks for
// deriving the URI of the wasm binary from whosonfirst.browser.uris
var whosonfirst = whosonfirst || {};
whosonfirst.browser = whosonfirst.browser || {};

whosonfirst.browser.validate = (function(){

    var self = {

	init: function(cb){

	    var pending = 1;
	    
	    return new Promise((resolve, reject) => {
		
		if (! WebAssembly.instantiateStreaming){
		    
		    WebAssembly.instantiateStreaming = async (resp, importObject) => {
			const source = await (await resp).arrayBuffer();
			return await WebAssembly.instantiate(source, importObject);
		    };
		}
		
		const export_go = new Go();
		
		let export_mod, export_inst;	

		var wasm_uri = whosonfirst.browser.uris.forCustomLabel("validate_wasm");
		
		WebAssembly.instantiateStreaming(fetch(wasm_uri), export_go.importObject).then(
		    
		    async result => {
			
			pending -= 1;
			
			if (pending == 0){
			    resolve();
			}
			
			export_mod = result.module;
			export_inst = result.instance;
			await export_go.run(export_inst);
		    }
		);
		
	    });
	    
	},	// init
	
    };

    return self;
    
})();
