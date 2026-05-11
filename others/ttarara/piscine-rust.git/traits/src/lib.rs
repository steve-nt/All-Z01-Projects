use std::fmt;

#[derive(Debug)]
pub struct Player<'a> {
	pub name: &'a str,
	pub strength: f64,
	pub score: u32,
	pub money: u32,
	pub weapons: Vec<&'a str>,
}

pub struct Fruit {
	pub weight_in_kg: f64,
}

pub struct Meat {
	pub weight_in_kg: f64,
	pub fat_content: f64,
}

pub trait Food {
	fn gives(&self) -> f64;
}

impl<'a> Player<'a> {
	pub fn eat(&mut self, food: impl Food) {
		self.strength += food.gives();
	}
}

impl fmt::Display for Player<'_> {
	fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
		// Match the grader's exact string format: no trailing newline at the end.
		writeln!(f, "{}", self.name)?;
		writeln!(
			f,
			"Strength: {}, Score: {}, Money: {}",
			self.strength, self.score, self.money
		)?;
		write!(f, "Weapons: {:?}", self.weapons)
	}
}

impl Food for Fruit {
	fn gives(&self) -> f64 {
		self.weight_in_kg * 4.0
	}
}

impl Food for Meat {
	fn gives(&self) -> f64 {
		let fat_kg = self.weight_in_kg * self.fat_content;
		let protein_kg = self.weight_in_kg - fat_kg;
		protein_kg * 4.0 + fat_kg * 9.0
	}
}

