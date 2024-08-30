<script lang="ts">
    export let suspectUUID: string;
    export let imgSrc: string;

    import { FreeSuspect } from '../wailsjs/go/main/App.js';

    let isFree = false; // innocent suspects can be free
    let fled = false; // if it was criminal

    async function selected() {
        if (isFree) {
            return;
        }

        try {
            const isInnocent = await FreeSuspect(suspectUUID);
            if (isInnocent) {
                isFree = true;
                console.log(`Suspect ${suspectUUID} freed`);
                return;
            }
            fled = true;
            console.log(`Criminal ${suspectUUID} released`);
            return;
        } catch (error) {
            console.error(`Failed to free suspect ${suspectUUID}:`, error);
        }
    }
</script>

<div 
    class="suspect {isFree ? 'free' : ''} {fled ? 'fled' : ''}"
    id={suspectUUID} 
    on:click={selected}
    on:keydown={selected}
    aria-disabled={isFree}
>
    <div class="suspect-image" style="background-image: url({imgSrc});"></div>
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

