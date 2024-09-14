FROM debian:stable-slim
RUN apt-get update && apt-get install -y ca-certificates
ENV PORT=8080 \
    SQLI_PATH=database.db \
    OAUTH_CALLBACK_URI=https://localhost:8080/auth/callback

ADD cfaSocietyScrapper /usr/bin/cfaSocietyScrapper
ADD database.db /usr/bin/database.db
ADD _dist /usr/bin/_dist
WORKDIR /usr/bin/
RUN openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 3650 -nodes -subj "/C=XX/ST=StateName/L=CityName/O=CompanyName/OU=CompanySectionName/CN=CommonNameOrHostname/CN=localhost"
CMD ["cfaSocietyScrapper"]

