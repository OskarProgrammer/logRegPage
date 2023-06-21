package main

import(
	"html/template"
	"log"
	"net/http"
	"os"
	"fmt"
	"bufio"
	"strings"
)
func errorCheck(err error){
	if err != nil{
		log.Fatal(err)
	}
}

func getStrings(filename string) []string {
	var lines []string
	file, err := os.Open(filename)

	if os.IsNotExist(err){
		return nil
	}
	errorCheck(err)

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	errorCheck(scanner.Err())

	return lines
}

func handleMain(writer http.ResponseWriter,
	request *http.Request){
		tmpl, err := template.ParseFiles("index.html")
		errorCheck(err)
		err = tmpl.Execute(writer, nil)
	}

func handleCreate(writer http.ResponseWriter,
	request *http.Request){
		tmpl, err := template.ParseFiles("create.html")
		errorCheck(err)

		nick := request.FormValue("nick")
		password := request.FormValue("password")
		passwordConfirm := request.FormValue("passwordConfirm")

		if password != nick && passwordConfirm == password {
				options := os.O_APPEND | os.O_WRONLY
	
				file, err := os.OpenFile("database.txt", options,os.FileMode(0600))
				errorCheck(err)

				_, err = fmt.Fprintf(file,"%s %s\n",nick,password)
				errorCheck(err)

				err = file.Close()
				errorCheck(err)

				http.Redirect(writer, request, "/successRegistered", http.StatusFound)
		}
		err = tmpl.Execute(writer, nil)
	}

func handleLogin(writer http.ResponseWriter,
	request *http.Request){
		tmpl, err := template.ParseFiles("login.html")
		errorCheck(err)
		
		nick := request.FormValue("nick")
		password := request.FormValue("password")
		
		loginVals := getStrings("database.txt")
		
		for _, v:= range loginVals{
			item := strings.Split(v, " ")
			if nick==item[0] && password==item[1]{
				fmt.Printf("%s with password %s logged in\n",nick,password)
				http.Redirect(writer, request, "/successLogin", http.StatusFound)
			}
			
		}
		// fmt.Printf("%#v\n", loginVals)
		err = tmpl.Execute(writer, nil)
	}

func handleSuccess(writer http.ResponseWriter,
	request *http.Request){
		tmpl, err := template.ParseFiles("success.html")
		errorCheck(err)
		err = tmpl.Execute(writer, nil)
	}

func handleSuccessReg(writer http.ResponseWriter,
	request *http.Request){
		tmpl, err := template.ParseFiles("successRegistered.html")
		errorCheck(err)
		err = tmpl.Execute(writer, nil)
	}

func main(){
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/createAccount", handleCreate)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/successLogin",handleSuccess)
	http.HandleFunc("/successRegistered",handleSuccessReg)

	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}