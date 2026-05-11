pub fn nbr_function(c: i32) -> (i32, f64, f64) {
    let x = c as f64;
    let abs_x = c.abs() as f64;
    (c, x.exp(), abs_x.ln())
}

pub fn vec_function(b: Vec<i32>) -> (Vec<i32>, Vec<f64>) {
    let logs = b.iter().map(|&n| (n as f64).ln()).collect();
    (b, logs)
}

pub fn str_function(a: String) -> (String, String) {
    let exp_values = a
        .split_whitespace()
        .map(|part| part.parse::<f64>().unwrap().exp().to_string())
        .collect::<Vec<String>>()
        .join(" ");

    (a, exp_values)
}