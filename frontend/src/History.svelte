<script lang="ts">
    import { database } from '../wailsjs/go/models';

    export let game: database.Game;
</script>

{#each [...game.investigation.rounds].reverse().slice(1) as round}
    <div>
        {round.Question.English} - {round.answer}
    </div>
{/each}

<style>

</style>