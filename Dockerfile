FROM debian:stable-slim
RUN apt-get update && apt-get install -y ca-certificates
ADD cfaSocietyScrapper /usr/bin/cfaSocietyScrapper
ADD _dist /usr/bin/_dist
WORKDIR /usr/bin/
RUN openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 3650 -nodes -subj "/C=XX/ST=StateName/L=CityName/O=CompanyName/OU=CompanySectionName/CN=CommonNameOrHostname/CN=localhost"
CMD ["cfaSocietyScrapper"]

