<script lang="ts">
    import { main } from '../wailsjs/go/models';
    import { NextLevel } from '../wailsjs/go/main/App.js';

    import Suspects from './Suspects.svelte'
    
    export let game: main.Game; // HOW THE FUCK import from models.ts?

    // HOME BUTTON
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();
    function goToMenu() {
        dispatch('message', { message: 'goToHome' });
    }

    // NEXT QUESTION
    async function nextRound() {
        try {
                game = await NextLevel();

            } catch (error) {
                console.log(`NewGame() has failed: ${error}`)
            }
    }
</script>


<button on:click={goToMenu}>Menu</button>
<h1>{game.question}</h1>
<Suspects suspects={game.suspects}/>
<button on:click={nextRound}>Next Question</button>

<style>

</style>