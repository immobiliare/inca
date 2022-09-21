#!/bin/bash

set -eu
set -o pipefail

name="${1}"

go run main.go gen --ca -o test -n "${name}"
go run main.go server -c test/inca.yml &!
until $(curl -sfI -o /dev/null 127.0.0.1:8080); do
    echo "waiting for server to come up..."
    sleep 1
done

curl "127.0.0.1:8080/ca/local" > "${name}.crt"
openssl x509 -in "${name}.crt" -text

curl "127.0.0.1:8080/test.${name}" > "test.${name}.crt"
openssl x509 -in "test.${name}.crt" -text
test -f "test/test.${name}/crt.pem"
curl "127.0.0.1:8080/enum" | jq -r '.results|length' | grep 1
curl "127.0.0.1:8080/test.${name}/key" > "test.${name}.key"
curl "127.0.0.1:8080/test.${name}/show" | jq

go run test/test.go "test.${name}.crt" "test.${name}.key" &!
echo -e "127.0.1.1\ttest.${name}" >> /etc/hosts
until $(curl --cacert "${name}.crt" -sfI -o /dev/null https://test.${name}:8081); do
    echo "waiting for server to come up..."
    sleep 1
done
curl --cacert "${name}.crt" "https://test.${name}:8081"
kill %2

curl "127.0.0.1:8080/test.${name}" -X DELETE
test -f "test/test.${name} /crt.pem" && exit 1 || echo -n
kill %1