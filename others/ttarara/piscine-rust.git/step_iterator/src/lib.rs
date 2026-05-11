use std::ops::Add;

pub struct StepIterator<T> {
    current: Option<T>,
    end: T,
    step: T,
}

impl<T: PartialOrd> StepIterator<T> {
    pub fn new(beg: T, end: T, step: T) -> Self {
        let current = if beg <= end { Some(beg) } else { None };
        StepIterator {
            current,
            end,
            step,
        }
    }
}

impl<T> Iterator for StepIterator<T>
where
    T: Copy + Add<Output = T> + PartialOrd,
{
    type Item = T;

    fn next(&mut self) -> Option<Self::Item> {
        let cur = self.current?;
        let nxt = cur + self.step;
        if nxt <= self.end {
            self.current = Some(nxt);
        } else {
            self.current = None;
        }
        Some(cur)
    }
}
