var whosonfirst = whosonfirst || {};

whosonfirst.render = (function(){

    var _custom_renderers = {
	'dict': function(d, ctx){ return null; },
	'text': function(d, ctx){ return null; },
    };

    var _exclusions = {
	'text': function(d, ctx){ return null; },
    };

    var self = {

	'enabled': function(bool){

	    if (typeof(bool) != "undefined"){
		if (bool){
		    _enabled = true;
		} else {
		    _enabled = false;
		}
	    }

	    return _enabled;
	},
	
	'set_custom_renderers': function(t, r){

	    if (! _custom_renderers[t]){
		return;
	    }

	    if (! r){
		return;
	    }

	    _custom_renderers[t] = r;
	},

	'get_custom_renderer': function(t, d, ctx){

	    if (! _custom_renderers[t]){
		return null;
	    }

	    var custom = _custom_renderers[t];
	    return custom(d, ctx);
	},

	'set_custom_exclusions': function(t, e){

	    if (! _exclusions[t]){
		return;
	    }

	    if ((! e) || (typeof(e) != "function")){
		return;
	    }

	    _exclusions[t] = e;
	},

	'get_custom_exclusion': function(t, d, ctx){

	    if (! _exclusions[t]){
		return null;
	    }

	    var exclude =  _exclusions[t];
	    return exclude(d, ctx);
	},
	
	'apply': function(data, target){
	    
	    var el = document.getElementById(target);
	    
		if (! el){
		return false;
	    }
	    
	    var pretty = self.render(data);
	    el.appendChild(pretty);

	    return true;
	},
	
	'render': function(props){

	    var pretty = document.createElement("div");
	    pretty.setAttribute("id", "render-pretty");

	    buckets = self.bucket_props(props);
	    
	    var namespaces = Object.keys(buckets);
	    namespaces = namespaces.sort();

	    var count_ns = namespaces.length;
	    
	    for (var i=0; i < count_ns; i++){
		var ns = namespaces[i];
		var dom = self.render_bucket(ns, buckets[ns]);
		pretty.appendChild(dom);
	    }

	    return pretty;				
	},

	'render_bucket': function(ns, bucket){

	    // var wrapper = document.createElement("div");

	    var details = document.createElement("details");	    
	    details.setAttribute("style", "margin-bottom:1rem;");
	    
	    if (ns != '_global_'){
		
		var summary = document.createElement("summary");
		summary.setAttribute("id", ns);
		
		var anchor = document.createElement("a");
		anchor.setAttribute("href", "#" + ns);
		anchor.setAttribute("class", "anchor");
		anchor.appendChild(document.createTextNode("¶"));

		// TO DO: write whosonfirst.sources.js and
		// assign the full label for 'ns' to 'content'
		// (20191214/straup)
		
		var content = document.createTextNode(ns);

		var header = document.createElement("h3");		
		header.setAttribute("style", "display:inline");
		
		header.appendChild(content);
		// header.appendChild(anchor);
		
		summary.appendChild(header);
		details.appendChild(summary);			
	    }

	    var menu = document.createElement("details-menu");
	    
	    var sorted = self.sort_bucket(bucket);
	    var body = self.render_data(sorted, ns);
	    
	    menu.appendChild(body);
	    details.appendChild(menu);
	    
	    return details;
	},
	
	'render_data': function(d, ctx){
	    
	    if (Array.isArray(d)){
		// console.log("render list for " + ctx);
		return self.render_list(d, ctx);
	    }
	    
	    else if (typeof(d) == "object"){
		// console.log("render dict for " + ctx);
		return self.render_dict(d, ctx);
	    }
	    
	    else {
		// console.log("render text for " + ctx);

		var wrapper = document.createElement("span");
		wrapper.setAttribute("class", "render-content");
		
		var content;
		
		var renderer = self.get_custom_renderer('text', d, ctx);
		// console.log("rendered for " + ctx + " : " + typeof(renderer), renderer);

		if (renderer){
		    // console.log("try to render for " + ctx + " with renderer");
		    try {
			content = renderer(d, ctx);
		    } catch (e) {
			// console.log("UNABLE TO RENDER " + ctx + " BECAUSE " + e);
		    }
		}

		else {
		    // console.log("render " + ctx + " as plain text");		    
		    content = self.render_text(d, ctx);
		}

		wrapper.appendChild(content);
		return wrapper;
	    }
	},
	
	'render_dict': function(d, ctx){
	    
	    var table = document.createElement("table");
	    table.setAttribute("class", "table table-sm");
	    
	    for (k in d){
		
		var row = document.createElement("tr");
		var label_text = k;

		var _ctx = (ctx) ? ctx + "." + k : k;

		var renderer = self.get_custom_renderer('dict', d, _ctx);

		if (renderer){
		    try {
			label_text = renderer(d, _ctx);
		    } catch (e) {
			// console.log("UNABLE TO RENDER " + _ctx + " BECAUSE " + e);
		    }
		}

		/*
		  unclear if the rule should just be only text (as it currently is)
		  or whether custom markup is allowed... smells like feature quicksand
		  so moving along for now (20160211/thisisaaronland)
		 */

		var header = document.createElement("th");
		var label = document.createTextNode(self.htmlspecialchars(label_text));
		header.appendChild(label);
		
		var content = document.createElement("td");

		var body = self.render_data(d[k], _ctx);		
		content.appendChild(body);

		row.appendChild(header);
		row.appendChild(content);
		
		table.appendChild(row);
	    }

	    var wrapper = document.createElement("div");
	    wrapper.setAttribute("class", "table-responsive");

	    wrapper.appendChild(table);
	    return wrapper;
	},
	
	'render_list': function(d, ctx){
	    
	    var count = d.length;
	    
	    if (count == 0){
		return self.render_text("–", ctx);
	    }
	    
	    if (count <= 1){
		return self.render_data(d[0], ctx);
	    }
	    
	    var list = document.createElement("ul");
	    
	    for (var i=0; i < count; i++){
		
		var item = document.createElement("li");
		var body = self.render_data(d[i], ctx + "#" + i);
		
		item.appendChild(body);
		list.appendChild(item);
	    }
	    
	    return list;
	},
	
	'render_text': function(d, ctx){
	    
	    var text = self.htmlspecialchars(d);
	    
	    var span = document.createElement("span");
	    span.setAttribute("id", ctx);
	    span.setAttribute("title", ctx);
	    span.setAttribute("class", "yesnofix-uoc");
	    	    
	    var el = document.createTextNode(text);
	    span.appendChild(el);

	    return span;
	},
	
	'render_link': function(link, text, ctx){

	    var anchor = document.createElement("a");
	    anchor.setAttribute("href", link);
	    anchor.setAttribute("target", "_wof");
	    var body = self.render_text(text, ctx);
	    anchor.appendChild(body);

	    return anchor;
	},

	'render_code': function(text, ctx){
	    
	    var code = document.createElement("code");
	    var body = self.render_text(text, ctx);
	    code.appendChild(body);
	    return code;
	},
	
	'render_timestamp': function(text, ctx){
	    var dt = new Date(parseInt(text) * 1000);
	    return self.render_text(dt.toISOString(), ctx);
	},
	
	'bucket_props': function(props){
	    
	    buckets = {};
	    
	    for (k in props){
		parts = k.split(":", 2);
		
		ns = parts[0];
		pred = parts[1];
		
		if (parts.length != 2){
		    ns = "_global_";
		    pred = k;
		}
		
		if (! buckets[ns]){
		    buckets[ns] = {};					
		}
		
		buckets[ns][pred] = props[k];
	    }
	    
	    return buckets;
	},
	
	'sort_bucket': function(bucket){
	    
	    var sorted = {};
	    
	    var keys = Object.keys(bucket);
	    keys = keys.sort();
	    
	    var count_keys = keys.length;
	    
	    for (var j=0; j < count_keys; j++){
		var k = keys[j];
		sorted[k] = bucket[k];
	    }
	    
	    return sorted;
	},

	'htmlspecialchars': function(string, quote_style, charset, double_encode){
	    //       discuss at: http://phpjs.org/functions/htmlspecialchars/
	    //      original by: Mirek Slugen
	    //      improved by: Kevin van Zonneveld (http://kevin.vanzonneveld.net)
	    //      bugfixed by: Nathan
	    //      bugfixed by: Arno
	    //      bugfixed by: Brett Zamir (http://brett-zamir.me)
	    //      bugfixed by: Brett Zamir (http://brett-zamir.me)
	    //       revised by: Kevin van Zonneveld (http://kevin.vanzonneveld.net)
	    //         input by: Ratheous
	    //         input by: Mailfaker (http://www.weedem.fr/)
	    //         input by: felix
	    // reimplemented by: Brett Zamir (http://brett-zamir.me)
	    //             note: charset argument not supported
	    //        example 1: htmlspecialchars("<a href='test'>Test</a>", 'ENT_QUOTES');
	    //        returns 1: '&lt;a href=&#039;test&#039;&gt;Test&lt;/a&gt;'
	    //        example 2: htmlspecialchars("ab\"c'd", ['ENT_NOQUOTES', 'ENT_QUOTES']);
	    //        returns 2: 'ab"c&#039;d'
	    //        example 3: htmlspecialchars('my "&entity;" is still here', null, null, false);
	    //        returns 3: 'my &quot;&entity;&quot; is still here'
	    
	    var optTemp = 0,
	    i = 0,
	    noquotes = false;
	    if (typeof quote_style === 'undefined' || quote_style === null) {
		quote_style = 2;
	    }
	    string = string.toString();
	    if (double_encode !== false) {
		// Put this first to avoid double-encoding
		string = string.replace(/&/g, '&amp;');
	    }
	    string = string.replace(/</g, '&lt;')
		.replace(/>/g, '&gt;');
	    
	    var OPTS = {
		'ENT_NOQUOTES'          : 0,
		'ENT_HTML_QUOTE_SINGLE' : 1,
		'ENT_HTML_QUOTE_DOUBLE' : 2,
		'ENT_COMPAT'            : 2,
		'ENT_QUOTES'            : 3,
		'ENT_IGNORE'            : 4
	    };
	    if (quote_style === 0) {
		noquotes = true;
	    }
	    if (typeof quote_style !== 'number') {
		// Allow for a single string or an array of string flags
		quote_style = [].concat(quote_style);
		for (i = 0; i < quote_style.length; i++) {
		    // Resolve string input to bitwise e.g. 'ENT_IGNORE' becomes 4
		    if (OPTS[quote_style[i]] === 0) {
			noquotes = true;
		    } else if (OPTS[quote_style[i]]) {
			optTemp = optTemp | OPTS[quote_style[i]];
		    }
		}
		quote_style = optTemp;
	    }
	    if (quote_style & OPTS.ENT_HTML_QUOTE_SINGLE) {
		string = string.replace(/'/g, '&#039;');
	    }
	    if (!noquotes) {
		string = string.replace(/"/g, '&quot;');
	    }
	    
	    return string;
	}
	
    }
    
    return self;
    
})();
