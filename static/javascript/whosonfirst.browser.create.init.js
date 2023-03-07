window.addEventListener("load", function load(event){

    var wasm_uri = whosonfirst.browser.uris.forCustomLabel("validate_wasm");
    var custom_wasm_uri = whosonfirst.browser.uris.forCustomLabel("custom_validate_wasm");

    sfomuseum.wasm.fetch(wasm_uri).then(rsp => {				
		    
	if (! custom_wasm_uri){
	    whosonfirst.browser.create.init();	    
	    return;
	}
	
	sfomuseum.wasm.fetch(custom_wasm_uri).then(rsp => {
	    whosonfirst.browser.create.init();	    	    
	    return;
	}).catch(err => {
	    console.log("Failed to initialize custom validation code", err);
	});
	
    }).catch(err => {
	console.log("Failed to initialize validation code", err);	
	reject(err);
    });
    
});
