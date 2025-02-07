# Go Auth Gate

Go Auth Gate is an advanced middleware authentication framework designed for Go applications, enabling robust
authentication mechanisms through Apple Sign-In, Google Sign-In, and JWT-based authorization. The framework is
engineered to seamlessly integrate with existing web architectures, providing a secure and scalable authentication
solution.

## Key Features

- Comprehensive authentication support for Apple and Google OAuth providers.
- Secure JWT-based authentication mechanism for robust session management.
- Middleware-based architecture facilitating modular integration within Go applications.
- Configurable parameters via YAML for flexible deployment and customization.
- Refresh mechanism for authentication tokens via `refresh.go`.

## Installation

To integrate Go Auth Gate into your project, execute the following command:

```sh
 go get github.com/iooojik/go-auth-gate
```

## Implementation Guide

The following example illustrates the integration of Go Auth Gate's authentication middleware within a Go application:

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	authgate "github.com/iooojik/go-auth-gate"
	"github.com/iooojik/go-auth-gate/internal/config"
)

func main() {
	router := mux.NewRouter()

	ctx := context.Background()

	cfg := config.Load("configs/config.yaml")

	authMiddleware := authgate.NewMiddleware(ctx, cfg)

	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.Use(authMiddleware.Auth)

	authRouter.HandleFunc("/test", TestHandler).Methods(http.MethodPost)

	loginRouter := router.PathPrefix("/login").Subrouter()
	loginRouter.Use(authMiddleware.Login)

	loginRouter.HandleFunc("/test", TestHandler).Methods(http.MethodPost)

	server := http.Server{
		Addr:              ":8001",
		Handler:           router,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(fmt.Errorf("listen: %w", err))
	}
}

func TestHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("OK"))
}
```

## Token Refresh Mechanism

A token refresh process has been introduced via `refresh.go`, ensuring that authentication tokens remain valid. The
function `RunRefresh(ctx, cfg)` facilitates automatic refresh of Apple and Google authentication tokens. Below is a
high-level overview of its functionality:

```go
func RunRefresh(ctx context.Context, cfg config.Config) error {
db, err := sqlx.ConnectContext(ctx, "mysql", cfg.SQL.SQLDsn)
if err != nil {
panic(err)
}

sessionsRepo := session.New(db)

appleSecretsFile, err := os.Open(cfg.AppleSignIn.KeyPath)
if err != nil {
panic(fmt.Errorf("open apple_sign_in_key_file: %w", err))
}
defer appleSecretsFile.Close()

appleSecretsContent, err := io.ReadAll(appleSecretsFile)
if err != nil {
panic(fmt.Errorf("read apple_sign_in_key_file: %w", err))
}

srv := authservice.New(
apple.New(
cfg.AppleSignIn,
apple.GenerateClientSecret(appleSecretsContent),
http.DefaultClient,
),
google.New(cfg.GoogleSignIn, http.DefaultClient),
sessionsRepo,
)

r := applerefresh.New(srv)
err = r.Run(ctx)
if err != nil {
return fmt.Errorf("run refresh: %w", err)
}

return nil
}
```

## Configuration Schema

Authentication and service configurations are specified via a YAML file. Below is an exemplary configuration file:

```yaml
appleSignIn:
  url: "https://appleid.apple.com"
  keyPath: "configs/apple_key.p8"
  token:
    clientID: ""
    teamID: ""
    keyID: ""
    audience: "https://appleid.apple.com"
    exp: 2592000

googleSignIn:
  url: "https://oauth2.googleapis.com"
  appID: ""

jwt:
  secretKey: ""
  domain: ""

sql:
  dsn: "vpn_root:vpn_pass@tcp(127.0.0.1:3306)/vpn"
```

## Execution Instructions

1. Duplicate the configuration template as `configs/config.yaml` and populate it with appropriate credentials and
   parameters.
2. Launch the authentication service by executing:

   ```sh
   go run main.go
   ```

3. Run the token refresh mechanism:
   ```sh
   go run refresh.go
   ```

4. Verify endpoint functionality:
   ```sh
   curl -X POST http://localhost:8001/auth/test
   ```

## Licensing

This software is distributed under the MIT License, permitting unrestricted use, modification, and distribution while
maintaining attribution to the original authorship.

