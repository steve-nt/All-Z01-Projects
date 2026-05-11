use std::io::{self, Write};

fn main() {
    let riddle = "I am the beginning of the end, and the end of time and space. I am essential to creation, and I surround every place. What am I?";
    let answer = "The letter e";
    let mut trials = 0;

    loop {
        println!("{}", riddle);

        trials += 1;

        let mut input = String::new();
        io::stdin()
            .read_line(&mut input)
            .expect("Failed to read line");

        if input.trim() == answer {
            println!("Number of trials: {}", trials);
            break; 
        }
    }
}