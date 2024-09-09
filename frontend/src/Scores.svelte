<script lang="ts">
    import { main } from '../wailsjs/go/models';
    import { GetScores, SaveScore } from '../wailsjs/go/main/App.js';
    import { createEventDispatcher } from 'svelte';

    export let game: main.Game;
    let name: string;

    const dispatch = createEventDispatcher();

    function closeScores() {
        let scoresVisible: boolean = false;
        dispatch('toggleScores', { scoresVisible });
    }

    function newGameDispatcher() {
        dispatch('newGame', { 'game_uuid': game.uuid });
    }

    function saveScore() {
        SaveScore(name, game.uuid);
    }
</script>

<div class="infobox">
    <h1>Game Over!</h1>
    <p>
        You've mistakenly released a criminal, while innocent suspects have been unjustly persecuted.
        Next time, try to delve deeper into the mindset of the AI during its interrogation.
    </p>

    <h2>High Scores</h2>
    {#await GetScores() }
        Loading High Scores...
    {:then scores}
        <div class="scores">
            {#each scores as score}
                {#if score.Position == 1}
                    <div class="score-item">🥇 {score.Score} {score.Investigator}</div>
                {:else if score.Position == 2}
                    <div class="score-item">🥈 {score.Score} {score.Investigator}</div>
                {:else if score.Position == 3}
                    <div class="score-item">🥉  {score.Score} {score.Investigator}</div>
                {:else if score.Position <= 10}
                    {#if score.GameUUID == game.uuid}
                        <div class="score-item" class:highlighted={score.GameUUID == game.uuid}>
                            {score.Position}. {score.Score} <input bind:value={name} placeholder="enter your name" />
                        </div>
                    {:else}
                        <div class="score-item">
                            {score.Position}.  {score.Score} {score.Investigator}
                        </div>
                    {/if}
                {:else if score.GameUUID === game.uuid}
                    <div class="score-item">...</div>
                    <div class="score-item" class:highlighted={score.GameUUID == game.uuid}>
                        {score.Position}. {score.Score}
                        <input bind:value={name} placeholder="enter your name"/>
                        <button on:click={saveScore}>Confirm</button>  
                    </div>
                {/if}
            {/each}
        </div>
    {/await }

    <button on:click={closeScores}>Close</button>
    <button on:click={newGameDispatcher}>New Game</button>  
</div>

<style>
.scores {
    margin-top: 20px;
}

.score-item {
    margin-bottom: 8px;
    font-size: 18px;
}

.score-item input {
    margin: 0 0 0 2rem;
}

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
.infobox p {
    padding: 0 2rem;
}
</style>
