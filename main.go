package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/connection"
	"personal-web/middleware"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type MetaData struct {
	Id        int
	Title     string
	IsLogin   bool
	Email     string
	FlashData string
}

var Data = MetaData{
	Title: "Personal Web",
}

type Project struct {
	Id          int
	ProjectName string
	SDate       time.Time
	StartDate   string
	EDate       time.Time
	EndDate     string
	Durasi      string
	Description string
	Node        string
	React       string
	TypeScript  string
	Next        string
	Image       string
	AuthId      int
	IsLogin     bool
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnection()
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	route.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/home", home).Methods("GET").Name("home")
	route.HandleFunc("/project", project).Methods("GET")
	route.HandleFunc("/project/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/project", middleware.UploadFile(addProject)).Methods("POST")
	route.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/edit-project/{id}", eProject).Methods("GET")
	route.HandleFunc("/edit-project/{id}", editProject).Methods("POST")
	route.HandleFunc("/register", register).Methods("GET")
	route.HandleFunc("/register", regProcess).Methods("POST")
	route.HandleFunc("/login", login).Methods("GET")
	route.HandleFunc("/login", loginProcess).Methods("POST")
	route.HandleFunc("/logout", logout).Methods("GET")

	fmt.Println("Server Is Running On Port 5000")
	http.ListenAndServe("localhost:5000", route)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message: " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Email = session.Values["Email"].(string)
	}
	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)

		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}
	Data.FlashData = strings.Join(flashes, "")

	rows, _ := connection.Conn.Query(context.Background(), "select tb_project.id, project_name, start_date, end_date, description, node, react, next, typescript, image, tb_user.id as author from tb_project left join tb_user on tb_user.id = tb_project.auth_id order by id desc")

	var result []Project
	for rows.Next() {
		var each = Project{}
		var err = rows.Scan(&each.Id, &each.ProjectName, &each.SDate, &each.EDate, &each.Description, &each.Node, &each.React, &each.Next, &each.TypeScript, &each.Image, &each.AuthId)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		each.StartDate = each.SDate.Format("2 Februari 2000")
		each.EndDate = each.EDate.Format("2 Februari 2000")
		if session.Values["IsLogin"] != true {
			each.IsLogin = false
		} else {
			each.IsLogin = session.Values["IsLogin"].(bool)
		}
		result = append(result, each)
	}

	respData := map[string]interface{}{
		"Data":     Data,
		"Projects": result,
	}
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/add-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message" + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Email = session.Values["Email"].(string)
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var tmpl, err = template.ParseFiles("views/detail-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message" + err.Error()))
		return
	}

	ProjectDetail := Project{}
	err = connection.Conn.QueryRow(context.Background(), "select id, project_name, start_date, end_date, description, node, react, next, typescript, image, auth_id from tb_project where id = $1", id).Scan(
		&ProjectDetail.Id, &ProjectDetail.ProjectName, &ProjectDetail.SDate, &ProjectDetail.EDate, &ProjectDetail.Description, &ProjectDetail.Node, &ProjectDetail.React, &ProjectDetail.Next, &ProjectDetail.TypeScript, &ProjectDetail.Image, &ProjectDetail.AuthId,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}
	ProjectDetail.StartDate = ProjectDetail.SDate.Format("2 Januari 2006")
	ProjectDetail.EndDate = ProjectDetail.EDate.Format("2 Januari 2006")

	resp := map[string]interface{}{
		"Data":     Data,
		"Projects": ProjectDetail,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	// fmt.Println("Title : " + r.PostForm.Get("name"))
	projectName := r.PostForm.Get("name")
	description := r.PostForm.Get("description")
	Node := r.PostForm.Get("node")
	React := r.PostForm.Get("react")
	Next := r.PostForm.Get("next")
	TypeScript := r.PostForm.Get("typescript")
	AuthId := session.Values["Id"].(int)
	fmt.Println(AuthId)

	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)
	_, err = connection.Conn.Exec(context.Background(), "insert into tb_project(project_name, description, node, react, next, typescript, image, auth_id) values($1, $2, $3, $4, $5, $6, $7, $8)", projectName, description, Node, React, Next, TypeScript, image, AuthId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content/type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	_, err := connection.Conn.Exec(context.Background(), "delete from tb_project where id = $1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message: " + err.Error()))
	}
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func eProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var tmpl, err = template.ParseFiles("views/edit-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message" + err.Error()))
		return
	}

	ProjectDetail := Project{}
	err = connection.Conn.QueryRow(context.Background(), "select id, project_name, start_date, end_date, description, node, react, next, typescript, image, auth_id from tb_project where id = $1", id).Scan(
		&ProjectDetail.Id, &ProjectDetail.ProjectName, &ProjectDetail.SDate, &ProjectDetail.EDate, &ProjectDetail.Description, &ProjectDetail.Node, &ProjectDetail.React, &ProjectDetail.Next, &ProjectDetail.TypeScript, &ProjectDetail.Image, &ProjectDetail.AuthId,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	ProjectDetail.AuthId = 2
	ProjectDetail.StartDate = ProjectDetail.SDate.Format("2 Januari 2006")
	ProjectDetail.EndDate = ProjectDetail.EDate.Format("2 Januari 2006")

	resp := map[string]interface{}{
		"Data":     Data,
		"Projects": ProjectDetail,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func editProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	projectName := r.PostForm.Get("name")
	// startDate := time.Now().String()
	// endDate := time.Now().String()
	// durasi := "10 Bulan"
	description := r.PostForm.Get("description")
	Node := r.PostForm.Get("node")
	React := r.PostForm.Get("react")
	Next := r.PostForm.Get("next")
	TypeScript := r.PostForm.Get("typescript")

	_, err = connection.Conn.Exec(context.Background(), "update tb_project set project_name = $1, description = $2, node = $3, react = $4, next = $5, typescript = $6 where id = $7", projectName, description, Node, React, Next, TypeScript, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact-form.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message" + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Email = session.Values["Email"].(string)
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/register.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message" + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func regProcess(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	passHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "insert into tb_user(name, email, password) values($1, $2, $3)", name, email, passHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message: " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	session.AddFlash("Successfully register!", "message")

	session.Save(r, w)

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/login.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message" + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func loginProcess(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	user := User{}
	err = connection.Conn.QueryRow(context.Background(), "select * from tb_user where email = $1", email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message: " + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message: " + err.Error()))
		return
	}

	session.Values["Id"] = user.Id
	session.Values["IsLogin"] = true
	session.Values["Email"] = user.Email
	session.Options.MaxAge = 10800
	session.AddFlash("Successfully Login", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("logout.")
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")
	session.Options.MaxAge = -1 // gak boleh kurang dari 0
	session.Save(r, w)

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
