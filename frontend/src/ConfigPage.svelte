<script lang="ts">
    import { database } from '../wailsjs/go/models';
    import { GetServices, SaveService, ActivateService, GetAllModels } from '../wailsjs/go/main/App.js';
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();

    let services: database.Service[] = [];
    let models: database.Model[] = [];
    let selectedService: string = ''; // Holds the selected AI service
    let selectedModel: database.Model; // Holds the selected model's name
    let hasErrors: boolean = false; // Track if there are validation errors

    // Fetch services when the component is loaded
    GetServices().then((fetchedServices) => {
        console.log("fetched services:", fetchedServices)
        services = fetchedServices;
        selectedService = services.find(s => s.Active).Name || services[0].Name;
        console.log("selected service:", selectedService)
    });

    // Fetch models when the component is loaded
    GetAllModels().then((fetchedModels) => {
        models = fetchedModels;
        selectedModel = models.find(m => m.Active) || models[0];
    });

    function goToMenu() {
        dispatch('message', { message: 'goToHome' });
    }

    async function showServiceDetail(event) {selectedService = event.target.value;}

    async function saveService() {
        const service = services.find(s => s.Name === selectedService);
        if (!service) {
            console.error("No service selected to save.");
            return;
        }

        try {
            await SaveService(service.Name, service.Model, service.Token, service.URL);
            console.log(`Service ${service.Name} saved successfully.`);
        } catch (error) {
            console.error(`Error saving service ${service.Name}:`, error);
        }
    }

    async function activateService() {
        const service = services.find(s => s.Name === selectedService);
        if (!service) {
            console.error("No service selected to save.");
            return;
        }

        await saveService();
        try {
            await ActivateService(service.Name);
            console.log(`Service ${service.Name} successfully activated.`);
            services = services.map(s => ({
                ...s,
                Active: s.Name === service.Name, // Only the selected service is active
            }));
        } catch (error) {
            console.error(`Error saving service ${service.Name}:`, error);
        }
    }


</script>

<h1>Game Configuration</h1>
<div class="info">
    {#if services.length > 0}
        {#each services as service}
            {#if service.Active}
                <p>
                    Game uses the active service <strong>{service.Name}</strong> with model <strong>{service.Model}</strong>.
                </p>
            {/if}
        {/each}
    {:else}
        <p>No services available.</p>
    {/if}
</div>

<!-- AI Services Tabs -->
<div class="service-config">
    <h2>AI Services</h2>
    {#if services.length === 0}
        Loading services...
    {:else}
        <div class="services">
            {#each services as service}
                <button on:click={showServiceDetail} value={service.Name} class:selected={selectedService === service.Name}>
                    {service.Name}
                </button>
            {/each}
        </div>
    {/if}

    <!-- Configuration for the selected AI service -->
    {#if selectedService}
        {#each services as service}
            {#if service.Name == selectedService}
                <div class="service-details">
                    {#if service.Type == "API"}
                        <div class="service-token">
                            <label for="token-{service.Name}">API token:</label>
                            <input id="token-{service.Name}" type="password" bind:value={service.Token} placeholder="Enter token" class:error={service.Token.trim() === ''}>
                        </div>
                    {/if}
                    {#if service.Type == "local"}
                        <div class="service-URL">
                            <label for="token-{service.Name}">URL:</label>
                            <input id="token-{service.Name}" bind:value={service.URL} placeholder="Enter local instance URL" class:error={service.URL.trim() === ''}>
                        </div>
                    {/if}
                    <div class="service-model">
                        <div class="service-model">
                            <label for="token-{service.Model}">Model:</label>
                            <input id="token-{service.Model}" bind:value={service.Model} placeholder="Name of model to use" class:error={service.Model.trim() === ''}>
                        </div>
                    </div>
                    <div class="actions">
                        <button on:click={saveService}>Save</button>
                        <button on:click={activateService}>Activate</button>
                    </div>
                </div>
            {/if}
        {/each}
    {/if}
</div>

<div class="back">
    <button on:click={goToMenu}>Back</button>
</div>

<style>
    .actions {
        margin: 1rem 0;
    }

    input.error {
        border: 1px solid red;
    }

    .service-details {
        margin: 1rem 0;
    }

    .back {
        margin: 8rem 0 0 0;
    }

</style>
