<script lang="ts">
  import logo from './assets/images/logo-universal.png'
  import Suspects from './Suspects.svelte'

  import {Greet} from '../wailsjs/go/main/App.js'

  let resultText: string = "Please enter your name below 👇"
  let name: string

  function greet(): void {
    Greet(name).then(result => resultText = result)
  }
</script>

<main>
  <Suspects />
</main>

<style>


</style>
