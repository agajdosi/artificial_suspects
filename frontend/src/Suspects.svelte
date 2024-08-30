<script lang="ts">
    import Suspect from './Suspect.svelte'

    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();
    function goToMenu() {
        dispatch('message', { message: 'goToHome' });
    }
</script>

<button on:click={goToMenu}>Menu</button>

<div class="suspects">
    <Suspect suspectUUID="sus1" imgSrc="src/assets/images/suspects/1.jpg" on:message />
    <Suspect suspectUUID="sus2" imgSrc="src/assets/images/suspects/2.jpg" on:message />
    <Suspect suspectUUID="sus3" imgSrc="src/assets/images/suspects/1.jpg" on:message />
    <Suspect suspectUUID="sus4" imgSrc="src/assets/images/suspects/2.jpg" on:message />
    <Suspect suspectUUID="sus5" imgSrc="src/assets/images/suspects/1.jpg" on:message />

    <Suspect suspectUUID="sus6" imgSrc="src/assets/images/suspects/2.jpg" on:message />
    <Suspect suspectUUID="sus7" imgSrc="src/assets/images/suspects/1.jpg" on:message />
    <Suspect suspectUUID="sus8" imgSrc="src/assets/images/suspects/2.jpg" on:message />
    <Suspect suspectUUID="sus9" imgSrc="src/assets/images/suspects/1.jpg" on:message />
    <Suspect suspectUUID="sus10" imgSrc="src/assets/images/suspects/2.jpg" on:message />
    
    <Suspect suspectUUID="sus11" imgSrc="src/assets/images/suspects/1.jpg" on:message />
    <Suspect suspectUUID="sus12" imgSrc="src/assets/images/suspects/2.jpg" on:message />
    <Suspect suspectUUID="sus13" imgSrc="src/assets/images/suspects/1.jpg" on:message />
    <Suspect suspectUUID="sus14" imgSrc="src/assets/images/suspects/2.jpg" on:message />
    <Suspect suspectUUID="sus15" imgSrc="src/assets/images/suspects/1.jpg" on:message />
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