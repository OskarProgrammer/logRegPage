package main

import(
	"html/template"
	"log"
	"net/http"
	"os"
	"fmt"
	"bufio"
	"strings"
	"regexp"
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
		match, err := regexp.MatchString("^[a-zA-Z0-9][^!@^&*#$%]{1,20}$", password)

		if password != nick && passwordConfirm == password && match == true{
				options := os.O_APPEND | os.O_WRONLY
	
				file, err := os.OpenFile("database.txt", options,os.FileMode(0600))
				errorCheck(err)

				_, err = fmt.Fprintf(file,"%s %s\n",nick,password)
				errorCheck(err)

				err = file.Close()
				errorCheck(err)

				http.Redirect(writer, request, "/createAccount/successRegistered", http.StatusFound)
		}else if (passwordConfirm != password || match == true) && nick != password{
			http.Redirect(writer, request, "/createAccount/failedRegistered", http.StatusFound)
		}else if (passwordConfirm == password || match != true) && nick != password{
			http.Redirect(writer, request, "/createAccount/failedRegistered", http.StatusFound)
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
				http.Redirect(writer, request, "/login/successLogin", http.StatusFound)
			}
		}
		var i int = 1
		if nick != password && i == 1{
			http.Redirect(writer, request, "/login/failedLogin", http.StatusFound)
		}
		err = tmpl.Execute(writer, nil)
	}

func handleFailedLogin(writer http.ResponseWriter,
	request *http.Request){
		tmpl, err := template.ParseFiles("failedLogin.html")
		errorCheck(err)
		err = tmpl.Execute(writer, nil)
	}

func handleFailed(writer http.ResponseWriter,
	request *http.Request){
		tmpl, err := template.ParseFiles("failed.html")
		errorCheck(err)
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
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/createAccount", handleCreate)

	http.HandleFunc("/login/successLogin",handleSuccess)
	http.HandleFunc("/createAccount/successRegistered",handleSuccessReg)

	http.HandleFunc("/login/failedLogin",handleFailedLogin)
	http.HandleFunc("/createAccount/failedRegistered",handleFailed)

	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}