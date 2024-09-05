<script lang="ts">
    import { main } from '../wailsjs/go/models';
    import { NextRound, FreeSuspect, GetGame } from '../wailsjs/go/main/App.js';
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
            game = await NextRound();
        } catch (error) {
            console.log(`NextRound() has failed: ${error}`)
        }
        console.log(`>>> NEW ROUND: ${game.investigation.rounds.at(-1)}`)
    }

    async function handleSuspectFreeing(event) {
        const { suspect } = event.detail;
        try {
            await FreeSuspect(suspect.UUID, game.investigation.rounds.at(-1).uuid);
        } catch (error) {
            console.error(`Failed to free suspect ${suspect.UUID}:`, error);
        }
        game = await GetGame();
        console.log(`GAME OVER: ${game.GameOver}`);
    }   
</script>

<button on:click={goToMenu}>Menu</button>
<h1>{game.investigation.rounds.at(-1).question} '{game.investigation.rounds.at(-1).answer}'</h1>


<Suspects suspects={game.investigation.suspects} gameOver={game.GameOver} on:suspect_freeing={handleSuspectFreeing} />

<button
on:click={nextRound}
disabled={!game.investigation.rounds.at(-1).Eliminations || game.GameOver}
aria-disabled="{!game.investigation.rounds.at(-1).Eliminations || game.GameOver ? 'true': 'false'}"
>
    Next Question
</button>

<style>
.disabled-suspects {
    pointer-events: none;
    opacity: 0.5; /* Optional: To make it look visually disabled */
}
</style>