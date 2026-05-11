use core::f64;

// use rand::Rng;
use rand::prelude::*;
use raster::{Color, Image};

#[derive(Debug, Clone, Copy)]
pub enum Colours {
    Red,
    Green,
    Blue,
    Yellow,
    Purple,
}

impl Colours {
    pub fn to_colour(&self) -> Color {
        match self {
            Colours::Red => Color {
                r: 255,
                g: 0,
                b: 0,
                a: 255,
            },
            Colours::Green => Color {
                r: 0,
                g: 255,
                b: 0,
                a: 255,
            },
            Colours::Blue => Color {
                r: 0,
                g: 0,
                b: 255,
                a: 255,
            },
            Colours::Yellow => Color {
                r: 255,
                g: 255,
                b: 0,
                a: 255,
            },
            Colours::Purple => Color {
                r: 128,
                g: 0,
                b: 128,
                a: 255,
            },
        }
    }

    pub fn randomize() -> Colours {
        let mut rndm_c = rand::rng();
        match rndm_c.random_range(0..5) {
            0 => Colours::Red,
            1 => Colours::Green,
            2 => Colours::Blue,
            3 => Colours::Yellow,
            4 => Colours::Purple,
            _ => {
                panic!("Pick an actual Colour bruh!")
            }
        }
    }
}

#[derive(Debug, Clone, Copy)]
pub struct Line {
    start: Point,
    end: Point,
    color: Colours,
} //two points

impl Line {
    //alex
    pub fn new(start: &Point, end: &Point) -> Line {
        //two points

        Line {
            start: *start,
            end: *end,
            color: Colours::randomize(),
        }
    }

    pub fn new_with_color(start: &Point, end: &Point, color: Colours) -> Line {
        Line {
            start: *start,
            end: *end,
            color,
        }
    }

    pub fn random(height: i32, width: i32) -> Line {
        let mut rng = rand::rng();
        let start_x = rng.random_range(0..width);
        let start_y = rng.random_range(0..height);
        let end_x = rng.random_range(0..width);
        let end_y = rng.random_range(0..height);
        Line::new(&Point::new(start_x, start_y), &Point::new(end_x, end_y))
    }
}

impl Drawable for Line {
    fn draw(&self, im: &mut Image) {
        let mut dist = self.start.dist(&self.end);
        let mut draw_vec = Vec2::new_p(&self.start);
        let mut travel_dir = Vec2::new_p(&self.end);
        travel_dir.sub(&Vec2::new_p(&self.start));
        travel_dir.norm();

        //to_point
        while dist >= 0.0 {
            let p = draw_vec.to_point();
            im.display(p.x, p.y, self.color());
            draw_vec.add(&travel_dir);
            dist -= 1.0;
        }
    }
    fn color(&self) -> Color {
        self.color.to_colour()
    }
}

// =========================================================
#[derive(Copy, Clone, Debug)]
pub struct Vec2 {
    x: f64,
    y: f64,
}

impl Vec2 {
    //new vector2d
    pub fn new(x: f64, y: f64) -> Vec2 {
        Vec2 { x: x, y: y }
    }

    //new vector2d using i32
    pub fn new_i(x: i32, y: i32) -> Vec2 {
        Vec2 {
            x: (x as f64),
            y: (y as f64),
        }
    }

    //new vector2d from Point
    pub fn new_p(p: &Point) -> Vec2 {
        Vec2 {
            x: (p.x as f64),
            y: (p.y as f64),
        }
    }

    //normalize
    pub fn norm(&mut self) {
        let dist = self.mag();
        self.x /= dist;
        self.y /= dist;
    }

    //multiply
    pub fn mult(&mut self, multiplier: f64) {
        self.x *= multiplier;
        self.y *= multiplier;
    }

    //add
    pub fn add(&mut self, v: &Vec2) {
        self.x += v.x;
        self.y += v.y;
    }

    //add two vecs and return a new one
    pub fn add_new(&self, v: Vec2) -> Vec2 {
        Vec2::new(self.x + v.x, self.y + v.y)
    }

    //subtract two vecs
    pub fn sub(&mut self, v: &Vec2) {
        self.x -= v.x;
        self.y -= v.y;
    }

    //subtract two vecs and return a new one
    pub fn sub_new(&self, v: &Vec2) -> Vec2 {
        Vec2::new(self.x - v.x, self.y - v.y)
    }

    //magnitude
    pub fn mag(&mut self) -> f64 {
        (self.x * self.x + self.y * self.y).sqrt().abs()
    }

    //rotate
    pub fn rot(&mut self, deg: f64) {
        let rads = deg.to_radians();
        let cos = rads.cos();
        let sin = rads.sin();
        let new_x = (cos * self.x) - (sin * self.y);
        self.y = (cos * self.y) + (sin * self.x);
        self.x = new_x;
    }

    pub fn to_point(&self) -> Point {
        Point::new(self.x.round() as i32, self.y.round() as i32)
    }
}

// =========================================================
#[derive(Debug, Clone, Copy)]
pub struct Circle {
    pub center: Point,
    pub radius: i32,
    color: Colours,
} //point, radius

impl Circle {
    //memo
    pub fn new(center: &Point, radius: i32) -> Circle {
        Circle {
            center: *center,
            radius: radius,
            color: Colours::randomize(),
        }
    }

    pub fn random(height: i32, width: i32) -> Circle {
        let mut a = rand::rng();
        let radius = a.random_range(50..250);
        Circle::new(&Point::random(height, width), radius)
    }
}

impl Drawable for Circle {
    fn draw(&self, im: &mut Image) {
        let v_cent = Vec2::new_p(&self.center);
        let edge_vec = Vec2::new_i(self.center.x, &self.center.y + self.radius);
        let mut rad_vec = edge_vec.sub_new(&v_cent);
        let per = f64::consts::PI * self.radius as f64 * 2.0;
        for _ in 0..=per as u32 {
            let new_edge = v_cent.add_new(rad_vec);
            let p = new_edge.to_point();
            im.display(p.x, p.y, self.color());
            rad_vec.rot(360.0 / per);
        }
    }

    fn color(&self) -> Color {
        self.color.to_colour()
    }
}

// =========================================================
#[derive(Debug, Clone, Copy)]
pub struct Point {
    pub x: i32,
    pub y: i32,
    color: Colours,
}

impl Point {
    pub fn new(x: i32, y: i32) -> Point {
        Point {
            x: x,
            y: y,
            color: Colours::randomize(),
        }
    }

    pub fn random(height: i32, width: i32) -> Point {
        let mut rng = rand::rng();
        let x = rng.random_range(0..width);
        let y = rng.random_range(0..height);
        Point::new(x, y)
    }

    pub fn dist(&self, p: &Point) -> f64 {
        let mut s = Vec2::new_p(&self);
        let t = Vec2::new_p(p);
        s.sub(&t);
        s.mag()
    }
}

impl Drawable for Point {
    fn draw(&self, im: &mut Image) {
        _ = im.display(self.x, self.y, self.color());
    }

    fn color(&self) -> Color {
        self.color.to_colour()
    }
}
// =========================================================

#[derive(Debug, Clone, Copy)]
pub struct Rectangle {
    top_left: Point,
    bottom_right: Point,
    color: Colours,
}

impl Rectangle {
    //theo

    pub fn new(top_left: &Point, bottom_right: &Point) -> Rectangle {
        let min_x = std::cmp::min(top_left.x, bottom_right.x);
        let max_x = std::cmp::max(top_left.x, bottom_right.x);
        let min_y = std::cmp::min(top_left.y, bottom_right.y);
        let max_y = std::cmp::max(top_left.y, bottom_right.y);
        Rectangle {
            top_left: Point::new(min_x, min_y),
            bottom_right: Point::new(max_x, max_y),
            color: Colours::randomize(),
        }
    }

    pub fn random(height: i32, width: i32) -> Rectangle {
        let mut rng = rand::rng();
        let x1 = rng.random_range(0..width);
        let y1 = rng.random_range(0..height);
        let x2 = rng.random_range(0..width);
        let y2 = rng.random_range(0..height);
        Rectangle::new(&Point::new(x1, y1), &Point::new(x2, y2))
    }

    pub fn top_left(&self) -> Point {
        self.top_left
    }

    pub fn bottom_right(&self) -> Point {
        self.bottom_right
    }
}


impl Drawable for Rectangle {
    /// Draws the rectangle as 4 lines.
    fn draw(&self, im: &mut Image) {
        let tl = self.top_left;
        let br = self.bottom_right;
        let tr = Point::new(br.x, tl.y);
        let bl = Point::new(tl.x, br.y);

        let color = self.color;
        Line::new_with_color(&tl, &tr, color).draw(im); // top
        Line::new_with_color(&tr, &br, color).draw(im); // right
        Line::new_with_color(&br, &bl, color).draw(im); // bottom
        Line::new_with_color(&bl, &tl, color).draw(im); // left
    }

    fn color(&self) -> Color {
        self.color.to_colour()
    }
}

// =========================================================
#[derive(Debug, Clone, Copy)]
pub struct Triangle {
    pub p1: Point,
    pub p2: Point,
    pub p3: Point,
    color: Colours,
}

impl Triangle {
    //memos
    //3 points
    pub fn new(p1: &Point, p2: &Point, p3: &Point) -> Triangle {
        Triangle {
            p1: *p1,
            p2: *p2,
            p3: *p3,
            color: Colours::randomize(),
        }
    }
}

impl Drawable for Triangle {
    fn draw(&self, im: &mut Image) {
        let color = self.color;
        let l1 = Line::new_with_color(&self.p1, &self.p2, color);
        let l2 = Line::new_with_color(&self.p2, &self.p3, color);
        let l3 = Line::new_with_color(&self.p3, &self.p1, color);

        l1.draw(im);
        l2.draw(im);
        l3.draw(im);
    }

    fn color(&self) -> Color {
        self.color.to_colour()
    }
}

pub trait Drawable {
    fn draw(&self, im: &mut Image);
    fn color(&self) -> Color;
}

pub trait Displayable {
    fn display(&mut self, x: i32, y: i32, color: Color);
}

#[cfg(test)]
mod tests {
    use crate::geometrical_shapes::*;
    use raster::Image;

    #[test]
    fn test_random_point() {
        const SIZE: i32 = 100;
        for _ in 0..10000 {
            let p = Point::random(SIZE, SIZE);
            assert!(p.x >= 0 && p.x <= SIZE && p.y >= 0 && p.y <= SIZE);
        }
    }

    #[test]
    fn test_random_line() {
        const SIZE: i32 = 100;
        for _ in 0..10000 {
            let l = Line::random(SIZE, SIZE);
            assert!(l.start.x >= 0 && l.start.x <= SIZE && l.start.y >= 0 && l.start.y <= SIZE);
            assert!(l.end.x >= 0 && l.end.x <= SIZE && l.end.y >= 0 && l.end.y <= SIZE);
        }
    }

    #[test]
    fn test_rectangle_new_and_draw() {
        let tl = Point::new(2, 3);
        let br = Point::new(8, 6);
        let rect = Rectangle::new(&tl, &br);
        let mut im = Image::blank(20, 20);
        rect.draw(&mut im);
    }

    #[test]
    fn test_rectangle_inverted_coordinates() {
        let rect = Rectangle::new(&Point::new(10, 8), &Point::new(2, 3));
        assert_eq!((rect.top_left().x, rect.top_left().y), (2, 3));
        assert_eq!((rect.bottom_right().x, rect.bottom_right().y), (10, 8));
    }

    // ================== Triangle Tests ==================

    #[test]
    fn test_triangle_new_and_draw() {
        let p1 = Point::new(1, 1);
        let p2 = Point::new(5, 1);
        let p3 = Point::new(3, 4);
        let tri = Triangle::new(&p1, &p2, &p3);
        let mut im = Image::blank(10, 10);
        tri.draw(&mut im); // draw without panic ???
    }

    #[test]
    fn test_triangle_collinear_points() {
        let p1 = Point::new(0, 0);
        let p2 = Point::new(2, 2);
        let p3 = Point::new(4, 4); 
        let tri = Triangle::new(&p1, &p2, &p3);
        let mut im = Image::blank(10, 10);
        tri.draw(&mut im); 
    }

    #[test]
    fn test_triangle_identical_points() {
        let p = Point::new(5, 5);
        let tri = Triangle::new(&p, &p, &p);
        let mut im = Image::blank(10, 10);
        tri.draw(&mut im); 
    }

    #[test]
    fn test_random_triangle_within_bounds() {
        const SIZE: i32 = 50;
        for _ in 0..1000 {
            let p1 = Point::random(SIZE, SIZE);
            let p2 = Point::random(SIZE, SIZE);
            let p3 = Point::random(SIZE, SIZE);
            let tri = Triangle::new(&p1, &p2, &p3);
            assert!(tri.p1.x >= 0 && tri.p1.x <= SIZE && tri.p1.y >= 0 && tri.p1.y <= SIZE);
            assert!(tri.p2.x >= 0 && tri.p2.x <= SIZE && tri.p2.y >= 0 && tri.p2.y <= SIZE);
            assert!(tri.p3.x >= 0 && tri.p3.x <= SIZE && tri.p3.y >= 0 && tri.p3.y <= SIZE);
        }
    }
}