use easy_traits::*;

fn main() {
    let mut s = "hello".to_owned();

    println!("Before append: {}", s);

    s.append_str(" there!");
    println!("After append: {}", s);

    s.remove_punctuation_marks();
    println!("After removing punctuation: {}", s);
}
