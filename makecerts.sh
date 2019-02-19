# Create root certificate
# openssl req -newkey rsa:2048 -new -nodes -x509 -days 3650 -out root.crt -keyout root.key -subj "/C=US/ST=California/L=Mountain View/O=Organisation/OU=Root/CN=localhost"

# install root certificates
# echo "Installing root certificate"
# sudo cp root.crt /etc/pki/ca-trust/source/anchors
# sudo update-ca-trust extract

# create Prometheus certificate & key
openssl req -newkey rsa:2048 -new -nodes -x509 -days 3650 -keyout client.key -subj "/C=US/ST=California/L=Mountain View/O=Organisation/OU=Prometheus/CN=localhost"
openssl req -new -key client.key -out client.csr -subj "/C=US/ST=California/L=Mountain View/O=Organisation/OU=Prometheus/CN=localhost"
openssl x509 -req -in client.csr -CA root.crt -CAkey root.key -CAcreateserial -out clientSigned.crt

# create Node certificate & key
openssl req -newkey rsa:2048 -new -nodes -x509 -days 3650 -keyout server.key -subj "/C=US/ST=California/L=Mountain View/O=Organisation/OU=NodeExporter/CN=localhost"
openssl req -new -key server.key -out server.csr -subj "/C=US/ST=California/L=Mountain View/O=Organisation/OU=NodeExporter/CN=localhost"
openssl x509 -req -in server.csr -CA root.crt -CAkey root.key -CAcreateserial -out serverSigned.crt