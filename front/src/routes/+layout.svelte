<script lang="ts">
    import { register, init, waitLocale } from 'svelte-i18n';
    import { onMount } from 'svelte';
    import OverlayError from '$lib/OverlayError.svelte';

    register('en', () => import('$lib/locales/en.json'));
    register('cz', () => import('$lib/locales/cz.json'));
    register('pl', () => import('$lib/locales/pl.json'));

    init({
        fallbackLocale: 'en',
        initialLocale: 'en'
    });

    let isLocaleLoaded = false;

    onMount(async () => {
        await waitLocale(); // Ensure locale is ready before rendering
        isLocaleLoaded = true;
        const player = localStorage.getItem('player');
        console.log("Current player:", player);
    });
</script>

{#if isLocaleLoaded}
    <slot />
{:else}
    <p>Loading translations...</p>
{/if}

<OverlayError/>

<style>
:global(:root){
       --bg-color: rgba(27, 38, 54, 1);
}

:global(html) {
    background-color: var(--bg-color);
    text-align: center;
    color: white;
}

:global(body) {
    margin: 0;
    color: white;
    font-family: "Nunito", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto",
    "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue",
    sans-serif;
}

@font-face {
    font-family: "Nunito";
    font-style: normal;
    font-weight: 400;
    src: local(""), url("/nunito-v16-latin-regular.woff2") format("woff2");
}

@font-face {
    font-family: "Bitcount";
    font-style: normal;
    src: local(""), url("/Bitcount-VariableFont_CRSV,ELSH,ELXP,slnt,wght.ttf") format("truetype");
}

:global(h1) {
    margin: 6rem 0 1rem 0;
    font-size: 6rem;
    font-family: "Bitcount";
    font-weight: 390;
    font-variation-settings:
        "slnt" -8,
        "CRSV" 0,
        "ELSH" 0,
        "ELXP" 0;
}

:global(#app) {
    height: 100vh;
    text-align: center;
}

:global(button) {
    display: inline-block;
    outline: 0;
    text-align: center;
    cursor: pointer;
    padding: 5px 10px;
    border: 0;
    color: #fff;
    font-size: 17.5px;
    border: 2px solid transparent;
    border-color: #ffffff;
    color: #ffffff;
    background: transparent;
    transition: background,color .1s ease-in-out;
}
                
:global(button:hover) {
    background-color: #ffffff;
    color: #000000;
}

:global(button:disabled) {
    color: #666666;
    border-color:#666666;
    background-color: unset;
    cursor: not-allowed;
}

:global(button.selected) {
    background: #007bff;
    color: white;
    border-color: #0056b3;
}

:global(header) {
    display: flex;
    justify-content: end;
    padding: 3px 7px 0 0;
}

/* Use display: flow-root to create new Block Formatting Context.
BFC prevents margins to overflow outside main and #app, breaking its 100vh.
*/
:global(main) {
    display: flow-root;
}

</style>