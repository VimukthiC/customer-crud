package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	_ "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	//_ "github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

var connection *sql.DB

type Customer struct {
	CusId string `json:"cusId"`
	Name string `json:"name"`
	Address string `json:"address"`
	Salary string `json:"salary"`
}

func CreateDBConnection()  {
	db, err := sql.Open("mysql", "root:ijse@tcp(127.0.0.1:3307)/customer_detail")
	if err!=nil{
		fmt.Println(err.Error())
		panic("Failed connect to database")
	}
	connection=db
}

func GetAllCustomer(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","application/json")
	var customers []Customer

	row, err := connection.Query("select * from customer")
	if err != nil {
		panic(err.Error())
	} else{
		for row.Next(){
			var cusId string
			var name string
			var address string
			var salary string
			err2 := row.Scan(&cusId, &name, &address, &salary)
			if err2 != nil{
				panic(err2.Error())
			}else{
				customer := Customer{
					CusId:cusId,
					Name:name,
					Address: address,
					Salary: salary,
				}
				customers = append(customers, customer)
			}
		}
	}
	defer row.Close()
	json.NewEncoder(w).Encode(customers)
}

func GetCustomer(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","application/json")

	var customerObject Customer
	params := mux.Vars(r)

	row, err := connection.Query("select * from customer where cusId=?", params["cusId"])
	if err != nil {
		panic(err.Error())
	}else{
		for row.Next(){
			var cusId string
			var name string
			var address string
			var salary string
			err2 := row.Scan(&cusId, &name, &address, &salary)
			if err2 != nil{
				panic(err2.Error())
			}else{
				customer := Customer{
					CusId:cusId,
					Name:name,
					Address: address,
					Salary: salary,
				}
				customerObject = customer
			}
		}
	}
	defer row.Close()
	json.NewEncoder(w).Encode(customerObject)
}

func Save(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","application/json")
	var customer Customer
	_=json.NewDecoder(r.Body).Decode(&customer)

	insert, err := connection.Query("insert into customer (cusId, name, address, salary) values (?, ?, ?, ?)", customer.CusId, customer.Name, customer.Address, customer.Salary)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	json.NewEncoder(w).Encode("Success")
}

func Delete(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","application/json")
	params := mux.Vars(r)
	delete, err := connection.Query("delete from customer where cusId=?", params["cusId"])
	fmt.Println( params["cusId"])
	if err != nil {
		panic(err.Error())
	}
	defer delete.Close()
	json.NewEncoder(w).Encode("Success")
}

func Update(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","application/json")
	params := mux.Vars(r)
	var customer Customer
	_= json.NewDecoder(r.Body).Decode(&customer)


	update, err := connection.Query("update customer set name=? , address=? , salary=? where cusId=? ", customer.Name, customer.Address,customer.Salary,params["cusId"])
	fmt.Println( params["cusId"])
	if err != nil {
		panic(err.Error())
	}
	defer update.Close()
	json.NewEncoder(w).Encode(customer)
}


func handleRequests(){
	router := mux.NewRouter().StrictSlash(true)
	headers:= handlers.AllowedHeaders([]string{"X-Requested-With","Content-Type","Authorization"})
	methods := handlers.AllowedMethods([]string{"GET","POST","PUT","DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})


	router.HandleFunc("/user",GetAllCustomer).Methods("GET")
	router.HandleFunc("/user/{cusId}",GetCustomer).Methods("GET")
	router.HandleFunc("/user",Save).Methods("POST")
	router.HandleFunc("/user/{cusId}",Delete).Methods("DELETE")
	router.HandleFunc("/user/{cusId}",Update).Methods("PUT")

	http.ListenAndServe(":3000", handlers.CORS(headers,methods,origins)(router))
}

func main()  {
	fmt.Println("Server Started")

	CreateDBConnection()

	handleRequests()
}
