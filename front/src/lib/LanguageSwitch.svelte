<script lang="ts">
    import {locale} from 'svelte-i18n';

    function changeLanguage(code: string) {
        $locale = code;
    }
</script>

<button on:click={() => changeLanguage('en')} class="langbtn" class:active={$locale === 'en'}>en</button>
<button on:click={() => changeLanguage('cz')} class="langbtn" class:active={$locale === 'cz'}>cz</button>
<button on:click={() => changeLanguage('pl')} class="langbtn" class:active={$locale === 'pl'}>pl</button>

<style>
.langbtn {
    all: unset;
    text-decoration: underline;
    min-width: 20px;
}
.langbtn:hover{
    cursor: pointer;
}
.langbtn.active {
    text-transform: uppercase;
}
</style>