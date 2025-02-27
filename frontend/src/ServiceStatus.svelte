<script lang="ts">
    import { serviceStatus } from './lib/stores';
    import { AIServiceIsReady } from '../wailsjs/go/main/App';
    import { onMount } from 'svelte';

    let interval: number; // Store the interval ID
    let overlayActive: boolean = false;
    function closePopup() {overlayActive = false;}
    function togglePopup() {overlayActive = !overlayActive;}


    async function getServiceStatus() {
        let status = await AIServiceIsReady();
        serviceStatus.set(status);
    }

    onMount(async () => {
        getServiceStatus();
        interval = setInterval(getServiceStatus, 10_000);
    })
</script>

<button on:click={togglePopup} class="status">
    {#if $serviceStatus?.Ready}
        🟢
    {:else}
        🔴
    {/if}
</button>

{#if overlayActive}
<div class="popup">
    <div>Service: {$serviceStatus.Service?.Name}</div>
    <div>Model: {$serviceStatus.Service?.VisualModel}</div>
    <div>Ready: {$serviceStatus.Ready}</div>
    <div>{$serviceStatus.Message}</div>
    <button on:click={getServiceStatus}>Refresh</button>
    <button on:click={closePopup}>Close</button>
</div>
{/if}

<style>
.status {
    border: none;
}
.status:hover {
    background-color: unset;
}

.status {
    position: fixed;
    bottom: 0;
    left: 0;
    padding: 0 0 0 4px;
}

.popup{
    position: fixed;
    bottom: 2rem;
    left: 0;
    padding: 0.4rem 0.8rem 0.6rem;
    text-align: left;
    background-color: rgba(0, 0, 0, 0.7);
    backdrop-filter: blur(5px);
}
</style>
