pub fn rotate(input: &str, key: i8) -> String {
    let shift = key.rem_euclid(26) as u8;

    input
        .chars()
        .map(|c| {
            if c.is_ascii_lowercase() {
                let base = b'a';
                let pos = c as u8 - base;
                (base + (pos + shift) % 26) as char
            } else if c.is_ascii_uppercase() {
                let base = b'A';
                let pos = c as u8 - base;
                (base + (pos + shift) % 26) as char
            } else {
                c
            }
        })
        .collect()
}

