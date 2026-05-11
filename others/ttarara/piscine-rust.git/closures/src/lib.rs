pub fn first_fifty_even_square() -> Vec<i32> {
    (2..)
        .step_by(2)
        .map(|x| x * x)
        .take(50)
        .collect()
}
