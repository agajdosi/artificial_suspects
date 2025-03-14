<script lang="ts">
    import { register, init, locale, waitLocale } from 'svelte-i18n';
    import { onMount } from 'svelte';
    import ErrorOverlay from '$lib/ErrorOverlay.svelte';
    import ServiceStatus from '$lib/ServiceStatus.svelte';

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
    });
</script>

{#if isLocaleLoaded}
    <slot />
{:else}
    <p>Loading translations...</p>
{/if}

<ErrorOverlay on:message={handleMessage}/>
<ServiceStatus/>

<style>
:global(html) {
    background-color: rgba(27, 38, 54, 1);
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
    src: local(""),
    url("/nunito-v16-latin-regular.woff2") format("woff2");
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

/* Use display: flow-root to create new Block Formatting Context.
BFC prevents margins to overflow outside main and #app, breaking its 100vh.
*/
:global(main) {
    display: flow-root;
}

</style>