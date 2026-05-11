pub fn edit_distance(source: &str, target: &str) -> usize {
    
    let s: Vec<char> = source.chars().collect();
    let t: Vec<char> = target.chars().collect();
    let m = s.len();
    let n = t.len();
    
    if m == 0 { return n; }
    if n == 0 { return m; }
    
    let mut row: Vec<usize> = (0..=n).collect();
    for i in 0..m {
        
        let mut prev_diag = row[0];
        
        
        row[0] = i + 1;
        for j in 0..n {
            
            
            let prev_above = row[j + 1];
            
            let cost = if s[i] == t[j] { 0 } else { 1 };
            
            row[j + 1] = std::cmp::min(
                row[j + 1] + 1, 
                std::cmp::min(
                    row[j] + 1, 
                    prev_diag + cost 
                )
            );
            
            prev_diag = prev_above;
        }
    }
    
    row[n]
}