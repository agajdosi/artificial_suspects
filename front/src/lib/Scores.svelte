<script lang="ts">
    import { currentGame, serviceStatus, hint } from '$lib/stores';
    import { NewGame, type FinalScore } from '$lib/main';
    import { GetScores, SaveScore } from '$lib/main';
    import { createEventDispatcher, onMount } from 'svelte';
    import { t } from 'svelte-i18n';

    let name: string;
    let scores: FinalScore[] = [];
    let loading: boolean = true;

    const dispatch = createEventDispatcher();

    // Fetch the scores when the component is mounted
    onMount(async () => {
        try {
            scores = await GetScores();
        } catch (error) {
            console.error('Error fetching scores:', error);
        } finally {
            loading = false;
        }
    });

    function closeScores() {
        dispatch('toggleScores', { scoresVisible: false });
    }

    function saveScore() {
        SaveScore(name, $currentGame.uuid);
    }

    // Helper function to return the position label (medal or rank)
    function getPositionLabel(position: number): string {
        if (position === 1) return '🥇';
        if (position === 2) return '🥈';
        if (position === 3) return '🥉';
        return `${position}.`;
    }

    // Function to check if the current score belongs to the current game
    function isCurrentGame(scoreUUID: string): boolean {
        return scoreUUID === $currentGame.uuid;
    }

    function getHintNewGame() {
        if (!$serviceStatus.ready) return hint.set("Cannot start new game, AI service is not ready!");
        return hint.set("Start a new game and try it again!");
    }
</script>

<div class="infobox">
    <h1>{$t('gameOver.gameOver')}</h1>
    <div class="riptext">{$t('gameOver.riptext')}</div>

    <h2>{$t('gameOver.highScores')}</h2>
    {#if loading}
        {$t('gameOver.loadingScores')}
    {:else}
        <div class="scores">
            {#each scores as score, index}
                {#if index < 10 || isCurrentGame(score.GameUUID)}                
                    {#if index >= 10 && isCurrentGame(score.GameUUID)}
                        <div class="score-item">...</div>
                    {/if}
                    <div class="score-item" class:highlighted={isCurrentGame(score.GameUUID)}>
                        <span class="position">{getPositionLabel(index + 1)}</span>
                        
                        {#if isCurrentGame(score.GameUUID)}
                            <span>
                                <input
                                    bind:value={name}
                                    on:mouseenter={() => hint.set("Inscribe your name to the leaderboards.")}
                                    on:mouseleave={() => hint.set("")}
                                    placeholder="{$t('gameOver.enterName')}"
                                />
                                <button
                                    on:click={saveScore}
                                    on:mouseenter={() => hint.set("Confirm your name and save.")}
                                    on:mouseleave={() => hint.set("")}
                                    >
                                    {$t('buttons.confirm')}
                                </button>
                            </span>
                        {:else}
                            {score.Investigator}
                        {/if}
                        <span class="score">{score.Score}</span>
                    </div>

                {/if}
            {/each}
        </div>
    {/if}

    <button
        on:click={closeScores}
        on:mouseenter={() => hint.set("Close scores window.")}
        on:mouseleave={() => hint.set("")}
        >
        {$t('buttons.close')}
    </button>
    <button
        on:click={NewGame}
        on:mouseenter={() => getHintNewGame()}
        on:mouseleave={() => hint.set("")}
        class:offline={!$serviceStatus.ready}
        disabled={!$serviceStatus.ready}
        >
        {$t('buttons.newGame')}
    </button>  
</div>

<style>
h1 {
    margin: 10px 0 0 0;
}

.scores {
    margin: 30px 0;
    display: flex;
    flex-direction: column;
    align-items: center; /* Center align items horizontally */
    text-align: center;  /* Center align text */
}

.score-item {
    max-width: 90%;
    margin-bottom: 8px;
    font-size: 18px;
    display: flex;
    justify-content: space-between;
    align-items: center;  /* Align items vertically in the center */
    width: 100%;  /* Make sure it spans the full width */
}

.score-item input {
    margin: 0 0 0 2rem;
}


.highlighted {
    background-color: rgb(255, 89, 0);
    padding: 5px;
    border-radius: 5px;
}

.infobox {
    position: absolute;
    left: 25vw;
    top: 10vh;
    background-color: grey;
    width: 50vw;
    height: 80vh;
}
.riptext {
    padding: 0 2rem;
}
.offline {
    cursor: wait;
}
</style>
