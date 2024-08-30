<script lang="ts">
  import GamePage from './GamePage.svelte'
  import HomePage from './HomePage.svelte'

  let currentScreen = 'game'; // State to track the current screen

  function handleMessage(event) {
    console.log(event)
    const { message } = event.detail;
    if (message === 'goToGame') {
      currentScreen = 'game';
    } else if (message === 'goToHome') {
      currentScreen = 'home';
    }
  }
</script>

<main>
  {#if currentScreen === 'home'}
    <HomePage on:message={handleMessage} />
  {:else if currentScreen === 'game'}
    <GamePage on:message={handleMessage} />
  {/if}
</main>
