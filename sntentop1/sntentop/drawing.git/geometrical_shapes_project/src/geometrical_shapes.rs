use raster::{Color, Image};
use rand::Rng;

pub struct Point {
    x: i32,
    y: i32,
}

  impl Point {
      pub fn new(x: i32, y: i32) -> Point {
          Point { x, y }
      }

      pub fn random(width: i32, height: i32) -> Point {
          let mut rng = rand::thread_rng();
          Point {
              x: rng.gen_range(0, width),
              y: rng.gen_range(0, height),
          }
      }
  }


pub struct Line {
    p1: Point,
    p2: Point,
}


 impl Line {
      pub fn new(p1: &Point, p2: &Point) -> Line {
          Line {
              p1: Point { x: p1.x, y: p1.y },
              p2: Point { x: p2.x, y: p2.y },
          }
      }

      pub fn random(width: i32, height: i32) -> Line {
          Line::new(
              &Point::random(width, height),
              &Point::random(width, height),
          )
      }
  }

pub struct Triangle {
    p1: Point,
    p2: Point,
    p3: Point,
}

 impl Triangle {                                                                                                                                          
      pub fn new(p1: &Point, p2: &Point,p3: &Point) -> Triangle {                                                                                                   
          Triangle {
              p1: Point { x: p1.x, y: p1.y },
              p2: Point { x: p2.x, y: p2.y },
              p3: Point { x: p3.x, y: p3.y },
          }                                                                                                                                            
      }
  } 



pub struct Rectangle {
    p1: Point,
    p2: Point,
}

                                                                                                                                                     
  impl Rectangle {                                                                                                                                     
      pub fn new(p1: &Point, p2: &Point) -> Rectangle {                                                                                                
          Rectangle {                                                                                                                                  
              p1: Point { x: p1.x, y: p1.y },                                                                                                        
              p2: Point { x: p2.x, y: p2.y },                                                                                                          
          }
      }                                                                                                                                                
  } 

pub struct Circle {
    center: Point,
    radius: i32,
}

  impl Circle {
      pub fn new(center: &Point, radius: i32) -> Circle {
          Circle {
              center: Point { x: center.x, y: center.y },
              radius,
          }
      }

      pub fn random(width: i32, height: i32) -> Circle {
          let mut rng = rand::thread_rng();
          Circle::new(
              &Point::random(width, height),
              rng.gen_range(10, 200),
          )
      }
  }

pub trait Drawable {
    fn draw(&self, image: &mut Image);
    fn color(&self) -> Color;
}

pub trait Displayable {
    fn display(&mut self, x: i32, y: i32, color: Color);
}

fn draw_line(image: &mut Image, p1: &Point, p2: &Point, color: Color) {
    let mut x0 = p1.x;
    let mut y0 = p1.y;
    let x1 = p2.x;
    let y1 = p2.y;

    let dx = (x1 - x0).abs();
    let dy = (y1 - y0).abs();
    let sx = if x0 < x1 { 1 } else { -1 };
    let sy = if y0 < y1 { 1 } else { -1 };
    let mut err = dx - dy;

    loop {
        image.display(x0, y0, color.clone());
        if x0 == x1 && y0 == y1 {
            break;
        }
        let e2 = 2 * err;
        if e2 > -dy {
            err -= dy;
            x0 += sx;
        }
        if e2 < dx {
            err += dx;
            y0 += sy;
        }
    }
}

impl Drawable for Point {
    fn color(&self) -> Color {
        Color::rgb(255, 0, 0)
    }

    fn draw(&self, image: &mut Image) {
        image.display(self.x, self.y, self.color());
    }
}

impl Drawable for Line {
    fn color(&self) -> Color {
        Color::rgb(0, 255, 0)
    }

    fn draw(&self, image: &mut Image) {
        draw_line(image, &self.p1, &self.p2, self.color());
    }
}

impl Drawable for Triangle {
    fn color(&self) -> Color {
        Color::rgb(0, 0, 255)
    }

    fn draw(&self, image: &mut Image) {
        let color = self.color();
        draw_line(image, &self.p1, &self.p2, color.clone());
        draw_line(image, &self.p2, &self.p3, color.clone());
        draw_line(image, &self.p3, &self.p1, color.clone());
    }
}

impl Drawable for Rectangle {
    fn color(&self) -> Color {
        Color::rgb(255, 165, 0)
    }

    fn draw(&self, image: &mut Image) {
        let color = self.color();
        let top_right = Point { x: self.p2.x, y: self.p1.y };
        let bottom_left = Point { x: self.p1.x, y: self.p2.y };
        draw_line(image, &self.p1, &top_right, color.clone());
        draw_line(image, &top_right, &self.p2, color.clone());
        draw_line(image, &self.p2, &bottom_left, color.clone());
        draw_line(image, &bottom_left, &self.p1, color.clone());
    }
}

impl Drawable for Circle {
    fn color(&self) -> Color {
        Color::rgb(255, 0, 255)
    }

    fn draw(&self, image: &mut Image) {
        let color = self.color();
        let mut x = self.radius;
        let mut y = 0;
        let mut err = 0;
        let cx = self.center.x;
        let cy = self.center.y;

        while x >= y {
            image.display(cx + x, cy + y, color.clone());
            image.display(cx + y, cy + x, color.clone());
            image.display(cx - y, cy + x, color.clone());
            image.display(cx - x, cy + y, color.clone());
            image.display(cx - x, cy - y, color.clone());
            image.display(cx - y, cy - x, color.clone());
            image.display(cx + y, cy - x, color.clone());
            image.display(cx + x, cy - y, color.clone());

            y += 1;
            if err <= 0 {
                err += 2 * y + 1;
            } else {
                x -= 1;
                err += 2 * (y - x) + 1;
            }
        }
    }
}