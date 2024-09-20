<script lang="ts">
    import { database } from '../wailsjs/go/models';
    import { NextRound, EliminateSuspect, GetGame, WaitForAnswer, NextInvestigation } from '../wailsjs/go/main/App.js';
    import Suspects from './Suspects.svelte';
    import History from './History.svelte';
    import Scores from './Scores.svelte';

    export let game: database.Game;
    let lastRoundUUID: string;
    let answerIsLoading: boolean;
    let answer: string;
    let hint: string = "hint...";  // TODO: capture hints
    let scoresVisible: boolean = true;

    const czech: string = "cz";
    const polish: string = "pl";
    const english: string = "en";
    let language: string = english; // TODO: probably should be set in App.svelte to make it global

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

    function changeLanguage(code: string) {
        language = code;
    }

</script>

<div class="top">
    <div class="top-left">
        <div class="main">
        {#if game.investigation.InvestigationOver}
            <div class="jailtime">Arrest the Perp!</div>
        {:else}
            <div class="question">
                {game.investigation.rounds.length}.
                {#if language == czech}
                    {game.investigation.rounds.at(-1).Question.Czech}
                {:else if language == polish}
                    {game.investigation.rounds.at(-1).Question.Polish}
                {:else}
                    {game.investigation.rounds.at(-1).Question.English}
                {/if}
            </div>
            {#if answerIsLoading}
                <div class="waiting">*thinking*</div>
            {:else}
                <div class="answer">{answer.toUpperCase()}!</div>
            {/if}
        {/if}
        </div>
        <div class="instruction">
            {#if !answerIsLoading}
                {#if answer == "yes"}Release those who aren't/doesn't.
                {:else}Release those who are/do.
                {/if}
            {:else}Waiting for the answer...
            {/if}
        </div>
    </div>
    <div class="top-right">
        <button on:click={() => changeLanguage('en')} class="langbtn" class:active={language === 'en'}>en</button>
        <button on:click={() => changeLanguage('cz')} class="langbtn" class:active={language === 'cz'}>cz</button>
        <button on:click={() => changeLanguage('pl')} class="langbtn" class:active={language === 'pl'}>pl</button>
    </div>
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

.middle {
    display: flex;
}
.left .actions {
    padding: 2rem 0;
}
.right {
    padding: 0.2rem 0 0 0;
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

.top {
    width: 100vw;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
}

.top .main {
    padding: 1rem 0 0 0;
    display: flex;
    flex-direction: row;
    gap: 1rem;
    margin: 0 0 0 1rem;
    font-size: 2rem;
}

.top .instruction {
    font-size: 1.2rem;
    display: flex;
    margin: 0 0 0 1.1rem;
}

.top-right {
    padding: 3px 7px 0 0;
}

.langbtn {
    all: unset;
    text-decoration: underline;
    min-width: 20px;
}
.langbtn:hover{
    cursor: pointer;
}
.langbtn.active {
    text-transform: uppercase;
}

.history {
    display: flex;
    flex-direction: column-reverse;
}

</style>
