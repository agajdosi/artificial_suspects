// Copyright (C) 2024 (Andreas Gajdosik) <andreas@gajdosik.org>
// This file is part of project.
//
// project is non-violent software: you can use, redistribute,
// and/or modify it under the terms of the CNPLv7+ as found
// in the LICENSE file in the source code root directory or
// at <https://git.pixie.town/thufie/npl-builder>.
//
// project comes with ABSOLUTELY NO WARRANTY, to the extent
// permitted by applicable law. See the CNPL for details.

package database

// Prepacked default questions for the game. This gets filled into
// the questions database.
var defaultQuestions = []Question{
	{
		English: "Does the suspect like pizza?",
		Czech:   "Má podezřelý rád pizzu?",
		Polish:  "Czy podejrzany lubi pizzę?",
		Topic:   "basic", Level: 1,
	},
	{
		English: "Is the suspect leftist?",
		Czech:   "Je podezřelý levičák?",
		Polish:  "Czy podejrzany jest lewicowy?",
		Topic:   "political", Level: 1,
	},
	{
		English: "Does the suspect have depressions?",
		Czech:   "Má podezřelý deprese?",
		Polish:  "Czy podejrzany ma depresje?",
		Topic:   "psychological", Level: 1,
	},
	{
		English: "Is the suspect a fan of social media?",
		Czech:   "Je podezřelý fanouškem sociálních sítí?",
		Polish:  "Czy podejrzany jest fanem mediów społecznościowych?",
		Topic:   "sociological", Level: 1,
	},
	{
		English: "Does the suspect enjoy traveling?",
		Czech:   "Má podezřelý rád cestování?",
		Polish:  "Czy podejrzany lubi podróże?",
		Topic:   "basic", Level: 1,
	},
	{
		English: "Is the suspect environmentally conscious?",
		Czech:   "Je podezřelý ohleduplný k životnímu prostředí?",
		Polish:  "Czy podejrzany ma świadomość ekologiczną?",
		Topic:   "political", Level: 1,
	},
	{
		English: "Does the suspect attend therapy?",
		Czech:   "Navštěvuje podezřelý terapii?",
		Polish:  "Czy podejrzany chodzi na terapię?",
		Topic:   "psychological", Level: 1},
	{
		English: "Does the suspect believe in traditional family?",
		Czech:   "Vyznává podezřelý tradiční rodinu?",
		Polish:  "Czy podejrzany wierzy w tradycyjną rodzinę?",
		Topic:   "sociological", Level: 1},
	{
		English: "Is the suspect vegetarian?",
		Czech:   "Je podezřelý vegetarián?",
		Polish:  "Czy podejrzany jest wegetarianinem?",
		Topic:   "basic", Level: 1},
	{
		English: "Is the suspect vegan?",
		Czech:   "Je podezřelý vegan?",
		Polish:  "Czy podejrzany jest weganinem?",
		Topic:   "basic", Level: 1},
	{
		English: "Does the suspect vote regularly?",
		Czech:   "Chodí podezřelý pravidelně k volbám?",
		Polish:  "Czy podejrzany regularnie głosuje?",
		Topic:   "political", Level: 1},
	{
		English: "Does the suspect struggle with anxiety?",
		Czech:   "Má podezřelý problémy s úzkostmi?",
		Polish:  "Czy podejrzany zmaga się z lękiem?",
		Topic:   "psychological", Level: 1},
	{
		English: "Is the suspect an extrovert?",
		Czech:   "Je podezřelý extrovert?",
		Polish:  "Czy podejrzany jest ekstrawertykiem?",
		Topic:   "sociological", Level: 1},
	{
		English: "Does the suspect sport regularly?",
		Czech:   "Sportuje podezřelý pravidelně?",
		Polish:  "Czy podejrzany regularnie uprawia sport?",
		Topic:   "basic", Level: 1},
	{
		English: "Does the suspect have strong political opinions?",
		Czech:   "Má podezřelý vyhraněné politické názory?",
		Polish:  "Czy podejrzany ma silne poglądy polityczne?",
		Topic:   "political", Level: 1},
	{
		English: "Does the suspect meditate?",
		Czech:   "Medituje podezřelý?",
		Polish:  "Czy podejrzany medytuje?",
		Topic:   "psychological", Level: 1},
	{
		English: "Is the suspect part of a secret community?",
		Czech:   "Je podezřelý členem tajné komunity?",
		Polish:  "Czy podejrzany jest częścią tajnej społeczności?",
		Topic:   "sociological", Level: 1},
	{
		English: "Does the suspect enjoy cooking?",
		Czech:   "Je podezřelý členem uzavřené komunity?",
		Polish:  "Czy podejrzany lubi gotować?",
		Topic:   "basic", Level: 1},
	{
		English: "Is the suspect involved in activism?",
		Czech:   "Je podezřelý zapojen do aktivismu?",
		Polish:  "Czy podejrzany jest zaangażowany w aktywizm?",
		Topic:   "political", Level: 1},
	{
		English: "Does the suspect have mood swings?",
		Czech:   "Má podezřelý výkyvy nálad?",
		Polish:  "Czy podejrzany miewa wahania nastroju?",
		Topic:   "psychological", Level: 1},
	{
		English: "Does the suspect follow trends?",
		Czech:   "Řídí se podezřelý trendy?",
		Polish:  "Czy podejrzany podąża za trendami?",
		Topic:   "sociological", Level: 1},
	{
		English: "Is the suspect a fan of sci-fi movies?",
		Czech:   "Je podezřelý fanouškem sci-fi filmů?",
		Polish:  "Czy podejrzany jest fanem filmów science-fiction?",
		Topic:   "basic", Level: 1},
	{
		English: "Does the suspect lean towards conservatism?",
		Czech:   "Přiklání se podezřelý ke konzervatismu?",
		Polish:  "Czy podejrzany skłania się ku konserwatyzmowi?",
		Topic:   "political", Level: 1},
	{
		English: "Does the suspect enjoy large social gatherings?",
		Czech:   "Má podezřelý rád velká společenská setkání?",
		Polish:  "Czy podejrzany lubi duże spotkania towarzyskie?",
		Topic:   "sociological", Level: 1},
	{
		English: "Does the suspect enjoy hiking?",
		Czech:   "Má podezřelý rád pěší turistiku?",
		Polish:  "Czy podejrzany lubi piesze wędrówki?",
		Topic:   "basic", Level: 1},
	{
		English: "Does the suspect have progressive views?",
		Czech:   "Má podezřelý pokrokové názory?",
		Polish:  "Czy podejrzany ma postępowe poglądy?",
		Topic:   "political", Level: 1},
	{
		English: "Does the suspect think they are member of a minority?",
		Czech:   "Myslí si podezřelý, že je menšinou?",
		Polish:  "Czy podejrzany uważa się za członka mniejszości?",
		Topic:   "sociological", Level: 1},
	{
		English: "Does the suspect enjoy reading books?",
		Czech:   "Čte podezřelý rád knihy?",
		Polish:  "Czy podejrzany lubi czytać książki?",
		Topic:   "basic", Level: 1},
	{
		English: "Is the suspect politically active online?",
		Czech:   "Je podezřelý politicky aktivní na internetu?",
		Polish:  "Czy podejrzany jest aktywny politycznie w Internecie?",
		Topic:   "political", Level: 1},
	{
		English: "Does the suspect have low self-esteem?",
		Czech:   "Má podezřelý nízké sebevědomí?",
		Polish:  "Czy podejrzany ma niską samoocenę?",
		Topic:   "psychological", Level: 1},
	{
		English: "Does the suspect belong to a religious organization?",
		Czech:   "Patří podezřelý k náboženské organizaci?",
		Polish:  "Czy podejrzany należy do organizacji religijnej?",
		Topic:   "sociological", Level: 1},
	{
		English: "Does the suspect discuss politics frequently?",
		Czech:   "Diskutuje podezřelý často o politice?",
		Polish:  "Czy podejrzany często dyskutuje o polityce?",
		Topic:   "political", Level: 1},
	{
		English: "Is the suspect involved in charity work?",
		Czech:   "Podílí se podezřelý na charitativní činnosti?",
		Polish:  "Czy podejrzany jest zaangażowany w działalność charytatywną?",
		Topic:   "sociological", Level: 1},
	{
		English: "Is the suspect a pet owner?",
		Czech:   "Má podezřelý domácího mazlíčka?",
		Polish:  "Czy podejrzany jest właścicielem zwierzęcia?",
		Topic:   "basic", Level: 1},
	{
		English: "Does the suspect experience panic attacks?",
		Czech:   "Mívá podezřelý panické ataky?",
		Polish:  "",
		Topic:   "psychological", Level: 1},
	{
		English: "Is the suspect active in a subculture?",
		Czech:   "Je podezřelý aktivní v některé subkultuře?",
		Polish:  "Czy podejrzany doświadcza ataków paniki?",
		Topic:   "sociological", Level: 1},
	{
		English: "Does the suspect enjoy classical music?",
		Czech:   "Má podezřelý rád klasickou hudbu?",
		Polish:  "Czy podejrzany lubi muzykę klasyczną?",
		Topic:   "basic", Level: 1},
	{
		English: "Does the suspect practice positive affirmations?",
		Czech:   "Praktikuje podezřelý pozitivní afirmace?",
		Polish:  "Czy podejrzany praktykuje pozytywne afirmacje?",
		Topic:   "psychological", Level: 1},
	{
		English: "Does the suspect have many friends?",
		Czech:   "Má podezřelý hodně přátel?",
		Polish:  "Czy podejrzany ma wielu przyjaciół?",
		Topic:   "sociological", Level: 1},
	{
		English: "Does the suspect enjoy fast food?",
		Czech:   "Má podezřelý rád rychlé občerstvení?",
		Polish:  "Czy podejrzany lubi fast foody?",
		Topic:   "basic", Level: 1},
	{
		English: "Does the suspect support LGBTQ+ rights?",
		Czech:   "Podporuje podezřelý práva LGBTQ+ lidí?",
		Polish:  "Czy podejrzany wspiera prawa osób LGBTQ+?",
		Topic:   "political", Level: 1},
	{
		English: "Does the suspect watch reality TV?",
		Czech:   "Sleduje podezřelý reality show?",
		Polish:  "Czy podejrzany ogląda reality TV?",
		Topic:   "basic", Level: 1},
	{
		English: "Does the suspect believe in socialism?",
		Czech:   "Je podezřelý zastáncem socialismu?",
		Polish:  "Czy podejrzany wierzy w socjalizm?",
		Topic:   "political", Level: 1},
	{
		English: "Does the suspect engage in self-care?",
		Czech:   "Pečuje podezřelý o sebe?",
		Polish:  "Czy podejrzany angażuje się w samoopiekę?",
		Topic:   "psychological", Level: 1},
	{
		English: "Is the suspect socially awkward?",
		Czech:   "Je podezřelý společensky neohrabaný?",
		Polish:  "Czy podejrzany jest niezręczny społecznie?",
		Topic:   "sociological", Level: 1},
	{
		English: "Does the suspect align with feminist ideals?",
		Czech:   "Je podezřelý v souladu s feministickými ideály?",
		Polish:  "Czy podejrzany jest zgodny z feministycznymi ideałami?",
		Topic:   "political", Level: 1},
	{
		English: "Does the suspect have trust issues?",
		Czech:   "Má podezřelý problémy s důvěrou?",
		Polish:  "Czy podejrzany ma problemy z zaufaniem?",
		Topic:   "psychological", Level: 1},
	{
		English: "Does the suspect drink alcohol?",
		Czech:   "Pije podezřelý alkohol?",
		Polish:  "Czy podejrzany pije alkohol?",
		Topic:   "basic", Level: 1},
	{
		English: "Does the suspect have anger issues?",
		Czech:   "Má podezřelý problémy se vztekem?",
		Polish:  "Czy podejrzany ma problemy z gniewem?",
		Topic:   "psychological", Level: 1},
	{
		English: "Does the suspect regularly attend social events?",
		Czech:   "Navštěvuje podezřelý pravidelně společenské akce?",
		Polish:  "Czy podejrzany regularnie uczestniczy w wydarzeniach towarzyskich?",
		Topic:   "sociological", Level: 1},
	{
		English: "Does the suspect enjoy gardening?",
		Czech:   "Pracuje podezřelý rád na zahradě?",
		Polish:  "Czy podejrzany lubi ogrodnictwo?",
		Topic:   "basic", Level: 1},
	{
		English: "Does the suspect have a fear of failure?",
		Czech:   "Má podezřelý strach ze selhání?",
		Polish:  "Czy podejrzany obawia się porażki?",
		Topic:   "psychological", Level: 1},
	{
		English: "Is the suspect a fan of horror movies?",
		Czech:   "Je podezřelý fanouškem hororů?",
		Polish:  "Czy podejrzany jest fanem horrorów?",
		Topic:   "basic", Level: 1},
	{
		English: "Does the suspect believe in capitalism?",
		Czech:   "Věří podezřelý v kapitalismus?",
		Polish:  "Czy podejrzany wierzy w kapitalizm?",
		Topic:   "political", Level: 1},
	{
		English: "Does the suspect enjoy fine dining?",
		Czech:   "Má podezřelý rád dobré jídlo?",
		Polish:  "Czy podejrzany lubi dobrze zjeść?",
		Topic:   "basic", Level: 1},
	{
		English: "Does the suspect support authoritarianism?",
		Czech:   "Podporuje podezřelý autoritářství?",
		Polish:  "Czy podejrzany sympatyzuje z autorytaryzmem?",
		Topic:   "political", Level: 1},
	{
		English: "Does the suspect feel isolated?",
		Czech:   "Cítí se podezřelý izolovaný?",
		Polish:  "Czy podejrzany czuje się odizolowany?",
		Topic:   "psychological", Level: 1},
	{
		English: "Does the suspect collect anything?",
		Czech:   "Sbírá podezřelý něco?",
		Polish:  "Czy podejrzany coś zbiera?",
		Topic:   "basic", Level: 1},
	{
		English: "Does the suspect support progressive taxation?",
		Czech:   "Podporuje podezřelý progresivní zdanění?",
		Polish:  "Czy podejrzany popiera progresywne opodatkowanie?",
		Topic:   "political", Level: 1},
	{
		English: "Does the suspect struggle with self-doubt?",
		Czech:   "Bojuje podezřelý s pochybnostmi o sobě samém?",
		Polish:  "Czy podejrzany zmaga się z wątpliwościami?",
		Topic:   "psychological", Level: 1},
	{
		English: "Does the suspect support universal basic income?",
		Czech:   "Podporuje podezřelý univerzální základní příjem?",
		Polish:  "",
		Topic:   "political", Level: 1},
	{
		English: "Does the suspect deal with imposter syndrome?",
		Czech:   "Trpí podezřelý syndromem podvodníka?",
		Polish:  "Czy podejrzany popiera uniwersalny dochód podstawowy?",
		Topic:   "psychological", Level: 1},
	{
		English: "Is the suspect well-connected in their neighborhood?",
		Czech:   "Má podezřelý ve svém okolí dobré kontakty?",
		Polish:  "Czy podejrzany ma dobre kontakty w swojej okolicy?",
		Topic:   "sociological", Level: 1},
	{
		English: "Does the suspect prefer cats?",
		Czech:   "Má podezřelý raději kočky?",
		Polish:  "Czy podejrzany preferuje koty?",
		Topic:   "basic", Level: 1},
	{
		English: "Does the suspect regularly journal?",
		Czech:   "Píše si podezřelý pravidelně deník?",
		Polish:  "Czy podejrzany regularnie prowadzi dziennik?",
		Topic:   "psychological", Level: 1},
	{
		English: "Is the suspect engaged in social justice movements?",
		Czech:   "Je podezřelý zapojen do hnutí za sociální spravedlnost?",
		Polish:  "Czy podejrzany jest zaangażowany w ruchy na rzecz sprawiedliwości społecznej?",
		Topic:   "sociological", Level: 1},
	{
		English: "Has the suspect ever tried drugs?",
		Czech:   "Zkusil podezřelý někdy drogy?",
		Polish:  "Czy podejrzany kiedykolwiek próbował narkotyków?",
		Topic:   "political", Level: 1},
	{
		English: "Does the suspect believe in global climate change?",
		Czech:   "Věří podezřelý v globální změnu klimatu?",
		Polish:  "Czy podejrzany wierzy w globalne zmiany klimatu?",
		Topic:   "ecology", Level: 1},
	{
		English: "Does the suspect like contemporary art?",
		Czech:   "Má rád podezřelý současné umění?",
		Polish:  "Czy podejrzany lubi sztukę współczesną?",
		Topic:   "art", Level: 1},
}
