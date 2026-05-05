# On the build host — produce a fresh binary.
cd api
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o aart_api .

cd ..

# Stage to a new release directory on the server (no sudo — moruogrsadmin owns /opt/aart-gr).
RELEASE=$(date -u +%Y-%m-%dT%H%MZ)
ssh -p 2222 moruogrsadmin@102.133.230.104 "mkdir -p /opt/aart-gr/releases/$RELEASE"

# 1. Binary.
rsync -av -e "ssh -p 2222" \
  api/aart_api \
  moruogrsadmin@102.133.230.104:/opt/aart-gr/releases/$RELEASE/aart_api

# 2. Migrations. Without this, the runner finds no migrations directory, logs a Warn,
#    and applies nothing — schema drift with no error. Always ship.
rsync -av --delete -e "ssh -p 2222" \
  api/migrations/ \
  moruogrsadmin@102.133.230.104:/opt/aart-gr/releases/$RELEASE/migrations/

# 3. Permissions, config symlink, writable-state symlinks (logs, tmp/uploads).
ssh -p 2222 moruogrsadmin@102.133.230.104 "
  chmod 0750 /opt/aart-gr/releases/$RELEASE/aart_api &&
  ln -sf  /opt/aart-gr/shared/config/config.json /opt/aart-gr/releases/$RELEASE/config.json &&
  ln -sfn /opt/aart-gr/shared/logs              /opt/aart-gr/releases/$RELEASE/logs &&
  mkdir -p /opt/aart-gr/releases/$RELEASE/tmp &&
  ln -sfn /opt/aart-gr/shared/uploads           /opt/aart-gr/releases/$RELEASE/tmp/uploads
"

# 4. Flip the symlink atomically. The new release is now live for the next process start,
#    but currently-running processes still hold the old binary in memory.
#    Guard: if $RELEASE is empty (lost variable, wrong shell), the link would resolve
#    to /opt/aart-gr/releases/ — a directory with no aart_api in it — and the next
#    service restart would fail status=203/EXEC. Abort before that can happen.
[ -n "$RELEASE" ] || { echo "ABORT: \$RELEASE not set in this shell; re-export it from §16 step 1"; exit 1; }
ssh -p 2222 moruogrsadmin@102.133.230.104 "
  [ -x /opt/aart-gr/releases/$RELEASE/aart_api ] || {
    echo 'ABORT: /opt/aart-gr/releases/$RELEASE/aart_api missing or not executable'; exit 1;
  } &&
  ln -sfn /opt/aart-gr/releases/$RELEASE /opt/aart-gr/current
"

# 5. Rolling restart. Instance 1 first (it applies any new migrations from step 2).
ssh -p 2222 moruogrsadmin@102.133.230.104 '
  sudo systemctl restart aart-gr-api@9091 &&
  for i in $(seq 1 30); do
    code=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:9091/health);
    [ "$code" = "200" ] && break;
    sleep 1;
  done &&
  echo "instance 9091: HTTP $code (after ${i}s)" &&
  [ "$code" = "200" ] || { echo "ABORT: 9091 not healthy, not restarting 9092"; exit 1; } &&
  sudo systemctl restart aart-gr-api@9092 &&
  for i in $(seq 1 30); do
    code=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:9092/health);
    [ "$code" = "200" ] && break;
    sleep 1;
  done &&
  echo "instance 9092: HTTP $code (after ${i}s)"
'

# 6. Verify migrations actually applied (only meaningful if your release contained any).
ssh -p 2222 moruogrsadmin@102.133.230.104 \
  "journalctl -u aart-gr-api@9091 --since '2 minutes ago' --no-pager | grep -iE 'migration|migrate' | tail -10"
# expect at least: "Running migrations on startup" and "Migrations completed successfully"
# expect for new ones: "Applied migration: <timestamp>_<name>.sql"

# Prune old releases (keep last 5). No sudo — moruogrsadmin owns the dirs.
ssh -p 2222 moruogrsadmin@102.133.230.104 '
  cd /opt/aart-gr/releases &&
  ls -1tr | head -n -5 | xargs --no-run-if-empty rm -rf
'
rm api/aart_api

