<script lang="ts">
    import { GetGame } from '../wailsjs/go/main/App.js';
    import { main } from '../wailsjs/go/models';

    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();
    
    function newGame() {dispatch('message', { message: 'newGame' });}
    function continueGame() {dispatch('message', { message: 'continueGame' });}

    async function getGame(){
        let game: main.Game
        try {
            game = await GetGame();
            } catch (error) {
                console.log(`NewGame() has failed: ${error}`)
            }
            console.log(game)
    }

</script>

<h1>Unusual Suspects</h1>

<div class="menu">
    <button on:click={newGame}>NEW GAME</button>
    <button on:click={continueGame}>CONTINUE GAME</button>
    <button on:click={getGame}>GET GAME</button>
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