class ExistentialFlag extends HTMLElement {
    constructor() {
	
	super();

	// Y U NO ATTRIBUTES???
	var id = this.getAttribute("id");
	console.log("ID", id);
	
	const shadow = this.attachShadow({mode: 'open'});
	
	// Create text node and add word count to it
	var select = document.createElement('select');

	const flags = {
	    "-1": "unknown",
	    "0": "false",
	    "1": "true"
	};

	for (var i in flags){

	    var option = document.createElement("option");
	    option.setAttribute("value", i);
	    option.appendChild(document.createTextNode(flags[i]));

	    select.appendChild(option);
	}

	shadow.appendChild(select);
    }
}

// Define the new element
customElements.define('existential-flag', ExistentialFlag);
