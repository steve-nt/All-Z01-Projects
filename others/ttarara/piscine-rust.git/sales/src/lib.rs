#[derive(Debug, Clone, PartialEq)]
pub struct Store {
    pub products: Vec<(String, f32)>,
}

impl Store {
    pub fn new(products: Vec<(String, f32)>) -> Store {
        Store { products }
    }
}

#[derive(Debug, Clone, PartialEq)]
pub struct Cart {
    pub items: Vec<(String, f32)>,
    pub receipt: Vec<f32>,
}

fn round_two_decimals(x: f32) -> f32 {
    (x * 100.0).round() / 100.0
}

impl Cart {
    pub fn new() -> Cart {
        Cart {
            items: Vec::new(),
            receipt: Vec::new(),
        }
    }

    pub fn insert_item(&mut self, s: &Store, ele: String) {
        if let Some((_, price)) = s.products.iter().find(|(n, _)| n == &ele) {
            self.items.push((ele, *price));
        }
    }

    pub fn generate_receipt(&mut self) -> Vec<f32> {
        let prices: Vec<f32> = self.items.iter().map(|(_, p)| *p).collect();
        let total: f32 = prices.iter().copied().sum();
        let free_total: f32 = prices
            .chunks(3)
            .filter(|c| c.len() == 3)
            .map(|chunk| chunk.iter().copied().fold(f32::INFINITY, f32::min))
            .sum();
        let factor = if total != 0.0 {
            (total - free_total) / total
        } else {
            1.0
        };
        let scaled = |p: f32| round_two_decimals(p * factor);
        let mut receipt: Vec<f32> = prices.iter().copied().map(scaled).collect();
        receipt.sort_by(|a, b| a.total_cmp(b));
        self.receipt.clone_from(&receipt);
        receipt
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn usage_three_products() {
        let store = Store::new(vec![
            (String::from("product A"), 1.23),
            (String::from("product B"), 23.1),
            (String::from("product C"), 3.12),
        ]);
        let mut cart = Cart::new();
        cart.insert_item(&store, String::from("product A"));
        cart.insert_item(&store, String::from("product B"));
        cart.insert_item(&store, String::from("product C"));
        let r = cart.generate_receipt();
        assert_eq!(r, vec![1.17, 2.98, 22.06]);
        assert_eq!(cart.receipt, r);
    }

    #[test]
    fn nine_items_promotion() {
        let store = Store::new(vec![
            (String::from("a"), 1.23),
            (String::from("b"), 23.1),
            (String::from("c"), 3.12),
            (String::from("d"), 9.75),
            (String::from("e"), 1.75),
            (String::from("f"), 23.75),
            (String::from("g"), 2.75),
            (String::from("h"), 1.64),
            (String::from("i"), 15.23),
        ]);
        let mut cart = Cart::new();
        for name in ["a", "b", "c", "d", "e", "f", "g", "h", "i"] {
            cart.insert_item(&store, String::from(name));
        }
        let r = cart.generate_receipt();
        assert_eq!(
            r,
            vec![1.16, 1.55, 1.65, 2.6, 2.94, 9.2, 14.38, 21.8, 22.42]
        );
    }

    #[test]
    fn seven_items_two_triples_one_leftover() {
        let store = Store::new(vec![
            (String::from("a"), 3.12),
            (String::from("b"), 9.75),
            (String::from("c"), 1.75),
            (String::from("d"), 23.75),
            (String::from("e"), 2.75),
            (String::from("f"), 1.64),
            (String::from("g"), 15.23),
        ]);
        let mut cart = Cart::new();
        for name in ["a", "b", "c", "d", "e", "f", "g"] {
            cart.insert_item(&store, String::from(name));
        }
        let r = cart.generate_receipt();
        assert_eq!(
            r,
            vec![1.54, 1.65, 2.59, 2.94, 9.18, 14.34, 22.36]
        );
    }
}
