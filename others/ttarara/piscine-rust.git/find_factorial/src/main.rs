use find_factorial::*;

fn main() {
    println!("The factorial of 0 = {}", factorial(0));
    println!("The factorial of 1 = {}", factorial(1));
    println!("The factorial of 5 = {}", factorial(5));
    println!("The factorial of 10 = {}", factorial(10));
    println!("The factorial of 19 = {}", factorial(19));
}

/*
And its output:

$ cargo run
The factorial of 0 = 1
The factorial of 1 = 1
The factorial of 5 = 120
The factorial of 10 = 3628800
The factorial of 19 = 121645100408832000
$
*/