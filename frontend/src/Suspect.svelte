<script lang="ts">
    import { main } from '../wailsjs/go/models';
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();

    export let suspect: main.Suspect;

    const imgDir: string = 'src/assets/images/suspects/';
    let isFree = false; // innocent suspects can be free
    let fled = false; // if it was criminal

    async function selected() {
        if (isFree) return;
        dispatch('suspect_freeing', { 'suspect': suspect });
    }
</script>

<div 
    class="suspect {suspect.Free ? 'free' : ''} {suspect.Fled ? 'fled' : ''}"
    id={suspect.UUID} 
    on:click={selected}
    on:keydown={selected}
    aria-disabled={suspect.Free || suspect.Fled}
>
    <div class="suspect-image" style="background-image: url({imgDir+suspect.Image});"></div>
</div>

<style>
    .suspect {
        position: relative;
        height: 21vh;
        width: 21vh;
        margin: 1%;
        cursor: pointer;
        transition: opacity 0.3s ease, filter 0.3s ease;
    }

    .suspect.free .suspect-image {
        opacity: 0.2;
        filter: grayscale(100%);
        cursor: not-allowed;
    }

    .suspect.fled::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-color: rgba(255, 0, 0, 0.5); /* Red overlay */
        pointer-events: none; /* Ensure the overlay doesn’t interfere with clicks */
        transition: background-color 0.3s ease; /* Transition for the red overlay */
        opacity: 0; /* Initially hidden */
    }

    .suspect.fled .suspect-image {
        opacity: 0.2;
        filter: grayscale(100%);
    }

    .suspect.fled::before {
        opacity: 1; /* Fade in the red overlay */
    }

    .suspect-image {
        height: 100%;
        width: 100%;
        background-size: cover;
        background-position: center;
        background-repeat: no-repeat;
        transition: opacity 0.3s ease, filter 0.3s ease;
    }
</style>

