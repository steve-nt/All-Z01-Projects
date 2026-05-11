use crate::RomanDigit::*;

#[derive(Copy, Clone, Debug, PartialEq, Eq)]
pub enum RomanDigit {
    Nulla,
    I,
    V,
    X,
    L,
    C,
    D,
    M,
}

#[derive(Clone, Debug, PartialEq, Eq)]
pub struct RomanNumber(pub Vec<RomanDigit>);

fn digit_value(d: RomanDigit) -> u32 {
    match d {
        Nulla => 0,
        I => 1,
        V => 5,
        X => 10,
        L => 50,
        C => 100,
        D => 500,
        M => 1000,
    }
}

impl RomanNumber {
    fn to_u32(&self) -> Option<u32> {
        let digits = &self.0;
        if digits.is_empty() {
            return None;
        }
        if digits == &[Nulla] {
            return Some(0);
        }
        let mut total = 0u32;
        let mut i = 0;
        while i < digits.len() {
            let v = digit_value(digits[i]);
            if i + 1 < digits.len() {
                let w = digit_value(digits[i + 1]);
                if v < w {
                    total = total.checked_add(w.checked_sub(v)?)?;
                    i += 2;
                    continue;
                }
            }
            total = total.checked_add(v)?;
            i += 1;
        }
        Some(total)
    }
}

impl From<u32> for RomanNumber {
    fn from(value: u32) -> Self {
        if value == 0 {
            return RomanNumber(vec![Nulla]);
        }

        let mut n = value;
        let mut out: Vec<RomanDigit> = Vec::new();

        let table: &[(u32, &[RomanDigit])] = &[
            (1000, &[M]),
            (900, &[C, M]),
            (500, &[D]),
            (400, &[C, D]),
            (100, &[C]),
            (90, &[X, C]),
            (50, &[L]),
            (40, &[X, L]),
            (10, &[X]),
            (9, &[I, X]),
            (5, &[V]),
            (4, &[I, V]),
            (1, &[I]),
        ];

        for (val, digits) in table {
            while n >= *val {
                for d in *digits {
                    out.push(*d);
                }
                n -= *val;
            }
        }

        RomanNumber(out)
    }
}

impl Iterator for RomanNumber {
    type Item = RomanNumber;

    fn next(&mut self) -> Option<Self::Item> {
        let n = self.to_u32()?;
        let next_n = n.checked_add(1)?;
        let new = RomanNumber::from(next_n);
        let out = new.clone();
        *self = new;
        Some(out)
    }
}
