#!/bin/bash

# Maestro Certificate Generation Script
# Generates self-signed certificates for development

set -e

CERT_DIR="$(dirname "$0")"
DAYS_VALID=365

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Maestro Certificate Generator${NC}"
echo "==============================="
echo ""

# Create certificate directory if it doesn't exist
mkdir -p "$CERT_DIR"

# Generate CA private key
echo -e "${YELLOW}Generating CA private key...${NC}"
openssl genrsa -out "$CERT_DIR/ca.key" 4096

# Generate CA certificate
echo -e "${YELLOW}Generating CA certificate...${NC}"
openssl req -new -x509 -days $DAYS_VALID -key "$CERT_DIR/ca.key" -out "$CERT_DIR/ca.crt" \
    -subj "/C=US/ST=State/L=City/O=Maestro/CN=Maestro CA"

# Generate server private key
echo -e "${YELLOW}Generating server private key...${NC}"
openssl genrsa -out "$CERT_DIR/server.key" 4096

# Generate server certificate request
echo -e "${YELLOW}Generating server certificate request...${NC}"
openssl req -new -key "$CERT_DIR/server.key" -out "$CERT_DIR/server.csr" \
    -subj "/C=US/ST=State/L=City/O=Maestro/CN=localhost"

# Sign server certificate with CA
echo -e "${YELLOW}Signing server certificate...${NC}"
openssl x509 -req -days $DAYS_VALID -in "$CERT_DIR/server.csr" \
    -CA "$CERT_DIR/ca.crt" -CAkey "$CERT_DIR/ca.key" -CAcreateserial \
    -out "$CERT_DIR/server.crt" \
    -extfile <(printf "subjectAltName=DNS:localhost,IP:127.0.0.1")

# Generate client private key
echo -e "${YELLOW}Generating client private key...${NC}"
openssl genrsa -out "$CERT_DIR/client.key" 4096

# Generate client certificate request
echo -e "${YELLOW}Generating client certificate request...${NC}"
openssl req -new -key "$CERT_DIR/client.key" -out "$CERT_DIR/client.csr" \
    -subj "/C=US/ST=State/L=City/O=Maestro/CN=maestro-client"

# Sign client certificate with CA
echo -e "${YELLOW}Signing client certificate...${NC}"
openssl x509 -req -days $DAYS_VALID -in "$CERT_DIR/client.csr" \
    -CA "$CERT_DIR/ca.crt" -CAkey "$CERT_DIR/ca.key" -CAcreateserial \
    -out "$CERT_DIR/client.crt"

# Generate MCP client certificate
echo -e "${YELLOW}Generating MCP client certificate...${NC}"
openssl genrsa -out "$CERT_DIR/mcp-client.key" 4096
openssl req -new -key "$CERT_DIR/mcp-client.key" -out "$CERT_DIR/mcp-client.csr" \
    -subj "/C=US/ST=State/L=City/O=Maestro/CN=mcp-client"
openssl x509 -req -days $DAYS_VALID -in "$CERT_DIR/mcp-client.csr" \
    -CA "$CERT_DIR/ca.crt" -CAkey "$CERT_DIR/ca.key" -CAcreateserial \
    -out "$CERT_DIR/mcp-client.crt"

# Clean up CSR files
rm -f "$CERT_DIR"/*.csr
rm -f "$CERT_DIR"/*.srl

# Set appropriate permissions
chmod 600 "$CERT_DIR"/*.key
chmod 644 "$CERT_DIR"/*.crt

echo ""
echo -e "${GREEN}✓ Certificate generation complete!${NC}"
echo ""
echo "Generated files:"
echo "  • CA Certificate: $CERT_DIR/ca.crt"
echo "  • CA Private Key: $CERT_DIR/ca.key"
echo "  • Server Certificate: $CERT_DIR/server.crt"
echo "  • Server Private Key: $CERT_DIR/server.key"
echo "  • Client Certificate: $CERT_DIR/client.crt"
echo "  • Client Private Key: $CERT_DIR/client.key"
echo "  • MCP Client Certificate: $CERT_DIR/mcp-client.crt"
echo "  • MCP Client Private Key: $CERT_DIR/mcp-client.key"
echo ""
echo "Certificates are valid for $DAYS_VALID days"
echo ""
echo -e "${YELLOW}Note: These are self-signed certificates for development only!${NC}"
