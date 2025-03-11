// MARK: CONSTANTS

export const API_URL = 'http://localhost:8080';

// MARK: TYPES

export interface ServiceStatus {
    Ready: boolean;
    Message: string;
    Service: Service;
}

export interface Service {
    Name: string;
    Type: string; // API or local
    Active: boolean;
    TextModel: string;
    VisualModel: string;
    Token: string;
    URL: string;
}

export interface Investigation {
    uuid: string;
    game_uuid: string;
    suspects: Suspect[];
    rounds: Round[];
    CriminalUUID: string;
    InvestigationOver: boolean;
    Timestamp: string;
}
 
export interface Game {
    uuid: string;
    investigation: Investigation;
    level: number;
    Score: number;
    GameOver: boolean;
    Investigator: string;
    Timestamp: string;
}

export interface Suspect {
    UUID: string;
    Image: string;
    Free: boolean;
    Fled: boolean;
    Timestamp: string;
}

export interface Round {
    uuid: string;
    InvestigationUUID: string;
    Question: Question;
    AnswerUUID: string;
    answer: string;
    Eliminations: Elimination[];
    Timestamp: string;
}

export interface Question {
    UUID: string;
    English: string;
    Czech: string;
    Polish: string;
    Topic: string;
    Level: number;
}

export interface Elimination {
    UUID: string;
    RoundUUID: string;
    SuspectUUID: string;
    Timestamp: string;
}

export interface ErrorMessage {
    Severity: string;
    Title: string;
    Message: string;
    Actions: string[];
}

// MARK: FUNCTIONS

export async function NewGame(): Promise<Game> {
    const response = await fetch(`${API_URL}/new_game`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    });
    
    if (!response.ok) {
        throw new Error('Failed to create new game');
    }
    
    return await response.json();
}

export async function GetGame(): Promise<Game> {
    const response = await fetch(`${API_URL}/get_game`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    });
    
    if (!response.ok) {
        throw new Error('Failed to fetch game');
    }
    
    return await response.json();
}


