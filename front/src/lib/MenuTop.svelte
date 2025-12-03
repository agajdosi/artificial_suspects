<script lang="ts">
    import {locale} from 'svelte-i18n';
    import { NewGame } from '$lib/main';
    import { goto } from '$app/navigation';

    let {overlayConfigVisible = $bindable(false)} = $props();

    function changeLanguage(code: string) {
        $locale = code;
    }

    async function toggleOverlayConfig() {
        console.log("toggleOverlayConfig", overlayConfigVisible);
        overlayConfigVisible = !overlayConfigVisible;
    }
</script>

<button onclick={() => goto('/')}>Home</button>
<button onclick={() => goto('/new_game')}>Restart</button>

<span class="separator"></span>

<button onclick={() => changeLanguage('en')} class:active={$locale === 'en'}>en</button>
<button onclick={() => changeLanguage('cz')} class:active={$locale === 'cz'}>cz</button>
<button onclick={() => changeLanguage('pl')} class:active={$locale === 'pl'}>pl</button>

<style>
button {
    all: unset;
    text-decoration: underline;
    min-width: 20px;
}
button:hover{
    cursor: pointer;
}
button.active {
    text-transform: uppercase;
}

.separator {
    margin: 0 1rem;
    background-color: var(--bg-color);
}
</style>