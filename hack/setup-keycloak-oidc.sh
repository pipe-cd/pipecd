#!/bin/bash

echo "Setting up Keycloak OIDC provider for testing..."

set -e

KEYCLOAK_VERSION="22.0.1"
KEYCLOAK_PORT=8091
KEYCLOAK_ADMIN="admin"
KEYCLOAK_ADMIN_PASSWORD="admin"
KEYCLOAK_CONTAINER_NAME="pipecd-keycloak"
KEYCLOAK_REALM="pipecd-test"
KEYCLOAK_CLIENT_ID="pipecd-client"
KEYCLOAK_CLIENT_SECRET="pipecd-client-secret"
PIPECD_REDIRECT_URI="http://localhost:8080/auth/callback"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Docker is not installed. Please install Docker first."
    exit 1
fi

# Stop and remove existing container if it exists
echo "Cleaning up any existing Keycloak container..."
docker stop $KEYCLOAK_CONTAINER_NAME 2>/dev/null || true
docker rm $KEYCLOAK_CONTAINER_NAME 2>/dev/null || true

# Create temp directory for configuration
TEMP_DIR=$(mktemp -d)
trap 'rm -rf $TEMP_DIR' EXIT

# Create realm configuration file
cat > $TEMP_DIR/realm.json <<EOF
{
  "realm": "$KEYCLOAK_REALM",
  "enabled": true,
  "accessTokenLifespan": 300,
  "roles": {
    "realm": [
      {
        "name": "viewer"
      },
      {
        "name": "editor"
      },
      {
        "name": "admin"
      }
    ]
  },
  "users": [
    {
      "username": "admin-user",
      "email": "admin@example.com",
      "enabled": true,
      "credentials": [
        {
          "type": "password",
          "value": "password",
          "temporary": false
        }
      ],
      "realmRoles": ["admin"],
      "attributes": {
        "groups": ["admin"]
      }
    },
    {
      "username": "editor-user",
      "email": "editor@example.com",
      "enabled": true,
      "credentials": [
        {
          "type": "password",
          "value": "password",
          "temporary": false
        }
      ],
      "realmRoles": ["editor"],
      "attributes": {
        "groups": ["editor"]
      }
    },
    {
      "username": "viewer-user",
      "email": "viewer@example.com",
      "enabled": true,
      "credentials": [
        {
          "type": "password",
          "value": "password",
          "temporary": false
        }
      ],
      "realmRoles": ["viewer"],
      "attributes": {
        "groups": ["viewer"]
      }
    }
  ],
  "clients": [
    {
      "clientId": "$KEYCLOAK_CLIENT_ID",
      "secret": "$KEYCLOAK_CLIENT_SECRET",
      "enabled": true,
      "redirectUris": ["$PIPECD_REDIRECT_URI"],
      "clientAuthenticatorType": "client-secret",
      "protocol": "openid-connect",
      "standardFlowEnabled": true,
      "implicitFlowEnabled": false,
      "directAccessGrantsEnabled": true,
      "serviceAccountsEnabled": false,
      "publicClient": false,
      "frontchannelLogout": false,
      "attributes": {
        "access.token.lifespan": "300"
      },
      "protocolMappers": [
        {
          "name": "groups",
          "protocol": "openid-connect",
          "protocolMapper": "oidc-usermodel-attribute-mapper",
          "consentRequired": false,
          "config": {
            "userinfo.token.claim": "true",
            "user.attribute": "groups",
            "id.token.claim": "true",
            "access.token.claim": "true",
            "claim.name": "groups",
            "jsonType.label": "String"
          }
        },
        {
          "name": "username",
          "protocol": "openid-connect",
          "protocolMapper": "oidc-usermodel-property-mapper",
          "consentRequired": false,
          "config": {
            "userinfo.token.claim": "true",
            "user.attribute": "username",
            "id.token.claim": "true",
            "access.token.claim": "true",
            "claim.name": "username",
            "jsonType.label": "String"
          }
        }
      ]
    }
  ]
}
EOF

echo "Starting Keycloak container..."
docker run -d \
  --name $KEYCLOAK_CONTAINER_NAME \
  -p $KEYCLOAK_PORT:8080 \
  -e KEYCLOAK_ADMIN=$KEYCLOAK_ADMIN \
  -e KEYCLOAK_ADMIN_PASSWORD=$KEYCLOAK_ADMIN_PASSWORD \
  -v $TEMP_DIR/realm.json:/opt/keycloak/data/import/realm.json \
  quay.io/keycloak/keycloak:$KEYCLOAK_VERSION \
  start-dev \
  --import-realm

# Wait for Keycloak to be ready
echo "Waiting for Keycloak to start..."
until curl -s http://localhost:$KEYCLOAK_PORT > /dev/null; do
  echo "Waiting for Keycloak to be ready..."
  sleep 3
done

echo "---------------------------------------------------------------------"
echo "Keycloak OIDC setup complete!"
echo "---------------------------------------------------------------------"
echo "Keycloak URL: http://localhost:$KEYCLOAK_PORT"
echo "Admin console: http://localhost:$KEYCLOAK_PORT/admin"
echo "Admin username: $KEYCLOAK_ADMIN"
echo "Admin password: $KEYCLOAK_ADMIN_PASSWORD"
echo ""
echo "Realm: $KEYCLOAK_REALM"
echo "Client ID: $KEYCLOAK_CLIENT_ID"
echo "Client Secret: $KEYCLOAK_CLIENT_SECRET"
echo "Redirect URI: $PIPECD_REDIRECT_URI"
echo ""
echo "Test users:"
echo "  admin-user / password (admin role)"
echo "  editor-user / password (editor role)"
echo "  viewer-user / password (viewer role)"
echo ""
echo "PipeCD Control Plane Configuration Example:"
echo "---------------------------------------------------------------------"
cat << EOF
apiVersion: "pipecd.dev/v1beta1"
kind: ControlPlane
spec:
  sharedSSOConfigs:
    - name: oidc
      provider: OIDC
      oidc:
        clientId: "$KEYCLOAK_CLIENT_ID"
        clientSecret: "$KEYCLOAK_CLIENT_SECRET"
        issuer: "http://localhost:$KEYCLOAK_PORT/realms/$KEYCLOAK_REALM"
        redirectUri: "$PIPECD_REDIRECT_URI"
        scopes:
          - openid
          - profile
EOF
echo "---------------------------------------------------------------------"
echo "To stop the Keycloak container: docker stop $KEYCLOAK_CONTAINER_NAME"