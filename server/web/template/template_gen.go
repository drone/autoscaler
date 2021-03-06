package template

import "html/template"

// list of embedded template files.
var files = []struct {
	name string
	data string
}{
	{
		name: "index.tmpl",
		data: index,
	}, {
		name: "logs.tmpl",
		data: logs,
	},
}

// T exposes the embedded templates.
var T *template.Template

func init() {
	T = template.New("_").Funcs(funcMap)
	for _, file := range files {
		T = template.Must(
			T.New(file.name).Parse(file.data),
		)
	}
}

//
// embedded template files.
//

// files/index.tmpl
var index = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<meta http-equiv="refresh" content="30">
<title>Dashboard</title>
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/style.css">
<link rel="icon" type="image/png" id="favicon" href="/static/favicon.png">
<script src="/static/timeago.js" type="text/javascript"></script>
</head>
<body>

<header class="navbar">
    <div class="logo">
        <svg viewBox="0 0 60 60" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"><defs><path d="M12.086 5.814l-.257.258 10.514 10.514C20.856 18.906 20 21.757 20 25c0 9.014 6.618 15 15 15 3.132 0 6.018-.836 8.404-2.353l10.568 10.568C48.497 55.447 39.796 60 30 60 13.434 60 0 46.978 0 30 0 19.903 4.751 11.206 12.086 5.814zm5.002-2.97C20.998 1.015 25.378 0 30 0c16.566 0 30 13.022 30 30 0 4.67-1.016 9.04-2.835 12.923l-9.508-9.509C49.144 31.094 50 28.243 50 25c0-9.014-6.618-15-15-15-3.132 0-6.018.836-8.404 2.353l-9.508-9.508zM35 34c-5.03 0-9-3.591-9-9s3.97-9 9-9c5.03 0 9 3.591 9 9s-3.97 9-9 9z" id="a"></path></defs><use fill="#fff" xlink:href="#a" fill-rule="evenodd"></use></svg>
    </div>
    <nav class="inline-nav">
        <ul>
            <li><a href="#" class="active">Servers</a></li>
            <li><a href="ui/logs">Logging</a></li>
        </ul>
    </nav>
</header>

<main>
    <section>
        <header>
            <h1>Servers</h1>
        </header>
        <article class="cards stages">
            {{ if not .Items }}
            <div class="card alert sleeping">
                <span>
                    <h1>There are no active servers</h1>
                    <p>The system will not provision instances when the queue is empty.</p>
                </span>
                <img src="/static/icons/server-list-empty.svg" />
            </div>
            {{ else }}
                {{ range .Items }}
                <div class="card instance">
                    <div class="icon">
                        <svg viewBox="0 0 24 24" class="icon-server icon-server-{{ .State }}"><path class="primary" d="M5 3h14a2 2 0 0 1 2 2v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5c0-1.1.9-2 2-2zm0 10h14a2 2 0 0 1 2 2v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4c0-1.1.9-2 2-2zm2 3a1 1 0 0 0 0 2h3a1 1 0 0 0 0-2H7z"/><rect width="5" height="2" x="6" y="6" class="secondary" rx="1"/></svg>
                    </div>
                    <div class="addr">{{ .Name }}</div>
                    <div class="id">{{ if .Address }}{{ .Address }}{{ else }}0.0.0.0{{ end }}</div>
                    <div class="state">
                        <span class="badge badge-{{ .State }}">
                            {{ if eq .State "error" }}
                            <svg viewBox="0 0 24 24" class="icon-close-circle"><path class="primary" d="M12 2a10 10 0 1 1 0 20 10 10 0 0 1 0-20z"/><path class="secondary" d="M12 18a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3zm1-5.9c-.13 1.2-1.88 1.2-2 0l-.5-5a1 1 0 0 1 1-1.1h1a1 1 0 0 1 1 1.1l-.5 5z"/></svg>
                            {{ else }}
                            <svg viewBox="0 0 24 24" class="icon-dot"><path d="M12,10A2,2 0 0,0 10,12C10,13.11 10.9,14 12,14C13.11,14 14,13.11 14,12A2,2 0 0,0 12,10Z" /></svg>
                            {{ end }}
                            {{ .State }}
                        </span>
                    </div>
                    <div class="region">{{ .Region }}</div>
                    <div class="image">{{ .Image }}</div>
                    <div class="size">{{ .Size }}</div>
                    <div class="time" datetime="{{ timestamp .Created }}"></div>
                </div>
                {{ end }}
            {{ end }}
        </article>
    </section>
</main>

<footer></footer>

<script>
timeago.render(document.querySelectorAll('.time'));
</script>
</body>
</html>`

// files/logs.tmpl
var logs = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Dashboard</title>
<link rel="stylesheet" type="text/css" href="/static/reset.css">
<link rel="stylesheet" type="text/css" href="/static/style.css">
<link rel="icon" type="image/png" id="favicon" href="/static/favicon.png">
</head>
<body>

<header class="navbar">
    <div class="logo">
        <svg viewBox="0 0 60 60" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"><defs><path d="M12.086 5.814l-.257.258 10.514 10.514C20.856 18.906 20 21.757 20 25c0 9.014 6.618 15 15 15 3.132 0 6.018-.836 8.404-2.353l10.568 10.568C48.497 55.447 39.796 60 30 60 13.434 60 0 46.978 0 30 0 19.903 4.751 11.206 12.086 5.814zm5.002-2.97C20.998 1.015 25.378 0 30 0c16.566 0 30 13.022 30 30 0 4.67-1.016 9.04-2.835 12.923l-9.508-9.509C49.144 31.094 50 28.243 50 25c0-9.014-6.618-15-15-15-3.132 0-6.018.836-8.404 2.353l-9.508-9.508zM35 34c-5.03 0-9-3.591-9-9s3.97-9 9-9c5.03 0 9 3.591 9 9s-3.97 9-9 9z" id="a"></path></defs><use fill="#FFF" xlink:href="#a" fill-rule="evenodd"></use></svg>
    </div>
    <nav class="inline-nav">
        <ul>
            <li><a href="../">Servers</a></li>
            <li><a href="#" class="active">Logging</a></li>
        </ul>
    </nav>
</header>

<main>
    <section>
        <header>
            <h1>Recent Logs</h1>
        </header>

        <div class="cards">
            {{ range .Entries }}
            <div class="card entry">
                <div class="level">
                    <span class="badge badge-{{ .Level }}">{{ .Level }}</span>
                </div>
                <div class="message">{{ .Message }}</div>
                <div class="fields">
                    {{ range $key, $val := .Data }}
                    <span><em>{{ $key }}</em>{{ $val }}</span>
                    {{ end }}
                </div>
                <div class="time" datetime="{{ timestamp .Unix }}">{{ timestamp .Unix }}</div>
            </div>
            {{ end }}
        </div>

    </section>
</main>

<footer></footer>
</body>
</html>`
