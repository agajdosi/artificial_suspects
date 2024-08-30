<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();

    function newGame() {
        dispatch('message', { message: 'goToGame' });
    }
</script>

<h1>Unusual Suspects</h1>
<button on:click={newGame}>NEW GAME</button>

<style>

</style>