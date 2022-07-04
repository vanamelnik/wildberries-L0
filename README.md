# wildberries-L0
_Task L0 for Wildberries interns_
### Summary
Task L0 consists of 2 applications:
 - **orderpub** publishes orders in JSON format to *nats-streaming-server* from the provided file or from the console.
 - **orderserver** - listens *nats-streaming-server* (subject *orders*) and stores incoming orders to the Postgresql database using in-memory cache.

#### Run
```bash
scripts/start_postgres
scripts/start_nats
./orderserver
```