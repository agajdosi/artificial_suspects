<script lang="ts">
    import { goto } from '$app/navigation';
    function newGame() {
        goto('/play');
    }    
</script>

<h1>Artificial Suspects</h1>

<p>
A game exposing how different AI models form biased interpretations of people.</p>
<p>
Using the familiar structure of Unusual Suspects, the game lets an AI take the role of the witness.
Every decision it makes highlights subtle biases that usually stay invisible behind the training data, weights and code.
</p>

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
