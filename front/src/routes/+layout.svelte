<script lang="ts">

</script>

<div>
    <slot />
</div>

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
    url("static/nunito-v16-latin-regular.woff2") format("woff2");
}

:global(#app) {
    height: 100vh;
    text-align: center;
}

</style>