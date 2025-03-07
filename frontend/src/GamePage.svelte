<script lang="ts">
    import { serviceStatus, hint } from './lib/stores';
    import { database } from '../wailsjs/go/models';
    import { NextRound, EliminateSuspect, GetGame, WaitForAnswer, NextInvestigation } from '../wailsjs/go/main/App.js';
    import Suspects from './Suspects.svelte';
    import History from './History.svelte';
    import Scores from './Scores.svelte';
    import Help from './Help.svelte';
    import IntroOverlay from './IntroOverlay.svelte';
    import { locale, t } from 'svelte-i18n';

    export let game: database.Game;
    let lastRoundUUID: string;
    let answerIsLoading: boolean;
    let answer: string;
    let scoresVisible: boolean = true;
    let helpVisible: boolean = false;

    // HOME BUTTON
    import { createEventDispatcher } from 'svelte';
    import LanguageSwitch from './LanguageSwitch.svelte';
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
    function getHintNextQuestion(){
        if (answerIsLoading) return hint.set("Wait for the AI to answer the question.")
        if (!game.investigation.rounds.at(-1).Eliminations) return hint.set("Eliminate at least 1 suspect before proceeding to next question.");
        return hint.set("Proceed to next question.");
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

    //HELP
    function toggleHelp() {
        helpVisible = !helpVisible;
        scoresVisible = false;
    }
    function handleToggleHelp(event) {
        helpVisible = event.detail.helpVisible;
    }

    //INTRO
    let introVisible: boolean = false;
    function handleToggleIntro(event) {
        introVisible = event.detail.introVisible;
    }

    let IdleTimer: number | null = null;
    window.addEventListener('mousemove', resetIdleTimer); // (re)sets IdleTimer
    function resetIdleTimer(): void {
        const msTimeout = 5 * 60 * 1000;
        if (IdleTimer) {
            clearTimeout(IdleTimer);
        }
        IdleTimer = setTimeout(userIsIdle, msTimeout);
    }
    function userIsIdle() {
        introVisible = true;
    }

</script>

<div class="top">
    <div class="top-left">
        <div class="main">
        {#if game.investigation.InvestigationOver}
            <div class="jailtime">
                {$t('arrest')}
            </div>
        {:else}
            <div
                class="question"
                on:mouseenter={() => hint.set("A question about the wanted person, answered by an AI witness.")}
                on:mouseleave={() => hint.set("")}
                >
                {game.investigation.rounds.length}.
                {#if $locale == "cz"}
                    {game.investigation.rounds.at(-1).Question.Czech}
                {:else if $locale == "pl"}
                    {game.investigation.rounds.at(-1).Question.Polish}
                {:else}
                    {game.investigation.rounds.at(-1).Question.English}
                {/if}
            </div>
            {#if answerIsLoading}
                <div class="waiting"
                    on:mouseenter={() => hint.set("Waiting for the AI witness to answer the question.")}
                    on:mouseleave={() => hint.set("")}
                    >
                    *{$t('thinking')}*
                </div>
            {:else}
                <div class="answer"
                    on:mouseenter={() => hint.set("The AI witness' response to the question about the wanted person.")}
                    on:mouseleave={() => hint.set("")}
                    >
                    {$t(answer.toLocaleLowerCase())}!
                </div>
            {/if}
        {/if}
        </div>
        <div class="instruction">
            {#if game.investigation.InvestigationOver}
                {$t('arrestInstruction')}
            {:else if !answerIsLoading}
                {#if answer.toLowerCase() == "yes"}{$t('release-no')}
                {:else}{$t('release-yes')}
                {/if}
            {:else}
                {$t('waiting')}...
            {/if}
        </div>
    </div>
    <div class="top-right"
        on:mouseenter={() => hint.set("Switch language of the user interface.")}
        on:mouseleave={() => hint.set("")}
        >
        <LanguageSwitch/>
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
                    <button
                        on:click={newGame}
                        on:mouseenter={() => hint.set("Start a new game and try it again!")}
                        on:mouseleave={() => hint.set("")}
                        class="{!$serviceStatus.Ready && 'offline'}"
                        disabled={!$serviceStatus.Ready}>
                        {$t('buttons.newGame')}
                    </button>
                {:else}
                <button
                    on:click={nextRound}
                    on:mouseenter={() => getHintNextQuestion()}
                    on:mouseleave={() => hint.set("")}
                    class="{!$serviceStatus.Ready && 'offline'}"
                    disabled={!game.investigation.rounds.at(-1).Eliminations || game.GameOver || !$serviceStatus.Ready}
                    aria-disabled="{!game.investigation.rounds.at(-1).Eliminations || game.GameOver || !$serviceStatus.Ready ? 'true': 'false'}"
                    >
                    {$t('buttons.nextQuestion')}
                </button>
                {/if}
            {/if}
        </div>
    </div>

    <div class="right">
        <History {game}/>
    </div>
</div>

<div class="bottom">
    <div class="help">
        <button on:click={toggleHelp} class="langbtn">{$t('buttons.help')}</button>
    </div>
    <div class="hint">
        <span>{$hint}</span>
    </div>
    <div class="stats">
        <div
            on:mouseenter={() => hint.set("Successfully finish the investigation to get into higher level.")}
            on:mouseleave={() => hint.set("")}
            >
            level: {game.level}
        </div>
        <div
            on:mouseenter={() => hint.set("Your current score. Free innocent suspects and finish the investigation to get more points.")}
            on:mouseleave={() => hint.set("")}
            >
            score: {game.Score}
        </div>
    </div>
</div>

{#if game.GameOver && scoresVisible}
    <Scores {game} on:toggleScores={handleToggleScores} on:newGame/>    
{/if}

{#if helpVisible}
    <Help on:toggleHelp={handleToggleHelp}/>    
{/if}

{#if introVisible}
    <IntroOverlay on:toggleIntro={handleToggleIntro}/>
{/if}

<style>
.middle {
    display: flex;
}
.left .actions {
    padding: 2rem 0;
}
.right {
    padding: 0.2rem 0 0 0;
    width: 100%;
    max-height: 73vh;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
}

.bottom {
    display: flex;
    justify-content: space-between;
    position: absolute;
    bottom: 0;
    width: calc(100vw - 2.5rem);
    padding: 0 0.5rem 0 2rem;
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

.answer {
    text-transform: uppercase;
}

.offline {
    cursor: wait;
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

</style>
