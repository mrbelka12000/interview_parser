export namespace app {
	
	export class APIKeyResult {
	    success: boolean;
	    message?: string;
	    apiKey?: string;
	    lastUpdated?: string;
	
	    static createFrom(source: any = {}) {
	        return new APIKeyResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.apiKey = source["apiKey"];
	        this.lastUpdated = source["lastUpdated"];
	    }
	}
	export class FileContent {
	    name: string;
	    path: string;
	    isDir: boolean;
	    size: number;
	    extension: string;
	    content: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new FileContent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.isDir = source["isDir"];
	        this.size = source["size"];
	        this.extension = source["extension"];
	        this.content = source["content"];
	        this.error = source["error"];
	    }
	}
	export class FileInfo {
	    name: string;
	    path: string;
	    isDir: boolean;
	    size: number;
	    extension: string;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.isDir = source["isDir"];
	        this.size = source["size"];
	        this.extension = source["extension"];
	    }
	}
	export class TranscriptionResult {
	    success: boolean;
	    message: string;
	    transcriptPath?: string;
	    analysisPath?: string;
	
	    static createFrom(source: any = {}) {
	        return new TranscriptionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.transcriptPath = source["transcriptPath"];
	        this.analysisPath = source["analysisPath"];
	    }
	}

}

