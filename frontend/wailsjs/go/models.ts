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
	export class DeviceResult {
	    success: boolean;
	    message?: string;
	    devices?: audiocapture.AudioDevice[];
	    device?: audiocapture.AudioDevice;
	
	    static createFrom(source: any = {}) {
	        return new DeviceResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.devices = this.convertValues(source["devices"], audiocapture.AudioDevice);
	        this.device = this.convertValues(source["device"], audiocapture.AudioDevice);
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
	export class RecordingResult {
	    success: boolean;
	    message: string;
	    filePath?: string;
	    duration?: number;
	    dataSize?: number;
	
	    static createFrom(source: any = {}) {
	        return new RecordingResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.filePath = source["filePath"];
	        this.duration = source["duration"];
	        this.dataSize = source["dataSize"];
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

export namespace audiocapture {
	
	export class AudioDevice {
	    id: string;
	    name: string;
	    isInput: boolean;
	    isOutput: boolean;
	    isDefault: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AudioDevice(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.isInput = source["isInput"];
	        this.isOutput = source["isOutput"];
	        this.isDefault = source["isDefault"];
	    }
	}

}

