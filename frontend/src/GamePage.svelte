<script lang="ts">
    import Suspects from './Suspects.svelte'

    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();
    function goToMenu() {
        dispatch('message', { message: 'goToHome' });
    }
</script>


<button on:click={goToMenu}>Menu</button>
<Suspects/>

<style>

</style>