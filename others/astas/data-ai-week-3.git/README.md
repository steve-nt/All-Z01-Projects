Use a Virtual Environment (Recommended)

This is the cleanest way, and avoids breaking your system Python.

Create a virtual environment:

python3 -m venv venv


This creates a folder venv in your project.

Activate it:

source venv/bin/activate


Your shell prompt should now show (venv).

Install packages inside the venv:

pip install pandas scikit-learn matplotlib seaborn
pip install scikit-learn pandas numpy



Run your scripts with the venv Python:

python script.py


All packages are isolated and won’t interfere with system Python.