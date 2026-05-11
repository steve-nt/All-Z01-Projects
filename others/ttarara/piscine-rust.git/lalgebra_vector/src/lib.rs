use std::ops::Add;

use lalgebra_scalar::Scalar;

#[derive(Debug, PartialEq)]
pub struct Vector<T: Scalar>(pub Vec<T>);

impl<T> Add for Vector<T>
where
	T: Scalar<Item = T>,
{
	type Output = Option<Self>;

	fn add(self, rhs: Self) -> Self::Output {
		if self.0.len() != rhs.0.len() {
			return None;
		}
		Some(Vector(
			self.0
				.into_iter()
				.zip(rhs.0)
				.map(|(a, b)| a + b)
				.collect(),
		))
	}
}

impl<T> Vector<T>
where
	T: Scalar<Item = T>,
{
	pub fn dot(self, rhs: Self) -> Option<T> {
		if self.0.len() != rhs.0.len() {
			return None;
		}
		let acc = self
			.0
			.into_iter()
			.zip(rhs.0.into_iter())
			.fold(T::zero(), |sum, (a, b)| sum + a * b);
		Some(acc)
	}
}

