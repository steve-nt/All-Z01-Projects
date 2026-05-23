## Carbon copy

### Video

Check out [this little video](https://youtu.be/ym-Qiio2ktI)!

### Instructions

Today is a big day: you're going to make your own webpage. Like a boss, yes.
If that one already sounds scary to you, don't worry, we've got your back ; you'll be provided some assets to complete this task.

The goal of this raid is to make a website that would look like this this webpage template.
![](https://github.com/01-edu/public/blob/master/subjects/carbon-copy/page-template.jpg)
You must customize the content of the site according to the theme that your group chose, trying to keep the structure of the placeholder version.  


The raid is divided in 3 phases:

- pure HTML structure
- custom CSS style
- JS interactions

> The code editor is not available for this Raid, so you'll learn to use an **IDE** (Integrated Development Environnement). We'll walk you through the basics.
 
To complete this project, participants will need to have Visual Studio Code installed. If it's not already installed, you can [Download IDE Visual Studio here](https://code.visualstudio.com/Download)

Next, download the zip file containing the three basic files for any website: an index.html, a styles.css file, and a script.js. You can download the zip file [here](https://we.tl/t-hbh3s9zRAf).

Once downloaded, you can put them all in a folder on your laptop and then open it using the "Open folder with Visual Studio" option.

Additionally, ensure you have the "Live Share" extension installed in VSC to collaborate simultaneously on the same project. You can install the Live Share extension from [this page](https://marketplace.visualstudio.com/items?itemName=MS-vsliveshare.vsliveshare).

To help you out, you can download [`carbon-copy.zip`](https://assets.01-edu.org/carbon-copy) to have at your disposal the following files:

- the CSS file `styles.css` containing the pre-styled elements
- the `images` folder containing the images to display in the webpage
- the `assets` folder containing the templates & wireframes images of the webpage

Coaches can help you turn in your work.


### Phase 1: HTML only _(mandatory)_

Create & write the HTML file `index.html` to build the structure of the page.
For this phase, we provide you a CSS file (`styles.css`) that you just have to [link](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/link) to your `index.html`, meaning you won't need to style anything. However, you can do it from scratch and write yourself the whole CSS if you feel up to it.

NB: Each item of the navbar menu - in the top right corner - has to bring to the corresponding section on click (clicking on "About" scrolls to the "About" section in the page); take a look at the [`href` anchor use](https://www.w3.org/TR/html401/struct/links.html#h-12.2.3).

Here is a [wireframe](https://en.wikipedia.org/wiki/Website_wireframe) of the webpage, showing the HTML tags you have to use:
![](https://raw.githubusercontent.com/01-edu/public/master/subjects/carbon-copy/page-wireframe.jpg)

### Phase 2: custom CSS _(mandatory)_

First of all, let's customize the color atmosphere of the webpage to your own taste: go to the CSS file `styles.css`, & replace the current blue & yellow with 2 new colors of your choice.

Now that the page is mainly built, you have to populate the "Dashboard" section with 3 new elements.
![](https://raw.githubusercontent.com/01-edu/public/master/subjects/carbon-copy/dashboard-template.jpg)

Those 3 cards have to display respectively one information with:

- a title
- a subtitle
- a text paragraph

For this phase, you'll have to make the whole HTML & CSS by yourself.

### Phase 3: JS interactions _(mandatory & optional)_

If you made it until here pretty fast, now the fun will begin! You're going to add a bunch of JS interactions to make elements appear / disappear / change in the HTML & CSS by linking a JS [script](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/script) to your `index.html`.

- Change the order of the pictures when clicking on the pictures' section, [toggling](https://css-tricks.com/snippets/javascript/the-classlist-api/) the pre-defined class `row-reverse` from the CSS file `styles.css`.

![](https://raw.githubusercontent.com/01-edu/public/master/subjects/carbon-copy/images-order.gif)

- Option 1: In the Contact section, when clicking on the "Introduce yourself" button, get the text typed in the `input` and display it in the middle of the following sentence: "Nice to meet you _[put here the input data]_ 👋! Thanks for introducing yourself." Also, the `<p>`, `<input>` & `<button>` elements have to disappear after the button has been clicked.

![](https://raw.githubusercontent.com/01-edu/public/master/subjects/carbon-copy/contact-input.gif)

- Option 2: When clicking on a card, open a modal window that will show the whole article ; the modal will be closed either when clicking on a "Close" button, or when the "Escape" key is pressed.

![](https://raw.githubusercontent.com/01-edu/public/master/subjects/carbon-copy/modale.gif)

- Option 3: In the modal article, create a widget that allows to change the text alignment ; on click on `left` or `center` buttons, the layout changes to the chosen justification, and the selected option's `font-weight` becomes `bold` whereas the other becomes `light`.

![](https://raw.githubusercontent.com/01-edu/public/master/subjects/carbon-copy/text-alignment.gif)

- Warrior option: set the `header` text content with a random quote every time the page is loaded, and then every 10 seconds. You can use this marvelous [Chuck Norris API](https://api.chucknorris.io/) to [fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch) his most inspiring sayings and display them in your own page!

![](https://raw.githubusercontent.com/01-edu/public/master/subjects/carbon-copy/fetch-quote.gif)

### Images

To add images to your site you can add the url of the image in an appropriate tag.  
Example:  

```HTML
<img src='https://i.postimg.cc/ygq96fdH/lou-batier-5-Eo-WFa-Htdo-unsplash.jpg' border='0' alt='lou-batier-5-Eo-WFa-Htdo-unsplash'/>
```  

Here are the urls of the images of the template.  
- [https://postimg.cc/34DBn4yK](https://i.postimg.cc/Kz0Hctxg/fabian-kozdon-5-Zeoo-CGNw3s-unsplash.jpg)
- [https://postimg.cc/3WXjNHmb](https://i.postimg.cc/wxPwWqvH/levi-midnight-DApw8e-Rf-R8-unsplash.jpg)
- [https://postimg.cc/ygq96fdH](https://i.postimg.cc/qRqGLSHz/lou-batier-5-Eo-WFa-Htdo-unsplash.jpg)

### Notions

- [HTML tags](https://developer.mozilla.org/en-US/docs/Web/HTML/Element)
- [CSS basics](https://developer.mozilla.org/en-US/docs/Learn/Getting_started_with_the_web/CSS_basics)
- [Css flexbox layout](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Flexible_Box_Layout/Basic_Concepts_of_Flexbox) & [Flexbox froggy](https://flexboxfroggy.com/) are always useful
- [Link a JS script](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/script)
- [`classList` / `toggle`](https://css-tricks.com/snippets/javascript/the-classlist-api/)
- [`getElementById`](https://developer.mozilla.org/en-US/docs/Web/API/Document/getElementById) or [`querySelector`](https://developer.mozilla.org/en-US/docs/Web/API/Element/querySelector)
- [`textContent`](https://developer.mozilla.org/en-US/docs/Web/API/Node/textContent)
- [`style`](https://developer.mozilla.org/en-US/docs/Web/API/ElementCSSInlineStyle/style)
- [`addEventListener`](https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener): [`click` event](https://developer.mozilla.org/en-US/docs/Web/API/Element/click_event) / [`keydown` event](https://developer.mozilla.org/en-US/docs/Web/API/Element/keydown_event) / [`load` event](https://developer.mozilla.org/en-US/docs/Web/API/Window/load_event)
- [`fetch`](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch)
- [`setInterval`](https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/setInterval)
