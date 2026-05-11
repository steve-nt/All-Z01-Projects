# Telecom Churn Prediction — Student Objectives

## 1. Prepare the Data

- Load the dataset and inspect its structure (shape, column types, sample rows).
- Identify the target variable (`churn`) and features.
- Check for missing values and handle them appropriately (e.g., impute median for numeric, most frequent for categorical).
- Drop irrelevant identifiers (`customer_id`).
- Summarize categorical distributions (e.g., `plan_type`, `region`) and numeric stats.

---

## 2. Feature Engineering

- Create **age groups** (young, adult, middle, senior) using bins.
- Create a flag for **heavy data users** (above-median `avg_data_gb`).
- Generate an **interaction feature** (e.g., `support_tickets × unpaid_bills`).
- Update feature lists to include new engineered features.

---

## 3. Preprocessing

- Build a `ColumnTransformer` pipeline:
  - Scale numeric features.
  - One-hot encode categorical features.
  - Impute missing values (median for numeric, most frequent for categorical).
- Apply preprocessing consistently to training and testing sets.

---

## 4. Train/Test Split

- Split the dataset into **train/test sets** using `train_test_split` with stratification on `churn`.
- Verify the size and distribution of the splits.

---

## 5. Model Building

- Train a **Logistic Regression** model (with preprocessing pipeline).
- Train a **Decision Tree** model (with preprocessing pipeline).
- Ensure both models use the same preprocessed inputs.

---

## 6. Model Evaluation

- Generate predictions on the test set.
- Compute evaluation metrics:
  - Accuracy
  - Precision
  - Recall
  - Confusion Matrix
- Compare the two models and discuss which performs better.

---

Compare the precision, recall, and F1-score of both models.
*
Discuss whether a simple model (like Logistic Regression) or a tree-based model performs better.
*
Reflect on what additional features or data sources could improve prediction accuracy.