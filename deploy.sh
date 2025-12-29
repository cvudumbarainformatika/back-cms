#!/usr/bin/env bash
set -u

###########################
# CONFIG
###########################
SERVER_USER="dev"
SERVER_HOST="202.10.48.229"
SERVER_DIR="/var/www/juwita"
IMAGE_NAME="starter-api"
IMAGE_TAG="prod"
FULL_IMAGE="${IMAGE_NAME}:${IMAGE_TAG}"
LOCAL_TMP_DIR="./.deploy_tmp"
REMOTE_IMAGES_DIR="${SERVER_DIR}/images"
REMOTE_RELEASES_DIR="${SERVER_DIR}/releases"
BACKEND_DIR="${SERVER_DIR}/backend"
COMPOSE_FILE="docker-compose.prod.yml"
ENV_FILE=".env.production"
LOGFILE="./deploy.log"
HEALTH_TIMEOUT=40
MAX_BACKUPS=2     # jumlah tar file backup yang disimpan
MAX_IMAGES=2      # jumlah image lama yang disimpan di Docker host
###########################

timestamp() { date +"%Y%m%d_%H%M%S"; }
now=$(timestamp)

echo "[$(date +"%Y-%m-%d %H:%M:%S")] START DEPLOY $now" | tee -a "$LOGFILE"

# 1) Build binary Go
echo "-> Build binary (GOOS=linux GOARCH=amd64)..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./main || {
  echo "❌ go build gagal" | tee -a "$LOGFILE"
  exit 1
}

# 2) Build docker image lokal
echo "-> Build docker image (platform linux/amd64) using buildx..."
docker buildx inspect >/dev/null 2>&1 || docker buildx create --use >/dev/null 2>&1
docker buildx build --platform linux/amd64 -t "${FULL_IMAGE}" --load . || {
  echo "❌ docker buildx gagal" | tee -a "$LOGFILE"
  exit 1
}

# 3) Save image ke tar
mkdir -p "$LOCAL_TMP_DIR"
IMAGE_TAR="${LOCAL_TMP_DIR}/${IMAGE_NAME}_${IMAGE_TAG}_${now}.tar"
echo "-> Save image to $IMAGE_TAR"
docker save -o "$IMAGE_TAR" "${FULL_IMAGE}" || {
  echo "❌ docker save gagal" | tee -a "$LOGFILE"
  exit 1
}

# 4) Upload tar ke server
echo "-> Uploading image to ${SERVER_USER}@${SERVER_HOST}:${REMOTE_IMAGES_DIR}/"
scp "$IMAGE_TAR" "${SERVER_USER}@${SERVER_HOST}:${REMOTE_IMAGES_DIR}/" || {
  echo "❌ scp gagal" | tee -a "$LOGFILE"
  exit 1
}

# 5) Remote deploy
SSH_CMD=$(cat <<EOF
set -u

IMAGES_DIR="${REMOTE_IMAGES_DIR}"
RELEASES_DIR="${REMOTE_RELEASES_DIR}"
BACKEND_DIR="${BACKEND_DIR}"
IMAGE_TAR="\$(basename "$IMAGE_TAR")"
IMAGE_NAME="${IMAGE_NAME}"
IMAGE_TAG="${IMAGE_TAG}"
FULL_IMAGE="\${IMAGE_NAME}:\${IMAGE_TAG}"
COMPOSE_FILE="${COMPOSE_FILE}"
ENV_FILE="${ENV_FILE}"
TS="${now}"
MAX_BACKUPS=${MAX_BACKUPS}
MAX_IMAGES=${MAX_IMAGES}

mkdir -p "\$RELEASES_DIR/${now}"
cd "\$IMAGES_DIR" || exit 1

echo "[\$(date '+%Y-%m-%d %H:%M:%S')] Remote deploy start (ts=${now})"

# Backup existing image
if docker image inspect "\$FULL_IMAGE" >/dev/null 2>&1; then
  echo "-> Found existing image \$FULL_IMAGE, saving backup..."
  BACKUP_TAR="\$RELEASES_DIR/${now}/\${IMAGE_NAME}_\${IMAGE_TAG}_backup_\$(date +%s).tar"
  docker save -o "\$BACKUP_TAR" "\$FULL_IMAGE" || echo "⚠ Gagal menyimpan backup image, lanjutkan..."
  echo "-> backup saved at \$BACKUP_TAR"
fi

# Load new image
echo "-> Loading new image: \$IMAGE_TAR"
docker load -i "\$IMAGES_DIR/\$IMAGE_TAR" || {
  echo "❌ docker load gagal"; exit 1
}

# Restart compose
cd "\$BACKEND_DIR" || exit 1
echo "-> docker compose down"
docker compose -f "\$COMPOSE_FILE" --env-file "\$ENV_FILE" down || echo "⚠ compose down non-fatal"
echo "-> docker compose up -d"
docker compose -f "\$COMPOSE_FILE" --env-file "\$ENV_FILE" up -d || {
  echo "❌ docker compose up gagal"; exit 1
}

# Healthcheck
CID="\$(docker compose -f "\$COMPOSE_FILE" ps -q app 2>/dev/null || true)"
if [ -z "\$CID" ]; then
  CID="\$(docker ps --filter ancestor=\$FULL_IMAGE --format '{{.ID}}' | head -n1 || true)"
fi

if [ -n "\$CID" ]; then
  echo "-> Waiting for container \$CID to be running (timeout ${HEALTH_TIMEOUT}s)..."
  for i in \$(seq 1 ${HEALTH_TIMEOUT}); do
    STATE="\$(docker inspect --format='{{json .State}}' \$CID 2>/dev/null || echo '')"
    RUNNING="\$(echo \"\$STATE\" | grep -o '\"Running\":true' || true)"
    if [ -n \"\$RUNNING\" ]; then
      echo "-> Container is running - assume OK"
      break
    fi
    sleep 1
  done
else
  echo "⚠ Tidak dapat menemukan container untuk healthcheck"
fi

# Cleanup old tar backups (keep MAX_BACKUPS)
echo "-> Cleanup old tar backups, keep \$MAX_BACKUPS latest"
ls -1t "\$RELEASES_DIR"/*/*.tar | tail -n +\$((MAX_BACKUPS+1)) | while read OLD_TAR; do
    echo "Deleting \$OLD_TAR"
    rm -f "\$OLD_TAR"
done

# Cleanup old Docker images (keep MAX_IMAGES)
echo "-> Cleanup old Docker images, keep \$MAX_IMAGES latest"
IMAGES_TO_DELETE=\$(docker images --format '{{.Repository}}:{{.Tag}} {{.CreatedAt}} {{.ID}}' | grep "\$IMAGE_NAME" | sort -rk2 | tail -n +\$((MAX_IMAGES+1)) | awk '{print \$3}')
for IMG in \$IMAGES_TO_DELETE; do
    echo "Deleting Docker image \$IMG"
    docker rmi -f \$IMG || echo "⚠ Gagal menghapus image \$IMG, lanjutkan..."
done

EOF
)

echo "-> Executing remote commands..."
ssh "${SERVER_USER}@${SERVER_HOST}" "$SSH_CMD"
REMOTE_EXIT=$?

if [ $REMOTE_EXIT -ne 0 ]; then
  echo "❌ Remote deploy reported failure (exit $REMOTE_EXIT). Check logs on server ${REMOTE_RELEASES_DIR}/${now}" | tee -a "$LOGFILE"
  exit 1
fi

# Cleanup lokal
echo "-> Cleaning local tmp"
rm -rf "$LOCAL_TMP_DIR"
echo "[$(date +"%Y-%m-%d %H:%M:%S")] DEPLOY COMPLETE" | tee -a "$LOGFILE"

exit 0
