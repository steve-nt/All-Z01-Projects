pub fn is_anagram(s1: &str, s2: &str) -> bool {
    // πρώτα κανονικοποιούμε και τα δύο strings σε lower case
    let n1: String = s1.chars().flat_map(|c| c.to_lowercase()).collect();
    let n2: String = s2.chars().flat_map(|c| c.to_lowercase()).collect();

    // αν τα μήκη σε χαρακτήρες διαφέρουν, δεν είναι anagram
    if n1.chars().count() != n2.chars().count() {
        return false;
    }

    let mut c1: Vec<char> = n1.chars().collect();
    let mut c2: Vec<char> = n2.chars().collect();

    c1.sort_unstable();
    c2.sort_unstable();

    c1 == c2
}