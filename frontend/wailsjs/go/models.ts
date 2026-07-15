export namespace main {
	
	export class AppInfo {
	    name: string;
	    tagline: string;
	    version: string;
	    commit: string;
	    buildTime: string;
	    runtime: string;
	    status: string;
	    error?: storage.SafeError;
	
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
	        this.error = this.convertValues(source["error"], storage.SafeError);
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

}

export namespace storage {
	
	export class SafeError {
	    code: string;
	    safeMessage: string;
	    retryable: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SafeError(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.safeMessage = source["safeMessage"];
	        this.retryable = source["retryable"];
	    }
	}

}
