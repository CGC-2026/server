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

## Database Migrations

We use [Alembic](https://alembic.sqlalchemy.org/) for database schema management.

### Quick Commands

**Generate migration automatically (when you add/modify models):**
```bash
alembic revision --autogenerate -m "Description of changes"
```

**Generate empty migration (if autogenerate fails):**
```bash
alembic revision -m "Description of changes"
```

**Apply migrations to database:**
```bash
alembic upgrade head
```

**Check migration status:**
```bash
alembic current
```

**Rollback last migration:**
```bash
alembic downgrade -1
```

### Workflow

1. **Add or modify models** in `app/models/`
2. **Import your models** in `alembic/env.py` (crucial step!)
   ```python
   # Import all your models here so they are registered with Base.metadata
   from app.models.user import User
   from app.models.your_new_model import YourNewModel  # Add new models here
   ```
3. **Generate migration** with `alembic revision --autogenerate -m "Description"`
4. **Review the generated migration** in `alembic/versions/`
5. **Apply to database** with `alembic upgrade head`

### Troubleshooting

If `--autogenerate` fails with database connection errors:

1. **Create empty migration:**
   ```bash
   alembic revision -m "Description of changes"
   ```

2. **Manually edit** the migration file in `alembic/versions/`:
   ```python
   def upgrade() -> None:
       op.create_table('your_table',
           sa.Column('id', sa.Integer(), nullable=False),
           # Add your columns here
           sa.PrimaryKeyConstraint('id')
       )

   def downgrade() -> None:
       op.drop_table('your_table')
   ```

### Important Notes

- **Models must be imported** in `alembic/env.py` to be detected
- **All models must inherit from `Base`** (imported from `app.db.base`)

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

---
