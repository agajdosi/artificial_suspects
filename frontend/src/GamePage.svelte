<script lang="ts">
    import { main } from '../wailsjs/go/models';
    import { NextRound, EliminateSuspect, GetGame, WaitForAnswer, NextInvestigation } from '../wailsjs/go/main/App.js';
    import Suspects from './Suspects.svelte';
    import History from './History.svelte';
    import Scores from './Scores.svelte';

    export let game: main.Game;
    let lastRoundUUID: string;
    let answerIsLoading: boolean;
    let answer: string;
    let hint: string = "hint...";  // TODO: capture hints
    let scoresVisible: boolean = true;

    // HOME BUTTON
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();
    function goToMenu() {dispatch('message', { message: 'goToHome' });}

    // NEXT QUESTION
    async function nextRound() {
        try {
            game = await NextRound();
        } catch (error) {
            console.log(`NextRound() has failed: ${error}`);
        }
        console.log(`>>> NEW ROUND: ${game.investigation.rounds.at(-1)}`);
    }

    async function handleSuspectFreeing(event) {
        const { suspect } = event.detail;
        try {
            await EliminateSuspect(suspect.UUID, game.investigation.rounds.at(-1).uuid, game.investigation.uuid);
        } catch (error) {
            console.error(`Failed to free suspect ${suspect.UUID}:`, error);
        }
        game = await GetGame();
        console.log(`GAME OVER: ${game.GameOver}`);
    }

    async function LoadAnswer(roundUUID: string) {
        answerIsLoading = true;
        try {
            answer = await WaitForAnswer(roundUUID);
        } catch (error) {
            console.error("Failed to get the answer:", error);
            answer = "Error fetching answer";
        } finally {
            answerIsLoading = false;
        }
    }

    // Run when ...rounds.at(-1).QuestionUUID changes
    $: if (game.investigation.rounds) {
        const currentRoundUUID = game.investigation.rounds.at(-1).uuid;
        if (currentRoundUUID !== lastRoundUUID) {
            lastRoundUUID = currentRoundUUID;
            LoadAnswer(currentRoundUUID)
        }
    }

    async function nextInvestigation(){
        game = await NextInvestigation();
        game.GameOver = false;
        console.log("GOT NEW INVESTIGATION", game)
    }
    
    function newGame() {
        scoresVisible = true;
        dispatch('newGame', { 'game_uuid': game.uuid });
    }

    // Scores
    function handleToggleScores(event) {
        scoresVisible = event.detail.scoresVisible;
    }
</script>

<div class="header">
    <button on:click={goToMenu}>Menu</button>
</div>

<div class="top">
    {#if game.investigation.InvestigationOver}
        <div class="jailtime">Arrest the Perp!</div>
    {:else}
        <div class="question">{game.investigation.rounds.at(-1).question}</div>
        {#if answerIsLoading}
            <div class="waiting">...waiting for answer</div>
        {:else}
            <div class="answer">{answer}</div>
        {/if}
    {/if}
</div>

<div class="middle">
    <div class="left">
        <Suspects
            suspects={game.investigation.suspects}
            gameOver={game.GameOver}
            investigationOver={game.investigation.InvestigationOver}
            {answerIsLoading}
            on:suspect_freeing={handleSuspectFreeing}
            on:suspect_jailing={nextInvestigation}
        />

        <div class="actions">
            {#if !game.investigation.InvestigationOver}
                {#if game.GameOver}
                    <button on:click={newGame}>New Game</button>
                {:else}
                <button
                        on:click={nextRound}
                        disabled={!game.investigation.rounds.at(-1).Eliminations || game.GameOver}
                        aria-disabled="{!game.investigation.rounds.at(-1).Eliminations || game.GameOver ? 'true': 'false'}"
                        >
                        Next Question
                    </button>
                {/if}
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

{#if game.GameOver && scoresVisible}
    <Scores {game} on:toggleScores={handleToggleScores} on:newGame/>    
{/if}

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
