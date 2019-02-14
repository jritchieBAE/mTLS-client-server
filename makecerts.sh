openssl req -newkey rsa:2048 -new -nodes -x509 -days 3650 -out ca.crt -keyout ca.key -subj "/C=US/ST=California/L=Mountain View/O=Organisation/OU=Unit/CN=localhost"
openssl req -new -key ca.key -out ca.csr -subj "/C=US/ST=California/L=Mountain View/O=Organisation/OU=Unit/CN=localhost"
openssl x509 -req -in ca.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out localhost.crt