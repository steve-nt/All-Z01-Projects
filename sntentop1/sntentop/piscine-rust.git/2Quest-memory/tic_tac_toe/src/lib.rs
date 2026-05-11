pub fn tic_tac_toe(table: [[char; 3]; 3]) -> String {
    
    if diagonals('O', table) || horizontal('O', table) || vertical('O', table) {
        return "player O won".to_string();
    }
    
    
    if diagonals('X', table) || horizontal('X', table) || vertical('X', table) {
        return "player X won".to_string();
    }
    
    
    "tie".to_string()
}

pub fn diagonals(player: char, table: [[char; 3]; 3]) -> bool {
    
    let diag1 = table[0][0] == player && table[1][1] == player && table[2][2] == player;
    
    
    let diag2 = table[0][2] == player && table[1][1] == player && table[2][0] == player;
    
    diag1 || diag2
}

pub fn horizontal(player: char, table: [[char; 3]; 3]) -> bool {
    
    for row in table.iter() {
        
        if row[0] == player && row[1] == player && row[2] == player {
            return true;
        }
    }
    false
}

pub fn vertical(player: char, table: [[char; 3]; 3]) -> bool {
    
    for col in 0..3 {
        
        if table[0][col] == player && table[1][col] == player && table[2][col] == player {
            return true;
        }
    }
    false
}