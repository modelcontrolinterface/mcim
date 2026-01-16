#!/usr/bin/env bash

CERT_DIR="certs"
KEY_FILE="$CERT_DIR/key.pem"
ERT_FILE="$CERT_DIR/cert.pem"

mkdir -p "$CERT_DIR"

openssl genrsa -out "$KEY_FILE" 2048
openssl req -new -x509 -key "$KEY_FILE" -out "$CERT_FILE" -days 365 -nodes -subj "/CN=localhost"
