package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"evolutionary-mcp/backend/internal/auth"
)

//go:generate oapi-codegen -generate server,types -o api.gen.go -package api ../../api/openapi.yaml

// swaggerHandler serves a simple Swagger UI page that points at the
// generated OpenAPI spec. The page uses the official CDN-hosted assets so we
// don't need to check any static files into version control. The UI is
// configured with OAuth2 settings so that users can "Authorize" using the
// same Okta tenant used by the application.
// SpecHandler serves the OpenAPI YAML spec with any runtime placeholders
// replaced. The file on disk still contains {oktaIssuer} so clients don't have
// to know the actual tenant or issuer URL; we substitute it here before returning.
func SpecHandler(oktaIssuer string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("api/openapi.yaml")
		if err != nil {
			http.Error(w, "failed to load spec", http.StatusInternalServerError)
			return
		}
		spec := strings.ReplaceAll(string(data), "{oktaIssuer}", oktaIssuer)
		w.Header().Set("Content-Type", "application/yaml")
		w.Write([]byte(spec))
	}
}

// SwaggerHandler returns an HTTP handler that serves the Swagger UI.
func SwaggerHandler(oktaDomain, swaggerClientID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		specURL := "/openapi.yaml"

		// Determine scheme
		scheme := "http"
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			scheme = "https"
		}

		// Construct absolute URL to ensure exact match with Okta config
		oauth2RedirectURL := fmt.Sprintf("%s://%s/docs/oauth2-redirect.html", scheme, r.Host)

		html := strings.ReplaceAll(swaggerHTML, "${SPEC_URL}", specURL)
		html = strings.ReplaceAll(html, "${OAUTH2_REDIRECT_URL}", oauth2RedirectURL)
		// OKTA_DOMAIN is really the issuer base URL; pass as-is
		html = strings.ReplaceAll(html, "${OKTA_DOMAIN}", oktaDomain)
		html = strings.ReplaceAll(html, "${CLIENT_ID}", swaggerClientID)
		html = strings.ReplaceAll(html, "${SCOPES}", strings.Join(auth.AllScopes, " "))
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}
}

// OAuth2RedirectHandler serves the static HTML file required by Swagger UI's
// OAuth2 flow. This file receives the token/code from the provider and sends
// it back to the main Swagger UI window.
func OAuth2RedirectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(oauth2RedirectHTML))
	}
}

const swaggerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <title>Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
  <style>
    html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
    *, *:before, *:after { box-sizing: inherit; }
    body { margin: 0; background: #fafafa; }
    /* Hide the top bar introduced by StandaloneLayout */
    .topbar { display: none !important; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-standalone-preset.js"></script>
  <script>
  window.onload = function() {
    const ui = SwaggerUIBundle({
      url: "${SPEC_URL}",
      dom_id: '#swagger-ui',
      presets: [
        SwaggerUIBundle.presets.apis,
        SwaggerUIStandalonePreset
      ],
      layout: "StandaloneLayout",
      oauth2RedirectUrl: "${OAUTH2_REDIRECT_URL}",
    });
    window.ui = ui;

    const clientId = "${CLIENT_ID}";
    console.log("Swagger OAuth configured with Client ID:", clientId);
    ui.initOAuth({
      clientId: clientId,
      usePkceWithAuthorizationCodeGrant: true,
      useBasicAuthenticationWithAccessCodeGrant: false,
      scopes: "${SCOPES}"
    });
  }
  </script>
</body>
</html>`

const oauth2RedirectHTML = `<!doctype html>
<html lang="en-US">
<head>
    <title>Swagger UI: OAuth2 Redirect</title>
</head>
<body>
<script>
    'use strict';
    function run () {
        var oauth2 = window.opener.swaggerUIRedirectOauth2;
        var sentState = oauth2.state;
        var redirectUrl = oauth2.redirectUrl;
        var isValid, qp, arr;

        if (/code|token|error/.test(window.location.hash)) {
            qp = window.location.hash.substring(1).replace('?', '&');
        } else {
            qp = location.search.substring(1);
        }

        arr = qp.split("&");
        arr.forEach(function (v,i,_arr) { _arr[i] = '"' + v.replace('=', '":"') + '"';});
        qp = qp ? JSON.parse('{' + arr.join() + '}',
                function (key, value) {
                    return key === "" ? value : decodeURIComponent(value)
                }
        ) : {};

        isValid = qp.state === sentState;

        if ((
          oauth2.auth.schema.get("flow") === "accessCode" ||
          oauth2.auth.schema.get("flow") === "authorizationCode" ||
          oauth2.auth.schema.get("flow") === "authorization_code"
        ) && !oauth2.auth.code) {
            if (!isValid) {
                oauth2.errCb({
                    authId: oauth2.auth.name,
                    source: "auth",
                    level: "warning",
                    message: "Authorization may be unsafe, passed state was changed in server. The passed state wasn't returned from auth server."
                });
            }

            if (qp.code) {
                delete oauth2.state;
                oauth2.auth.code = qp.code;
                oauth2.callback({auth: oauth2.auth, redirectUrl: redirectUrl});
            } else {
                let oauthErrorMsg;
                if (qp.error) {
                    oauthErrorMsg = "["+qp.error+"]: " +
                        (qp.error_description ? qp.error_description+ ". " : "no accessCode received from the server. ") +
                        (qp.error_uri ? "More info: "+qp.error_uri : "");
                }

                oauth2.errCb({
                    authId: oauth2.auth.name,
                    source: "auth",
                    level: "error",
                    message: oauthErrorMsg || "[Authorization failed]: no accessCode received from the server."
                });
            }
        } else {
            oauth2.callback({auth: oauth2.auth, token: qp, isValid: isValid, redirectUrl: redirectUrl});
        }
        window.close();
    }

    if (document.readyState !== 'loading') {
        run();
    } else {
        document.addEventListener('DOMContentLoaded', function () {
            run();
        });
    }
</script>
</body>
</html>`
