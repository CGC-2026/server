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

4. **Database Management Tool**
   PgAdmin is Recommended: [https://www.pgadmin.org/](https://www.pgadmin.org/)
   Or the Visual studio extension managed by Microsoft: [PostGreSQL](https://marketplace.visualstudio.com/items?itemName=ms-ossdata.vscode-pgsql)

---

## Development Setup

1. **Create a virtual environment**

   - macOS/Linux: `python3 -m venv .venv`
   - Windows (cmd): `python -m venv .venv`

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
   Mac/Unix: `cp .env.example .env`
   Windows: `copy .env.example .env`

6. **Start PostgreSQL (via Docker)**
   **\_Note:** It is reccomended to run the docker containers in a standalone terminal, as the VSCode terminal may cause performance issues\_
   From the root directory:

   ```bash
   docker-compose up -d
   ```

7. **Run the server**
   From the root directory:

   ```bash
   fastapi dev app/main.py
   ```

---

## Connecting DBMS to PostGresDb

1. Ensure docker container is running
2. In your Database management tool, select add new server
3. Fill in settings from your .env file, default settings are as follows:
   - Host: localhost
   - Port: 5432
   - database: mydatabase
   - username: myuser
   - password: mypassword
4. Save settings and test connection, you should now be able to view and manage the database through your tool of choice.

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

## Running Tests

We use [pytest](https://docs.pytest.org/).

**Run all tests**
pytest

**Verbose output**
pytest -v

**Run a specific file or test**
pytest app/tests/test_main.py
pytest -k "test_read_root"

**Coverage**
pytest --cov=app

---

## Documentation

- **FastAPI**: [https://fastapi.tiangolo.com](https://fastapi.tiangolo.com)
- **Alembic**: [https://alembic.sqlalchemy.org/](https://alembic.sqlalchemy.org/)

---
