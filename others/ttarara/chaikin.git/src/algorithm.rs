use macroquad::prelude::*;

// Builds step 0 through step max_steps.
// Step 0 is the original control polygon.
pub fn build_steps(control_points: &[Vec2], max_steps: usize) -> Vec<Vec<Vec2>> {
    let mut steps = Vec::new();
    steps.push(control_points.to_vec());

    let mut current = control_points.to_vec();

    for _ in 0..max_steps {
        current = chaikin_step(&current);
        steps.push(current.clone());
    }

    steps
}

// Applies one iteration of the Chaikin algorithm.
// For each pair of adjacent points, creates two new points:
// - One at 1/4 of the distance from the first point
// - One at 3/4 of the distance from the first point
// Preserves the first and last points to keep the curveδ anchored
fn chaikin_step(points: &[Vec2]) -> Vec<Vec2> {
    // Need at least 2 points to create a line segment
    if points.len() < 2 {
        return points.to_vec();
    }

    let mut new_points = Vec::new();

    // Always keep the first point
    new_points.push(points[0]);

    // For each segment between adjacent points
    for segment_index in 0..points.len() - 1 {
        let start_point = points[segment_index];
        let end_point = points[segment_index + 1];

        // First cut point at 1/4 of the way from start to end
        let quarter_point = start_point * 0.75 + end_point * 0.25;

        // Second cut point at 3/4 of the way from start to end
        let three_quarter_point = start_point * 0.25 + end_point * 0.75;

        new_points.push(quarter_point);
        new_points.push(three_quarter_point);
    }

    // Always keep the last point
    new_points.push(points[points.len() - 1]);

    new_points
}

#[cfg(test)]
mod tests {
    use super::*;
    use macroquad::math::vec2;

    /// Tolerance for float equality in tests.
    fn assert_vec2_close(a: Vec2, b: Vec2) {
        const EPS: f32 = 1e-5;
        assert!(
            (a.x - b.x).abs() < EPS && (a.y - b.y).abs() < EPS,
            "expected {:?} ≈ {:?}",
            a,
            b
        );
    }

    /// Verifies that after one Chaikin iteration, new points lie at 1/4 and 3/4 along each edge
    /// (corner cuts at 25% and 75% of the segment from the first point toward the second),
    /// and that endpoints are preserved for a two-point input.
    #[test]
    fn one_iteration_cut_points_at_quarter_and_three_quarter_along_segment() {
        let p0 = vec2(0.0, 0.0);
        let p1 = vec2(4.0, 0.0);
        let out = chaikin_step(&[p0, p1]);

        assert_eq!(out.len(), 4);
        assert_vec2_close(out[0], p0);
        // 1/4 along p0→p1: 0.75*p0 + 0.25*p1
        assert_vec2_close(out[1], vec2(1.0, 0.0));
        // 3/4 along p0→p1: 0.25*p0 + 0.75*p1
        assert_vec2_close(out[2], vec2(3.0, 0.0));
        assert_vec2_close(out[3], p1);
    }

    /// Verifies that the number of points grows as expected per iteration for an open polyline
    /// with preserved first and last control points: N input points (N ≥ 2) yield 2N points
    /// after one step (not 2N−2, which applies to variants that omit the original endpoints).
    #[test]
    fn one_iteration_point_count_open_curve_with_preserved_endpoints() {
        // Course notes often give 2N−2 when the new polyline does not repeat the original
        // endpoints as vertices. This implementation keeps the first and last input points,
        // so for N ≥ 2 the length is 2N.
        for n in 2..=12 {
            let points: Vec<Vec2> = (0..n)
                .map(|i| vec2(i as f32, (i * i) as f32 * 0.1))
                .collect();
            let out = chaikin_step(&points);
            assert_eq!(
                out.len(),
                2 * n,
                "N={}: expected 2N points with anchored endpoints",
                n
            );
        }
    }

    /// Ensures the algorithm does not panic on degenerate inputs: empty list, a single point,
    /// or two points; checks sensible outputs (empty, unchanged singleton, four-point polyline).
    #[test]
    fn edge_cases_zero_one_two_points_no_panic() {
        let empty: Vec<Vec2> = Vec::new();
        assert!(chaikin_step(&empty).is_empty());

        let one = [vec2(10.0, -3.0)];
        let o1 = chaikin_step(&one);
        assert_eq!(o1.len(), 1);
        assert_vec2_close(o1[0], one[0]);

        let two = [vec2(0.0, 0.0), vec2(8.0, 6.0)];
        let o2 = chaikin_step(&two);
        assert_eq!(o2.len(), 4);
        assert_vec2_close(o2[0], two[0]);
        assert_vec2_close(o2[3], two[1]);
    }
}
