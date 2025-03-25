<script lang="ts">
    import { errorMessage } from '$lib/stores';
    import type { ErrorMessage } from '$lib/main';
    import { createEventDispatcher } from 'svelte';

    const dispatch = createEventDispatcher();

    function clearError() {
        const em: ErrorMessage = {
            Title: "",
            Message: "",
            Actions: [],
            Severity: ""
        };
        errorMessage.set(em);
    }

    function enterConfigDispatcher() {
        dispatch('message', { message: 'goToConfig' });
        clearError();
    }

    function closeOverlay() {
        clearError();
    }

</script>

{#if $errorMessage.Title}
<div class="overlay">
    <div class="content">
        <h1>{$errorMessage.Title}</h1>
        <p>{$errorMessage.Message}</p>
        <div class="actions">
            {#if $errorMessage.Actions?.includes("goToConfig")}
            <button on:click={enterConfigDispatcher}>Configuration</button>
            {/if}
            {#if $errorMessage.Actions?.includes("close")}
            <button on:click={closeOverlay}>Close</button>
            {/if}
        </div>
    </div>
</div>
{/if}


<style>
.overlay{
    position: fixed;
    top: 0;
    left: 0;
    height: 100%;
    width: 100vw;
    background-color: rgba(0, 0, 0, 0.7);
    backdrop-filter: blur(5px);
}

.content {
    margin: 15vh auto 0 auto;
}
</style>