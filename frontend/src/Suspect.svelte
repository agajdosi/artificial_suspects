<script lang="ts">
    export let susID: string;
    export let imgSrc: string;

    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();
    function selected() {
        dispatch('message', {'susID': susID});
    }

</script>

<div class="suspect" id={susID} on:click={selected} on:keydown={selected} style="background-image: url({imgSrc});"></div>

<style>
    .suspect {
        height: 21vh;
        width: 21vh;
        margin: 1%;
        background-size: cover; /* Cover to ensure the image covers the entire div */
        background-position: center; /* Center the image */
        background-repeat: no-repeat; /* Prevent the image from repeating */
    }

</style>