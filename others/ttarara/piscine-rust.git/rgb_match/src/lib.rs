#[derive(Debug, Clone, Copy, PartialEq)]
pub struct Color {
    pub r: u8,
    pub g: u8,
    pub b: u8,
    pub a: u8,
}

impl Color {
    pub fn swap(mut self, first: u8, second: u8) -> Color {
        let mut channels = [self.r, self.g, self.b, self.a];

        let i = channels.iter().position(|&v| v == first);
        let j = channels.iter().position(|&v| v == second);

        if let (Some(i), Some(j)) = (i, j) {
            channels.swap(i, j);
        }

        self.r = channels[0];
        self.g = channels[1];
        self.b = channels[2];
        self.a = channels[3];
        self
    }
}

