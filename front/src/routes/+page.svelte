<script lang="ts">
    import { serviceStatus } from '$lib/stores';
    import { goto } from '$app/navigation';
    
    function newGame() {
        goto('/game');
    }
    
    function continueGame() {
        goto('/game');
    }
    
    function enterConfig() {
        goto('/game');
    }
</script>

<h1>Artificial Suspects</h1>

<div class="menu">
    <button disabled={!$serviceStatus.ready} on:click={newGame}>New Game</button>
    <button disabled={!$serviceStatus.ready} on:click={continueGame}>Continue Game</button>
    <button on:click={enterConfig}>Configuration</button>
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
