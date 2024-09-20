<script lang="ts">
    import { onMount } from 'svelte';
    import GamePage from './GamePage.svelte'
    import HomePage from './HomePage.svelte'
    import ConfigPage from './ConfigPage.svelte'
    import { GetGame, NewGame } from '../wailsjs/go/main/App.js';
    import { database } from '../wailsjs/go/models';

    let screen = 'home'; // State to track the current screen
    let game: database.Game;

    function handleKeyDown (event: KeyboardEvent) {
        if ((event.ctrlKey || event.metaKey) && event.shiftKey && event.key === 'm') {
            console.log('Ctrl+Shift+M or Cmd+Shift+M was pressed');
            event.preventDefault();
            screen = 'home';
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

    async function enterConfigHandler(event) {
        screen = "config";
    }

    async function handleMessage(event) {
        console.log(event)
        const { message } = event.detail;
        if (message === 'continueGame') {
            try {
                game = await GetGame();
            } catch (error) {
                console.log(`GetGame() has failed: ${error}`)
            }
            console.log(game)
            screen = 'game'
            return
        } else if (message === 'goToHome') {
            screen = 'home';
            return
        }
    }
</script>

<main>
    {#if screen === 'home'}
        <HomePage on:message={handleMessage} on:newGame={newGameHandler} on:enterConfig={enterConfigHandler}/>
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
</style>
