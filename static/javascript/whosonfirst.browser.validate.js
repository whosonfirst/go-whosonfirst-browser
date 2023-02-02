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

		this.fetchWasm(wasm_uri).then(rsp => {
		    
		    if (! custom_wasm_uri){
			resolve();
			return;
		    }

		    this.fetchWasm(custom_wasm_uri).then(rsp => {
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

	fetchWasm: function(wasm_uri){
	    
	    var pending = 1;
	    console.log("Fetch WASM ", wasm_uri);
	    
	    return new Promise((resolve, reject) => {
		
		if (! WebAssembly.instantiateStreaming){
		    
		    WebAssembly.instantiateStreaming = async (resp, importObject) => {
			const source = await (await resp).arrayBuffer();
			return await WebAssembly.instantiate(source, importObject);
		    };
		}
		
		const export_go = new Go();
		
		let export_mod, export_inst;	
		
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
