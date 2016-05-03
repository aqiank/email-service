package main

import (
    "flag"
    "net/http"
    "os"

    "github.com/codegangsta/negroni"
    "github.com/mailgun/mailgun-go"
)

var mg mailgun.Mailgun

func mailHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "POST":
        from := r.FormValue("from")
        if from == "" {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        subject := r.FormValue("text")
        if subject == "" {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        text := r.FormValue("text")
        if subject == "" {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        if to, ok := r.Form["to"]; !ok {
            w.WriteHeader(http.StatusBadRequest)
            return
        } else {
            for _, v := range to {
                if v == "" {
                    w.WriteHeader(http.StatusBadRequest)
                    return
                }
            }

            message := mailgun.NewMessage(from, subject, text, to...)
            if _, _, err := mg.Send(message); err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                return
            }
        }

        w.WriteHeader(http.StatusOK)

    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}

func init() {
    domain := os.Getenv("DOMAIN")
    apiKey := os.Getenv("API_KEY")
    publicAPIKey := os.Getenv("PUBLIC_KEY")
    mg = mailgun.NewMailgun(domain, apiKey, publicAPIKey)
}

func main() {
    port := flag.String("port", "8080", "server port")

    flag.Parse()

    n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
    n.UseHandlerFunc(mailHandler)
    n.Run(":" + *port)
}
