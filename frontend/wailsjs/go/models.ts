export namespace database {
	
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
	export class FinalScore {
	    Score: number;
	    Position: number;
	    Investigator: string;
	    GameUUID: string;
	    Timestamp: string;
	
	    static createFrom(source: any = {}) {
	        return new FinalScore(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Score = source["Score"];
	        this.Position = source["Position"];
	        this.Investigator = source["Investigator"];
	        this.GameUUID = source["GameUUID"];
	        this.Timestamp = source["Timestamp"];
	    }
	}
	export class Question {
	    UUID: string;
	    English: string;
	    Czech: string;
	    Polish: string;
	    Topic: string;
	    Level: number;
	
	    static createFrom(source: any = {}) {
	        return new Question(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.UUID = source["UUID"];
	        this.English = source["English"];
	        this.Czech = source["Czech"];
	        this.Polish = source["Polish"];
	        this.Topic = source["Topic"];
	        this.Level = source["Level"];
	    }
	}
	export class Round {
	    uuid: string;
	    InvestigationUUID: string;
	    Question: Question;
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
	        this.Question = this.convertValues(source["Question"], Question);
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
	    InvestigationOver: boolean;
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
	        this.InvestigationOver = source["InvestigationOver"];
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
	    Investigator: string;
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
	        this.Investigator = source["Investigator"];
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
	
	export class Model {
	    Name: string;
	    Service: string;
	    Visual: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Model(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Service = source["Service"];
	        this.Visual = source["Visual"];
	    }
	}
	
	
	export class Service {
	    Name: string;
	    Type: string;
	    Active: boolean;
	    TextModel: string;
	    VisualModel: string;
	    Token: string;
	    URL: string;
	
	    static createFrom(source: any = {}) {
	        return new Service(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Active = source["Active"];
	        this.TextModel = source["TextModel"];
	        this.VisualModel = source["VisualModel"];
	        this.Token = source["Token"];
	        this.URL = source["URL"];
	    }
	}

}

