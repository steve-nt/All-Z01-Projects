use fibonacci2::*;

fn main() {
    println!(
        "The element in the position {} in fibonacci series is {}",
        2,
        fibonacci(2)
    );
    println!(
        "The element in the position {} in fibonacci series is {}",
        4,
        fibonacci(4)
    );
    println!(
        "The element in the position {} in fibonacci series is {}",
        22,
        fibonacci(22)
    );
    println!(
        "The element in the position {} in fibonacci series is {}",
        20,
        fibonacci(20)
    );
}


/*And its output:

$ cargo run
The element in the position 2 in fibonacci series is 1
The element in the position 4 in fibonacci series is 3
The element in the position 22 in fibonacci series is 17711
The element in the position 20 in fibonacci series is 6765
$
*/