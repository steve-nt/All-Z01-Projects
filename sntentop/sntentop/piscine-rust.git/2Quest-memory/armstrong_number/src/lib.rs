pub fn is_armstrong_number(nb: u32) -> Option<u32> {
    let nb_str = nb.to_string();
    let power = nb_str.len() as u32;

    
    
    
    let sum: u64 = nb_str
        .chars()
        .filter_map(|c| c.to_digit(10))
        .map(|d| (d as u64).pow(power))
        .sum();

    if sum == (nb as u64) {
        Some(nb)
    } else {
        None
    }
}