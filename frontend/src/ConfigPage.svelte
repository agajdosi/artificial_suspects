<script lang="ts">
    import { database } from '../wailsjs/go/models';
    import { GetServices, SaveToken } from '../wailsjs/go/main/App.js';
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();

    let services: database.Service[] = [];
    let hasErrors: boolean = false; // Track if there are validation errors

    // Fetch services when the component is loaded
    GetServices().then((fetchedServices) => {
        services = fetchedServices;
    });

    function goToMenu() {
        dispatch('message', { message: 'goToHome' });
    }

    // Validate inputs and check if any tokens are empty
    function validateInputs() {
        hasErrors = services.some(service => service.Token.trim() === '');
        return !hasErrors;
    }

    // Save the updated tokens for each service
    async function saveAndClose() {
        if (!validateInputs()) {
            return; // Stop if validation fails
        }

        for (const service of services) {
            await SaveToken(service.Name, service.Token);
        }
        goToMenu(); // Close after saving
    }
</script>

<h1>Game Configuration</h1>

{#if services.length === 0}
    Loading services...
{:else}
    {#each services as service, i}
        <div class="service">
            <label for="token-{i}">{service.Name}</label> token:
            <!-- Apply the 'error' class if the token is empty -->
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

<div class="actions">
    <button on:click={saveAndClose}>Save & Close</button>
    <button on:click={goToMenu}>Close</button>
</div>

<style>
    .actions {
        margin: 4rem;
    }

    /* Add red border to input when there's an error */
    input.error {
        border: 1px solid red;
    }
</style>
