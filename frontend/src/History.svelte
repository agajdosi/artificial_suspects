<script lang="ts">
    import { database } from '../wailsjs/go/models';

    export let game: database.Game;
    export let language: string;
</script>

{#each [...game.investigation.rounds].reverse().slice(1).reverse() as round, index}
    <div class="round">
        <div class="question">
            {index+1}.
            {#if language == "cz"}
                {round.Question.Czech}
            {:else if language == "pl"}
                {round.Question.Polish}
            {:else}
                {round.Question.English}
            {/if}
        </div>
        <div class="answer">{round.answer}</div>
    </div>
{/each}

<style>
.round {
    display: flex;
}

.question, .answer {
    padding: 10px;
    border-radius: 10px;
    margin: 5px 0;
    position: relative;
    font-size: 16px;
    width: fit-content;
    max-width: 100%;
}

.question {
    background-color: #343563;
    align-self: flex-start;
    border-bottom-left-radius: 0;
}

.answer {
    background-color: #3c1c54;
    align-self: flex-end;
    border-bottom-right-radius: 0;
}
</style>
