fn spell_under_20(n: u16) -> &'static str {
    match n {
        0 => "zero",
        1 => "one",
        2 => "two",
        3 => "three",
        4 => "four",
        5 => "five",
        6 => "six",
        7 => "seven",
        8 => "eight",
        9 => "nine",
        10 => "ten",
        11 => "eleven",
        12 => "twelve",
        13 => "thirteen",
        14 => "fourteen",
        15 => "fifteen",
        16 => "sixteen",
        17 => "seventeen",
        18 => "eighteen",
        19 => "nineteen",
        _ => unreachable!(),
    }
}

fn spell_tens(n: u16) -> String {
    if n < 20 {
        return spell_under_20(n).to_string();
    }

    let tens = n / 10;
    let ones = n % 10;

    let tens_word = match tens {
        2 => "twenty",
        3 => "thirty",
        4 => "forty",
        5 => "fifty",
        6 => "sixty",
        7 => "seventy",
        8 => "eighty",
        9 => "ninety",
        _ => "",
    };

    if ones == 0 {
        tens_word.to_string()
    } else {
        format!("{}-{}", tens_word, spell_under_20(ones))
    }
}

fn spell_under_1000(n: u16) -> String {
    let hundreds = n / 100;
    let rem = n % 100;
    let mut parts = Vec::new();

    if hundreds > 0 {
        parts.push(format!("{} hundred", spell_under_20(hundreds)));
    }
    if rem > 0 {
        parts.push(spell_tens(rem));
    }

    if parts.is_empty() {
        "zero".to_string()
    } else {
        parts.join(" ")
    }
}

pub fn spell(n: u64) -> String {
    if n == 0 {
        return "zero".to_string();
    }
    if n == 1_000_000 {
        return "one million".to_string();
    }

    let thousands = (n / 1000) as u16;
    let rem = (n % 1000) as u16;

    let mut parts = Vec::new();

    if thousands > 0 {
        parts.push(format!("{} thousand", spell_under_1000(thousands)));
    }
    if rem > 0 {
        parts.push(spell_under_1000(rem));
    }

    parts.join(" ")
}


