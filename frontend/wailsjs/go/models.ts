export namespace cloudflare {
	
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

