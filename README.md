# Nos Notification Service

## Launching

The program is located under this path and takes no arguments:

    $ go run ./cmd/notification-service

There is a Dockerfile available.

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