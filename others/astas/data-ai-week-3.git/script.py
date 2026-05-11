import pandas as pd
import numpy as np
import sys # for redirecting print output to a log file
from sklearn.compose import ColumnTransformer
from sklearn.preprocessing import StandardScaler, OneHotEncoder
from sklearn.impute import SimpleImputer
from sklearn.pipeline import Pipeline
from sklearn.model_selection import train_test_split
from sklearn.linear_model import LogisticRegression
from sklearn.tree import DecisionTreeClassifier
from sklearn.metrics import classification_report, confusion_matrix, accuracy_score
import time
import joblib
from sklearn.metrics import accuracy_score, precision_score, recall_score, confusion_matrix, classification_report

# -----------------------------
# Setup logging to file
# -----------------------------
log_file = "step1_data_inspection.log"
sys.stdout = open(log_file, "w")  # redirect all prints to log file

# -----------------------------
# 1) Load the dataset
# -----------------------------
data_path = "telecom_churn_dataset.csv"  # <-- replace with your dataset
df = pd.read_csv(data_path)
print("="*80)
print("STEP 1: Loaded dataset")
print("="*80, "\n")

# -----------------------------
# 2) Inspect dataset structure
# -----------------------------
print(">>> Dataset shape (rows, columns):", df.shape)
print("\n>>> Column data types:")
print(df.dtypes)
print("\n>>> First 5 rows:")
print(df.head())
print("\n>>> DataFrame info:")
df.info()
print("\n" + "-"*80 + "\n")

# -----------------------------
# 3) Separate numeric and categorical columns
# -----------------------------
numeric_cols = df.select_dtypes(include=[np.number]).columns.tolist()
categorical_cols = df.select_dtypes(exclude=[np.number]).columns.tolist()
print(f"Numeric columns: {numeric_cols}")
print(f"Categorical columns: {categorical_cols}")
print("\n" + "-"*80 + "\n")

# -----------------------------
# 4) Check for missing values and handle them
# -----------------------------

missing_counts = df.isnull().sum()

if missing_counts.any():
    print(">>> Missing values detected. Handling missing values...\n")
    print("Missing values per column:")
    print(missing_counts[missing_counts > 0])

    # Handle numeric columns -> fill missing with median
    numeric_missing = [col for col in numeric_cols if df[col].isnull().any()]
    if numeric_missing:
        df[numeric_cols] = df[numeric_cols].fillna(df[numeric_cols].median())
        print(f"\nNumeric columns imputed with median: {numeric_missing}")

    # Handle categorical columns -> fill missing with mode
    categorical_missing = [col for col in categorical_cols if df[col].isnull().any()]
    for col in categorical_missing:
        mode_val = df[col].mode(dropna=True).iloc[0] if not df[col].mode(dropna=True).empty else 'Unknown'
        df[col] = df[col].fillna(mode_val)
        print(f"Categorical column '{col}' imputed with mode: {mode_val}")

    print("\n" + "-"*80 + "\n")

else:
    print(">>> No missing values detected. Skipping imputation.\n" + "-"*80 + "\n")

# -----------------------------
# 5) Drop irrelevant identifier
# -----------------------------
id_col = 'customer_id'
if id_col in df.columns:
    df = df.drop(columns=[id_col])
    print(f"Identifier column '{id_col}' found and dropped.")
else:
    print(f"No identifier column '{id_col}' found, nothing dropped.")

print("\n" + "-"*80 + "\n")

# -----------------------------
# 6) Define target and features
# -----------------------------
target_col = 'churn'
if target_col not in df.columns:
    raise ValueError(f"Target column '{target_col}' not found.")

features = [col for col in df.columns if col != target_col]

print(f"Target variable: '{target_col}'")
print(f"Number of features: {len(features)}")
print("Feature columns:")
print(features)
print("\n" + "-"*80 + "\n")


# -----------------------------
# 7) Numeric summaries
# -----------------------------
numeric_cols_final = df.select_dtypes(include=[np.number]).columns.tolist()
if numeric_cols_final:
    print(">>> Numeric summary statistics:")
    print(df[numeric_cols_final].describe().transpose())
else:
    print("No numeric columns to summarize.")

print("\n" + "-"*80 + "\n")

# -----------------------------
# 8) Categorical summaries
# -----------------------------
categorical_cols_final = df.select_dtypes(exclude=[np.number]).columns.tolist()
if categorical_cols_final:
    print(">>> Categorical value counts (top 10 per column):")
    for col in categorical_cols_final:
        print(f"\n--- {col} ---")
        print(df[col].value_counts(dropna=False).head(10))
else:
    print("No categorical columns to summarize.")

print("\n" + "="*80)
print("STEP 1 completed successfully. All detailed logs saved in", log_file)
print("="*80, "\n")


# -----------------------------
# Optional: Churn distribution
# -----------------------------
print("\n>>> Churn distribution:")
churn_counts = df[target_col].value_counts()
churn_percent = df[target_col].value_counts(normalize=True) * 100
for label, count in churn_counts.items():
    percent = churn_percent[label]
    print(f"{label}: {count} customers ({percent:.2f}%)")

print("\n" + "-"*80 + "\n")

# -----------------------------
# Terminal output for user
# -----------------------------
# Close the log file so we can print a short summary to terminal
sys.stdout.close()
sys.stdout = sys.__stdout__

print("STEP 1 completed.")
print(f"Detailed logs are available in '{log_file}'.")
print(f"Cleaned dataset shape: {df.shape}")


# -----------------------------
# STEP 2: Feature Engineering
# -----------------------------
log_file_step2 = "step2_feature_engineering.log"
sys.stdout = open(log_file_step2, "w")  # redirect prints to Step 2 log file

print("\n" + "="*80)
print("STEP 2: Feature Engineering")
print("="*80 + "\n")

# 1) Age groups
age_bins = [0, 29, 44, 59, 120]  # bin edges
age_labels = ['young', 'adult', 'middle', 'senior']
df['age_group'] = pd.cut(df['age'], bins=age_bins, labels=age_labels, right=True)
print("Created 'age_group' feature based on age bins.\n")
print(">>> Age group distribution:")
print(df['age_group'].value_counts())
print("\n" + "-"*40 + "\n")

# 2) Heavy data users flag
median_data = df['avg_data_gb'].median()
df['heavy_data_user'] = (df['avg_data_gb'] > median_data).astype(int)
print(f"Created 'heavy_data_user' flag (1 if avg_data_gb > {median_data:.2f}).\n")
print(">>> Heavy data user distribution:")
print(df['heavy_data_user'].value_counts())
print("\n" + "-"*40 + "\n")

# 3) Interaction feature: support_tickets * unpaid_bills
df['tickets_x_bills'] = df['support_tickets'] * df['unpaid_bills']
print("Created interaction feature 'tickets_x_bills' (support_tickets × unpaid_bills).\n")
print(">>> tickets_x_bills summary statistics:")
print(df['tickets_x_bills'].describe())
print("\n" + "-"*40 + "\n")

# 4) Update feature list
features += ['age_group', 'heavy_data_user', 'tickets_x_bills']
print(f"Updated feature list to include engineered features ({len(features)} total features).")
print("Current features:")
print(features)
print("\n" + "-"*40 + "\n")

# 5) Show sample of the dataset with new features
print(">>> Sample of dataset with new features (first 5 rows):")
print(df.head()[features + [target_col]])
print("\n" + "="*80)
print("STEP 2 completed successfully. All logs saved in", log_file_step2)
print("="*80 + "\n")

# Restore terminal output
sys.stdout.close()
sys.stdout = sys.__stdout__

print("STEP 2 completed.")
print(f"Detailed logs are available in '{log_file_step2}'.")


# -----------------------------
# STEP 3: Preprocessing Pipeline
# -----------------------------
log_file_step3 = "step3_preprocessing_pipeline.log"
sys.stdout = open(log_file_step3, "w")  # redirect logs to Step 3 log file

print("\n" + "="*80)
print("STEP 3: Preprocessing Pipeline")
print("="*80 + "\n")

# 1) Split features and target
X = df[features]
y = df[target_col]
print("Split dataset into X (features) and y (target).")

# 2) Split into training and testing sets
X_train, X_test, y_train, y_test = train_test_split(
    X, y, test_size=0.3, random_state=42
)
print(f"Training set shape: {X_train.shape}")
print(f"Testing set shape: {X_test.shape}\n")

# 3) Identify numeric and categorical columns
numeric_features = X.select_dtypes(include=[np.number]).columns.tolist()
categorical_features = X.select_dtypes(exclude=[np.number]).columns.tolist()
print(f"Numeric features: {numeric_features}")
print(f"Categorical features: {categorical_features}\n")

# 4) Define preprocessing steps
numeric_transformer = Pipeline(steps=[
    ('imputer', SimpleImputer(strategy='median')),  # handle missing numeric values
    ('scaler', StandardScaler())                    # scale numeric features
])

categorical_transformer = Pipeline(steps=[
    ('imputer', SimpleImputer(strategy='most_frequent')),  # handle missing categorical
    ('onehot', OneHotEncoder(handle_unknown='ignore'))     # one-hot encode categories
])

# 5) Combine transformers into ColumnTransformer
preprocessor = ColumnTransformer(
    transformers=[
        ('num', numeric_transformer, numeric_features),
        ('cat', categorical_transformer, categorical_features)
    ]
)

# 6) Fit preprocessor on training data and transform both sets
X_train_processed = preprocessor.fit_transform(X_train)
X_test_processed = preprocessor.transform(X_test)
print("Preprocessing pipeline applied: numeric features scaled, categorical features one-hot encoded.\n")

# 7) Show shapes of processed data
print(f"Processed X_train shape: {X_train_processed.shape}")
print(f"Processed X_test shape: {X_test_processed.shape}")

print("\n" + "="*80)
print("STEP 3 completed successfully. All logs saved in", log_file_step3)
print("="*80 + "\n")

# Restore terminal output
sys.stdout.close()
sys.stdout = sys.__stdout__

print("STEP 3 completed.")
print(f"Detailed logs are available in '{log_file_step3}'.")


# -----------------------------
# STEP 4: Train/Test Split with Stratification
# -----------------------------
log_file_step4 = "step4_train_test_split.log"
sys.stdout = open(log_file_step4, "w")  # redirect logs to Step 4 log file

print("\n" + "="*80)
print("STEP 4: Train/Test Split with Stratification")
print("="*80 + "\n")

# Split the dataset with stratification on churn
X = df[features]
y = df[target_col]

X_train, X_test, y_train, y_test = train_test_split(
    X, y, 
    test_size=0.3, 
    random_state=42, 
    stratify=y  # ensure churn ratio is consistent across sets
)

print("Dataset successfully split with stratification on 'churn'.\n")
print(f"Training set shape: X_train={X_train.shape}, y_train={y_train.shape}")
print(f"Testing set shape:  X_test={X_test.shape}, y_test={y_test.shape}")
print("\n" + "-"*40 + "\n")

# Verify churn distribution consistency
train_dist = y_train.value_counts(normalize=True) * 100
test_dist = y_test.value_counts(normalize=True) * 100
overall_dist = y.value_counts(normalize=True) * 100

print(">>> Churn distribution (%):")
print(f"Overall dataset: {overall_dist.to_dict()}")
print(f"Training set:   {train_dist.to_dict()}")
print(f"Testing set:    {test_dist.to_dict()}")

# Check that proportions are similar
diff = (train_dist - test_dist).abs().max()
print(f"\nMaximum churn proportion difference between train and test: {diff:.4f}%")

print("\n" + "="*80)
print("STEP 4 completed successfully. All logs saved in", log_file_step4)
print("="*80 + "\n")

# Restore terminal output
sys.stdout.close()
sys.stdout = sys.__stdout__

print("STEP 4 completed.")
print(f"Detailed logs are available in '{log_file_step4}'.")


# -----------------------------
# STEP 5: Model training & evaluation (Logistic Regression and Decision Tree)
# -----------------------------

log_file_step5 = "step5_model_training.log"
sys.stdout = open(log_file_step5, "w")  # redirect detailed prints to Step 5 log file

print("\n" + "="*80)
print("STEP 5: Model training & evaluation")
print("="*80 + "\n")

# NOTE: we reuse the existing `preprocessor` ColumnTransformer defined earlier
# This ensures the SAME preprocessing is applied to both models.
print("Using the same preprocessor for both models (ColumnTransformer pipeline).")
print("Preprocessor details:")
print(preprocessor)    # prints transformer structure

# Build pipelines that chain preprocessor + estimator (ensures fit on train only)
lr_pipeline = Pipeline([
    ('preprocessor', preprocessor),
    ('classifier', LogisticRegression(max_iter=1000, random_state=42))
])

dt_pipeline = Pipeline([
    ('preprocessor', preprocessor),
    ('classifier', DecisionTreeClassifier(max_depth=5, random_state=42))
])

# Fit Logistic Regression
t0 = time.time()
lr_pipeline.fit(X_train, y_train)
t1 = time.time()
print(f"Trained Logistic Regression in {t1 - t0:.3f} seconds.")

# Evaluate Logistic Regression
y_pred_lr = lr_pipeline.predict(X_test)
print("\nLogistic Regression performance (on test set):")
print(f"Accuracy: {accuracy_score(y_test, y_pred_lr):.4f}")
print("Confusion matrix:")
print(confusion_matrix(y_test, y_pred_lr))
print("\nClassification report:")
print(classification_report(y_test, y_pred_lr))

print("\n" + "-"*80 + "\n")

# Fit Decision Tree
t0 = time.time()
dt_pipeline.fit(X_train, y_train)
t1 = time.time()
print(f"Trained Decision Tree in {t1 - t0:.3f} seconds.")

# Evaluate Decision Tree
y_pred_dt = dt_pipeline.predict(X_test)
print("\nDecision Tree performance (on test set):")
print(f"Accuracy: {accuracy_score(y_test, y_pred_dt):.4f}")
print("Confusion matrix:")
print(confusion_matrix(y_test, y_pred_dt))
print("\nClassification report:")
print(classification_report(y_test, y_pred_dt))

print("\n" + "="*80)
print("STEP 5 completed: trained and evaluated Logistic Regression and Decision Tree.")
print("Detailed results and diagnostics saved in", log_file_step5)
print("="*80 + "\n")

# Restore terminal output
sys.stdout.close()
sys.stdout = sys.__stdout__

# Short terminal summary
print("STEP 5 completed.")
print(f"Detailed logs are available in '{log_file_step5}'.")
print("Models trained: LogisticRegression, DecisionTreeClassifier. Evaluations saved to log.")



# -----------------------------
# STEP 6: Model Training & Evaluation
# -----------------------------


log_file_step5 = "step6_model_evaluation.log"
sys.stdout = open(log_file_step5, "w")  # redirect output to log file

print("\n" + "="*80)
print("STEP 6: Model Training & Evaluation")
print("="*80 + "\n")

# 1) Train Logistic Regression
lr_model = LogisticRegression(max_iter=1000, random_state=42)
lr_model.fit(X_train_processed, y_train)
y_pred_lr = lr_model.predict(X_test_processed)

print(">>> Logistic Regression model trained successfully.\n")

# 2) Train Decision Tree
dt_model = DecisionTreeClassifier(max_depth=5, random_state=42)
dt_model.fit(X_train_processed, y_train)
y_pred_dt = dt_model.predict(X_test_processed)

print(">>> Decision Tree model trained successfully.\n")

# 3) Evaluate both models
def evaluate_model(name, y_true, y_pred):
    print(f"\n{'-'*60}")
    print(f"Model: {name}")
    print("-"*60)
    print(f"Accuracy:  {accuracy_score(y_true, y_pred):.4f}")
    print(f"Precision: {precision_score(y_true, y_pred):.4f}")
    print(f"Recall:    {recall_score(y_true, y_pred):.4f}")
    print("\nConfusion Matrix:")
    print(confusion_matrix(y_true, y_pred))
    print("\nDetailed Classification Report:")
    print(classification_report(y_true, y_pred))

evaluate_model("Logistic Regression", y_test, y_pred_lr)
evaluate_model("Decision Tree", y_test, y_pred_dt)

# 4) Compare models (simple summary)
lr_acc = accuracy_score(y_test, y_pred_lr)
dt_acc = accuracy_score(y_test, y_pred_dt)

print("\n" + "="*80)
print("MODEL COMPARISON SUMMARY")
print("="*80)
print(f"Logistic Regression Accuracy: {lr_acc:.4f}")
print(f"Decision Tree Accuracy:       {dt_acc:.4f}")

if lr_acc > dt_acc:
    print("\n✅ Logistic Regression performed better overall.")
elif dt_acc > lr_acc:
    print("\n✅ Decision Tree performed better overall.")
else:
    print("\n🤝 Both models performed equally well.")

print("\n" + "="*80)
print("STEP 6 completed successfully. All logs saved in", log_file_step5)
print("="*80 + "\n")

# Optional — Save trained models
joblib.dump(lr_model, "logistic_regression_model.pkl")
joblib.dump(dt_model, "decision_tree_model.pkl")
print("\nTrained models have been saved for future use.")

# Restore terminal output
sys.stdout.close()
sys.stdout = sys.__stdout__

print("STEP 6 completed.")
print(f"Detailed logs are available in '{log_file_step5}'.")
print("Trained models have been saved as .pkl files.")
