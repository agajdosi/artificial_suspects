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
            console.log(`NewGame() has failed: ${error}`)
        }
        console.log("GOT ROUNDS", game.investigation.rounds)
    }

    async function handleSuspectFreeing(event) {
        const { suspect } = event.detail;
        try {
            const isInnocent = await FreeSuspect(suspect.UUID, game.investigation.rounds.at(-1).uuid);
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
<h1>{game.investigation.rounds.at(-1).question} '{game.investigation.rounds.at(-1).answer}'</h1>
<Suspects suspects={game.investigation.suspects} on:suspect_freeing={handleSuspectFreeing} />

<button on:click={nextRound} disabled={!game.investigation.rounds.at(-1).Eliminations} aria-disabled="{!game.investigation.rounds.at(-1).Eliminations ? 'true': 'false'}">
    Next Question
</button>

<style>

</style>