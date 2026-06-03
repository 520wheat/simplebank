package swagger

import (
	_ "embed"
	"html/template"
	"net/http"
)

//go:embed simple_bank.swagger.json
var swaggerJSON []byte

func UIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("swagger").Parse(swaggerHTML))
		tmpl.Execute(w, nil)
	}
}

func JSONHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(swaggerJSON)
	}
}

const swaggerHTML = `<!DOCTYPE html>
  <html>
  <head>
      <title>SimpleBank API</title>
      <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css">
  </head>
  <body>
      <div id="swagger-ui"></div>
      <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
      <script>
          SwaggerUIBundle({
              url: "/swagger/simple_bank.swagger.json",
              dom_id: "#swagger-ui"
          });
      </script>
  </body>
  </html>`
