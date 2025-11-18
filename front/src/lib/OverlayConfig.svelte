<script lang="ts">
    import { activeService, services } from '$lib/stores';
    let { overlayConfigVisible = $bindable(false) } = $props();

    let selectedService = $state($activeService); // Holds the AI service selected in the UI

    async function startGame(event: Event) {
        const target = event.target as HTMLButtonElement;
        const name = target.value;
        selectedService = name;
        console.log("Starting game with service:", $services[name]);
        activeService.set(name);
        overlayConfigVisible = false;
    }
</script>

{#if overlayConfigVisible}
<div class="overlay">

<h1>Choose AI coponent</h1>

<div class="services">
    {#if Object.keys($services).length === 0}
        Loading services...
    {:else}
        {#each Object.entries($services) as [name, service]}
                <button onclick={startGame} value={service.Name}>{service.TextModel}</button>
        {/each}
    {/if}
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
    .services {
        display: flex;
        flex-direction: column;
        align-items: center;
        margin-top: 5rem;
        gap: 10px;
    }
</style>
