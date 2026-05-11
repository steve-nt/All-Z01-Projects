pub fn fibonacci(n: u32) -> u32 {

    if n == 0 {
        return 0;
    }

    if n == 1 {
        return 1;
    }

    let mut previous = 0;
    let mut current = 1;

    for _i in 2..=n{
        let next = previous + current;
        previous = current; 
        current = next;
    }

    return current;
}
