use std::ops::{Add, Mul, Sub};

use lalgebra_scalar::Scalar;
use matrix::Matrix;

#[derive(Debug, Eq, PartialEq, Clone, Copy)]
pub struct Wrapper<const W: usize, const H: usize, T>(pub Matrix<W, H, T>);

impl<const W: usize, const H: usize, T> From<[[T; W]; H]> for Wrapper<W, H, T>
where
    T: Scalar<Item = T> + Copy,
{
    fn from(value: [[T; W]; H]) -> Self {
        Self(Matrix(value))
    }
}

impl<const W: usize, const H: usize, T> Add for Wrapper<W, H, T>
where
    T: Scalar<Item = T> + Copy,
{
    type Output = Self;

    fn add(self, rhs: Self) -> Self::Output {
        let mut data = [[T::zero(); W]; H];
        for (r, row) in data.iter_mut().enumerate() {
            for (c, cell) in row.iter_mut().enumerate() {
                *cell = self.0 .0[r][c] + rhs.0 .0[r][c];
            }
        }
        Self(Matrix(data))
    }
}

impl<const W: usize, const H: usize, T> Sub for Wrapper<W, H, T>
where
    T: Scalar<Item = T> + Copy,
{
    type Output = Self;

    fn sub(self, rhs: Self) -> Self::Output {
        let mut data = [[T::zero(); W]; H];
        for (r, row) in data.iter_mut().enumerate() {
            for (c, cell) in row.iter_mut().enumerate() {
                *cell = self.0 .0[r][c] - rhs.0 .0[r][c];
            }
        }
        Self(Matrix(data))
    }
}

impl<const S: usize, T> Mul for Wrapper<S, S, T>
where
    T: Scalar<Item = T> + Copy,
{
    type Output = Self;

    fn mul(self, rhs: Self) -> Self::Output {
        let mut data = [[T::zero(); S]; S];

        for (r, row) in data.iter_mut().enumerate() {
            for (c, cell) in row.iter_mut().enumerate() {
                let mut acc = T::zero();
                for k in 0..S {
                    acc = acc + self.0 .0[r][k] * rhs.0 .0[k][c];
                }
                *cell = acc;
            }
        }

        Self(Matrix(data))
    }
}
