<script lang="ts">
    import { goto } from '$app/navigation';
    import { page } from '$app/state';
</script>

<div class="navigation">
    <button onclick={() => goto('/')} class="navigation" disabled={page.url.pathname==='/'}>↵</button>
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
</style>