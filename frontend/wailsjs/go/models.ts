export namespace main {
	
	export class Suspect {
	    uuid: string;
	    imageSource: string;
	    free: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Suspect(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uuid = source["uuid"];
	        this.imageSource = source["imageSource"];
	        this.free = source["free"];
	    }
	}
	export class Game {
	    suspects: Suspect[];
	    level: number;
	    question: string;
	    gameUUID: string;
	
	    static createFrom(source: any = {}) {
	        return new Game(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.suspects = this.convertValues(source["suspects"], Suspect);
	        this.level = source["level"];
	        this.question = source["question"];
	        this.gameUUID = source["gameUUID"];
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

