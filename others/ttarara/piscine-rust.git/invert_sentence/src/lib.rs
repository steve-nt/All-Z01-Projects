pub fn invert_sentence(string: &str) -> String {
    #[derive(Debug)]
    struct Token {
        is_ws: bool,
        text: String,
    }

    let mut tokens: Vec<Token> = Vec::new();
    let mut buf = String::new();
    let mut current_is_ws: Option<bool> = None;

    for ch in string.chars() {
        let is_ws = ch.is_whitespace();
        match current_is_ws {
            None => {
                current_is_ws = Some(is_ws);
                buf.push(ch);
            }
            Some(cur) if cur == is_ws => {
                buf.push(ch);
            }
            Some(cur) => {
                tokens.push(Token {
                    is_ws: cur,
                    text: std::mem::take(&mut buf),
                });
                current_is_ws = Some(is_ws);
                buf.push(ch);
            }
        }
    }

    if let Some(cur) = current_is_ws {
        tokens.push(Token { is_ws: cur, text: buf });
    }

    let mut words: Vec<String> = tokens
        .iter()
        .filter(|t| !t.is_ws)
        .map(|t| t.text.clone())
        .collect();
    words.reverse();

    let mut word_iter = words.into_iter();
    let mut out = String::with_capacity(string.len());
    for tok in tokens {
        if tok.is_ws {
            out.push_str(&tok.text);
        } else {
            out.push_str(&word_iter.next().unwrap_or_default());
        }
    }

    out
}

#[cfg(test)]
mod tests {
    use super::invert_sentence;

    #[test]
    fn examples() {
        assert_eq!(invert_sentence("Rust is Awesome"), "Awesome is Rust");
        assert_eq!(invert_sentence("    word1     word2  "), "    word2     word1  ");
        assert_eq!(invert_sentence("Hello, World!"), "World! Hello,");
    }

    #[test]
    fn preserves_tabs_and_newlines() {
        assert_eq!(
            invert_sentence("a\tb\nc  d"),
            "d\tc\nb  a"
        );
    }

    #[test]
    fn empty_and_whitespace_only() {
        assert_eq!(invert_sentence(""), "");
        assert_eq!(invert_sentence("   \n\t "), "   \n\t ");
    }
}
