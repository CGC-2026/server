# 🚀 CGC Server

CGC's main **FastAPI** service.

---

## Prerequisites

1. **Python (v3.13+)**  
   Download: [python.org/downloads](https://www.python.org/downloads/)

2. **VSCode Python Extension**  
   Install: [marketplace.cursorapi.com](https://marketplace.cursorapi.com/items/?itemName=ms-python.python)

3. **Docker Desktop**  
   Download: [docker.com/products/docker-desktop](https://www.docker.com/products/docker-desktop)  
   → Make sure Docker is running in the background

---

## Development Setup

1. **Create a virtual environment**
   - macOS: `python3 -m venv .venv`
   - Windows: `python -m venv .venv`

2. **Activate the environment**
   - macOS/Linux: `source .venv/bin/activate`
   - Windows (cmd): `.venv\Scripts\activate`

3. **Set Python interpreter in VSCode**  
   Guide: [Select and Activate Environment](https://code.visualstudio.com/docs/python/environments#_select-and-activate-an-environment)

4. **Install dependencies**
   ```bash
   pip install -r requirements.txt
   ```

5. **Create environment file**
   ```bash
   cp .env.example .env
   ```

6. **Start PostgreSQL (via Docker)**
   ```bash
   docker-compose up -d
   ```

7. **Run the server**
   ```bash
   fastapi dev app/main.py
   ```

---

## Adding Packages

To add a new Python package:

1. Install the package:
   ```bash
   pip install package-name
   ```

2. Save to `requirements.txt`:
   ```bash
   pip freeze > requirements.txt
   ```

3. Commit the updated file to version control.

---

## Documentation

- **FastAPI**: [https://fastapi.tiangolo.com](https://fastapi.tiangolo.com)

---
