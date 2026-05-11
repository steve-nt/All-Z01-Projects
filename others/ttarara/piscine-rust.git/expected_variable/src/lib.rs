pub fn edit_distance(source: &str, target: &str) -> usize {
    let source_chars: Vec<char> = source.chars().collect();
    let target_chars: Vec<char> = target.chars().collect();
    let source_len = source_chars.len();
    let target_len = target_chars.len();

    let mut table = vec![vec![0; target_len + 1]; source_len + 1];

    for i in 0..=source_len {
        table[i][0] = i;
    }
    for j in 0..=target_len {
        table[0][j] = j;
    }

    for i in 1..=source_len {
        for j in 1..=target_len {
            if source_chars[i - 1] == target_chars[j - 1] {
                table[i][j] = table[i - 1][j - 1];
            } else {
                let substitute = table[i - 1][j - 1] + 1;
                let delete = table[i - 1][j] + 1;
                let insert = table[i][j - 1] + 1;
                table[i][j] = substitute.min(delete).min(insert);
            }
        }
    }

    table[source_len][target_len]
}

pub fn expected_variable(compared: &str, expected: &str) -> Option<String> {
    let is_snake = compared.contains('_') && compared.chars().all(|c| c.is_alphanumeric() || c == '_');
    let is_camel = !compared.contains('_') && !compared.is_empty() && compared.chars().all(|c| c.is_alphanumeric());

    if !is_snake && !is_camel {
        return None;
    }

    let compared_lower = compared.to_lowercase();
    let expected_lower = expected.to_lowercase();
    
    let distance = edit_distance(&compared_lower, &expected_lower);
    let expected_len = expected_lower.len();
    
    if expected_len == 0 || distance > expected_len {
        return None;
    }
    
    let similarity = ((expected_len - distance) as f64 / expected_len as f64) * 100.0;
    
    if similarity > 50.0 {
        Some(format!("{}%", similarity.round() as usize))
    } else {
        None
    }
}


