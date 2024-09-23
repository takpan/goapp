package httpsrv

import (
	"goapp/pkg/util"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
)

const wsPath = "/goapp/ws"

type TemplateData struct {
	WsUrl     string
	CSRFToken template.HTML
}

func (s *Server) handlerHome(w http.ResponseWriter, r *http.Request) {
	// Create new session
	session, err := s.cookieStore.Get(r, "ws-session")
	if err != nil {
		http.Error(w, "Unable to retrieve session", http.StatusInternalServerError)
		return
	}

	// Generate a CSRF token
	csrfToken, err := util.GenerateKey(32)
	if err != nil {
		http.Error(w, "Error generating CSRF token", http.StatusInternalServerError)
		return
	}

	// Set session options
	session.Options = &sessions.Options{
		Path:     wsPath,
		MaxAge:   1800, // 30 min
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	// Store the CSRF token in the session cookie
	session.Values["csfr_token"] = csrfToken
	if err = session.Save(r, w); err != nil {
		http.Error(w, "Error saving session", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{ .WsUrl }}?csrf_token=" + encodeURIComponent({{ .CSRFToken }}));
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="{}">
<button id="send">Reset</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))

	data := TemplateData{
		WsUrl:     "ws://" + r.Host + wsPath,
		CSRFToken: template.HTML(csrfToken),
	}

	tmpl.Execute(w, data)
}
