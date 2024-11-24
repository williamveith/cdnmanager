export namespace cloudflare {
	
	export class ResponseInfo {
	    code: number;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ResponseInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	    }
	}
	export class Response {
	    success: boolean;
	    errors: ResponseInfo[];
	    messages: ResponseInfo[];
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.errors = this.convertValues(source["errors"], ResponseInfo);
	        this.messages = this.convertValues(source["messages"], ResponseInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class StorageKey {
	    name: string;
	    expiration: number;
	    metadata: any;
	
	    static createFrom(source: any = {}) {
	        return new StorageKey(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.expiration = source["expiration"];
	        this.metadata = source["metadata"];
	    }
	}

}

export namespace session {
	
	export class Entry {
	    Name: string;
	    Metadata: any;
	    Value: string;
	
	    static createFrom(source: any = {}) {
	        return new Entry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Metadata = source["Metadata"];
	        this.Value = source["Value"];
	    }
	}

}

