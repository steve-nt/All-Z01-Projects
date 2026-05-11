pub fn num_to_ordinal(x: u32) -> String {
    let suffix = match x % 100 {
        11 | 12 | 13 => "th",
        _ => match x % 10 {
            1 => "st",
            2 => "nd",
            3 => "rd",
            _ => "th",
        },
    };

    format!("{x}{suffix}")
}

#[cfg(test)]
mod tests {
    use super::num_to_ordinal;

    #[test]
    fn examples() {
        assert_eq!(num_to_ordinal(1), "1st");
        assert_eq!(num_to_ordinal(22), "22nd");
        assert_eq!(num_to_ordinal(43), "43rd");
        assert_eq!(num_to_ordinal(47), "47th");
    }

    #[test]
    fn teen_exceptions() {
        assert_eq!(num_to_ordinal(11), "11th");
        assert_eq!(num_to_ordinal(12), "12th");
        assert_eq!(num_to_ordinal(13), "13th");
    }

    #[test]
    fn last_digit_rules() {
        assert_eq!(num_to_ordinal(21), "21st");
        assert_eq!(num_to_ordinal(32), "32nd");
        assert_eq!(num_to_ordinal(53), "53rd");
        assert_eq!(num_to_ordinal(24), "24th");
    }
}

