pub fn delete_and_backspace(s: &mut String) {
    let chars: Vec<char> = s.chars().collect();
    let mut removed = vec![false; chars.len()];

    for i in 0..chars.len() {
        match chars[i] {
            '-' => {
                removed[i] = true;

                let mut j = i;
                while j > 0 {
                    j -= 1;
                    if !removed[j] && chars[j] != '-' && chars[j] != '+' {
                        removed[j] = true;
                        break;
                    }
                }
            }
            '+' => {
                removed[i] = true;

                let mut j = i + 1;
                while j < chars.len() {
                    if !removed[j] && chars[j] != '-' && chars[j] != '+' {
                        removed[j] = true;
                        break;
                    }
                    j += 1;
                }
            }
            _ => {}
        }
    }

    *s = chars
        .iter()
        .enumerate()
        .filter(|(i, ch)| !removed[*i] && **ch != '-' && **ch != '+')
        .map(|(_, ch)| *ch)
        .collect();
}

pub fn do_operations(v: &mut [String]) {
    for expr in v.iter_mut() {
        if let Some((left, right)) = expr.split_once('+') {
            let a: i32 = left.parse().unwrap();
            let b: i32 = right.parse().unwrap();
            *expr = (a + b).to_string();
        } else if let Some((left, right)) = expr.split_once('-') {
            let a: i32 = left.parse().unwrap();
            let b: i32 = right.parse().unwrap();
            *expr = (a - b).to_string();
        }
    }
}