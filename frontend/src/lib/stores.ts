import { database } from '../../wailsjs/go/models';
import { writable } from 'svelte/store';

const defaultServiceStatus = new database.ServiceStatus();
defaultServiceStatus.Ready = false;
defaultServiceStatus.Message = '';
defaultServiceStatus.Service = null;

export const serviceStatus = writable<database.ServiceStatus>(defaultServiceStatus);
