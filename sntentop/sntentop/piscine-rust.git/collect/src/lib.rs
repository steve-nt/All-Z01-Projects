pub fn bubble_sort(arr: &mut [i32]) {
    let mut n = arr.len();
    let mut swapped = true;

    
    while swapped {
        swapped = false;
        
        
        for i in 1..n {
            if arr[i - 1] > arr[i] {
                
                arr.swap(i - 1, i);
                swapped = true;
            }
        }
        
        
        
        if n > 0 {
            n -= 1;
        }
    }
}