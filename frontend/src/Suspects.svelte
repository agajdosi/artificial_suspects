<script lang="ts">
    import { main } from '../wailsjs/go/models';
    import Suspect from './Suspect.svelte'
    
    export let suspects: main.Suspect[];
</script>

<div class="suspects">
    {#each suspects as suspect}
        <Suspect suspect={suspect} on:message />
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