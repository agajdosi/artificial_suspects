<script lang="ts">
    import { currentGame, serviceStatus, hint } from '$lib/stores';
    import { NextRound, EliminateSuspect, GetGame, NextInvestigation } from '$lib/main';
    import Suspects from '$lib/Suspects.svelte';
    import History from '$lib/History.svelte';
    import Scores from '$lib/Scores.svelte';
    import Help from '$lib/Help.svelte';
    import IntroOverlay from '$lib/IntroOverlay.svelte';
    import { locale, t } from 'svelte-i18n';
    import { createEventDispatcher } from 'svelte';
    import LanguageSwitch from '$lib/LanguageSwitch.svelte';

    let scoresVisible: boolean = true;
    let helpVisible: boolean = false;

    // HOME BUTTON

    const dispatch = createEventDispatcher();
    function goToMenu() {dispatch('message', { message: 'goToHome' });}


    function getHintNextQuestion(){
        if ($currentGame.investigation?.rounds?.at(-1)?.answer != "") return hint.set("Wait for the AI to answer the question.")
        if (!$currentGame.investigation?.rounds?.at(-1)?.Eliminations) return hint.set("Eliminate at least 1 suspect before proceeding to next question.");
        return hint.set("Proceed to next question.");
    }

    async function handleSuspectFreeing(event) {
        console.log("FREEING SUSPECT", event)
        const { suspect } = event.detail;
        try {
            await EliminateSuspect(suspect.UUID, $currentGame.investigation?.rounds?.at(-1)?.uuid, $currentGame.investigation?.uuid);
        } catch (error) {
            console.error(`Failed to free suspect ${suspect.UUID}:`, error);
        }
        const game = await GetGame();
        currentGame.set(game);
        console.log(`GAME OVER: ${game.GameOver}`);
    }

    async function nextInvestigation(){
        let game = await NextInvestigation();
        game.GameOver = false;
        currentGame.set(game);
        console.log("GOT NEW INVESTIGATION", game)
    }
    
    function newGame() {
        scoresVisible = true;
        dispatch('newGame', { 'game_uuid': $currentGame.uuid });
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

    let IdleTimer: NodeJS.Timeout | null = null;
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
        {#if $currentGame.investigation?.InvestigationOver}
            <div class="jailtime">
                {$t('arrest')}
            </div>
        {:else}
            <div
                class="question"
                role="tooltip"
                on:mouseenter={() => hint.set("A question about the wanted person, answered by an AI witness.")}
                on:mouseleave={() => hint.set("")}
                >
                {$currentGame.investigation?.rounds?.length}.
                {#if $locale == "cz"}
                    {$currentGame.investigation?.rounds?.at(-1)?.Question?.Czech}
                {:else if $locale == "pl"}
                    {$currentGame.investigation?.rounds?.at(-1)?.Question?.Polish}
                {:else}
                    {$currentGame.investigation?.rounds?.at(-1)?.Question?.English}
                {/if}
            </div>
            {#if $currentGame.investigation?.rounds?.at(-1)?.answer == ""}
                <div class="waiting"
                    role="tooltip"
                    on:mouseenter={() => hint.set("Waiting for the AI witness to answer the question.")}
                    on:mouseleave={() => hint.set("")}
                    >
                    *{$t('thinking')}*
                </div>
            {:else}
                <div class="answer"
                    role="tooltip"
                    on:mouseenter={() => hint.set("The AI witness' response to the question about the wanted person.")}
                    on:mouseleave={() => hint.set("")}
                    >
                    {$t($currentGame.investigation?.rounds?.at(-1)?.answer?.toLocaleLowerCase())}!
                </div>
            {/if}
        {/if}
        </div>
        <div class="instruction">
            {#if $currentGame.investigation?.InvestigationOver}
                {$t('arrestInstruction')}
            {:else if $currentGame.investigation?.rounds?.at(-1)?.answer != ""}
                {#if $currentGame.investigation?.rounds?.at(-1)?.answer?.toLowerCase() == "yes"}{$t('release-no')}
                {:else}{$t('release-yes')}
                {/if}
            {:else}
                {$t('waiting')}...
            {/if}
        </div>
    </div>
    <div class="top-right"
        role="tooltip"
        on:mouseenter={() => hint.set("Switch language of the user interface.")}
        on:mouseleave={() => hint.set("")}
        >
        <LanguageSwitch/>
    </div>
</div>

<div class="middle">
    <div class="left">
        <Suspects
            suspects={$currentGame.investigation?.suspects || []}
            gameOver={$currentGame.GameOver}
            investigationOver={$currentGame.investigation?.InvestigationOver}
            answerIsLoading={$currentGame.investigation?.rounds?.at(-1)?.answer == ""}
            on:suspect_freeing={handleSuspectFreeing}
            on:suspect_jailing={nextInvestigation}
        />

        <div class="actions">
            {#if !$currentGame.investigation?.InvestigationOver}
                {#if $currentGame.GameOver}
                    <button
                        on:click={newGame}
                        on:mouseenter={() => hint.set("Start a new game and try it again!")}
                        on:mouseleave={() => hint.set("")}
                        class="{!$serviceStatus.ready && 'offline'}"
                        disabled={!$serviceStatus.ready}>
                        {$t('buttons.newGame')}
                    </button>
                {:else}
                <button
                    on:click={NextRound}
                    on:mouseenter={() => getHintNextQuestion()}
                    on:mouseleave={() => hint.set("")}
                    class="{!$serviceStatus.ready && 'offline'}"
                    disabled={!$currentGame.investigation?.rounds?.at(-1)?.Eliminations || $currentGame.GameOver || !$serviceStatus.ready}
                    aria-disabled="{!$currentGame.investigation?.rounds?.at(-1)?.Eliminations || $currentGame.GameOver || !$serviceStatus.ready ? 'true': 'false'}"
                    >
                    {$t('buttons.nextQuestion')}
                </button>
                {/if}
            {/if}
        </div>
    </div>

    <div class="right">
        <History/>
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
            role="tooltip"
            on:mouseenter={() => hint.set("Successfully finish the investigation to get into higher level.")}
            on:mouseleave={() => hint.set("")}
            >
            level: {$currentGame.level}
        </div>
        <div
            role="tooltip"
            on:mouseenter={() => hint.set("Your current score. Free innocent suspects and finish the investigation to get more points.")}
            on:mouseleave={() => hint.set("")}
            >
            score: {$currentGame.Score}
        </div>
    </div>
</div>

{#if $currentGame.GameOver && scoresVisible}
    <Scores on:toggleScores={handleToggleScores} on:newGame/>    
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

</style>
