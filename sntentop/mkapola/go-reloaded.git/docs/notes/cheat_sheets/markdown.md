# 🧾 Markdown Basics Cheat Sheet

Markdown is a simple way to format text using symbols — perfect for notes, READMEs, and documentation!

---

🧠 Quick Summary

| What You Want | What to Type               |     |     |   |
| ------------- | -------------------------- | --- | --- | - |
| Title         | `#`                        |     |     |   |
| Bold          | `**bold**`                 |     |     |   |
| Italic        | `*italic*`                 |     |     |   |
| List          | `- item` or `1. item`      |     |     |   |
| Quote         | `> text`                   |     |     |   |
| Code          | `` `code` `` or code block |     |     |   |
| Link          | `[text](url)`              |     |     |   |
| Image         | `![alt](image.png)`        |     |     |   |
| Line          | `---`                      |     |     |   |
| Table         | `                          | col | col | ` |
| Checkbox      | `- [ ]`                    |     |     |   |
| Escape symbol | `\*text\*`                 |     |     |   |

💡 Tip

Markdown works best on:

GitHub

VS Code

Notion

Discord

and many note-taking apps!

You can mix text and code easily, and everything stays neat, light, and readable.

## 🏷️ 1. Headings (Titles)

Use `#` at the start of a line to make titles.


👉 Example:

# Heading 1  
## Heading 2  
### Heading 3

---

## 💪 2. Emphasis (Bold, Italic, Bold+Italic)


👉 Example:  
*italic* **bold** ***bold and italic***

---

## 📋 3. Lists

### Unordered List (• bullets)
Use `-`, `*`, or `+` followed by a space.


👉 Example:
- Apple  
- Banana  
- Cherry

### Ordered List (1, 2, 3)
Use numbers with a dot.


👉 Example:
1. First  
2. Second  
3. Third

You can mix lists too


---

## 💬 4. Blockquotes (Quotes or Notes)

Use `>` at the beginning of a line.


👉 Example:
> This is a quote or note.

You can also nest quotes:

---

## 💻 5. Code

### Inline code (inside text)
Use one backtick `` ` ``.


👉 Example:  
Use `go run main.go` to start the program.

### Code block (for larger code)
Use three backticks (\`\`\`) before and after.

<pre>

package main

import "fmt"

func main() {
fmt.Println("Hello World")
}

</pre>

👉 Example:

package main

import "fmt"

func main() {
fmt.Println("Hello World")
}


You can also specify a language for syntax coloring:
<pre>
```go
fmt.Println("Hello World")

## Links 
[Link text](https://example.com)

## Images 

![Alt text](image.png)

## Horizontal Line

---
***
___

## Tables

| Name  | Age | Country |
|-------|-----|----------|
| Anna  | 25  | Greece   |
| John  | 30  | USA      |


📎 10. Inline HTML (optional)

<b>Bold</b> <i>Italic</i> <br> (I didn't get this, search other time)

🗒️ 11. Task Lists (Checkboxes)

- [x] Done task
- [ ] Not done yet

📂 12. Escaping Characters

If you want to show symbols instead of using them (like * or _), use a backslash \.
\*This will not be italic\*

👉 Example:
*This will not be italic*

📘 13. Line Breaks

To make a line break, end a line with two spaces or a blank line.
First line.  
Second line. (didn't get this also)


