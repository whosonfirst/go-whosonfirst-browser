class WhosOnFirstPlacetype extends HTMLElement {
    
    constructor() {
	super();
    }

    connectedCallback() {

	var wasm_uri = whosonfirst.browser.uris.forCustomLabel("placetypes_wasm");

	whosonfirst.browser.wasm.fetch(wasm_uri).
		    then(() => {

			console.log("WOO")
			whosonfirst_placetype_descendants("continent","common,common_optional,optional").
											  then((data) => {
											      console.log(data);
											  });
			
		    }).
		    catch((err) => {
			console.log("Failed to fetch ", wasm_uri, err);
		    });

	/*
	var id = this.getAttribute("data-id");
	
	const shadow = this.attachShadow({mode: 'open'});
	
	var select = document.createElement('select');
	select.setAttribute("class", "form-select wof-property");

	// Get placetypes here...
	
	select.onchange = function(){
	    var textarea = document.getElementById(id);
	    textarea.value = parseInt(select.value);
	};
	
	shadow.appendChild(select);
	*/
    }
}

// Define the new element
customElements.define('whosonfirst-placetype', WhosOnFirstPlacetype);
