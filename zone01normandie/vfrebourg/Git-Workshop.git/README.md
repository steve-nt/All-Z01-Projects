# Projet fictif : Site professionnel
## Introduction
Ce projet est un modèle de page Web pour un site pro que vous pouvez utiliser 
comme base pour votre propre projet. Il contient plusieurs sections "types" que l'on retrouve
généralement sur ce genre de site. 

Le projet fournit est un projet qui n'est pas fini, il contient
les bases du site pro et n'a pas toutes les parties nécessaires à un bon site professionnel
(services/produits vendues, formulaire de contact, avis, etc.). L'idée étant de vous appuyer sur 
l'existant pour ajouter ces parties manquantes.

## Objectifs
L'objectif de cet atelier est de fournir un projet fictif sur lequel vous pourrez travailler les notions de Git en collaboration, telles que les branches, les merge et les conflits. <br>

Vous serez invités à ajouter de nouvelles fonctionnalités (ici, ce sera des sections) à la plateforme dans des branches distinctes, à les fusionner avec la branche principale et à gérer les conflits qui pourraient survenir.

## Instructions
Voici les instructions pour utiliser ce projet dans votre propre dépôt Git :

<ol>
<li>Créer un dossier en local et aller dedans.</li>
<li>Récuperer ce projet via un git pull sur votre ordinateur :</li>

```git
git init
git pull https://zone01normandie.org/git/vfrebourg/Git-Workshop
```

<li>Créer un nouveau dépôt sur votre Gitea.</li>
<li>Ajouter les membres de votre équipe en tant que collaborateur sur Gitea.</li>
<li>Pousser le code du projet dans votre nouveau dépôt :</li>

```git
git remote add origin https://zone01normandie.org/git/votre-pseudo/repository-name
git push -u origin master
```

<li>Demandez à chaque membre de votre groupe de récupérer le dépôt et de créer une nouvelle branche chacun pour ajouter une des fonctionnalités de la liste fournie ci-dessous.</li>

<li>Une fois que chaque membre a terminé, vous devrez ajouter votre fonctionnalité à la branche principale (master ou main), <b>sans en récupérer le contenu auparavant</b>, pour justement avoir des conflits.

<li>Une fois que vous aurez ces conflits, les membres peuvent travailler ensemble pour les résoudre et fusionner leurs modifications.</li>

</ol>

## Fonctionnalités à ajouter

<i>NB: Les fonctionnalités sont plus ou moins identiques pour une bonne raison, l'idée est que vous travaillez sur le même fichier/la même partie pour justement avoir des conflits à régler.</i>

De manière générale, vous pouvez ajouter des "sections" à ce site pour créer 
les parties importantes d'un site de ce style : une partie contact, une partie 
avec des avis clients, une partie sur les produits, etc.

Si vous souhaitez faire vos propres sections, vous pouvez prendre ce bout de code qui permet de 
faire une section.

```html
<section class="u-align-center u-clearfix u-container-align-center u-grey-5 u-section-9" id="carousel_1208">
    <div class="u-clearfix u-sheet u-valign-middle u-sheet-1">
        <!-- Votre code ici -->
    </div>
</section>
```

Concernant l'emplacement de ces modifications, vous aurez un commentaire dans le fichier index.html vous indiquant où vous pouvez les mettre.

```html
...
    </section>
    <!-- Rajoutez vos sections avant ou après les sections ci-dessus -->
    
    <footer class="u-align-center u-clearfix u-footer u-grey-80 u-footer" id="sec-008b">
...
```

Pas de panique, le but de cet atelier n'étant pas la création web en soit mais plus la collaboration sur un projet,
vous avez déjà des sections préparées ci-dessous que vous pouvez copier-coller:

<details>
<summary>Liste des services/produits </summary>
<li>Voir le code dans le fichier <b>templates-part/services.html</b></li>
</details>

<details>
<summary>Inscription Newsletter</summary>
<li>Voir le code dans le fichier <b>templates-part/newsletter.html</b></li>
</details>

<details>
<summary>Un Call-to-Action </summary>
<li>Voir le code dans le fichier <b>templates-part/cta.html</b></li>
</details>

<details>
<summary>Une partie "Chiffres" </summary>
<li>Voir le code dans le fichier <b>templates-part/results.html</b></li>
</details>

<details>
<summary>Un formulaire de contact</summary>
<li>Voir le code dans le fichier <b>templates-part/contact.html</b></li>
</details>

<details>
<summary>Une partie "retour clients"</summary>
<li>Voir le code dans le fichier <b>templates-part/testimonials.html</b></li>
</details>

<details>
<summary>Une bannière "type" pour présenter</summary>
<li>Voir le code dans le fichier <b>templates-part/whatwedo.html</b></li>
</details>


## Notions
Cela vous permettra de voir et d'apprendre les concepts suivants :

* [Création de branches.](https://git-scm.com/docs/git-branch)
* [Fusion de branches.](https://git-scm.com/docs/git-merge) 
* [Résolution de conflits.](https://git-scm.com/docs/git-merge#_how_to_resolve_conflicts)