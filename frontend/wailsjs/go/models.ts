export namespace main {
	
	export class Round {
	    uuid: string;
	    question: string;
	    answer: string;
	
	    static createFrom(source: any = {}) {
	        return new Round(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uuid = source["uuid"];
	        this.question = source["question"];
	        this.answer = source["answer"];
	    }
	}
	export class Suspect {
	    UUID: string;
	    Image: string;
	    Free: boolean;
	    Timestamp: string;
	
	    static createFrom(source: any = {}) {
	        return new Suspect(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.UUID = source["UUID"];
	        this.Image = source["Image"];
	        this.Free = source["Free"];
	        this.Timestamp = source["Timestamp"];
	    }
	}
	export class Investigation {
	    uuid: string;
	    game_uuid: string;
	    suspects: Suspect[];
	    level: number;
	    rounds: Round[];
	
	    static createFrom(source: any = {}) {
	        return new Investigation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uuid = source["uuid"];
	        this.game_uuid = source["game_uuid"];
	        this.suspects = this.convertValues(source["suspects"], Suspect);
	        this.level = source["level"];
	        this.rounds = this.convertValues(source["rounds"], Round);
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
	export class Game {
	    uuid: string;
	    level: number;
	    investigation: Investigation;
	
	    static createFrom(source: any = {}) {
	        return new Game(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uuid = source["uuid"];
	        this.level = source["level"];
	        this.investigation = this.convertValues(source["investigation"], Investigation);
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

