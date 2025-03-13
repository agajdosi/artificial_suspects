import { writable } from 'svelte/store';
import type { ServiceStatus, ErrorMessage, Service } from './main';

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

// ActiveService
let storedActiveService: string = JSON.parse(localStorage.getItem('activeServiceName') || 'dadaista');
if (storedActiveService === '') {
    storedActiveService = "dadaisto";
}
export const activeService = writable<string>(storedActiveService);
activeService.subscribe((value) => {
    localStorage.setItem('activeServiceName', JSON.stringify(value));
});

// Services
const supportedServices: Service[] = [
    {
        Name: "ollama",
        Type: "local",
        Active: true,
        TextModel: "llama3",
        VisualModel: "llama3",
        Token: "",
        URL: "",
    },
    {
        Name: "openai",
        Type: "API",
        Active: false,
        TextModel: "gpt-4o",
        VisualModel: "gpt-4o",
        Token: "",
        URL: "",
    },
];
let storedServices: Service[] = JSON.parse(localStorage.getItem('services') || '[]');
if (storedServices.length === 0) {
    storedServices = supportedServices;
}
export const services = writable<Service[]>(storedServices);
services.subscribe((value) => {
    localStorage.setItem('services', JSON.stringify(value));
});
