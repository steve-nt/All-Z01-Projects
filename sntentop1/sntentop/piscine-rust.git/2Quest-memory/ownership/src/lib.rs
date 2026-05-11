pub fn first_subword(mut s: String) -> String {
    let mut split_index = s.len();
    for (i, c) in s.char_indices() {
        if c == '_' || (i > 0 && c.is_uppercase()) {
            split_index = i;
            break;
        }
    }
    s.truncate(split_index);

    s
}