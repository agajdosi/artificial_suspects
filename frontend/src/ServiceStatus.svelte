<script lang="ts">
    import { serviceStatus, errorMessage } from './lib/stores';
    import { AIServiceIsReady } from '../wailsjs/go/main/App';
    import { main } from './../wailsjs/go/models';
    import { onMount } from 'svelte';

    let interval: number; // Store the interval ID
    let overlayActive: boolean = false;
    function closeOverlay() {overlayActive = false;}
    function openOverlay() {overlayActive = true;}

    async function getServiceStatus() {
        let status = await AIServiceIsReady();
        if (!status.Ready) {
            const em = new main.ErrorMessage();
            em.Severity = "warning";
            em.Title = `${status.Service.VisualModel} from provider ${status.Service.Name} is not accessible`;
            em.Message = `Bla bla bla`;
            em.Actions = ["goToConfig"]
            errorMessage.set(em);
        }
        console.log("Service status is:", status);
        serviceStatus.set(status);
    }

    onMount(async () => {
        getServiceStatus();
        interval = setInterval(getServiceStatus, 10_000);
    })
</script>

<div class="status">
    {#if $serviceStatus?.Ready}
        🟢
    {:else}
        🔴
    {/if}
</div>

{#if overlayActive}
<div class="overlay">
    <div>Service: {$serviceStatus.Service?.Name}</div>
    <div>Model: {$serviceStatus.Service?.VisualModel}</div>
    <div>Ready: {$serviceStatus.Ready}</div>
    <div>Message: {$serviceStatus.Message}</div>
    <button on:click={getServiceStatus}>Refresh</button>
    <button on:click={closeOverlay}>Close</button>
</div>
{/if}

<style>
.status {
    position: fixed;
    bottom: 0;
    left: 100;
    padding: 0 0 0 4px;
}

.overlay{
    position: fixed;
    top: 0;
    height: 100%;
    width: 100vw;
    background-color: rgba(0, 0, 0, 0.7);
    backdrop-filter: blur(5px);
}
</style>
