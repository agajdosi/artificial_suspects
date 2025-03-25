<script lang="ts">
    import { goto } from '$app/navigation';
    function newGame() {
        goto('/play');
    }    
</script>

<h1>Artificial Suspects</h1>

<div class="menu">
    <button on:click={newGame}>Start Game</button>
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
