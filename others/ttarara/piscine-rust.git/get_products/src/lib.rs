pub fn get_products(arr: Vec<usize>) -> Vec<usize> {
    let n = arr.len();
    if n <= 1 {
        return Vec::new();
    }
    let mut prefix = vec![1; n];
    for i in 1..n {
        prefix[i] = prefix[i - 1] * arr[i - 1];
    }
    let mut suffix = vec![1; n];
    for i in (0..n - 1).rev() {
        suffix[i] = suffix[i + 1] * arr[i + 1];
    }
    (0..n).map(|i| prefix[i] * suffix[i]).collect()
}
