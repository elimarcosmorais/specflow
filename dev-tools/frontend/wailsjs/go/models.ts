export namespace desktop {
	
	export class BootstrapInfo {
	    version: string;
	    apiBaseUrl: string;
	    token: string;
	    apiStatus: string;
	    apiPort: number;
	    databaseStatus: string;
	    logsDir: string;
	    theme: string;
	    language: string;
	
	    static createFrom(source: any = {}) {
	        return new BootstrapInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.apiBaseUrl = source["apiBaseUrl"];
	        this.token = source["token"];
	        this.apiStatus = source["apiStatus"];
	        this.apiPort = source["apiPort"];
	        this.databaseStatus = source["databaseStatus"];
	        this.logsDir = source["logsDir"];
	        this.theme = source["theme"];
	        this.language = source["language"];
	    }
	}

}

