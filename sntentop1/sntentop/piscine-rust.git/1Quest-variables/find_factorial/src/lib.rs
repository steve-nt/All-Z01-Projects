pub fn factorial(num: u64) -> u64 {
    match num {
        0 | 1 => 1,                       
        _ => num * factorial(num - 1),   
    }
}

/*
ALTERNATIVE WAY 
pub fn factorial(num: u64) -> u64 {
    // Creates an inclusive range from 1 to num, and multiplies them all.
    // If num is 0, the range is empty, and .product() automatically returns 1.
    (1..=num).product()
}
*/