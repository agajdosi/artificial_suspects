<script lang="ts">
    import { onMount } from 'svelte';
    import GamePage from './GamePage.svelte'
    import HomePage from './HomePage.svelte'
    import ConfigPage from './ConfigPage.svelte'
    import { GetGame, NewGame } from '../wailsjs/go/main/App.js';
    import { database } from '../wailsjs/go/models';
    import { register, init} from 'svelte-i18n';
    import { ToggleFullscreen } from '../wailsjs/go/main/App';

    register('en', () => import('./assets/locales/en.json'));
    register('cz', () => import('./assets/locales/cz.json'));
    register('pl', () => import('./assets/locales/pl.json'));
    init({
        fallbackLocale: 'en',
        initialLocale: 'en'
    });

    let screen = 'home'; // State to track the current screen
    let game: database.Game;

    async function handleKeyDown (event: KeyboardEvent) {
        // HOME
        if ((event.ctrlKey || event.metaKey) && event.shiftKey && (event.key === 'h' || event.key === 'm')) {
            console.log('Ctrl/Cmd+Shift+H/M was pressed');
            event.preventDefault();
            screen = 'home';
        }
        if (event.key === 'Escape') {
            console.log('Escape has been pressed');
            event.preventDefault();
            screen = 'home';
        }

        // GAME
        if ((event.ctrlKey || event.metaKey) && event.shiftKey && event.key === 'g') {
            console.log('Ctrl/Cmd+Shift+G was pressed');
            event.preventDefault();
            game = await GetGame();
            screen = 'game';
        }

        // CONFIG
        if ((event.ctrlKey || event.metaKey) && event.shiftKey && event.key === 'o') {
            console.log('Ctrl/Cmd+Shift+O was pressed');
            event.preventDefault();
            screen = 'config';
        }

        // FULLSCREEN
        if ((event.ctrlKey || event.metaKey) && event.shiftKey && event.key === 'f') {
            console.log('Ctrl/Cmd+Shift+F was pressed');
            event.preventDefault();
            ToggleFullscreen();
        }
    };

    onMount(() => {
        window.addEventListener('keydown', handleKeyDown);
        return () => {window.removeEventListener('keydown', handleKeyDown);};
    });

    async function newGameHandler(event) {
        try {
            game = await NewGame();
        } catch (error) {
            console.log(`NewGame() has failed: ${error}`)
        }
        console.log(game)
        screen = 'game';
        return
    }

    async function handleMessage(event) {
        console.log("handleMessage:", event)
        const { message } = event.detail;
        if (message === 'continueGame') {
            try {
                game = await GetGame();
            } catch (error) {
                console.log(`GetGame() has failed: ${error}`)
            }
            console.log(game)
            screen = 'game';
            return
        }

        if (message === 'goToHome') {
            screen = 'home';
            return
        }
         
        if (message === 'goToConfig') {
            screen = 'config';
            return
        }
    }
</script>

<main>
    {#if screen === 'home'}
        <HomePage on:message={handleMessage} on:newGame={newGameHandler}/>
    {:else if screen === 'game'}
        <GamePage on:message={handleMessage} on:newGame={newGameHandler} {game}/>
    {:else if screen === 'config'}
        <ConfigPage on:message={handleMessage}/>
    {/if}
</main>

<style>
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
