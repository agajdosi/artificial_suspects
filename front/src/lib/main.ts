import { currentGame, currentPlayer } from '$lib/stores';
import { get } from 'svelte/store';
import { generateAnswer } from '$lib/intelligence';

// MARK: CONSTANTS

const API_URL = import.meta.env.PROD ? 'https://artsus.lab.gajdosik.org' : 'http://localhost:8080';
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

export interface Answer {
    uuid: string;
    answer: string;
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

export interface FinalScore {
    GameUUID: string;
    Score: number;
    Investigator: string; // AKA player name
}

export interface Player {
    UUID: string;
    Name: string;
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
    TextModel: string;
    VisualModel: string;
    Token: string;
    URL: string;
}

export interface ServiceStatus {
    ready: boolean;
    message: string;
    service: Service;
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
    console.log("NEW GAME HANDLER");
    let newGame: Game;
    try {
        const player = get(currentPlayer);
        const response = await fetch(`${API_URL}/new_game?player_uuid=${player.UUID}`, initGET);
        if (!response.ok) {
            throw new Error('Failed to create new game');
        }
        newGame = await response.json();
        console.log(`NewGame() response: ${newGame}`);
        currentGame.set(newGame);
    } catch (error) {
        console.log(`NewGame() has failed: ${error}`);
        throw error;
    }

    const answer = await generateAnswer(
        newGame.investigation.rounds.at(-1)?.uuid,
        newGame.investigation.rounds.at(-1)?.Question,
        newGame.investigation.CriminalUUID
    );

    await saveAnswer(answer.answer, newGame.investigation.rounds.at(-1)?.uuid);
    if (newGame.investigation.rounds.at(-1)) {
        newGame.investigation.rounds.at(-1).answer = answer.answer;
        newGame.investigation.rounds.at(-1).AnswerUUID = answer.uuid;
    }
    currentGame.set(newGame);
    return newGame;
}

export async function GetGame(): Promise<Game> {
    const player = get(currentPlayer);
    const response = await fetch(`${API_URL}/get_game?player_uuid=${player.UUID}`, initGET);
    if (!response.ok) {
        throw new Error('Failed to fetch game');
    }
    
    return await response.json();
}

export async function NextRound() {
    // FIRST GET THE NEW ROUND`
    const response = await fetch(`${API_URL}/next_round`, initGET);
    if (!response.ok) {
        throw new Error('Failed to fetch next round');
    }

    let game: Game = await response.json();
    console.log(`>>> NEW ROUND: ${game.investigation.rounds.at(-1)}`);
    currentGame.set(game);

    // THEN GENERATE ANSWER
    const answer = await generateAnswer(
        game.investigation.rounds.at(-1).uuid,
        game.investigation.rounds.at(-1).Question,
        game.investigation.CriminalUUID
    );

    await saveAnswer(answer.answer, game.investigation.rounds.at(-1).uuid);
    game.investigation.rounds.at(-1).answer = answer.answer;
    game.investigation.rounds.at(-1).AnswerUUID = answer.uuid;
    currentGame.set(game);
}

export async function NextInvestigation(): Promise<Game> {
    const response = await fetch(`${API_URL}/next_investigation`, initGET);

    if (!response.ok) {
        throw new Error('Failed to fetch next investigation');
    }

    return await response.json();
}

export async function EliminateSuspect(suspectUUID: string, roundUUID: string, investigationUUID: string): Promise<void> {
    const response = await fetch(`${API_URL}/eliminate_suspect?suspect_uuid=${suspectUUID}&round_uuid=${roundUUID}&investigation_uuid=${investigationUUID}`, initPOST);
    if (!response.ok) {
        throw new Error('Failed to eliminate suspect');
    }
}

export async function WaitForAnswer(roundUUID: string): Promise<string> {
    const response = await fetch(`${API_URL}/wait_for_answer?round_uuid=${roundUUID}`, initGET);
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

export async function SaveScore(playerName: string, gameUUID: string) {
    const response = await fetch(`${API_URL}/save_score?player_name=${playerName}&game_uuid=${gameUUID}`, initPOST);
    if (!response.ok) {
        throw new Error('Failed to save score');
    }
}


export async function ListModelsOllama(): Promise<Model[]> {
    let models: Model[] = [];
    return models;
}


export async function getDescriptionsForSuspect(suspectUUID: string, serviceName: string, modelName: string): Promise<string[]> {
    let descriptions: string[] = [];
    return descriptions;
}

export async function getQuestion(questionUUID: string): Promise<Question> {
    let question: Question;
    return question;
}


export async function saveAnswer(answer: string, roundUUID: string): Promise<void> {
    const response = await fetch(`${API_URL}/save_answer?answer=${answer}&round_uuid=${roundUUID}`, initPOST);
    if (!response.ok) {
        throw new Error('Failed to save answer');
    }
}
