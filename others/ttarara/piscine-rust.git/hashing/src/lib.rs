use std::collections::HashMap;

pub fn mean(list: &[i32]) -> f64 {
    let sum: i32 = list.iter().sum();
    sum as f64 / list.len() as f64
}

pub fn median(list: &[i32]) -> i32 {
    let mut values = list.to_vec();
    values.sort_unstable();

    let mid = values.len() / 2;

    if values.len() % 2 == 0 {
        (values[mid - 1] + values[mid]) / 2
    } else {
        values[mid]
    }
}

pub fn mode(list: &[i32]) -> i32 {
    let mut counts = HashMap::new();

    for &value in list {
        *counts.entry(value).or_insert(0) += 1;
    }

    let mut most_common = list[0];
    let mut max_count = 0;

    for (&value, &count) in &counts {
        if count > max_count {
            max_count = count;
            most_common = value;
        }
    }

    most_common
}