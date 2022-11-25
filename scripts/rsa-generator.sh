#!/bin/bash
set -e

cd $(dirname $0)

openssl genrsa -out ../assets/private.example.pem 2048
openssl rsa -in ../assets/private.example.pem -pubout > ../assets/public.example.pem
