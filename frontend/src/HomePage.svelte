<script lang="ts">
    import { serviceStatus } from './lib/stores';
    import { createEventDispatcher } from 'svelte';
    import { QuitApplication} from '../wailsjs/go/main/App';
    const dispatch = createEventDispatcher();
    
    function newGameDispatcher() {dispatch('newGame', {message: 'new game'});}
    function continueGameDispatcher() {dispatch('message', { message: 'continueGame' });}
    function enterConfigDispatcher() {dispatch('message', { message: 'goToConfig' });}
</script>

<h1>Artificial Suspects</h1>

<div class="menu">
    <button disabled={!$serviceStatus.Ready} on:click={newGameDispatcher}>New Game</button>
    <button disabled={!$serviceStatus.Ready} on:click={continueGameDispatcher}>Continue Game</button>
    <button on:click={enterConfigDispatcher}>Configuration</button>
    <button on:click={QuitApplication}>Exit</button>
</div>

<style>

.menu {
    display: flex;
    flex-direction: column;
    align-items: center;
}

button {
    width: 400px;
    margin: 2rem;
}

button:disabled{
    cursor: wait;
}

</style>