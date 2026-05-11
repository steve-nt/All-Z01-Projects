#[derive(Debug, PartialEq)]
pub struct CipherError {
    pub expected: String,
}

fn atbash_char(c: char) -> char {
    if c.is_ascii_lowercase() {
        // 'a' ↔ 'z', 'b' ↔ 'y', ...
        (b'z' - (c as u8 - b'a')) as char
    } else if c.is_ascii_uppercase() {
        // 'A' ↔ 'Z', 'B' ↔ 'Y', ...
        (b'Z' - (c as u8 - b'A')) as char
    } else {
        c
    }
}

pub fn cipher(original: &str, ciphered: &str) -> Result<(), CipherError> {
    let expected: String = original.chars().map(atbash_char).collect();

    if expected == ciphered {
        Ok(())
    } else {
        Err(CipherError { expected })
    }
}