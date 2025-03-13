// MARK: CONSTANTS

const API_URL = 'http://localhost:8080';
const initGET = {
    method: 'GET',
    headers: {
        'Content-Type': 'application/json',
    },
}
const initPOST = {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
}

// MARK: TYPES

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

export interface FinalScore {
    GameUUID: string;
    Score: number;
    Investigator: string; // AKA player name
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

export interface Investigation {
    uuid: string;
    game_uuid: string;
    suspects: Suspect[];
    rounds: Round[];
    CriminalUUID: string;
    InvestigationOver: boolean;
    Timestamp: string;
}

export interface Model {
    Name: string;
    Visual: boolean;
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

export interface Service {
    Name: string;
    Type: string; // API or local
    Active: boolean;
    TextModel: string;
    VisualModel: string;
    Token: string;
    URL: string;
}

export interface ServiceStatus {
    Ready: boolean;
    Message: string;
    Service: Service;
}

export interface Suspect {
    UUID: string;
    Image: string;
    Free: boolean;
    Fled: boolean;
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


// MARK: FUNCTIONS

export async function NewGame(): Promise<Game> {
    const response = await fetch(`${API_URL}/new_game`, initGET);
    
    if (!response.ok) {
        throw new Error('Failed to create new game');
    }
    
    return await response.json();
}

export async function GetGame(): Promise<Game> {
    const response = await fetch(`${API_URL}/get_game`, initGET);
    
    if (!response.ok) {
        throw new Error('Failed to fetch game');
    }
    
    return await response.json();
}

export async function NextRound(): Promise<Game> {
    const response = await fetch(`${API_URL}/next_round`, initGET);

    if (!response.ok) {
        throw new Error('Failed to fetch next round');
    }

    return await response.json();
}

export async function NextInvestigation(): Promise<Game> {
    const response = await fetch(`${API_URL}/next_investigation`, initGET);

    if (!response.ok) {
        throw new Error('Failed to fetch next investigation');
    }

    return await response.json();
}

export async function EliminateSuspect(suspectUUID: string, roundUUID: string, investigationUUID: string): Promise<void> {
    const response = await fetch(`${API_URL}/eliminate_suspect`, initPOST);

    if (!response.ok) {
        throw new Error('Failed to eliminate suspect');
    }
}

export async function WaitForAnswer(roundUUID: string): Promise<string> {
    const response = await fetch(`${API_URL}/wait_for_answer`, initGET);

    if (!response.ok) {
        throw new Error('Failed to wait for answer');
    }

    return await response.json();
}

export async function GetScores(): Promise<FinalScore[]> {
    const response = await fetch(`${API_URL}/get_scores`, initGET);

    if (!response.ok) {
        throw new Error('Failed to fetch scores');
    }

    return await response.json();
}

export async function SaveScore(name: string, gameUUID: string): Promise<void> {
    const body = {
        investigator: name,
        game_uuid: gameUUID,
    }
    const response = await fetch(`${API_URL}/save_score`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(body),
    });

    if (!response.ok) {
        throw new Error('Failed to save score');
    }
}



// MARK: AI SERVICES - dummy for now
// 

export async function GetServices(): Promise<Service[]> {
    let services: Service[] = [];
    return services;
}

export async function GetActiveService(): Promise<Service> {
    let service: Service = {
        Name: 'OLLAMA',
        Type: 'local',
        TextModel: 'llama3',
        VisualModel: 'llama3',
        Token: '',
        URL: 'localhost:12345',
        Active: true,
    };
    return service;
}

export async function SaveService(service: Service): Promise<void> {
    return;
}

export async function ActivateService(service: Service): Promise<void> {
    return;
}

export async function GetDefaultModels(): Promise<Model[]> {
    let models: Model[] = [];
    return models;
}

export async function ListModelsOllama(): Promise<Model[]> {
    let models: Model[] = [];
    return models;
}
