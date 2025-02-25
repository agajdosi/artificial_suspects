<script lang="ts">
    import { serviceStatus, errorMessage } from './lib/stores';
    import { AIServiceIsReady } from '../wailsjs/go/main/App';
    import { main } from './../wailsjs/go/models';
    import { onMount } from 'svelte';

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

    onMount(async () => {getServiceStatus();})
</script>

<details class="service_status">
    {#if !$serviceStatus}
        <!-- No status at all -->
        <summary>AI not ready</summary>
    {:else if $serviceStatus.Ready}
        <summary>AI ready</summary>
    {:else}
        <summary>AI not ready</summary>
    {/if}
    <div>Service: {$serviceStatus.Service?.Name}</div>
    <div>Model: {$serviceStatus.Service?.VisualModel}</div>
    <div>Ready: {$serviceStatus.Ready}</div>
    <div>Message: {$serviceStatus.Message}</div>

    <button on:click={getServiceStatus}>Refresh</button>
</details>

<style>
.service_status {
    text-align: left;
}
</style>
