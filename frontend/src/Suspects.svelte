<script lang="ts">
    import { database } from '../wailsjs/go/models';
    import Suspect from './Suspect.svelte'
    
    export let suspects: database.Suspect[];
    export let gameOver: boolean;
    export let investigationOver: boolean;
    export let answerIsLoading: boolean;
</script>

<div class="suspects">
    {#each suspects as suspect}
        <Suspect suspect={suspect} {gameOver} {investigationOver} {answerIsLoading} on:suspect_freeing on:suspect_jailing />
    {/each}
</div>

<style>
    .suspects {
        justify-content: flex-start;
        display: flex;
        flex-wrap: wrap;
        width: 120vh;
        max-width: 100vw;
        max-height: 80vh;
        flex-shrink: 1;
    }
</style>