# Server
CGC's main FastAPI service

### Prerequistes

1. Download Python (v3.13+) for your OS: https://www.python.org/downloads/
2. Download the VSCode Python Extension: https://marketplace.cursorapi.com/items/?itemName=ms-python.python
3. Download Docker Desktop: https://www.docker.com/products/docker-desktop/ (make sure the docker daemon is running aka docker desktop is open)

## Development Set up

1. Create python virtual enviroment.
    - MacOS: `python3 -m venv .venv`
    - Windows: `python -m venv .venv`
2. Activate the virtual enviroment: `source .venv/bin/activate`
3. For VSCode set the interepreter path: https://code.visualstudio.com/docs/python/environments#_select-and-activate-an-environment
3. Install dependencies: `pip install -r requirements.txt`
4. Copy .env values: `touch .env && cp .env.example .env`
5. Spin up PostgreSQL Docker Container:  `docker-compose up -d`
6. Run the server `fastapi dev app/main.py`


## Adding Packages:


## Documentation:
 1. FastAPI: https://fastapi.tiangolo.com/