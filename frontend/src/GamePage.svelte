<script lang="ts">
    import { main } from '../wailsjs/go/models';
    import { NextRound, FreeSuspect, GetGame } from '../wailsjs/go/main/App.js';
    import Suspects from './Suspects.svelte'
    import History from './History.svelte';

    export let game: main.Game;
    let hint: string = "hint..."; // TODO: capture hints

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
    
    function newGame() {dispatch('message', { message: 'newGame' });}
</script>

<div class="header">
    <button on:click={goToMenu}>Menu</button>
</div>

<div class="top">
    <div class="question">{game.investigation.rounds.at(-1).question}</div>
    {#if game.investigation.rounds.at(-1).answer != "" }
        <div>{game.investigation.rounds.at(-1).answer}</div>
    {:else}
        <div class="answer">💭</div>
    {/if}
</div>

<div class="middle">
    <div class="left">
        <Suspects suspects={game.investigation.suspects} gameOver={game.GameOver} on:suspect_freeing={handleSuspectFreeing} />

        <div class="actions">
            {#if !game.GameOver }
                <button
                    on:click={nextRound}
                    disabled={!game.investigation.rounds.at(-1).Eliminations || game.GameOver}
                    aria-disabled="{!game.investigation.rounds.at(-1).Eliminations || game.GameOver ? 'true': 'false'}"
                    >
                    Next Question
                </button>
            {:else}
                <button on:click={newGame}>
                    New Game
                </button>
            {/if}
        </div>
    </div>

    <div class="right">
        <div class="history"><History {game}/></div>
    </div>
</div>

<div class="bottom">
    <div class="hint">{hint}</div>
    <div class="stats">
        <div>level: {game.level}</div>
        <div>score: {game.Score}</div>
    </div>
</div>

<style>

.header {
    display: flex;
    justify-content: right;
}

.top {
    width: 100vw;
    display: flex;
    gap: 2rem;
    padding: 0.5rem 0 0 0;
    justify-content: left;
    font-size: 2rem;
}

.middle {
    display: flex;
}
.left .actions {
    padding: 2rem 0;
}
.right {
    padding: 2rem 0 0 0;
    width: 100%;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
}

.bottom {
    display: flex;
    justify-content: space-between;
    position: absolute;
    bottom: 0;
    width: calc(100vw - 1rem);
    padding: 0 0.5rem;
}
.stats {
    display: flex;
    gap: 1rem;
}

</style>