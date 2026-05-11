pub fn capitalize_first(input: &str) -> String {
    let mut chars = input.chars();

    match chars.next() {
        Some(first) => {
            let mut result = String::new();
            result.push(first.to_uppercase().next().unwrap());
            result.extend(chars);
            result
        }
        None => String::new(),
    }
}

pub fn title_case(input: &str) -> String {
    let mut result = String::new();
    let mut capitalize_next = true;

    for ch in input.chars() {
        if ch.is_whitespace() {
            result.push(ch);
            capitalize_next = true;
        } else if capitalize_next {
            result.push(ch.to_uppercase().next().unwrap());
            capitalize_next = false;
        } else {
            result.push(ch);
        }
    }

    result
}

pub fn change_case(input: &str) -> String {
    let mut result = String::new();

    for ch in input.chars() {
        if ch.is_lowercase() {
            result.push(ch.to_uppercase().next().unwrap());
        } else if ch.is_uppercase() {
            result.push(ch.to_lowercase().next().unwrap());
        } else {
            result.push(ch);
        }
    }

    result
}