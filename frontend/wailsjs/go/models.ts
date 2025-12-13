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
	export class CallAnalysisResult {
	    success: boolean;
	    message: string;
	    analysisPath?: string;
	
	    static createFrom(source: any = {}) {
	        return new CallAnalysisResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.analysisPath = source["analysisPath"];
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

export namespace models {
	
	export class AnalyzeInterview {
	    id: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new AnalyzeInterview(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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
	export class QuestionAnswer {
	    id: number;
	    interview_id: number;
	    question: string;
	    full_answer: string;
	    accuracy: number;
	    reason_unanswered: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new QuestionAnswer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.interview_id = source["interview_id"];
	        this.question = source["question"];
	        this.full_answer = source["full_answer"];
	        this.accuracy = source["accuracy"];
	        this.reason_unanswered = source["reason_unanswered"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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
	export class AnalyzeInterviewWithQA {
	    id: number;
	    qa: QuestionAnswer[];
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new AnalyzeInterviewWithQA(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.qa = this.convertValues(source["qa"], QuestionAnswer);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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
	export class Call {
	    id: number;
	    transcript: string;
	    analysis: number[];
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Call(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.transcript = source["transcript"];
	        this.analysis = source["analysis"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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
	export class GlobalAnalytics {
	    totalInterviews: number;
	    totalQuestions: number;
	    totalAnswered: number;
	    totalUnanswered: number;
	    globalAnsweredPercent: number;
	    globalAverageAccuracy: number;
	    globalAnsweredAccuracy: number;
	    bestInterviewID: number;
	    bestInterviewScore: number;
	    worstInterviewID: number;
	    worstInterviewScore: number;
	    // Go type: time
	    lastUpdated: any;
	
	    static createFrom(source: any = {}) {
	        return new GlobalAnalytics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalInterviews = source["totalInterviews"];
	        this.totalQuestions = source["totalQuestions"];
	        this.totalAnswered = source["totalAnswered"];
	        this.totalUnanswered = source["totalUnanswered"];
	        this.globalAnsweredPercent = source["globalAnsweredPercent"];
	        this.globalAverageAccuracy = source["globalAverageAccuracy"];
	        this.globalAnsweredAccuracy = source["globalAnsweredAccuracy"];
	        this.bestInterviewID = source["bestInterviewID"];
	        this.bestInterviewScore = source["bestInterviewScore"];
	        this.worstInterviewID = source["worstInterviewID"];
	        this.worstInterviewScore = source["worstInterviewScore"];
	        this.lastUpdated = this.convertValues(source["lastUpdated"], null);
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
	export class InterviewAnalytics {
	    id: number;
	    totalQuestions: number;
	    answeredQuestions: number;
	    unansweredQuestions: number;
	    answeredPercentage: number;
	    unansweredPercentage: number;
	    averageAccuracy: number;
	    averageAnsweredAccuracy: number;
	    highConfidenceQuestions: number;
	    mediumConfidenceQuestions: number;
	    lowConfidenceQuestions: number;
	    questionsWithReason: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new InterviewAnalytics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.totalQuestions = source["totalQuestions"];
	        this.answeredQuestions = source["answeredQuestions"];
	        this.unansweredQuestions = source["unansweredQuestions"];
	        this.answeredPercentage = source["answeredPercentage"];
	        this.unansweredPercentage = source["unansweredPercentage"];
	        this.averageAccuracy = source["averageAccuracy"];
	        this.averageAnsweredAccuracy = source["averageAnsweredAccuracy"];
	        this.highConfidenceQuestions = source["highConfidenceQuestions"];
	        this.mediumConfidenceQuestions = source["mediumConfidenceQuestions"];
	        this.lowConfidenceQuestions = source["lowConfidenceQuestions"];
	        this.questionsWithReason = source["questionsWithReason"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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

