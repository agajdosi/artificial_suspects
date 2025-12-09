<svelte:head>
    <title>Artificial Witness | How AI Interprets People</title>
</svelte:head>

<script lang="ts">
    import { goto } from '$app/navigation';
    import { t } from 'svelte-i18n';

    function newGame() {        
        goto('/new_game');
    }

</script>

<h1>Artificial Witness</h1>

<div class="menu">

<p>{$t('home.intro')}</p>

<button on:click={newGame}>{$t('buttons.newGame')}</button>
<button disabled on:click={newGame}>{$t('buttons.continueGame')}</button>


</div>

<style>

.menu {
    display: flex;
    flex-direction: column;
    align-items: center;
}
.menu > * {
    max-width: 910px;
}
button {
    width: 490px;
    margin: 1rem;
}
p {
    margin: 0 0 6rem 0;
}

</style>
