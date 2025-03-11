import { writable } from 'svelte/store';
import type { ServiceStatus, ErrorMessage } from './main';

// ServiceStatus - is AI accessible?
const defaultServiceStatus: ServiceStatus = {
    Ready: false,
    Message: '',
    Service: {
        Name: '',
        Type: '',
        Active: false,
        TextModel: '',
        VisualModel: '',
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
