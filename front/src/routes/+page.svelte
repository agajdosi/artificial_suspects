<svelte:head>
    <title>Artificial Witness | How AI Interprets People</title>
</svelte:head>

<script lang="ts">
    import { goto } from '$app/navigation';
    function newGame() {        
        goto('/new_game');
    }    
</script>

<div class="menu">

<h1>Artificial Witness</h1>

<p>
A playful investigation into how machines judge people.
Question the AI witness, catch the criminal, and uncover the distortions embedded in its data, its design, and the preferences of its builders and owners.
</p>

<button on:click={newGame}>Start Game</button>
<button disabled on:click={newGame}>Continue Game</button>


</div>

<style>

.menu {
    display: flex;
    flex-direction: column;
    align-items: center;
}
.menu > * {
    max-width: 500px;
}
button {
    width: 300px;
    margin: 2rem;
}


</style>
