use binary_search::binary_search;

fn main() {
    let sorted_list = vec![1, 3, 5, 7, 9, 11, 13];
    let target = 7;
    
    match binary_search(&sorted_list, target) {
        Some(index) => println!("{} found at index {}", target, index),
        None => println!("{} not found in the list", target),
    }
}