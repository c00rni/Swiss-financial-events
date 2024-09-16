# Unofficial swiss CFA society API

Unffocial API for [CFA society](https://cfasocietyswitzerland.org/) events in Swizterland. A helpful API for any developer who wants to develop application around CFA society events in Swizterland.

Navigate to [the app website](https://swiss-cfa-api-530073081731.europe-west6.run.app/) to create a free API key. The API can be test through the [SwaggerHub](https://app.swaggerhub.com/apis/AntonyGandonoumigan/unofficial-cfa_society_switzerland_api/1.0.0#/default/get_api_v1_events) documentation.

## Usage

Login to the webapp to generate an API Key and request the API. Accounts are **limited to 100 API calls** per month.

```bash
curl -X 'GET' \
  'https://swiss-cfa-api-530073081731.europe-west6.run.app/api/v1/categories' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer <YOUR_API_KEY>'
```

See the [documentation](https://app.swaggerhub.com/apis/AntonyGandonoumigan/unofficial-cfa_society_switzerland_api/1.0.0) to learn more about the API.

## Install locally

### Docker

The API use oauth2 to manage authentifiaction. Create a project and [ID client Oauth](https://console.cloud.google.com/apis/credentials) on Google Platfome and run the commands below. 

The Go server use libsql driver to connect to the Database. You can create a database on [Turso](https://turso.tech/) and get and migrate the database schema by executing `./_scripts/migrateup.sh` from the root.

```text
docker pull corni/cfasociety:latest
docker run --rm -it -e GCP_CLIENT_ID=<YOUR_PROJECT_ID> -e GCP_CLIENT_SECRET=<YOUR_KEY> -e PORT=8080 -e DATABASE_URL=<YOUR_DATABASE_LINK> -e OAUTH_CALLBACK_URI=<YOUR_OAUTH_CALLBACK_URI> -p 8080:8080 corni/cfasociety:latest
```

## 🤝 Contributing

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.
