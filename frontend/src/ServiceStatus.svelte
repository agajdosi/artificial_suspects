<script lang="ts">
    import { serviceStatus } from './lib/stores';
    import { AIServiceIsReady } from '../wailsjs/go/main/App';
    import { onMount } from 'svelte';

    async function getServiceStatus() {
        let status = await AIServiceIsReady();
        console.log("Service status is:", status);
        serviceStatus.set(status);
    }

    onMount(async () => {
        try {
            let status = await AIServiceIsReady();
            serviceStatus.set(status);
        } catch (error) {
            console.error('Error fetching status:', error);
        }
    });
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
