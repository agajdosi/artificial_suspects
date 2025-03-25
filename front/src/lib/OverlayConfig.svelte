<script lang="ts">
    import { activeService, services } from '$lib/stores';

    let { configOverlayVisible = false } = $props();

    let selectedService = $state($activeService); // Holds the AI service selected in the UI

    async function showServiceDetail(event: Event) {
        const target = event.target as HTMLButtonElement;
        selectedService = target.value;
    }

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

        configOverlayVisible = false;
    }

    function quitConfig() {
        configOverlayVisible = false;
    }
</script>

{#if configOverlayVisible}
<div class="overlay">

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
    <h2>1. Choose AI Service</h2>
    {#if Object.keys($services).length === 0}
        Loading services...
    {:else}
        <div class="services">
            {#each Object.entries($services) as [name, service]}
                <button onclick={showServiceDetail} value={name} class:selected={selectedService === name}>
                    {name}
                </button>
            {/each}
        </div>
    {/if}
</div>

<div>
    <h2>2. Configure AI Service</h2>

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
                        

                    </div>
                </div>
            {/if}
        {/each}
    {/if}
</div>

<div class="back">
    <button onclick={saveService}>Save</button>
    <button onclick={quitConfig}>Quit</button>
</div>

</div>
{/if}

<style>
    .overlay {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-color: var(--bg-color);
    }
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
