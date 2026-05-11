
## Steps left TO-DO:

> 1. **Print to the terminal**
  > * Take input from http request 
  > * Take banner input 
  > * Combine 1 - 2  
  > - add ascii convertion

> 2. **Return the output as a requestWriter**

> 3. **input validation**

> 4. **Http error posting**

> 5. **loger**

> 7. **Error handling**











----

## How the Reference Works in Go templates

The reference isn't explicitly declared in the code. Instead, it happens because:

  >  All Templates Are Parsed Together\
        When you call template.ParseGlob("templates/*.html") in the init() function, Go loads and parses all templates in the templates directory (both base.html and home.html in this case).
        The parsed templates are stored in the *template.Template object (tmpl).

  > Named Blocks ({{define "content"}}) in home.html\
        In home.html, the {{define "content"}} block defines a "named template" called "content". This block is made available to any other template that uses it.

  > Placeholders in base.html\
        In base.html, the {{template "content" .}} directive serves as a placeholder for the "content" block.
       When you execute "home.html", Go automatically merges base.html (because it has the layout) with home.html (because it defines "content").

  > Execution via tmpl.ExecuteTemplate\
        When tmpl.ExecuteTemplate(w, "home.html", data) is called, Go resolves the dependencies between the templates:
        It sees that home.html includes "content", and it knows base.html expects "content".
        The "content" block defined in home.html is inserted into the {{template "content" .}} placeholder in base.html.

  * Go looks for the home.html template in the parsed templates.
  * It sees that home.html defines a content block via {{define "content"}}.
  * It checks if any other template (like base.html) uses {{template "content"}}.
  * If so, it renders base.html first, then injects the content block from home.html into the {{template "content" .}} placeholder.

  This design makes templates modular and reusable without explicitly linking them in the code!


        <!-- How base.html and anyPage.html Work Together
    base.html as the Layout Template
        base.html provides the overall structure for the page, such as the header, footer, and a placeholder ({{template "content" .}}) for the main content.
        The content placeholder is filled dynamically by templates like anyPage.html.

    Template Inheritance Using {{deaasfine "content"}}
        Go templates allow defining named blocks (e.g., content) that can be included into a parent template like base.html. -->
   

