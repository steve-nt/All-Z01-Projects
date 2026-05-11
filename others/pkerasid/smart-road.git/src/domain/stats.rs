/// Collected at runtime, displayed on exit.
#[derive(Debug, Clone)]
pub struct Statistics {
    pub vehicles_completed: u32,
    pub max_speed: f32,
    pub min_speed: f32,
    pub max_time_to_pass: f32,
    pub min_time_to_pass: f32,
    pub close_calls: u32,
    speed_samples: u32,
    completion_samples: u32,
}

impl Statistics {
    pub fn record_speed_sample(&mut self, speed: f32) {
        if self.speed_samples == 0 {
            self.max_speed = speed;
            self.min_speed = speed;
        } else {
            self.max_speed = self.max_speed.max(speed);
            self.min_speed = self.min_speed.min(speed);
        }
        self.speed_samples += 1;
    }

    pub fn record_completion(&mut self, time_to_pass: f32) {
        if self.completion_samples == 0 {
            self.max_time_to_pass = time_to_pass;
            self.min_time_to_pass = time_to_pass;
        } else {
            self.max_time_to_pass = self.max_time_to_pass.max(time_to_pass);
            self.min_time_to_pass = self.min_time_to_pass.min(time_to_pass);
        }

        self.vehicles_completed += 1;
        self.completion_samples += 1;
    }

    pub fn record_close_call(&mut self) {
        self.close_calls += 1;
    }
}

impl Default for Statistics {
    fn default() -> Self {
        Self {
            vehicles_completed: 0,
            max_speed: 0.0,
            min_speed: 0.0,
            max_time_to_pass: 0.0,
            min_time_to_pass: 0.0,
            close_calls: 0,
            speed_samples: 0,
            completion_samples: 0,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::Statistics;

    #[test]
    fn first_speed_sample_sets_both_extrema() {
        let mut stats = Statistics::default();
        stats.record_speed_sample(42.0);

        assert!((stats.max_speed - 42.0).abs() < f32::EPSILON);
        assert!((stats.min_speed - 42.0).abs() < f32::EPSILON);
    }

    #[test]
    fn completion_updates_count_and_time_bounds() {
        let mut stats = Statistics::default();
        stats.record_completion(3.5);
        stats.record_completion(1.5);

        assert_eq!(stats.vehicles_completed, 2);
        assert!((stats.max_time_to_pass - 3.5).abs() < f32::EPSILON);
        assert!((stats.min_time_to_pass - 1.5).abs() < f32::EPSILON);
    }
}
