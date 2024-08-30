<script lang="ts">
  import Suspects from './Suspects.svelte'
  import Home from './Home.svelte'

  let currentScreen = 'suspects'; // State to track the current screen

  function handleMessage(event) {
    console.log(event)
    const { message } = event.detail;
    if (message === 'goToSuspects') {
      currentScreen = 'suspects';
    } else if (message === 'goToHome') {
      currentScreen = 'home';
    }
  }
</script>

<main>
  {#if currentScreen === 'home'}
    <Home on:message={handleMessage} />
  {:else if currentScreen === 'suspects'}
    <Suspects on:message={handleMessage} />
  {/if}
</main>
