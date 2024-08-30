<script lang="ts">
    import { main } from '../wailsjs/go/models';
    import Suspects from './Suspects.svelte'
    export let game: main.Game; // HOW THE FUCK import from models.ts?

    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();
    function goToMenu() {
        dispatch('message', { message: 'goToHome' });
    }
</script>


<button on:click={goToMenu}>Menu</button>
<Suspects suspects={game.suspects}/>

<style>

</style>