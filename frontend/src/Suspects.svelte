<script lang="ts">
    import Suspect from './Suspect.svelte'
</script>

<div class="suspects">
    <Suspect susID="sus1" imgSrc="src/assets/images/suspects/1.jpg"/>
    <Suspect susID="sus2" imgSrc="src/assets/images/suspects/2.jpg"/>
    <Suspect susID="sus3" imgSrc="src/assets/images/suspects/1.jpg"/>
    <Suspect susID="sus4" imgSrc="src/assets/images/suspects/2.jpg"/>
    <Suspect susID="sus5" imgSrc="src/assets/images/suspects/1.jpg"/>

    <Suspect susID="sus6" imgSrc="src/assets/images/suspects/2.jpg"/>
    <Suspect susID="sus7" imgSrc="src/assets/images/suspects/1.jpg"/>
    <Suspect susID="sus8" imgSrc="src/assets/images/suspects/2.jpg"/>
    <Suspect susID="sus9" imgSrc="src/assets/images/suspects/1.jpg"/>
    <Suspect susID="sus10" imgSrc="src/assets/images/suspects/2.jpg"/>
    
    <Suspect susID="sus11" imgSrc="src/assets/images/suspects/1.jpg"/>
    <Suspect susID="sus12" imgSrc="src/assets/images/suspects/2.jpg"/>
    <Suspect susID="sus13" imgSrc="src/assets/images/suspects/1.jpg"/>
    <Suspect susID="sus14" imgSrc="src/assets/images/suspects/2.jpg"/>
    <Suspect susID="sus15" imgSrc="src/assets/images/suspects/1.jpg"/>
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