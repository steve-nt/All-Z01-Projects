#[derive(Debug, PartialEq, Eq)]
pub enum Suit {
    Heart,
    Diamond,
    Spade,
    Club,
}

#[derive(Debug, PartialEq, Eq)]
pub enum Rank {
    Ace,
    Number(u8),
    Jack,
    Queen,
    King,
}


impl Suit {
    pub fn random() -> Suit {
        
        
        let random_val = (rand::random::<u8>() % 4) + 1;
        Suit::translate(random_val)
    }

    pub fn translate(value: u8) -> Suit {
        match value {
            1 => Suit::Heart,
            2 => Suit::Diamond,
            3 => Suit::Spade,
            _ => Suit::Club, 
        }
    }
}

impl Rank {
    pub fn random() -> Rank {
        
        let random_val = (rand::random::<u8>() % 13) + 1;
        Rank::translate(random_val)
    }

    pub fn translate(value: u8) -> Rank {
        match value {
            1 => Rank::Ace,
            2..=10 => Rank::Number(value), 
            11 => Rank::Jack,
            12 => Rank::Queen,
            13 => Rank::King,
            _ => Rank::Ace, 
        }
    }
}

#[derive(Debug, PartialEq, Eq)]
pub struct Card {
    pub suit: Suit,
    pub rank: Rank,
}

pub fn winner_card(card: &Card) -> bool {
    
    
    matches!(
        card,
        Card {
            suit: Suit::Spade,
            rank: Rank::Ace
        }
    )
}