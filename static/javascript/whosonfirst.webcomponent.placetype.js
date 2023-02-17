class WhosOnFirstPlacetype extends HTMLElement {
    
    constructor() {
	super();
    }

    connectedCallback() {
	
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
    }
}

// Define the new element
customElements.define('whosonfirst-placetype', WhosOnFirstPlacetype);
