window.addEventListener("load", function load(event){

    whosonfirst.browser.validate.init().then(rsp => {
	whosonfirst.browser.create.init();
    }).catch(err => {
	console.log("Failed to initialize validation code", err);
    });

});
