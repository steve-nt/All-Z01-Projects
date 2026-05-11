pub fn get_diamond(c: char) -> Vec<String> {
    // Special case: just 'A'
    if c == 'A' {
        return vec!["A".to_string()];
    }

    let n = (c as u8 - b'A') as usize; // 0 for A, 1 for B, ..., n for c
    let width = 2 * n + 1;

    // Helper to build one line for a given letter index i (0 = A, ..., n)
    fn line_for(i: usize, n: usize, width: usize) -> String {
        if i == 0 {
            // 'A' row: centered single A
            let mut s = String::new();
            let lead = n;
            s.push_str(&" ".repeat(lead));
            s.push('A');
            s.push_str(&" ".repeat(lead));
            return s;
        }

        let ch = (b'A' + i as u8) as char;
        let lead = n - i;
        let inner = width - 2 - 2 * lead;

        let mut s = String::new();
        s.push_str(&" ".repeat(lead));
        s.push(ch);
        s.push_str(&" ".repeat(inner));
        s.push(ch);
        s.push_str(&" ".repeat(lead));
        s
    }

    let mut rows: Vec<String> = Vec::new();

    // Top half including middle row
    for i in 0..=n {
        rows.push(line_for(i, n, width));
    }
    // Bottom half (mirror without the middle row)
    for i in (0..n).rev() {
        rows.push(line_for(i, n, width));
    }

    rows
}
