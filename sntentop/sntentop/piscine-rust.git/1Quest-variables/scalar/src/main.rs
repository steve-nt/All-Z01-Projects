use scalar::*;

fn main() {
    // sum
    println!("sum: {}", sum(234, 2)); 
    // println!("sum: {}", sum(1, 255)); // Uncommenting this will cause a panic!

    // diff
    println!("diff: {}", diff(234, 2)); 
    // println!("diff: {}", diff(-32768, 32766)); // Uncommenting this will cause a panic!

    // product
    println!("pro: {}", pro(23, 2)); 
    // println!("pro: {}", pro(-128, 2)); // Uncommenting this will cause a panic!

    // quotient
    println!("quo: {}", quo(22.0, 2.0)); 
    println!("quo: {}", quo(-128.23, 2.0)); 

    // remainder
    println!("rem: {}", rem(22.0, 2.0)); 
    println!("rem: {}", rem(-128.23, 2.0)); 
}