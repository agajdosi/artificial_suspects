<script lang="ts">
    import { main } from '../wailsjs/go/models';

    export let game: main.Game;
</script>

{#each [...game.investigation.rounds].reverse() as round}
    <div>
        "{round.question}" - "{round.answer}"
    </div>
{/each}

<style>

</style>