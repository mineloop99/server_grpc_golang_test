
#Generate ca.crt
openssl genrsa -passout pass:171099 -aes256 -out ca.key 4096
openssl req -passin pass:171099 -new -x509 -days 365 -key ca.key -out ca.crt -subj "/C=VN/ST=HCM/L=THUDUC/O=Mine Loop, WanatabeYuu/OU=IT/CN=127.0.0.1"

#Generate server.key
openssl genrsa -passout pass:171099 -aes256 -out server.key 4096

#Generate server.csr
openssl req -passin pass:171099 -new -key server.key -out server.csr -subj "/C=VN/ST=HCM/L=THUDUC/O=Mine Loop, WanatabeYuu/OU=IT/CN=127.0.0.1"

#Generate server.crt
openssl x509 -req -passin pass:171099 -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt

#Convert the server certificate to .pem format (server.pem) - usable by gRPC

openssl pkcs8 -topk8 -nocrypt -passin pass:171099 -in server.key -out server.pem

