var whosonfirst = whosonfirst || {};
whosonfirst.browser = whosonfirst.browser || {};

whosonfirst.browser.feedback = (function(){
    
    var self = {

	emit: function(){

	    // To do: Add code to finesse each element in arguments
	    // in to a string, for example "error" objects or dictionaries.
	    
	    var msg = Array.from(arguments).join(" ");

	    var dt = new Date();
	    console.log(dt.toUTCString(), msg);
	    
	    var feedback_el = document.getElementById("feedback");

	    if (feedback_el){
		feedback.innerText = msg;
	    }
	},

	clear: function(){

	    var feedback_el = document.getElementById("feedback");

	    if (feedback_el){
		feedback_el.innerHTML = '';
	    }
	    
	}
	
    };

    return self;

})();    
