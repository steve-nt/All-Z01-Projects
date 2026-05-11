pub fn arrange_phrase(phrase: &str) -> String {
    
    
    let mut words = Vec::with_capacity(phrase.split_whitespace().count());
    words.extend(phrase.split_whitespace());

    
    words.sort_by_key(|word| {
        
        
        word.chars()
            .filter(|c| c.is_ascii_digit())
            .fold(0, |acc, c| acc * 10 + c.to_digit(10).unwrap() as usize)
    });

    
    
    let mut result = String::with_capacity(phrase.len());
    
    for (i, word) in words.iter().enumerate() {
        if i > 0 {
            result.push(' ');
        }
        
        
        for c in word.chars() {
            if !c.is_ascii_digit() {
                result.push(c);
            }
        }
    }

    result
}