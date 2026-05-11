/// Returns `true` if `num` is equal to the sum of its digits each raised
/// to the power of the number of digits.
///
/// Example:
/// - 153 = 1^3 + 5^3 + 3^3
/// - 154 != 1^3 + 5^3 + 4^3
pub fn number_logic(num: u32) -> bool {
    let n_digits: u32 = if num == 0 {
        1
    } else {
        let mut tmp = num;
        let mut count = 0;
        while tmp > 0 {
            count += 1;
            tmp /= 10;
        }
        count
    };

    let mut tmp = num;
    let mut sum: u64 = 0;

    // Special-case: if num == 0, it has one digit '0'.
    if tmp == 0 {
        return num as u64 == 0u64.pow(n_digits);
    }

    while tmp > 0 {
        let digit = tmp % 10;
        sum += (digit as u64).pow(n_digits);
        tmp /= 10;
    }

    sum == num as u64
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn examples() {
        assert_eq!(number_logic(9), true);
        assert_eq!(number_logic(10), false);
        assert_eq!(number_logic(153), true);
        assert_eq!(number_logic(154), false);
    }
}
