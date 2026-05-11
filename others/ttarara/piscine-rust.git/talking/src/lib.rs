pub fn talking(text: &str) -> &str {
    let trimmed = text.trim();
    if trimmed.is_empty() {
        return "Just say something!";
    }

    let mut has_letter = false;
    let mut all_letters_uppercase = true;

    for ch in trimmed.chars() {
        if ch.is_ascii_alphabetic() {
            has_letter = true;
            if !ch.is_ascii_uppercase() {
                all_letters_uppercase = false;
                break;
            }
        }
    }

    let yelling = has_letter && all_letters_uppercase;
    let is_question = trimmed.ends_with('?');

    match (yelling, is_question) {
        (true, true) => "Quiet, I am thinking!",
        (true, false) => "There is no need to yell, calm down!",
        (false, true) => "Sure.",
        (false, false) => "Interesting",
    }
}               