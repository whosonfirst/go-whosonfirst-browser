class WhosOnFirstExistentialFlag extends HTMLElement {
    
    constructor() {
	super();
    }

    connectedCallback() {
	
	var id = this.getAttribute("data-id");
	
	const shadow = this.attachShadow({mode: 'open'});
	
	var select = document.createElement('select');
	select.setAttribute("class", "form-select wof-property");

	// I don't understand why this is necessary; it appears that shadow DOM elements don't get CSS?
	select.setAttribute("style", "display: block; padding: .3rem; font-size: 1rem; width: 100%;");
	
	const flags = [ "-1", "0" , "1" ];
	
	const labels = {
	    "-1": "unknown",
	    "0": "false",
	    "1": "true"
	};

	for (var i in flags){
	    var option = document.createElement("option");
	    option.setAttribute("value", flags[i]);
	    option.appendChild(document.createTextNode(labels[flags[i]]));
	    select.appendChild(option);
	}

	select.onchange = function(){
	    var textarea = document.getElementById(id);
	    textarea.value = parseInt(select.value);
	};
	
	shadow.appendChild(select);
    }
}

// Define the new element
customElements.define('whosonfirst-existential-flag', WhosOnFirstExistentialFlag);
