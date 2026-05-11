pub fn arrange_phrase(phrase: &str) -> String {
    let mut words: Vec<(usize, &str)> = phrase
        .split_whitespace()
        .map(|word| {
            let position = word
                .chars()
                .find(|ch| ch.is_ascii_digit())
                .unwrap()
                .to_digit(10)
                .unwrap() as usize;

            (position, word)
        })
        .collect();

    words.sort_unstable_by_key(|&(position, _)| position);

    let mut result = String::with_capacity(phrase.len());

    for (i, (_, word)) in words.iter().enumerate() {
        if i > 0 {
            result.push(' ');
        }

        for ch in word.chars() {
            if !ch.is_ascii_digit() {
                result.push(ch);
            }
        }
    }

    result
}