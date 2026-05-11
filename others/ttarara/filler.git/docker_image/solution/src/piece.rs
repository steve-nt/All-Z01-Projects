#[derive(Debug, Clone)]
pub struct Piece {
    pub w: u16,
    pub h: u16,
    pub hashes: Vec<(u16, u16)>,
}

impl Piece {
    pub fn new(w: u16, h: u16, hashes: Vec<(u16, u16)>) -> Self {
        Self { w, h, hashes }
    }
}
