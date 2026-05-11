use logic_number::*;

fn main() {
    let array = [9, 10, 153, 154];
    for pat in &array {
        if number_logic(*pat) {
            println!(
                "this number returns {} because the number {} obey the rules of the sequence",
                number_logic(*pat),
                pat
            );
        } else {
            println!(
                "this number returns {} because the number {} does not obey the rules of the sequence",
                number_logic(*pat),
                pat
            );
        }
    }
}

