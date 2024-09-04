<script lang="ts">
    import GamePage from './GamePage.svelte'
    import HomePage from './HomePage.svelte'
    import { GetGame, NewGame } from '../wailsjs/go/main/App.js';
    import { main } from '../wailsjs/go/models';

    let currentScreen = 'home'; // State to track the current screen
    let game: main.Game;

    async function handleMessage(event) {
        console.log(event)
        const { message } = event.detail;
        if (message === 'newGame') {
            try {
                game = await NewGame();
            } catch (error) {
                console.log(`NewGame() has failed: ${error}`)
            }
            console.log(game)
            currentScreen = 'game';
            return
        } else if (message === 'continueGame') {
            try {
                game = await GetGame();
            } catch (error) {
                console.log(`GetGame() has failed: ${error}`)
            }
            console.log(game)
            currentScreen = 'game'
            return
        } else if (message === 'goToHome') {
            currentScreen = 'home';
            return
        }
    }
</script>

<main>
    {#if currentScreen === 'home'}
        <HomePage on:message={handleMessage} />
    {:else if currentScreen === 'game'}
        <GamePage on:message={handleMessage} {game}/>
    {/if}
</main>
