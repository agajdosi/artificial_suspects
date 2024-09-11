<script lang="ts">
    import GamePage from './GamePage.svelte'
    import HomePage from './HomePage.svelte'
    import ConfigPage from './ConfigPage.svelte'
    import { GetGame, NewGame } from '../wailsjs/go/main/App.js';
    import { main } from '../wailsjs/go/models';

    let screen = 'home'; // State to track the current screen
    let game: main.Game;

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
