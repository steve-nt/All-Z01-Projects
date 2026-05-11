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

impl From<u32> for RomanNumber {
	fn from(value: u32) -> Self {
		if value == 0 {
			return RomanNumber(vec![Nulla]);
		}

		let mut n = value;
		let mut out: Vec<RomanDigit> = Vec::new();

		// (value, digits) in descending subtractive order
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

