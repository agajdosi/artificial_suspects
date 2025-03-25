import { writable } from 'svelte/store';
import type { ServiceStatus, ErrorMessage, Service, Game } from '$lib/main';


// GAME STATE
const storedGame = localStorage.getItem('currentGame');
const defaultGame: Game = {
    uuid: '',
    level: 0,
    Score: 0,
    investigation: {
        uuid: '',
        game_uuid: '',
        suspects: [],
        rounds: [],
        CriminalUUID: '',
        InvestigationOver: false,
        Timestamp: ''
    },
    GameOver: false,
    Investigator: '',
    Timestamp: ''
};
export const currentGame = writable<Game>(storedGame ? JSON.parse(storedGame) : defaultGame);
currentGame.subscribe((value) => {
    localStorage.setItem('currentGame', JSON.stringify(value));
});

// ServiceStatus - is AI accessible?
const defaultServiceStatus: ServiceStatus = {
    ready: true,
    message: '',
    service: { // Dummy service - until the real fetching of service is implemented
        Name: 'Dummy',
        Type: 'local',
        TextModel: 'llava:latest',
        VisualModel: 'llava:latest',
        Token: '',
        URL: ''
    }
};
export const serviceStatus = writable<ServiceStatus>(defaultServiceStatus);

// ErrorMessage
const defaultErrorMessage: ErrorMessage = {
    Severity: '',
    Title: '',
    Message: '',
    Actions: []
};
export const errorMessage = writable<ErrorMessage>(defaultErrorMessage);

// Hint
export const hint = writable<string>("");

// ActiveService
const activeServiceName = localStorage.getItem('activeServiceName');
let storedActiveService: string = activeServiceName ? JSON.parse(activeServiceName) : '';
if (storedActiveService === '') {
    storedActiveService = "ollama";
}
export const activeService = writable<string>(storedActiveService);
activeService.subscribe((value) => {
    localStorage.setItem('activeServiceName', JSON.stringify(value));
});

// Services
const supportedServices: Record<string, Service> = {
    "ollama": {
        Name: "ollama",
        Type: "local", 
        TextModel: "llama3",
        VisualModel: "llama3",
        Token: "",
        URL: "",
    },
    "openai": {
        Name: "openai",
        Type: "API",
        TextModel: "gpt-4o",
        VisualModel: "gpt-4o",
        Token: "",
        URL: "",
    }
};

let storedServices: Record<string, Service>;
try {
    const stored = JSON.parse(localStorage.getItem('services') || '{}');
    // Convert array to object if needed
    storedServices = Array.isArray(stored) 
        ? stored.reduce((obj, service) => ({...obj, [service.Name]: service}), {})
        : stored;
} catch (e) {
    storedServices = {};
}

if (Object.keys(storedServices).length === 0) {
    storedServices = supportedServices;
}

export const services = writable<Record<string, Service>>(storedServices);
services.subscribe((value) => {
    // Ensure we're always storing an object
    localStorage.setItem('services', JSON.stringify(value));
});
