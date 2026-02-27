package api

import (
    "net/http"
    "os"
    "strings"
)

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
func SwaggerHandler(oktaDomain, clientID string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        specURL := "/openapi.yaml"
        oauth2Redirect := r.URL.Scheme + "://" + r.Host + "/docs/oauth2-redirect.html"
        // r.URL.Scheme may be empty (Go's request only populates it when
        // behind proxy), so derive from header if necessary
        if oauth2Redirect == "://"+r.Host+"/docs/oauth2-redirect.html" {
            scheme := "http"
            if r.TLS != nil {
                scheme = "https"
            }
            oauth2Redirect = scheme + "://" + r.Host + "/docs/oauth2-redirect.html"
        }

        html := strings.ReplaceAll(swaggerHTML, "${SPEC_URL}", specURL)
        html = strings.ReplaceAll(html, "${OAUTH2_REDIRECT}", oauth2Redirect)
        // OKTA_DOMAIN is really the issuer base URL; pass as-is
        html = strings.ReplaceAll(html, "${OKTA_DOMAIN}", oktaDomain)
        html = strings.ReplaceAll(html, "${CLIENT_ID}", clientID)
        w.Header().Set("Content-Type", "text/html")
        w.Write([]byte(html))
    }
}

// OAuthRedirectHandler serves the OAuth2 redirect page used by Swagger UI
func OAuthRedirectHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte(oauthRedirectHTML))
}

const swaggerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <title>Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
  <script>
  window.onload = function() {
    const ui = SwaggerUIBundle({
      url: "${SPEC_URL}",
      dom_id: '#swagger-ui',
      presets: [SwaggerUIBundle.presets.apis],
      layout: "BaseLayout",
      oauth2RedirectUrl: "${OAUTH2_REDIRECT}",
      // you can also set clientId here but we'll init explicitly below
    });
    window.ui = ui;

    // initialize OAuth settings with client id (no secret)
    ui.initOAuth({
      clientId: "${CLIENT_ID}",
      usePkceWithAuthorizationCodeGrant: true,
      // leaving clientSecret unset since PKCE is used
    });

    // hide only the client_id field; secret remains visible but we mark it optional
    const style = document.createElement('style');
    style.textContent =
      "/* Swagger UI modal uses .dialog-ux class for the form */\n" +
      " .dialog-ux input[name=\"client_id\"],\n" +
      " .dialog-ux label[for=\"client_id\"] {\n" +
      "     display: none !important;\n" +
      " }\n";
    document.head.appendChild(style);

    // once the modal appears, prefill client_id and annotate the secret field
    const observer = new MutationObserver(() => {
      const cidInput = document.querySelector('.dialog-ux input[name="client_id"]');
      if (cidInput) {
        cidInput.value = "${CLIENT_ID}";
      }
      const secretInput = document.querySelector('.dialog-ux input[name="client_secret"]');
      if (secretInput) {
        secretInput.placeholder = "optional – PKCE is used, leave blank";
        secretInput.disabled = true;
      }
    });
    observer.observe(document.body, { childList: true, subtree: true });

    tokenBox.placeholder = 'Bearer token will appear here after authorization';
    const container = document.createElement('div');
    container.style.margin = '10px 0';
    container.appendChild(tokenBox);
    document.body.insertBefore(container, document.getElementById('swagger-ui'));

    // poll for token after auth (Swagger UI stores it internally)
    function updateToken() {
      try {
        const at = ui.authActions && ui.authActions.getAccessToken && ui.authActions.getAccessToken();
        if (at && Object.keys(at).length) {
          // pick first token value
          const val = Object.values(at)[0];
          if (val) {
            tokenBox.value = val;
          }
        }
      } catch (e) {
        // ignore until ui is ready
      }
    }
    // run periodically for a short while
    const interval = setInterval(() => {
      updateToken();
    }, 1000);
    setTimeout(() => clearInterval(interval), 60000);
  }
  </script>
</body>
</html>`

const oauthRedirectHTML = `<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"/><title>OAuth2 Redirect</title></head>
<body>
<script>
if (window.opener && window.opener.swaggerUIRedirectCallback) {
  window.opener.swaggerUIRedirectCallback(window.location.href);
}
</script>
</body>
</html>`
