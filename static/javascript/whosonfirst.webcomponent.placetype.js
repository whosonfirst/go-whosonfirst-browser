class WhosOnFirstPlacetype extends HTMLElement {
    
    constructor() {
	super();
    }

    connectedCallback() {

	var wasm_uri = whosonfirst.browser.uris.forCustomLabel("placetypes_wasm");
	var _self = this;

	// Once we have placetypes, replace textarea with select menu
	
	var placetypes_cb = function(placetypes){

	    var count = placetypes.length;
	    
	    var id = _self.getAttribute("data-id");
	
	    const shadow = _self.attachShadow({mode: 'open'});
	    
	    var select = document.createElement('select');
	    select.setAttribute("class", "form-select wof-property");

	    // I don't understand why this is necessary; it appears that shadow DOM elements don't get CSS?
	    select.setAttribute("style", "display: block; padding: .3rem; font-size: 1rem; width: 100%;");

	    var opt = document.createElement("option");
	    select.appendChild(opt);
	    
	    for (var i=0; i < count; i++){
		var pt = placetypes[i];
		var name = pt["name"];

		var opt = document.createElement("option");
		opt.value = name;
		opt.appendChild(document.createTextNode(name));
		select.appendChild(opt);
	    }
	    
	    select.onchange = function(){
		var textarea = document.getElementById(id);
		textarea.value = parseInt(select.value);
	    };
	    
	    shadow.appendChild(select);
	};

	// Once we've loaded the whosonfirst/go-whosonfirst-placetypes-wasm WASM binary
	// fetch all the descendants for the 'planet' placetype
	
	var wasm_cb = function(){

	    whosonfirst_placetypes_descendants("planet","common,common_optional,optional").then((data) => {
		var placetypes = JSON.parse(data);
		placetypes_cb(placetypes);
	    }).catch((err)=> {
		console.log(err);
	    });
			
	};

	// Fetch the whosonfirst/go-whosonfirst-placetypes-wasm WASM binary

	sfomuseum.wasm.fetch(wasm_uri). then(() => {	    
	    wasm_cb();
	}).catch((err) => {
	    console.log("Failed to fetch ", wasm_uri, err);
	});

    }
}

// Define the new element
customElements.define('whosonfirst-placetype', WhosOnFirstPlacetype);
