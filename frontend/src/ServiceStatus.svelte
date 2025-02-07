<script lang="ts">
    import { onMount } from 'svelte';
    import {AIServiceIsReady} from '../wailsjs/go/main/App';
    import { database } from '../wailsjs/go/models';

    let serviceStatus: database.ServiceStatus;

    async function getServiceStatus() {
        serviceStatus = await AIServiceIsReady();
        console.log("Service status is:", serviceStatus)
    }

    onMount(async () => {
        try {
            serviceStatus = await AIServiceIsReady();
        } catch (error) {
            console.error('Error fetching scores:', error);
        }
    });

</script>

<details class="service_status">
    {#if serviceStatus === undefined}
        <summary>AI not ready</summary>
    {:else}
        {#if serviceStatus.Ready === true}
            <summary>AI ready</summary>
        {:else}
            <summary>AI not ready</summary>
        {/if}
        <div>Service: {serviceStatus.Service.Name}</div>
        <div>Model: {serviceStatus.Service.VisualModel}</div>
        <div>Ready: {serviceStatus.Ready}</div>
        <div>Message: {serviceStatus.Message}</div>
    {/if}
    <button on:click={getServiceStatus}>Refresh</button>
</details>


<style>

.service_status {
    text-align: left;
}

</style>

