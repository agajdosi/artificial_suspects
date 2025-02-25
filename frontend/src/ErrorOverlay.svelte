<script lang="ts">
    
</script>

<div class="overlay">
    <div class="content">
        <h1>AI is not accessible.</h1>
        <p>Please check that the Ollama service is running on your machine and selected model is available.</p>
        <div class="options">
            <button>Check Again</button>
            <button>Configuration</button>
        </div>
    </div>
</div>

<style>
.overlay{
    position: fixed;
    top: 0;
    left: 0;
    height: 100%;
    width: 100vw;
    background-color: rgba(0, 0, 0, 0.7);
    backdrop-filter: blur(5px);
}

.content {
    margin: 15vh auto 0 auto;
}
</style>