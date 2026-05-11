pub fn nbr_function(c: i32) -> (i32, f64, f64) {

    let val = c as f64;
    
    let exp_val = val.exp();
    let ln_val = val.abs().ln();

    (c, exp_val, ln_val)
}

pub fn str_function(a: String) -> (String, String) {

    let exp_str = a
        .split_whitespace()
        .map(|num_str| {

            let val: f64 = num_str.parse().unwrap();
            val.exp().to_string()
        })
        .collect::<Vec<String>>()
        .join(" "); 
    (a, exp_str)
}

pub fn vec_function(b: Vec<i32>) -> (Vec<i32>, Vec<f64>) {
    let ln_vec: Vec<f64> = b
        .iter()
        .map(|&val| (val as f64).abs().ln())
        .collect();

    (b, ln_vec)
}