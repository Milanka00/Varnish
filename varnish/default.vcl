vcl 4.1;
import std;

backend default {
    .host = "router";
    .port = "9096";

}

acl purge {
    "localhost";
   # "router";
}

sub vcl_recv {

     set req.http.Host = req.http.x-wso2-actual-host;
     set req.url = req.http.x-wso2-request-path;
     set req.url = std.querysort(req.url);

     if (req.method == "PURGE") {
		# check if the client is allowed to purge content
		if (!client.ip ~ purge) {
			return(synth(405,"Not allowed."));
		}
		return (purge);
	}
}


sub vcl_miss {

    set req.http.x-wso2-cluster-header = req.http.x-wso2-actual-cluster; 
    return (fetch);
    
}

sub vcl_backend_fetch {
    if (bereq.http.x-wso2-request-path) {
        
        set bereq.url = bereq.http.x-wso2-request-path;
        set bereq.url = std.querysort(bereq.url);
        unset bereq.http.x-wso2-request-path;
        unset bereq.http.x-wso2-actual-cluster;
    }
}

sub vcl_backend_response {

    # Don't cache 404 responses
    if (beresp.status == 404) {
        set beresp.uncacheable = true;
    }

    if (bereq.http.x-cache-default-ttl) {
        set beresp.ttl = std.duration(bereq.http.x-cache-default-ttl + "s", 120s);
    }


    # Determine the appropriate storage
    if (bereq.http.x-cache-partition == "org1") {
        set beresp.storage = storage.org1;
    } else if (bereq.http.x-cache-partition == "org2") {
        set beresp.storage = storage.org2;
    } else {
        set beresp.storage = storage.default;
    }
    set beresp.http.x-storage = bereq.http.x-cache-partition;

}


sub vcl_deliver {
    if (obj.hits > 0) {
        set resp.http.X-Cached-By = "Varnish";
        set resp.http.X-Cache-Info = "Cached under host: " + req.http.Host + "; Request URI: " + req.url;
    }
}



