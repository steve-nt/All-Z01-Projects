pub fn stars(n: u32) -> String {
    // 2^n stars
    let count = 1usize
        .checked_shl(n)
        .expect("n too large to allocate");
    "*".repeat(count)
}

#[cfg(test)]
mod tests {
    use super::stars;

    #[test]
    fn examples() {
        assert_eq!(stars(1), "**");
        assert_eq!(stars(4), "*".repeat(16));
        assert_eq!(stars(5), "*".repeat(32));
    }

    #[test]
    fn n_zero() {
        assert_eq!(stars(0), "*");
    }
}

