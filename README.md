# Unofficial swiss CFA society API

Unffocial API for [CFA society](https://cfasocietyswitzerland.org/) events in Swizterland. A helpful API for any developer or service who wants to develop application around CFA society events in Swizterland.

Navigate to [SwaggerHub](https://app.swaggerhub.com/apis/AntonyGandonoumigan/unofficial-cfa_society_switzerland_api/1.0.0#/default/get_api_v1_events) for API documentation.

## Install locally

### Docker
Create a new GCP project, a new [ID client Oauth](https://console.cloud.google.com/apis/credentials) and run the commands below.

```text
docker pull corni/cfasociety:latest
docker run --rm -e GCP_CLIENT_ID=<YOUR_PROJRCT_ID> -e GCP_CLIENT_SECRET=<YOUR_KEY> -p 8080:8080 corni/cfasociety
```

## 🤝 Contributing

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.
