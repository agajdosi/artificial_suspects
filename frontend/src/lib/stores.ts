import { writable } from 'svelte/store';
import type { ServiceStatus, ErrorMessage } from './main';

// ServiceStatus - is AI accessible?
const defaultServiceStatus: ServiceStatus = {
    Ready: true,
    Message: '',
    Service: { // Dummy service - until the real fetching of service is implemented
        Name: 'Dummy',
        Type: 'local',
        Active: true,
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
