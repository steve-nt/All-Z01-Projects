#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub enum Direction {
    North,
    South,
    East,
    West,
}

impl Direction {
    pub const ALL: [Self; 4] = [Self::North, Self::South, Self::East, Self::West];

    #[must_use]
    pub const fn index(self) -> usize {
        match self {
            Self::North => 0,
            Self::South => 1,
            Self::East => 2,
            Self::West => 3,
        }
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub enum Route {
    Right,
    Straight,
    Left,
}

impl Route {
    pub const ALL: [Self; 3] = [Self::Right, Self::Straight, Self::Left];
}
