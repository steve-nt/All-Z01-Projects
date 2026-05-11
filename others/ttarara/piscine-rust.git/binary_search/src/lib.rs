pub fn binary_search(sorted_list: &[i32], target: i32) -> Option<usize> {
    let mut left = 0;
    let mut right = sorted_list.len();

    while left < right {
        let mid = left + (right - left) / 2;
        let value = sorted_list[mid];

        if value == target {
            return Some(mid);
        } else if value < target {
            left = mid + 1;
        } else {
            right = mid;
        }
    }

    None
}