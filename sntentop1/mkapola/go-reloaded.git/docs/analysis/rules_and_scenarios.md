# Goals

1. In this exercise, we will use **functions from the piscine repo**.
2. With these functions, we will create a **tool** that will **complete, process, and correct plain text**.
3. We will be **tested by other members of our cohort**, and each of us will test others — in other words, we will **audit each other**.
4. It is recommended to create **our own tests** to check both ourselves and those we will be auditing.

# Instructions and Rules

* The program must be written in **Go**.
* It must follow **good practices** [good practices](https://platform.zone01.gr/api/content/root/public/subjects/good-practices/README.md) ---> tip: give the link to ChatGPT and ask it to format it properly <3 Krystalenia <3
* It is recommended to have test files for [unit testing](https://go.dev/doc/tutorial/add-a-test)

---

### 1. **(hex)**

**RULE -->** Convert the **previous** word to its **decimal form**. In this case, the previous word will **always be a hexadecimal number**.

Examples

* *"**1E (hex)** files were added" -> "**30** files were added"*
* *"I am **0x23 (hex)** years old" -> "I am **35** years old"*
* *"Dunbar's number is **0x96**" -> "Dunbar's number is **150**"*

**Information about hexadecimal numbers:**

1. (base-16): uses 16 symbols →
   0–9 (for values 0–9) and A–F (for values 10–15).
2. The hexadecimal order is:

   * 0 1 2 3 4 5 6 7 8 9  A  B  C  D  E  F
   * 0 1 2 3 4 5 6 7 8 9 10  11 12 13 14 15
3. To indicate that a number is hexadecimal, we write **0x** in front of it, because sometimes it contains only numeric digits, so it’s impossible to tell if it’s decimal or hexadecimal.
   Example:
   96 -> decimal
   0x96 -> hexadecimal
4. Each hex digit = 4 binary bits

* 1 hex = 4 bits = “nibble”
* 2 hex = 1 byte (8 bits)

Example: A3₁₆ = 1010 0011₂

5. Computers often use hex as a compact way to read or write binary data.

**Conversion Method**

1. 🔢 Conversion Logic (Hex → Decimal)

Steps:
Write down the hex number.
Multiply each digit by 16ⁿ, counting positions from right (0).
Add up all results.

Example:
2F₁₆ = (2×16¹) + (F×16⁰)
= (2×16) + (15×1)
= 32 + 15 = 47₁₀

2. 💡 Quick Pattern

Moving one digit left → multiply by 16.
Add the next digit’s value.

Example:
1A₁₆
Start 0
→ (0×16)+1 = 1
→ (1×16)+10 = 26₁₀

---

### 2. **(bin)**

**RULE -->** Convert the **previous** word to its **decimal form**. In this case, the previous word will **always be a binary number**.

Examples

* *"It has been **10 (bin)** years" -> "It has been **2** years"*
* *"I am **100011 (bin)** years old" -> "I am **35** years old"*
* *"Dunbar's number is **10010110**" -> "Dunbar's number is **150**"*

**Information about binary numbers:**

1. It’s a system that uses only 0 and 1 (base-2).
2. Binary system: Base-2 / each position = a power of 2 (1, 2, 4, 8, 16, …), counting from right to left starting at 0.
3. Bit: single binary digit (0 or 1)
4. Byte: 8 bits
5. MSB = Most Significant Bit (leftmost)
6. LSB = Least Significant Bit (rightmost)
7. This is the system used by computers to store and process data.

**Conversion Method**

1. Basic

🔢 Conversion Logic (Binary → Decimal)

Step-by-step rule:
Write down the binary number.
Multiply each bit by 2ⁿ, counting positions from right to left (starting at 0).
Add all the results.

Example:
1011₂ = (1×2³) + (0×2²) + (1×2¹) + (1×2⁰)
= 8 + 0 + 2 + 1 = 11₁₀

2. Mental Shortcut

Each time you move one digit to the left → multiply the previous total by 2 and add the current bit.

Example:
Binary 1011:
Start at 0
→ (0×2)+1 = 1
→ (1×2)+0 = 2
→ (2×2)+1 = 5
→ (5×2)+1 = 11

---

### 3. **(up)**

**RULE -->** Convert the **previous** word to **UPPERCASE**.

Examples

* *"Ready, set, **go (up)** !" -> "Ready, set, **GO**!"*
* *"**stop (up)**" -> "**STOP**"*
* *"The new game of sucker punch productions is **amazing (up)** !" -> "The new game of sucker punch productions is **AMAZING**!"*

---

### 4. **(low)**

**RULE -->** Convert the **previous** word to **lowercase**.

Examples

* *"I should stop **SHOUTING (low)**" -> "I should stop **shouting**"*
* *"I need to **REMEMBER (low)** that..." -> "I need to **remember** that..."*
* *"It's so **RELAXING (low)** here" -> "It's so **relaxing** here"*

---

### 5. **(cap)**

**RULE -->** Change the **previous** word so that **its first letter is capitalized**.

Examples

* *"Welcome to the Brooklyn bridge (cap)" -> "Welcome to the Brooklyn Bridge"*
* *"The new game of **sucker (cap) punch (cap) productions (cap)** is AMAZING!" -> "The new game of **Sucker Punch Productions** is AMAZING!"*
* *"The name of **the (cap)** woman is **irene (cap) adler (cap)**" -> "The name of **The** woman is **Irene Adler**"*

**In cases of (low), (up), (cap) ONLY**, the **number of words to be affected** can be specified when a comma and a number are added next to the parenthesis, like this:
**(low, <number>)**

Examples

* *"This is **so exciting (up, 2)**" -> "This is **SO EXCITING**"*
* *"The new game of **sucker punch productions (cap, 3)** is AMAZING!" -> "The new game of **Sucker Punch Productions** is AMAZING!"*
* *"The name of the (cap) woman is **irene adler (cap, 2)**" -> "The name of The woman is **Irene Adler**"*

---

### 6. **PUNCTUATION MARKS: . | , | ! | ? | : | ; |**

**RULE -->** Attach each punctuation mark **directly to the previous word** and add **one space after it**, regardless of what follows (word, letter, symbol, or punctuation).

Examples

* *"I was sitting over there ,and then **BAMM !!**" -> "I was sitting over there, and then **BAMM!!**"*

---

### **-> EXCEPTION <-**

Groups of punctuation marks: **...** or **!?**

**RULE -->** In these two cases, move the entire punctuation group **directly after the previous word**, without any spaces between the word and the punctuation marks. Then, add **one space** after it before anything that follows.

Examples

* *"I was **thinking ...** You were right" -> "I was **thinking...** You were right"*

---

### 7. **QUOTES**

**RULE -->** Whenever there is a single quote `'`, look for the next one. Identify the word between them and remove any spaces **immediately after the first quote** and **immediately before the second**, so that the quotes **touch** the word(s) inside without spaces.

Examples

* *"I am exactly how they describe me: **' awesome '**" -> "I am exactly how they describe me: **'awesome'**"*
* *"My dog thinks the word **' sit '** is optional." -> "My dog thinks the word **'sit'** is optional."*
* *"His plan sounded **' brilliant '**. Until we actually tried it." -> "His plan sounded **'brilliant'**. Until we actually tried it."*

If there are **multiple words** between the quotes, remove spaces only before the first word’s first letter and after the last word’s last letter, so that each quote touches the respective edge word.

Examples

* *"As Elton John said: **' I am the most well-known homosexual in the world '**" -> "As Elton John said: **'I am the most well-known homosexual in the world'**"*
* *"My cat looked at me like I was **' beneath her '** — and honestly, she’s probably right." -> "My cat looked at me like I was **'beneath her'** — and honestly, she’s probably right."*
* *"He said I'm **' too dramatic '** just because I cried over a pizza ad" -> "He said I'm **'too dramatic'** just because I cried over a pizza ad."*

---

### 8. **INDEFINITE ARTICLE a - an**

**RULE -->** Wherever there is the indefinite article **a**, change it to **an** when the following word begins with a vowel (a, e, i, o, u) or with **h**.

Examples

* *"There it was. A amazing rock!" -> "There it was. An amazing rock!"*
* *"That's **a awesome** bike!" -> "That's **an awesome** bike!"*
* *"**A hour** is definitely not enough for this" -> "**An hour** is definitely not enough for this"*

**NOTE**

In real English grammar, the rule for **a/an** concerning **h** and **y**, as well as some other cases, is more complex. Specifically, **a** changes to **an** when the next word begins with a **vowel sound**, not necessarily a vowel letter.

Examples of words starting with a vowel letter but a consonant sound:

* a home
* a university
* a European
* a young man
* a year

Examples of words starting with a consonant letter but a vowel sound:

* an yttrium sample
* an xylophone
* an MBA

This occurs rarely in everyday English, mainly in scientific or foreign words.
