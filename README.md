# Nos Notification Service

## How to do local development

1. Get an APNs certificate.
   1. Obtain a certificate from Apple (see ["Obtain a provider certificate from Apple"][get-apns-cert]).
   2. Export the certificate from your keychain in the PKCS#12 format.
2. Start the docker daemon.
3. Run `make recreate-emulator` to start the Firestore emulator using Docker compose.
4. Run the following command changing `NOTIFICATIONS_APNS_CERTIFICATE_PATH` and `NOTIFICATIONS_APNS_CERTIFICATE_PASSWORD`:

```
NOTIFICATIONS_APNS_CERTIFICATE_PATH=/path/to/your/apns/cert.p12 \
NOTIFICATIONS_APNS_CERTIFICATE_PASSWORD="your cert password if you set one" \
FIRESTORE_EMULATOR_HOST=localhost:8200 \
NOTIFICATIONS_FIRESTORE_PROJECT_ID=test-project-id \
NOTIFICATIONS_APNS_TOPIC=com.verse.Nos \
NOTIFICATIONS_ENVIRONMENT=DEVELOPMENT \
go run ./cmd/notification-service
```


## Building and running

Buid the program like so:

    $ go build -o notification-service ./cmd/notification-service
    $ ./notification-service

The program takes no arguments. There is a Dockerfile available.

## Configuration

Configuration is performed using environment variables. This is also the case
for the Dockerfile.

### `NOTIFICATIONS_NOSTR_LISTEN_ADDRESS`

Listen address for the websocket connections in the format accepted by the
standard library.

Optional, defaults to `:8008` if empty.

### `NOTIFICATIONS_FIRESTORE_PROJECT_ID`

Your Firestore project id.

Required.

### `NOTIFICATIONS_FIRESTORE_CREDENTIALS_JSON_PATH`

Path to your Firestore credentials JSON file.

Required if you are not using the emulator (`FIRESTORE_EMULATOR_HOST` is not set).

### `NOTIFICATIONS_APNS_TOPIC`

Topic on which APNs notifications will be sent. Probably your iOS app id.

Required.

### `NOTIFICATIONS_APNS_CERTIFICATE_PATH`

Path to your APNs certificate file in the PKCS#12 format. They normally come in
a different format I think so you need to presumably export this from your
keychain.

Required.

### `NOTIFICATIONS_APNS_CERTIFICATE_PASSWORD`

Password to your APNs certificate file.

Optional, leave empty if the certificate doesn't have a password.

### `NOTIFICATIONS_ENVIRONMENT`

Execution environment. Affects:
- whether testing or production APNs server is used

Optional, can be set to `PRODUCTION` or `DEVELOPMENT`. Defaults to `PRODUCTION`.


### `FIRESTORE_EMULATOR_HOST`

Optional, this is used by the Firestore libraries and can be useful for testing
but you shouldn't ever have to set this in production.


[get-apns-cert]: https://developer.apple.com/documentation/usernotifications/setting_up_a_remote_notification_server/establishing_a_certificate-based_connection_to_apns#2947597