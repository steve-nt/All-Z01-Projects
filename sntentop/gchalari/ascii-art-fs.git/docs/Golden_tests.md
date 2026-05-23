# GOLDEN TESTS — ascii-art-output

Golden tests are fixed input/output pairs taken directly from the audit. The program output must match exactly.

Note: the audit uses `cat -e` to display file contents, which adds `$` at the end of each line to show line endings. The actual expected file content does not contain `$`.

---

```bash
go run . "banana" standard abc

go run . "hello" standard | cat -e

go run . "hello world" shadow | cat -e

go run . "nice 2 meet you" thinkertoy | cat -e

go run . "you & me" standard | cat -e

go run . "123" shadow | cat -e

go run . "/(\")" thinkertoy | cat -e

go run . "ABCDEFGHIJKLMNOPQRSTUVWXYZ" shadow | cat -e

go run . "\"#$%&/()*+,-./" thinkertoy | cat -e

go run . "It's Working" thinkertoy | cat -e

go run . "HeLLo WoRLD" standard | cat -e

go run . "2025" standard | cat -e

go run . "@#$%^&*" shadow | cat -e

go run . "GoLang 123 !@#" thinkertoy | cat -e

go run . "hello\nworld" standard | cat -e

go run . "" standard | cat -e

go run . "   " thinkertoy | cat -e

go run . "hello" invalidbanner
```