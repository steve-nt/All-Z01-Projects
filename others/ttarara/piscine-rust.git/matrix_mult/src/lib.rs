use std::ops::{Add, Mul};

#[derive(Debug, PartialEq, Eq, Clone)]
pub struct Matrix<T>(pub Vec<Vec<T>>);

impl<T> Matrix<T>
where
    T: Clone,
{
    pub fn number_of_cols(&self) -> usize {
        self.0.first().map_or(0, Vec::len)
    }

    pub fn number_of_rows(&self) -> usize {
        self.0.len()
    }

    pub fn row(&self, n: usize) -> Vec<T> {
        self.0.get(n).cloned().unwrap_or_default()
    }

    pub fn col(&self, n: usize) -> Vec<T> {
        self.0
            .iter()
            .filter_map(|row| row.get(n).cloned())
            .collect()
    }
}

impl<T> Mul for Matrix<T>
where
    T: Copy + Clone + Default + Add<Output = T> + Mul<Output = T>,
{
    type Output = Option<Self>;

    fn mul(self, rhs: Self) -> Self::Output {
        let left_rows = self.number_of_rows();
        let left_cols = self.number_of_cols();
        let right_rows = rhs.number_of_rows();
        let right_cols = rhs.number_of_cols();

        if left_cols == 0 || right_rows == 0 || left_cols != right_rows {
            return None;
        }

        if self.0.iter().any(|row| row.len() != left_cols)
            || rhs.0.iter().any(|row| row.len() != right_cols)
        {
            return None;
        }

        let mut out = vec![vec![T::default(); right_cols]; left_rows];
        for (r, out_row) in out.iter_mut().enumerate() {
            for (c, out_cell) in out_row.iter_mut().enumerate() {
                let mut acc = T::default();
                for k in 0..left_cols {
                    acc = acc + self.0[r][k] * rhs.0[k][c];
                }
                *out_cell = acc;
            }
        }

        Some(Matrix(out))
    }
}
