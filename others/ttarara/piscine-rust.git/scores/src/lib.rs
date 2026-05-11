
pub fn score(string: &str) -> u64 {
    let mut total: u64 = 0;

    for ch in string.chars() {
        // Map only ASCII letters; everything else scores 0.
        let up = ch.to_ascii_uppercase();

        total += match up {
            'A' | 'E' | 'I' | 'O' | 'U' | 'L' | 'N' | 'R' | 'S' | 'T' => 1,
            'D' | 'G' => 2,
            'B' | 'C' | 'M' | 'P' => 3,
            'F' | 'H' | 'V' | 'W' | 'Y' => 4,
            'K' => 5,
            'J' | 'X' => 8,
            'Q' | 'Z' => 10,
            _ => 0,
        };
    }

    total
}

