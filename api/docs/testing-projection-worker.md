Testing the projection job concurrency claims

Goal: Verify that only one instance can claim a queued job at a time and prevent duplicate processing when multiple instances of the app share the same database.

Pre-requisites
- Database engine that supports SKIP LOCKED.
  - PostgreSQL 9.5+ or MySQL 8.0+ are recommended. Ensure the app is configured accordingly (config.json or environment variables).
- Buildable application environment (Go 1.20+).

Test scenarios
1) Manual multi-instance test (recommended)
   This test simulates two application instances connected to the same DB.
   Steps:
   - Ensure DB is running and schema is set up (the app will set up tables on first run if configured to do so).
   - Seed 1 or more queued jobs:
     SQL (example):
       INSERT INTO projection_jobs (jobs_template_id, run_job_id, product_id, creation_date, run_time, run_type, run_date, run_name, run_description, status, status_error, total_points, points_done, shock_setting_id, aggregation_period, yield_curve_basis, yield_curve_month, mp_version, run_basis, user_name, user_email)
       VALUES (0, 0, 0, NOW(), 0, 0, TO_CHAR(NOW(), 'YYYY-MM-DD'), 'Test Run', 'Manual concurrency test', 'Queued', '', 0, 0, 0, 0, '', 0, '', '', 'tester', 'tester@example.com');
     Notes:
     - Use NOW() for MySQL/Postgres; adjust for MSSQL if needed.
     - Ensure status='Queued'.
   - Start Instance A of the app (terminal 1):
       go run ./main.go
   - Start Instance B of the app (terminal 2) pointing to the SAME DB configuration:
       go run ./main.go
   - Observe logs: only one instance should log that it claimed the job, e.g.:
       Claimed and processing queued job ID: <id>, Name: Test Run
     The other instance should not claim the same job; it will either report no queued jobs or that there is an in-progress job and return.
   - Optional: Seed multiple queued jobs and observe that instances will pick jobs one-by-one without duplicate claims.

2) Using the helper tool (two processes)
   A small helper tool is provided at tools/projection_worker_tester to drive the claim loop. Run two OS processes concurrently to simulate two instances.
   Steps:
   - Build the helper:
       go build -o ./bin/projection_worker_tester ./tools/projection_worker_tester
   - Seed a few queued jobs as above (3–5 jobs recommended).
   - Run two processes concurrently with different instance names:
       ./bin/projection_worker_tester -instance=A
       ./bin/projection_worker_tester -instance=B
   - Watch the output. Each process will attempt to claim periodically. You should see claims alternating without any job being claimed twice. If a race occurs, the guarded update will result in 0 rows affected and the tool will log "lost claim" and skip.

3) Observability via SQL
   While tests are running, you can inspect the job statuses with SQL:
     -- Currently in progress
     SELECT id, run_name, status FROM projection_jobs WHERE status='In Progress';
     -- Completed
     SELECT id, run_name, status FROM projection_jobs WHERE status='Complete';
     -- Still queued
     SELECT id, run_name, status FROM projection_jobs WHERE status='Queued';
   There should never be two rows with the same id claimed by different instances; the same id should not appear twice in logs as being claimed simultaneously.

4) Optional: Go test (e2e/manual DB)
   Because this feature requires actual DB-level row locking, a unit test with a mocked DB is insufficient. You can craft an end-to-end test backed by a real DB, but it requires environment configuration and is therefore not enabled by default. If you’d like, we can add an e2e test with a build tag that:
   - seeds jobs
   - spawns concurrent goroutines calling ProcessQueuedProjectionJobs
   - polls the DB to assert that at no point the in-progress count exceeds 1 and that all jobs are eventually processed exactly once.

Notes about environments
- PostgreSQL: The code uses FOR UPDATE SKIP LOCKED which is supported and recommended.
- MySQL 8+: SKIP LOCKED is supported when using InnoDB and appropriate transaction isolation; the raw SQL used includes SKIP LOCKED.
- MSSQL: SKIP LOCKED support differs (READPAST hints). The current code path is optimized for Postgres/MySQL; if using MSSQL we’d need to adjust the locking SQL.

Troubleshooting
- If you see errors about the SELECT ... FOR UPDATE SKIP LOCKED query, confirm the DB type in config.json is set correctly and the DB version supports SKIP LOCKED.
- If jobs appear to get stuck in "In Progress" (e.g., due to a crash), the startup path calls RecoverStalledJobs to re-queue them and clear partial outputs.
