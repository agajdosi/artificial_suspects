<script lang="ts">
    import { serviceStatus, activeService, services, hint } from './lib/stores';
    import type { ServiceStatus } from './lib/main';
    import { checkServiceStatusOllama } from './lib/intelligence';
    import { onMount } from 'svelte';

    let interval: NodeJS.Timeout; // Store the interval ID
    let overlayActive: boolean = false;
    function closePopup() {overlayActive = false;}
    function togglePopup() {overlayActive = !overlayActive;}


    async function getServiceStatus() {
        console.log("getServiceStatus()")
        let status: ServiceStatus;
        if ($activeService.toLowerCase() == "openai") {
            // TODO: get status from openai
        } else {
            status = await checkServiceStatusOllama($services[$activeService]);
        }

        serviceStatus.set(status);
    }

    onMount(async () => {
        getServiceStatus();
        interval = setInterval(getServiceStatus, 10_000);
    })
</script>

<button
    class="status"
    on:click={togglePopup}
    on:mouseenter={() => hint.set("Status of the AI service. Click to show the details.")}
    on:mouseleave={() => hint.set("")}
    >
    {#if $serviceStatus?.ready}
        🟢
    {:else}
        🔴
    {/if}
</button>

{#if overlayActive}
<div class="popup">
    <div>Service: {$serviceStatus.service?.Name}</div>
    <div>Model: {$serviceStatus.service?.VisualModel}</div>
    <div>Ready: {$serviceStatus.ready}</div>
    <div>{$serviceStatus.message}</div>
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
