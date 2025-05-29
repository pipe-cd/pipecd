#!/bin/bash
set -eu

# function to get the private IP address
get_private_ip() {
  # macOS
  if [[ "$OSTYPE" == "darwin"* ]]; then
    ip=$(ipconfig getifaddr en0 2>/dev/null)
    if [[ -z "$ip" ]]; then
      ip=$(ipconfig getifaddr en1 2>/dev/null)
    fi
    echo "$ip"
    return
  fi

  # Linux
  if command -v ip >/dev/null 2>&1; then
    for iface in wlan0 eth0 enp0s3; do
      ip=$(ip -4 addr show "$iface" 2>/dev/null | awk '/inet / {print $2}' | cut -d/ -f1)
      if [[ -n "$ip" && "$ip" != "127.0.0.1" ]]; then
        echo "$ip"
        return
      fi
    done

    ip=$(ip route get 1.1.1.1 2>/dev/null | awk '/src/ {print $NF; exit}')
    if [[ -n "$ip" && "$ip" != "127.0.0.1" ]]; then
      echo "$ip"
      return
    fi
  fi

  # ifconfig fallback (BusyBox etc.)
  if command -v ifconfig >/dev/null 2>&1; then
    ip=$(ifconfig | awk '/inet / && $2 != "127.0.0.1" {print $2; exit}' | sed 's/addr://')
    echo "$ip"
    return
  fi

  echo "Could not determine private IP address" >&2
  return 1
}

# ====== Defaults ======
PIPECD_CONTROL_PLANE_ADDRESS=${PIPECD_CONTROL_PLANE_ADDRESS:-http://localhost:8080}
LOCAL_KEYCLOAK_ADDRESS=${LOCAL_KEYCLOAK_ADDRESS:-"http://$(get_private_ip):8081"}

# ====== File Setup ======
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SOURCE_FILE="$SCRIPT_DIR/realm.json"
WORKING_FILE="$SCRIPT_DIR/realm.local.json"

if [ ! -f "$SOURCE_FILE" ]; then
  echo "❌ $SOURCE_FILE not found. Aborting."
  exit 1
fi

echo "📁 Copying $SOURCE_FILE to $WORKING_FILE"
cp "$SOURCE_FILE" "$WORKING_FILE"

echo "🔧 Replacing environment variables in $WORKING_FILE"
export PIPECD_CONTROL_PLANE_ADDRESS LOCAL_KEYCLOAK_ADDRESS
envsubst <"$WORKING_FILE" >"${WORKING_FILE}.tmp" && mv "${WORKING_FILE}.tmp" "$WORKING_FILE"

# ====== Docker Compose Up ======
echo ""
echo "🚀 Starting Keycloak via docker compose"
cd "$SCRIPT_DIR"
docker compose up -d --force-recreate --wait

# ====== Output Preview ======
echo ""
echo "✅ Generated $WORKING_FILE successfully"
echo ""
echo "🔍 Configuration (for your reference):"
echo ""
cat <<EOF
apiVersion: "pipecd.dev/v1beta1"
kind: ControlPlane
spec:
  stateKey: test
  sharedSSOConfigs:
  - name: oidc
    provider: OIDC
    oidc:
      clientId: pipecd-keycloak
      clientSecret: L702XzhZfnYDzcA0dbtoaMRt3ZUTw9m9
      issuer: ${LOCAL_KEYCLOAK_ADDRESS}/realms/pipecd-test-realm
      redirectUri: ${PIPECD_CONTROL_PLANE_ADDRESS}/auth/callback
      usernameClaimKey: name
      rolesClaimKey: realm_roles
      scopes:
      - openid
      - profile
EOF
