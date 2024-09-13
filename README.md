# Swiss Financial Events

## Usage

Create a new GCP project, a new [ID client Oauth](https://console.cloud.google.com/apis/credentials) and run the commands below.

```text
docker pull corni/cfasociety:latest
docker run --rm -e GCP_CLIENT_ID=<YOUR_PROJRCT_ID> -e GCP_CLIENT_SECRET=<YOUR_KEY> -p 8080:8080 corni/cfasociety
```


