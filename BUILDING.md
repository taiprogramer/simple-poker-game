# Building Simple poker game

This document provides instructions for building & running an application. If you read this document but can't build an application, don't hesitate to [file an issue](https://github.com/taiprogramer/simple-poker-game/issues).

## Prerequisite

This project is a cross-platform application which means you can run it on every OS
you want.

### Server

- go (^1.18.1): go executable.
- gow (any): go wrapper that watches for changes and reruns go command.

### Client

### Mobile app: Android

- Open JDK (java-11-openjdk).
- flutter (^2.10.4): flutter executable.
- adb: Android Debug Bridge (connecting real device for running an app).

## Building

### Server

- Setting up dot env

```sh
cp backend/env backend/.env
```

| Variable              | Description                        |
| --------------------- | ---------------------------------- |
| HMAC_SECRET_KEY=      | Use for bcrypt password encryption |
| SERVER_MODE=          | "production" or "development"      |
| JWT_EXPIRE_IN_MINUTE= | JWT expiration time in minute      |

Right now, server only supports development mode (SERVER_MODE).

- Running the server and watch for changes

```sh
cd backend
gow run main.go
```

### Client

#### Mobile app: Android

- Connecting your Android phone using adb.
- Running the application on connected phone

```sh
cd app_frontend/simple_poker_game
flutter run
```
