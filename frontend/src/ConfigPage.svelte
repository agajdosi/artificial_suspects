<script lang="ts">
    import { ListModelsOllama } from './lib/main';
    import { activeService, services } from './lib/stores';
    import { createEventDispatcher } from 'svelte';
    import LanguageSwitch from './LanguageSwitch.svelte';
    const dispatch = createEventDispatcher();

    let selectedService: string = $activeService; // Holds the AI service selected in the UI

    async function goToMenu() {dispatch('message', { message: 'goToHome' });}

    async function showServiceDetail(event) {selectedService = event.target.value;}

    async function saveService() {
        const service = $services[selectedService];
        if (!service) {
            console.error("No service selected to save.");
            return;
        }
        const URL = service.Type === "local" ? (document.getElementById('url-' + selectedService) as HTMLInputElement).value : "";
        const Token = service.Type === "API" ? (document.getElementById('token-' + selectedService) as HTMLInputElement).value : "";
        const VisualModel = (document.getElementById('model-' + selectedService) as HTMLInputElement).value;
    
        console.log(`Saving service ${selectedService} with URL ${URL}, Token ${Token}, VisualModel ${VisualModel}`);
        $services = {
            ...$services,
            [selectedService]: {
                ...service,
                URL,
                Token,
                VisualModel
            }
        };
    }

    async function activateService() {
        const service = $services[selectedService];
        if (!service) {
            console.error("No service selected to save.");
            return;
        }
        console.log(`Activating service ${service.Name}`);
        activeService.set(service.Name);
    }

    async function listModelsOllama() {
        const res = await ListModelsOllama();
    }
</script>

<h1>Game Configuration</h1>
<div class="info">
    {#if Object.keys($services).length > 0}
        {#each Object.entries($services) as [name, service]}
            {#if service.Name === $activeService}
                <p>
                    Game uses the active service <strong>{name}</strong> with visual LLM <strong>{service.VisualModel}</strong>.
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
    {#if Object.keys($services).length === 0}
        Loading services...
    {:else}
        <div class="services">
            {#each Object.entries($services) as [name, service]}
                <button on:click={showServiceDetail} value={name} class:selected={selectedService === name}>
                    {name}
                </button>
            {/each}
        </div>
    {/if}

    <!-- Configuration for the selected AI service -->
    {#if selectedService}
        {#each Object.entries($services) as [name, service]}
            {#if name === selectedService}
                <div class="service-details">
                    {#if service.Type === "API"}
                        <div class="service-token">
                            <label for="token-{name}">API token:</label>
                            <input id="token-{name}" type="password" bind:value={service.Token} placeholder="Enter token" class:error={service.Token.trim() === ''}>
                        </div>
                    {/if}
                    {#if service.Type === "local"}
                        <button on:click={listModelsOllama}>List models</button>
                        <div class="service-URL">
                            <label for="url-{name}">URL:</label>
                            <input id="url-{name}" bind:value={service.URL} placeholder="Enter local instance URL" class:error={service.URL.trim() === ''}>
                        </div>
                    {/if}
                    <div class="service-model">
                        <div class="service-model">
                            <label for="model-{name}">Visual LLM:</label>
                            <input id="model-{name}" bind:value={service.VisualModel} placeholder="Name of model to use" class:error={service.VisualModel.trim() === ''}>
                        </div>
                    </div>
                    <div class="actions">
                        <button on:click={saveService}>Save</button>
                        {#if service.Name === $activeService}
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
