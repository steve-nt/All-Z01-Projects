pub fn fibonacci(n: u32) -> u32 {
    match n {
        0 => 0,
        1 => 1,
        // Calculate the sum of the two preceding numbers
        _ => fibonacci(n - 1) + fibonacci(n - 2), 
    }
}

/*
Alternative

pub fn fibonacci(n: u32) -> u32 {
    if n == 0 {
        return 0;
    }
    
    let mut a = 0;
    let mut b = 1;
    
    // Loop starting from 2 up to n (inclusive)
    for _ in 2..=n {
        let next = a + b;
        a = b;
        b = next;
    }
    
    b
}
*/