use std::collections::HashMap;

pub fn mean(list: &[i32]) -> f64 {
    if list.is_empty() {
        return 0.0;
    }
    
    
    
    let sum: f64 = list.iter().map(|&x| x as f64).sum();
    
    sum / (list.len() as f64)
}

pub fn median(list: &[i32]) -> i32 {
    if list.is_empty() {
        return 0;
    }
    
    
    let mut sorted_list = list.to_vec();
    sorted_list.sort();
    
    let mid = sorted_list.len() / 2;
    
    if sorted_list.len() % 2 == 0 {
        
        
        (sorted_list[mid - 1] + sorted_list[mid]) / 2
    } else {
        
        sorted_list[mid]
    }
}

pub fn mode(list: &[i32]) -> i32 {
    let mut counts = HashMap::new();

    
    for &num in list {
        *counts.entry(num).or_insert(0) += 1;
    }

    
    counts
        .into_iter()
        .max_by_key(|&(_, count)| count)
        .map(|(val, _)| val)
        .unwrap_or(0)
}