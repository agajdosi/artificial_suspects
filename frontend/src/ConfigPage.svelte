<script lang="ts">
    import { database } from '../wailsjs/go/models';
    import { GetServices, SaveToken, GetAllModels, SetActiveModel } from '../wailsjs/go/main/App.js';
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();

    let services: database.Service[] = [];
    let models: database.Model[] = [];
    let selectedModel: string = ''; // Holds the selected model's name
    let hasErrors: boolean = false; // Track if there are validation errors

    // Fetch services when the component is loaded
    GetServices().then((fetchedServices) => {
        services = fetchedServices;
    });

    // Fetch models when the component is loaded
    GetAllModels().then((fetchedModels) => {
        models = fetchedModels;
        // Set the first model as selected by default
        selectedModel = models.find(m => m.Active)?.Name || models[0].Name;
    });

    function goToMenu() {
        dispatch('message', { message: 'goToHome' });
    }

    // Validate inputs and check if any tokens are empty
    function validateInputs() {
        hasErrors = services.some(service => service.Token.trim() === '');
        return !hasErrors;
    }

    // Save the updated tokens for each service and set the selected model as active
    async function saveAndClose() {
        if (!validateInputs()) {
            return; // Stop if validation fails
        }

        for (const service of services) {
            await SaveToken(service.Name, service.Token);
        }

        // Set the selected model as active
        await SetActiveModel(selectedModel);

        goToMenu(); // Close after saving
    }
</script>

<h1>Game Configuration</h1>

<!-- Show services and tokens -->
{#if services.length === 0}
    Loading services...
{:else}
    {#each services as service, i}
        <div class="service">
            <label for="token-{i}">{service.Name}</label> token:
            <input 
                id="token-{i}" 
                type="password" 
                bind:value={service.Token} 
                placeholder="Enter token"
                class:error={service.Token.trim() === ''} 
            >
        </div>
    {/each}
{/if}

<!-- Dropdown for models -->
<div class="models-dropdown">
    <label for="models">Select Active Model:</label>
    {#if models.length === 0}
        Loading models...
    {:else}
        <select id="models" bind:value={selectedModel}>
            {#each models as model}
                <option value={model.Name} selected={model.Active}>{model.Name} ({model.Service})</option>
            {/each}
        </select>
    {/if}
</div>

<div class="actions">
    <button on:click={saveAndClose}>Save & Close</button>
    <button on:click={goToMenu}>Close</button>
</div>

<style>
    .actions {
        margin: 4rem;
    }

    input.error {
        border: 1px solid red;
    }

    .models-dropdown {
        margin: 2rem 0;
    }
</style>
