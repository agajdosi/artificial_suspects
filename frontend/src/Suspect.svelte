<script lang="ts">
    export let susID: string;
    export let imgSrc: string;

    import logo from './assets/images/logo-universal.png'
    
    function selected() {
        console.log("Suspect selected!")
    }
</script>

<div class="suspect" id={susID} on:click={selected} on:keydown={selected}>
    <img src={imgSrc} alt="Suspect">
</div>

<style></style>