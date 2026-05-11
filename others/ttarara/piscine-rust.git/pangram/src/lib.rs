pub fn is_pangram(s: &str) -> bool {
    let mut seen = [false; 26];

    for ch in s.chars() {
        let lower = ch.to_ascii_lowercase();
        if ('a'..='z').contains(&lower) {
            let idx = (lower as u8 - b'a') as usize;
            seen[idx] = true;
        }
    }

    seen.into_iter().all(|b| b)
}

