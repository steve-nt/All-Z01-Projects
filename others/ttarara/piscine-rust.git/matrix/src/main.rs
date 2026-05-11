use matrix::*;

fn main() {
    let m = Matrix([[0; 4]; 3]);
    println!("{:?}", m);
    println!("{:?}", Matrix::<4, 4, u32>::identity());
    println!("{:?}", Matrix::<3, 4, f64>::zero());
}