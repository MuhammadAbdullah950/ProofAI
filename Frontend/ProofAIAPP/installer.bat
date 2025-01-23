# Step 1: Install Python and Add It to PATH

# Download Python installer
Invoke-WebRequest -Uri "https://www.python.org/ftp/python/3.10.5/python-3.10.5-amd64.exe" -OutFile "python-installer.exe"

# Install Python silently with pip and Add to PATH
Start-Process -FilePath "python-installer.exe" -ArgumentList "/quiet InstallAllUsers=1 PrependPath=1" -Wait

# Verify Python and pip installation
python --version
pip --version

# Step 2: Create a Virtual Environment

# Navigate to your project directory
cd C:\Users\Abdullah

# Create a virtual environment
python -m venv Env

# Activate the virtual environment
.\Env\Scripts\activate

# Step 3: Install Dependencies

# Inside the virtual environment, install dependencies from requirements.txt
pip install -r requirements.txt

# Confirm the installed packages
pip list

# Step 4: Run Your Python Script

# Run your script or application inside the virtual environment
python your_script.py

# Optional: Clean Up

# Deactivate the virtual environment after you're done
deactivate
