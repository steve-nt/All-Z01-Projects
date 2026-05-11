pub fn first_subword(mut s: String) -> String {
    if let Some(pos) = s.find('_') {
        s.truncate(pos);
        return s;
    }

    for (i, ch) in s.char_indices().skip(1) {
        if ch.is_uppercase() {
            s.truncate(i);
            return s;
        }
    }

    s
}