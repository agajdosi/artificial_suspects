<script lang="ts">
    import { main } from '../wailsjs/go/models';
    import { NextLevel, FreeSuspect, GetGame } from '../wailsjs/go/main/App.js';
    import Suspects from './Suspects.svelte'
    
    export let game: main.Game;

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

    async function handleSuspectFreeing(event) {
        const { suspect } = event.detail;
        console.log("Freeing suspect:", suspect)

        try {
            const isInnocent = await FreeSuspect(suspect.UUID, game.investigation.rounds[0].uuid);
            if (isInnocent) {
                console.log(`Suspect ${suspect.UUID} freed`);
            } else {
                console.log(`Criminal ${suspect.UUID} released`);
            }
        } catch (error) {
            console.error(`Failed to free suspect ${suspect.UUID}:`, error);
        }

        game = await GetGame();
        console.log("GOT GAME:", game)
    }   
</script>


<button on:click={goToMenu}>Menu</button>
<h1>{game.investigation.rounds[0].question}</h1>
<Suspects suspects={game.investigation.suspects} on:suspect_freeing={handleSuspectFreeing} />
<button on:click={nextRound}>Next Question</button>

<style>

</style>