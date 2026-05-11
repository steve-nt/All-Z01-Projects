fn parse_number(token: &str) -> u32 {
    let token = token.trim();
    let (num_str, factor) = if let Some(stripped) = token.strip_suffix('k') {
        (stripped, 1000.0)
    } else if let Some(stripped) = token.strip_suffix('K') {
        (stripped, 1000.0)
    } else {
        (token, 1.0)
    };

    let v: f64 = num_str
        .parse()
        .unwrap_or_else(|_| panic!("invalid number token: {token:?}"));

    let scaled = v * factor;
    if !scaled.is_finite() || scaled < 0.0 || scaled > u32::MAX as f64 {
        panic!("number out of range for u32: {token:?}");
    }

    scaled.round() as u32
}

pub fn parse_into_boxed(s: String) -> Vec<Box<u32>> {
    s.split_whitespace()
        .map(|token| Box::new(parse_number(token)))
        .collect()
}

pub fn into_unboxed(a: Vec<Box<u32>>) -> Vec<u32> {
    a.into_iter().map(|b| *b).collect()
}
