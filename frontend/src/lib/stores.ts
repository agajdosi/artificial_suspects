import { database, main } from '../../wailsjs/go/models';
import { writable } from 'svelte/store';

// ServiceStatus - is AI accessible?
const defaultServiceStatus = new database.ServiceStatus();
//defaultServiceStatus.Ready = false;
//defaultServiceStatus.Message = '';
//defaultServiceStatus.Service = null;
export const serviceStatus = writable<database.ServiceStatus>(defaultServiceStatus);

// ErrorMessage
const defaultErrorMessage = new main.ErrorMessage();
//defaultErrorMessage.Severity = "error";
//defaultErrorMessage.Title = ""
//defaultErrorMessage.Message = "Something went very, very wrong.";
export const errorMessage = writable<main.ErrorMessage>(defaultErrorMessage);

// Hint
export const hint = writable<string>("");
