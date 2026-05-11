pub fn delete_and_backspace(s: &mut String) {
    let mut result = String::new();
    let mut skip_count = 0;

    for c in s.chars() {
        if c == '+' {
            
            skip_count += 1;
        } else if skip_count > 0 {
            
            skip_count -= 1;
        } else if c == '-' {
            
            result.pop();
        } else {
            
            result.push(c);
        }
    }
    
    
    *s = result;
}

pub fn do_operations(v: &mut [String]) {
    
    for s in v.iter_mut() {
        if let Some(pos) = s.find('+') {
            
            let left: i32 = s[..pos].trim().parse().unwrap();
            let right: i32 = s[pos + 1..].trim().parse().unwrap();
            
            
            *s = (left + right).to_string();
            
        } else if let Some(pos) = s.find('-') {
            
            let left: i32 = s[..pos].trim().parse().unwrap();
            let right: i32 = s[pos + 1..].trim().parse().unwrap();
            
            
            *s = (left - right).to_string();
        }
    }
}