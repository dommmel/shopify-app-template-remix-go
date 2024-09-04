## What is this?

This is a template for a Shopify app using Remix and Golang. The Remix frontend *does not* use the [shopify-app-remix/shopify-app-js](https://github.com/Shopify/shopify-app-js) packages nor the [Shopify remix template](https://github.com/Shopify/shopify-app-template-remix). This is made possible in part by Shopify's managed app installations. Shopify-related operations are handled in the backend.

The Shopify-related code is located in
- Backend: [golang-backend/pkg/shopify](golang-backend/pkg/shopify) and also [here](golang-backend/user/user_repository.go#L76)
- Frontend: [remix-frontend/app/.server/](remix-frontend/app/.server/)

## Key Characteristics:
- [gRPC](https://grpc.io) for communication between the frontend and backend
- No dependencies on Shopify libraries (Go or JavaScript)
- SQLite is used by default but can be replaced with any other database

## Prerequisites:
- `protoc` binary: See [protoc installation guide](https://grpc.io/docs/protoc-installation/)
- Shopify CLI
- Go (Golang)
- Node.js

## Running the Development Servers

Run the following command:

```bash
shopify app dev
```

On the first run, this will create a `shopify.app.toml` file in the root of the project. Be sure to enable [managed installs](https://shopify.dev/docs/apps/build/authentication-authorization/app-installation).

This command will:
- Install dependencies
- Run `codegen.sh` (for gRPC) in both the backend and frontend folders
- Start the Go backend with live reload (using `wgo`)
- Start the Remix dev server


## Things Not Included:
- Encryption at rest
- Webhooks
- Graphql Shopify API clients
