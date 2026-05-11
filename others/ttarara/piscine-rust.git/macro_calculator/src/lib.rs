use json::object;

pub struct Food {
    pub name: String,
    pub calories: (String, String), // (kJ, kcal)
    pub fats: f64,
    pub carbs: f64,
    pub proteins: f64,
    pub nbr_of_portions: f64,
}

fn round2(x: f64) -> f64 {
    (x * 100.0).round() / 100.0
}

pub fn calculate_macros(foods: &[Food]) -> json::JsonValue {
    let mut total_cals = 0.0;
    let mut total_carbs = 0.0;
    let mut total_proteins = 0.0;
    let mut total_fats = 0.0;

    for food in foods {
        let portions = food.nbr_of_portions;

        // Parse kcal string like "510kcal" or "358.65kcal"
        let kcal_str = food.calories.1.trim_end_matches("kcal");
        let kcal_val: f64 = kcal_str.parse().unwrap_or(0.0);

        total_cals += kcal_val * portions;
        total_carbs += food.carbs * portions;
        total_proteins += food.proteins * portions;
        total_fats += food.fats * portions;
    }

    let total_cals = round2(total_cals);
    let total_carbs = round2(total_carbs);
    let total_proteins = round2(total_proteins);
    let total_fats = round2(total_fats);

    object! {
        "cals" => total_cals,
        "carbs" => total_carbs,
        "proteins" => total_proteins,
        "fats" => total_fats
    }
}