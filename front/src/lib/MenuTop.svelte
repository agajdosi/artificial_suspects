<script lang="ts">
    import { locale } from 'svelte-i18n';
    import { goto } from '$app/navigation';
    import { page } from '$app/state';  

    let {overlayConfigVisible = $bindable(false)} = $props();

    function changeLanguage(code: string) {
        $locale = code;
    }
</script>

<button onclick={() => changeLanguage('en')} class:active={$locale === 'en'}>en</button>
<button onclick={() => changeLanguage('cz')} class:active={$locale === 'cz'}>cz</button>
<button onclick={() => changeLanguage('pl')} class:active={$locale === 'pl'}>pl</button>

<style>
button {
    all: unset;
    text-decoration: underline;
    min-width: 25px;
}
button:hover{
    all: unset;
    min-width: 25px;
    cursor: pointer;
}
button.active {
    text-transform: uppercase;
}
</style>