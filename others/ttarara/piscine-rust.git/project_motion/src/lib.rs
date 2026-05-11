const G: f32 = 9.8;

fn round_tenth(x: f32) -> f32 {
    (x * 10.0).round() / 10.0
}

#[derive(Debug, Clone, PartialEq)]
pub struct Object {
    pub x: f32,
    pub y: f32,
}

#[derive(Debug, Clone, PartialEq)]
pub struct ThrowObject {
    pub init_position: Object,
    pub init_velocity: Object,
    pub actual_position: Object,
    pub actual_velocity: Object,
    pub time: f32,
}

impl ThrowObject {
    pub fn new(init_position: Object, init_velocity: Object) -> ThrowObject {
        ThrowObject {
            init_position: init_position.clone(),
            init_velocity: init_velocity.clone(),
            actual_position: init_position,
            actual_velocity: init_velocity,
            time: 0.0,
        }
    }
}

impl Iterator for ThrowObject {
    type Item = ThrowObject;

    fn next(&mut self) -> Option<Self::Item> {
        let t = self.time + 1.0;
        let x0 = self.init_position.x;
        let y0 = self.init_position.y;
        let vx0 = self.init_velocity.x;
        let vy0 = self.init_velocity.y;

        let x = x0 + vx0 * t;
        let y = y0 + vy0 * t - 0.5 * G * t * t;
        let vx = vx0;
        let vy = vy0 - G * t;

        if y <= 0.0 {
            return None;
        }

        self.actual_position = Object {
            x: round_tenth(x),
            y: round_tenth(y),
        };
        self.actual_velocity = Object {
            x: round_tenth(vx),
            y: round_tenth(vy),
        };
        self.time = t;
        Some(self.clone())
    }
}
