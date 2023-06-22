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
var logged bool = false
var finished bool = false

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
		tmpl, err := template.ParseFiles("htmls/index.html")
		errorCheck(err)
		err = tmpl.Execute(writer, nil)
	}

func handleCreate(writer http.ResponseWriter,
	request *http.Request){
		tmpl, err := template.ParseFiles("htmls/create.html")
		errorCheck(err)

		nick := request.FormValue("nick")
		password := request.FormValue("password")
		passwordConfirm := request.FormValue("passwordConfirm")
		match, err := regexp.MatchString("^[a-zA-Z0-9][^!@^&*#$%]{1,20}$", password)

		if password != nick && passwordConfirm == password && match == true{
				options := os.O_APPEND | os.O_WRONLY
	
				file, err := os.OpenFile("htmls/database.txt", options,os.FileMode(0600))
				errorCheck(err)

				_, err = fmt.Fprintf(file,"%s %s\n",nick,password)
				errorCheck(err)

				err = file.Close()
				errorCheck(err)

				finished = true

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
		tmpl, err := template.ParseFiles("htmls/login.html")
		errorCheck(err)
		
		nick := request.FormValue("nick")
		password := request.FormValue("password")
		
		loginVals := getStrings("htmls/database.txt")
		
		for _, v:= range loginVals{
			item := strings.Split(v, " ")
			if nick==item[0] && password==item[1]{
				fmt.Printf("%s with password %s logged in\n",nick,password)
				logged = true
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
		tmpl, err := template.ParseFiles("htmls/failedLogin.html")
		errorCheck(err)
		err = tmpl.Execute(writer, nil)
	}

func handleFailed(writer http.ResponseWriter,
	request *http.Request){
		tmpl, err := template.ParseFiles("htmls/failed.html")
		errorCheck(err)
		err = tmpl.Execute(writer, nil)
	}

func handleErrorWithLogging(writer http.ResponseWriter,
	request *http.Request){
		tmpl, err := template.ParseFiles("htmls/notlogin.html")
		errorCheck(err)
		err = tmpl.Execute(writer, nil)
	}

func handleNotFinished(writer http.ResponseWriter,
	request *http.Request){
		tmpl, err := template.ParseFiles("htmls/notFinished.html")
		errorCheck(err)
		err = tmpl.Execute(writer, nil)
	}

func handleSuccess(writer http.ResponseWriter,
	request *http.Request){
		if logged == false {
			http.Redirect(writer, request, "/login/notLoggedIn", http.StatusFound)
		}else if logged == true {
			tmpl, err := template.ParseFiles("htmls/success.html")
			errorCheck(err)
			err = tmpl.Execute(writer, nil)
		}
	}

func handleSuccessReg(writer http.ResponseWriter,
	request *http.Request){
		if finished == false {
			http.Redirect(writer, request, "/createAccount/notFinishedForm",http.StatusFound)
		}else if finished == true {
			tmpl, err := template.ParseFiles("htmls/successRegistered.html")
			errorCheck(err)
			err = tmpl.Execute(writer, nil)
		}
	}

func handleLogout(writer http.ResponseWriter,
	request *http.Request){
		logged = false
		fmt.Printf("Logging out...")
		http.Redirect(writer, request, "/",http.StatusFound)
	}

func handleReset(writer http.ResponseWriter,
	request *http.Request){
		finished = false
		http.Redirect(writer, request, "/",http.StatusFound)
	}
func main(){
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/createAccount", handleCreate)

	http.HandleFunc("/login/successLogin",handleSuccess)
	http.HandleFunc("/createAccount/successRegistered",handleSuccessReg)

	http.HandleFunc("/login/failedLogin",handleFailedLogin)
	http.HandleFunc("/createAccount/failedRegistered",handleFailed)
	http.HandleFunc("/login/notLoggedIn",handleErrorWithLogging)
	http.HandleFunc("/createAccount/notFinishedForm", handleNotFinished)

	http.HandleFunc("/login/logout", handleLogout)
	http.HandleFunc("/createAccount/successRegistered/reset",handleReset)

	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}