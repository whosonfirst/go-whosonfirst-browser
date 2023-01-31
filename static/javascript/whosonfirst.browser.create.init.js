window.addEventListener("load", function load(event){

    // whosonfirst.browser.create.init();
    
    whosonfirst.validate.feature.init()
	       .then(rsp => {
		   whosonfirst.browser.create.init();
	       })
	       .catch(err => {
		   console.log("Failed to initialize validation code", err);
	       });

});
