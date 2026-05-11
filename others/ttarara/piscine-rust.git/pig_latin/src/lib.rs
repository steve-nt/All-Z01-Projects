fn is_vowel(c: char) -> bool {
    matches!(c, 'a' | 'e' | 'i' | 'o' | 'u')
}

pub fn pig_latin(text: &str) -> String {
    if text.is_empty() {
        return "ay".to_string();
    }

    let chars: Vec<char> = text.chars().collect();

    // Rule 1: starts with a vowel -> just add "ay"
    if is_vowel(chars[0]) {
        return format!("{text}ay");
    }

    let mut split = 0usize;
    while split < chars.len() && !is_vowel(chars[split]) {
      
        if split + 2 < chars.len() && chars[split + 1] == 'q' && chars[split + 2] == 'u' {
            split += 3;
            break;
        }
        split += 1;
    }

    let head: String = chars[..split].iter().collect();
    let tail: String = chars[split..].iter().collect();
    format!("{tail}{head}ay")
}
