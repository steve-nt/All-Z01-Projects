pub fn factorial(num: u64) -> u64 {

    if num == 0 || num == 1{
        return 1;
    }

    let mut result = 1;
    for i in 1..=num {
        result *= i;
    }                    

    result
}