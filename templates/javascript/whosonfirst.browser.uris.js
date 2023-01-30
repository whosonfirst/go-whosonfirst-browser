{{ define "uris" -}}
var whosonfirst = whosonfirst || {};
whosonfirst.browser = whosonfirst.browser || {};

whosonfirst.browser.uris = (function(){

    var uris = {{ .URIs }};
    
    var self = {

	forLabel: function(label){
	    return uris[label];
	},

	forCustomLabel: function(label){

	    if (!uris['custom']){
		return null;
	    }

	    return uris['custom'][label];
	},
    };

    return self;
    
})();
{{ end }}
