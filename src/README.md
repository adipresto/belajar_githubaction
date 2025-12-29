# go-fitbyte
Go app for ProjectSprint project

# How to Run
1. Clone the repo

   ```bash
    git clone https://github.com/arieffadhlan/go-fitbyte.git
    cd go-fitbyte
   ```

2. Create `.env` file

  Can copy from `.env-example` but adjust the value
   ```bash
    cp .env-example .env
   ```

3. Create database `go-fitbyte`

4. Run the migration

  ```bash
    make migrate-up
  ```

5. Run the app

  ```bash
  make run
  ```

6. Health check

```bash
curl http://localhost:{APP_PORT}/health-check
```
