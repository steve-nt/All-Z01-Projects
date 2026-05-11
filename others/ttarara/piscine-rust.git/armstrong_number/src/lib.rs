pub fn is_armstrong_number(nb: u32) -> Option<u32> {
    let digits = nb.to_string();
    let power = digits.len() as u32;

    let sum: u32 = digits
        .chars()
        .map(|ch| ch.to_digit(10).unwrap().pow(power))
        .sum();

    if sum == nb {
        Some(nb)
    } else {
        None
    }
}