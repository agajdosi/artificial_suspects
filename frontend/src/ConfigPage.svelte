<script lang="ts">
    import { GetServices, SaveService, ActivateService, GetDefaultModels, ListModelsOllama } from './lib/main';
    import type { Model, Service } from './lib/main';
    import { createEventDispatcher } from 'svelte';
    import LanguageSwitch from './LanguageSwitch.svelte';
    const dispatch = createEventDispatcher();

    let services: Service[] = [];
    let models: Model[] = [];
    let selectedService: string = ''; // Holds the selected AI service ()
    let activeService: {Name: string, TextModel: string, VisualModel: string} = {Name: "", TextModel: "", VisualModel: ""};

    // Fetch services when the component is loaded
    GetServices().then((fetchedServices) => {
        services = fetchedServices;
        const active = services.find(s => s.Active);
        if (active) {
            selectedService = active.Name;
            activeService = { Name: active.Name, TextModel: active.TextModel, VisualModel: active.VisualModel };
        } else {
            selectedService = services.length > 0 ? services[0].Name : '';
        }
    });

    // Fetch models when the component is loaded
    GetDefaultModels().then((fetchedModels) => {
        models = fetchedModels;
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
            await SaveService(service.Name, service.TextModel, service.VisualModel, service.Token, service.URL);
            console.log(`Service ${service.Name} saved successfully.`);
            activeService = { Name: service.Name, TextModel: service.TextModel, VisualModel: service.VisualModel }; // Update after saving
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
            activeService = { Name: service.Name, TextModel: service.TextModel, VisualModel: service.VisualModel }; // Update displayed info
        } catch (error) {
            console.error(`Error saving service ${service.Name}:`, error);
        }
    }

    async function listModelsOllama() {
        const res = await ListModelsOllama();
    }
</script>

<h1>Game Configuration</h1>
<div class="info">
    {#if services.length > 0}
        {#each services as service}
            {#if service.Active}
                <p>
                    Game uses the active service <strong>{activeService.Name}</strong> with visual LLM <strong>{activeService.VisualModel}</strong>.
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
                        <button on:click={listModelsOllama}>List models</button>
                        <div class="service-URL">
                            <label for="token-{service.Name}">URL:</label>
                            <input id="token-{service.Name}" bind:value={service.URL} placeholder="Enter local instance URL" class:error={service.URL.trim() === ''}>
                        </div>
                    {/if}
                    <div class="service-model">
                        <div class="service-model">
                            <label for="token-{service.VisualModel}">Visual LLM:</label>
                            <input id="token-{service.VisualModel}" bind:value={service.VisualModel} placeholder="Name of model to use" class:error={service.VisualModel.trim() === ''}>
                        </div>
                    </div>
                    <div class="actions">
                        <button on:click={saveService}>Save</button>
                        {#if service.Active}
                            <button on:click={activateService}>Active</button>
                        {:else}
                            <button on:click={activateService}>Activate</button>
                        {/if}
                    </div>
                </div>
            {/if}
        {/each}
    {/if}
</div>

<h2>Language</h2>
<LanguageSwitch/>

<div class="back">
    <button on:click={goToMenu}>Back</button>
</div>

<style>
    h2 {
        margin: 2rem 0 0.2rem 0;
    }
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
