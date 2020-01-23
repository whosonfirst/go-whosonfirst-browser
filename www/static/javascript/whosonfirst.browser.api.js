var whosonfirst = whosonfirst || {};
whosonfirst.browser = whosonfirst.browser || {};

whosonfirst.browser.api = (function(){
    
    var map;
    
    var self = {

	'init': function(){

	},

	'on_success_with_json': function(on_success, on_error){

	    var cb = function(raw){

		try {
		    var data = JSON.parse(raw);
		    on_success(data);
		}

		catch (e){
		    on_error(e);
		}
	    };

	    return cb;
	},
	
	// curl -s -F 'edtf:cessation=2020-01' http://localhost:8080/cessate/572206957

	'cessate': function(id, args, on_success, on_error){

	    var endpoint = "/cessate/" + id;
	    return self.execute_method('POST', endpoint, args, on_success, on_error);	    
	},

	// curl -s -F 'edtf:deprecated=2001-05-21' http://localhost:8080/deprecate/1108962799

	'deprecate': function(id, args, on_success, on_error){
	    var endpoint = "/deprecate/" + id;
	    return self.execute_method('DELETE', endpoint, args, on_success, on_error);
	},

	'update': function(id, args, on_success, on_error){

	    var endpoint = "/update/" + id;
	    return self.execute_method('POST', endpoint, args, on_success, on_error);	    	    
	},

	'execute_method': function(method, endpoint, args, on_success, on_error){

	    var form_data = new FormData();
		
	    for (key in args){
		form_data.append(key, args[key]);
	    }
	    
	    var on_load = function(rsp){

		var target = rsp.target;

		if (target.readyState != 4){
		    return;
		}

		var status_code = target['status'];
		var status_text = target['statusText'];

		if ((status_code < 200) || (status_code > 299)){
		    on_error({'code': status_code, 'message': status_text});
		    return;
		}

		var raw = target['responseText'];
		on_success(raw);
	    };

	    var on_failed = function(e){
		on_error(e);
	    };

	    var on_abort = function(){
		on_error(e);
	    };

	    var uri = location.origin + endpoint;

	    // console.log("EXECUTE", uri, args, form_data);
	    
	    try {
		var req = new XMLHttpRequest();

		req.addEventListener("load", on_load);
		req.addEventListener("error", on_failed);
		req.addEventListener("abort", on_abort);

		req.open(method, uri, true);
		req.send(form_data);
		
	    } catch (e) {
		on_error(e);
		return false;
	    }

	}
    };

    return self;
    
})();
    
