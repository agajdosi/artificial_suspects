<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();
    
    function newGameDispatcher() {dispatch('newGame', {message: 'new game'});}
    function continueGameDispatcher() {dispatch('message', { message: 'continueGame' });}
    function enterConfigDispatcher() {dispatch('enterConfig', { message: 'enterConfig' });}
</script>

<h1>Unusual Suspects</h1>

<div class="menu">
    <button on:click={newGameDispatcher}>NEW GAME</button>
    <button on:click={continueGameDispatcher}>CONTINUE GAME</button>
    <button on:click={enterConfigDispatcher}>CONFIGURATION</button>
</div>

<style>

.menu {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.menu button {
    width: 400px;
    margin: 2rem;
}

</style>