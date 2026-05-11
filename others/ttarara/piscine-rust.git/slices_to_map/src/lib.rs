use std::collections::HashMap;
use std::hash::Hash;

pub fn slices_to_map<'a, T: Eq + Hash, U>(keys: &'a [T], values: &'a [U]) -> HashMap<&'a T, &'a U> {
    keys.iter().zip(values.iter()).collect()
}
