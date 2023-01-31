var whosonfirst = whosonfirst || {};
whosonfirst.validate = whosonfirst.validate || {};

whosonfirst.validate.feature = (function(){

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
		
		WebAssembly.instantiateStreaming(fetch("wasm/validate_feature.wasm"), export_go.importObject).then(
		    
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
