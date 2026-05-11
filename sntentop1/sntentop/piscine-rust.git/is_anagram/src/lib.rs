pub fn is_anagram(s1: &str, s2: &str) -> bool {
    
    
    let clean_and_sort = |s: &str| {
        let mut chars: Vec<char> = s
            .chars()
            .filter(|c| !c.is_whitespace()) 
            .flat_map(|c| c.to_lowercase()) 
            .collect();
        
        
        chars.sort_unstable();
        
        chars
    };

    
    clean_and_sort(s1) == clean_and_sort(s2)
}