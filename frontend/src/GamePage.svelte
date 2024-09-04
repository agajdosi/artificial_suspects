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
        try {
            const isInnocent = await FreeSuspect(suspect.UUID, game.investigation.rounds[0].uuid);
            if (isInnocent) {
                console.log(`Eliminated Suspect ${suspect.UUID}`);
            } else {
                console.log(`Eliminated Criminal ${suspect.UUID}`);
            }
        } catch (error) {
            console.error(`Failed to free suspect ${suspect.UUID}:`, error);
        }

        game = await GetGame();
    }   
</script>

<button on:click={goToMenu}>Menu</button>
<h1>{game.investigation.rounds[0].question} '{game.investigation.rounds[0].answer}'</h1>
<Suspects suspects={game.investigation.suspects} on:suspect_freeing={handleSuspectFreeing} />

<button on:click={nextRound} disabled={!game.investigation.rounds[0].Eliminations} aria-disabled="{!game.investigation.rounds[0].Eliminations ? 'true': 'false'}">
    Next Question
</button>

<style>

</style>