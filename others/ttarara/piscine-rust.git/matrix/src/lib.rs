use lalgebra_scalar::Scalar;

/// `W` = αριθμός στηλών, `H` = αριθμός γραμμών.
#[derive(Debug, Eq, PartialEq, Clone, Copy)]
pub struct Matrix<const W: usize, const H: usize, T>(pub [[T; W]; H]);

impl<const W: usize, const H: usize, T> Matrix<W, H, T>
where
	T: Scalar<Item = T> + Copy,
{
	pub fn zero() -> Self {
		Self([[T::zero(); W]; H])
	}
}

impl<const S: usize, T> Matrix<S, S, T>
where
	T: Scalar<Item = T> + Copy,
{
	pub fn identity() -> Self {
		let mut data = [[T::zero(); S]; S];
		for i in 0..S {
			data[i][i] = T::one();
		}
		Self(data)
	}
}
