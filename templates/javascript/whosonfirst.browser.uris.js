{{ define "uris" -}}
var whosonfirst = whosonfirst || {};
whosonfirst.browser = whosonfirst.browser || {};

whosonfirst.browser.uris = (function(){

    var uris = {{ .URIs }};
    
    var self = {

	forLabel: function(label){
	    return uris[label];
	},
		   
    };

    return self;
    
})();
{{ end }}
