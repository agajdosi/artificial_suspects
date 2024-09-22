<script lang="ts">
    import { database } from '../wailsjs/go/models';
    import { GetScores, SaveScore } from '../wailsjs/go/main/App.js';
    import { createEventDispatcher, onMount } from 'svelte';
    import { t } from 'svelte-i18n';

    export let game: database.Game;
    let name: string;
    let scores: database.FinalScore[] = [];
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

    function newGameDispatcher() {
        dispatch('newGame', { 'game_uuid': game.uuid });
    }

    function saveScore() {
        SaveScore(name, game.uuid);
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
        return scoreUUID === game.uuid;
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
                            <input bind:value={name} placeholder="{$t('gameOver.enterName')}" />
                            <button on:click={saveScore}>{$t('buttons.confirm')}</button>
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

    <button on:click={closeScores}>{$t('buttons.close')}</button>
    <button on:click={newGameDispatcher}>{$t('buttons.newGame')}</button>  
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
</style>
