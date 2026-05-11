pub fn initials(names: Vec<&str>) -> Vec<String> {
    names
        .into_iter()
        .map(|name| {
            let mut result = String::new();

            for (i, part) in name.split_whitespace().enumerate() {
                if i > 0 {
                    result.push(' ');
                }

                result.push(part.chars().next().unwrap());
                result.push('.');
            }

            result
        })
        .collect()
}