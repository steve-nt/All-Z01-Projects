pub fn capitalize_first(input: &str) -> String {
    let mut chars = input.chars();
    match chars.next() {
        
        None => String::new(),
        
        Some(first_char) => {
            let mut result = String::new();
            
            
            result.extend(first_char.to_uppercase()); 
            result.push_str(chars.as_str());
            result
        }
    }
}

pub fn title_case(input: &str) -> String {
    let mut result = String::with_capacity(input.len());
    let mut capitalize_next = true;

    for c in input.chars() {
        if c.is_whitespace() {
            result.push(c);
            capitalize_next = true;
        } else if capitalize_next {
            result.extend(c.to_uppercase());
            capitalize_next = false;
        } else {
            result.push(c);
        }
    }

    result
}

pub fn change_case(input: &str) -> String {
    let mut result = String::with_capacity(input.len());
    
    for c in input.chars() {
        if c.is_uppercase() {
            result.extend(c.to_lowercase());
        } else {
            result.extend(c.to_uppercase());
        }
    }
    
    result
}