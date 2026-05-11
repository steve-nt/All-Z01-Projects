pub fn scytale_cipher(message: &str, i: usize) -> String {
    if i == 0 {
        return String::new();
    }

    let chars: Vec<char> = message.chars().collect();
    let len = chars.len();
    if len == 0 {
        return String::new();
    }

    let rows = len.div_ceil(i);
    let padded_len = rows * i;

    
    let mut padded = chars;
    padded.resize(padded_len, ' ');

    let mut out = String::with_capacity(padded_len);
    for col in 0..i {
        for row in 0..rows {
            out.push(padded[row * i + col]);
        }
    }
    out.trim_end().to_string()
}


