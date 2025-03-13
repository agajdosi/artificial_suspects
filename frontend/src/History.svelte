<script lang="ts">
    import { hint } from './lib/stores';
    import { GetActiveService } from './lib/main';
    import type { Game, Service } from './lib/main';
    import { locale, t } from 'svelte-i18n';

    export let game: Game;
    let activeService: Service | null = null;

    GetActiveService().then(service => {
        activeService = service;
    }).catch(err => {
        console.error('Error fetching active service:', err);
    });
</script>

<div class="history"
    on:mouseenter={() => hint.set("History of previous questions and their answers in current investigation.")}
    on:mouseleave={() => hint.set("")}
>
    {#each [...game.investigation.rounds].reverse().slice(1).reverse() as round, index}
        <div class="round">
            <div class="question">
                {index+1}.
                {#if $locale == "cz"}
                    {round.Question.Czech}
                {:else if $locale == "pl"}
                    {round.Question.Polish}
                {:else}
                    {round.Question.English}
                {/if}
            </div>
            <div class="answer">
                {$t(round.answer.toLocaleLowerCase())}!
            </div>
        </div>
    {/each}
</div>

<div class="roles">
    <div class="model"
        on:mouseenter={() => hint.set("An AI model that acts as a witness and responds to questions. You can change it in the game configuration.")}
        on:mouseleave={() => hint.set("")}
        >
        {$t("interrogated")}:
        {#if activeService === null}{$t("loading")}
        {:else} {activeService.VisualModel}
        {/if}
    </div>
</div>

<style>
.history {
    display: flex;
    flex-direction: column-reverse;
}

.roles {
    display: flex;
    justify-content: space-between;
    padding: 0 0 2rem 0;
}

.round {
    display: flex;
}

.question, .answer {
    padding: 10px;
    border-radius: 10px;
    margin: 5px 0;
    position: relative;
    font-size: 16px;
    width: fit-content;
    max-width: 100%;
}

.question {
    background-color: #343563;
    align-self: flex-start;
    border-bottom-left-radius: 0;
}

.answer {
    background-color: #3c1c54;
    align-self: flex-end;
    border-bottom-right-radius: 0;
    text-transform: capitalize;
}
</style>
