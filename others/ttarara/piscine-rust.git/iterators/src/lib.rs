#[derive(Copy, Clone)]
pub struct Collatz {
    pub v: u64,
}

impl Iterator for Collatz {
    type Item = Collatz;

    fn next(&mut self) -> Option<Self::Item> {
        if self.v == 0 {
            return None;
        }
        let current = self.v;
        if current == 1 {
            self.v = 0;
            return None;
        }
        self.v = if current % 2 == 0 {
            current / 2
        } else {
            current * 3 + 1
        };
        Some(Collatz { v: current })
    }
}

impl Collatz {
    pub fn new(n: u64) -> Self {
        Collatz { v: n }
    }
}

pub fn collatz(n: u64) -> usize {
    Collatz::new(n).count()
}
