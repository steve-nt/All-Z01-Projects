pub fn initials(names: Vec<&str>) -> Vec<String> {

    let mut result = Vec::with_capacity(names.len());

    for name in names {
        let words = name.split_whitespace();

        let mut init = String::with_capacity(words.clone().count() * 3);

        let mut is_first = true;

        for word in words {

            if let Some(c) = word.chars().next() {
                if !is_first {
                    init.push(' ');
                }

                init.push(c);
                init.push('.');
                
                is_first = false;
            }
        }
        
        result.push(init);
    }

    result
}