use lalgebra_vector::*;

fn main() {
    println!("{:?}", Vector(vec![1, 3, -5]).dot(Vector(vec![4, -2, -1])));
    println!("{:?}", Vector(vec![1, 3, -5]) + Vector(vec![4, -2, -1]));
}
