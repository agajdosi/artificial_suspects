<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { QuitApplication, AIServiceIsReady} from '../wailsjs/go/main/App';
    import ServiceStatus from './ServiceStatus.svelte';
    const dispatch = createEventDispatcher();
    
    function newGameDispatcher() {dispatch('newGame', {message: 'new game'});}
    function continueGameDispatcher() {dispatch('message', { message: 'continueGame' });}
    function enterConfigDispatcher() {dispatch('enterConfig', { message: 'enterConfig' });}
</script>

<h1>Artificial Suspects</h1>

<div class="menu">
    <button on:click={newGameDispatcher}>New Game</button>
    <button on:click={continueGameDispatcher}>Continue Game</button>
    <button on:click={enterConfigDispatcher}>Configuration</button>
    <button on:click={QuitApplication}>Exit</button>
</div>

<ServiceStatus />

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