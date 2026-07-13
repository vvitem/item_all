export namespace main {
	
	export class AppInfo {
	    name: string;
	    tagline: string;
	    version: string;
	    commit: string;
	    buildTime: string;
	    runtime: string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new AppInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.tagline = source["tagline"];
	        this.version = source["version"];
	        this.commit = source["commit"];
	        this.buildTime = source["buildTime"];
	        this.runtime = source["runtime"];
	        this.status = source["status"];
	    }
	}

}

