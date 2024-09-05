export namespace main {
	
	export class Elimination {
	    UUID: string;
	    RoundUUID: string;
	    SuspectUUID: string;
	    Timestamp: string;
	
	    static createFrom(source: any = {}) {
	        return new Elimination(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.UUID = source["UUID"];
	        this.RoundUUID = source["RoundUUID"];
	        this.SuspectUUID = source["SuspectUUID"];
	        this.Timestamp = source["Timestamp"];
	    }
	}
	export class Round {
	    uuid: string;
	    InvestigationUUID: string;
	    QuestionUUID: string;
	    question: string;
	    AnswerUUID: string;
	    answer: string;
	    Eliminations: Elimination[];
	    Timestamp: string;
	
	    static createFrom(source: any = {}) {
	        return new Round(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uuid = source["uuid"];
	        this.InvestigationUUID = source["InvestigationUUID"];
	        this.QuestionUUID = source["QuestionUUID"];
	        this.question = source["question"];
	        this.AnswerUUID = source["AnswerUUID"];
	        this.answer = source["answer"];
	        this.Eliminations = this.convertValues(source["Eliminations"], Elimination);
	        this.Timestamp = source["Timestamp"];
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
	export class Suspect {
	    UUID: string;
	    Image: string;
	    Free: boolean;
	    Fled: boolean;
	    Timestamp: string;
	
	    static createFrom(source: any = {}) {
	        return new Suspect(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.UUID = source["UUID"];
	        this.Image = source["Image"];
	        this.Free = source["Free"];
	        this.Fled = source["Fled"];
	        this.Timestamp = source["Timestamp"];
	    }
	}
	export class Investigation {
	    uuid: string;
	    game_uuid: string;
	    suspects: Suspect[];
	    rounds: Round[];
	    CriminalUUID: string;
	    Timestamp: string;
	
	    static createFrom(source: any = {}) {
	        return new Investigation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uuid = source["uuid"];
	        this.game_uuid = source["game_uuid"];
	        this.suspects = this.convertValues(source["suspects"], Suspect);
	        this.rounds = this.convertValues(source["rounds"], Round);
	        this.CriminalUUID = source["CriminalUUID"];
	        this.Timestamp = source["Timestamp"];
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
	    investigation: Investigation;
	    level: number;
	    Score: number;
	    GameOver: boolean;
	    Timestamp: string;
	
	    static createFrom(source: any = {}) {
	        return new Game(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uuid = source["uuid"];
	        this.investigation = this.convertValues(source["investigation"], Investigation);
	        this.level = source["level"];
	        this.Score = source["Score"];
	        this.GameOver = source["GameOver"];
	        this.Timestamp = source["Timestamp"];
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

