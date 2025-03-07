<script lang="ts">
  import { onMount } from "svelte";
  import Typed from "typed.js";

  let typedInstance: Typed | null = null;
  let showOverlay = true;

  let tutorialSteps = [
    "Welcome, Investigator! <br>Your mission is to find the guilty criminal among 15 suspects.",
    "Each round, the Witness will answer a question about the criminal.",
    "You must eliminate suspects who do NOT match the answer.",
    "For example, if the question is 'Does the suspect like reading books?' and I say 'Yes', remove those who likely don’t.",
    "Continue eliminating suspects round by round until only one remains.",
    "If the last suspect standing is the correct one, you win!<br>Now, let’s begin…"
  ];

  onMount(() => {
    typedInstance = new Typed("#typed", {
      strings: tutorialSteps,
      typeSpeed: 60,
      backSpeed: 10,
      fadeOut: true,
      loop: false,
      cursorChar: "|",
    });
  });

  function closeOverlay() {
    showOverlay = false;
    if (typedInstance) typedInstance.destroy(); // Stop the typing animation when closing
  }
</script>

{#if showOverlay}
  <div class="overlay">
    <p id="content">
      <span id="typed"></span>
    </p>
    <button on:click={closeOverlay}>Let's Play</button>
  </div>
{/if}

<style>
  .overlay {
    position: absolute;
    z-index: 100;
    top: 0;
    left: 0;
    height: 100vh;
    width: 100vw;
    background-color: rgba(27, 38, 54, 1);
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    color: white;
    text-align: center;
    font-size: 1.5rem;
  }

  #content {
    min-height: 20rem;
  }

</style>
