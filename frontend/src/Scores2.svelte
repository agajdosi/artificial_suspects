<script lang="ts">
    import { main } from '../wailsjs/go/models';
    import { GetScores } from '../wailsjs/go/main/App.js';

    export let game: main.Game;
</script>

<h2>High Scores</h2>
{#await GetScores() }
    Loading High Scores...
{:then scores}
    <div class="scores">
        {#each scores as score}
            {#if score.Position == 1}
                <div class="score-item">🥇 {score.Investigator} - {score.Score}</div>
            {:else if score.Position == 2}
                <div class="score-item">🥈 {score.Investigator} - {score.Score}</div>
            {:else if score.Position == 3}
                <div class="score-item">🥉 {score.Investigator} - {score.Score}</div>
            {:else if score.Position <= 10}
                {#if score.GameUUID == game.uuid}
                    <div class="score-item" class:highlighted={score.GameUUID == game.uuid}>
                        {score.Position}. ENTER YOUR NAME - {score.Score}
                    </div>
                {:else}
                    <div class="score-item">
                        {score.Position}. {score.Investigator} - {score.Score}
                    </div>
                {/if}
            {:else if score.GameUUID === game.uuid}
                <div class="score-item">...</div>
                <div class="score-item" class:highlighted={score.GameUUID == game.uuid}>
                    {score.Position}. ENTER YOUR NAME - {score.Score}
                </div>
            {/if}
        {/each}
    </div>
{/await }

<style>
.scores {
    margin-top: 20px;
}

.score-item {
    margin-bottom: 8px;
    font-size: 18px;
}

/* Style for the first, second, and third places */
.scores .score-item:nth-child(1) {
    font-weight: bold;
    font-size: 1.6rem;
}

.scores .score-item:nth-child(2) {
    font-weight: bold;
    font-size: 1.4rem;
}

.scores .score-item:nth-child(3) {
    font-weight: bold;
    font-size: 1.3rem;
}

/* Style for highlighted current player's score */
.highlighted {
    background-color: rgb(255, 89, 0);
    padding: 5px;
    border-radius: 5px;
}
</style>
