<script lang="ts">
    import { database } from '../wailsjs/go/models';

    export let game: database.Game;
</script>

{#each [...game.investigation.rounds].reverse().slice(1) as round}
    <div>
        {round.question} - {round.answer}
    </div>
{/each}

<style>

</style>