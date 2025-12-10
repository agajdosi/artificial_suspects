<script lang="ts">
    import { goto } from '$app/navigation';
    import { page } from '$app/state';
</script>

<div>
    <button onclick={() => goto('/')} class="navigation" disabled={page.url.pathname==='/'}>↵</button>
    {#if !page.url.pathname.startsWith('/play')}
        <button onclick={() => goto('/new_game')} class="navigation">?</button>
    {/if}
</div>

<style>
button {
    all: unset;
    min-width: 25px;
}
button:hover{
    all: unset;
    min-width: 25px;
    cursor: pointer;
}
button:disabled{
    color: #666666;
}
.separator {
    margin: 0 1rem;
    background-color: var(--bg-color);
}
</style>