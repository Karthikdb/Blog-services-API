package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
)

//local
// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "harish"
// 	password = "root"
// 	dbname   = "soa"
// )

// var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
// 	"password=%s dbname=%s sslmode=disable",
// 	host, port, user, password, dbname)

//Online
const (
	host     = "ec2-174-129-253-113.compute-1.amazonaws.com"
	port     = 5432
	user     = "txhwozuaedengx"
	password = "504cc96c5899d58072600873b4a013333a8721cebf89800fabd69ba6b796268d"
	dbname   = "dbe7el2krjpufv"
)

var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=require",
	host, port, user, password, dbname)

func register(user_name string, email string, password string, mobile string, dob string, occ string, sub string, hint string) []byte {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	stmt, err := db.Prepare("insert into users(user_name,user_mail_id,password,user_mobile_number,user_dob,occupation_id,sub_occupation_id,hint) values ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING user_id;")
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	fmt.Print(user_name, " ", email, " ", password, " ", mobile, " ", dob, " ", occ, " ", sub, " ", hint)
	var id string
	err = stmt.QueryRow(user_name, email, password, mobile, dob, occ, sub, hint).Scan(&id)
	// if err == nil {
	// 	return []byte("failed")
	// }
	obj := []byte(`{ "Id": ` + id + `, "message": "Created successfully"}`)
	fmt.Println(id)
	defer db.Close()
	return obj
}

func add_post(user_id string, category_id string, private_public string, heading string, content string) []byte {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	var stmt string
	stmt = "insert into posts(user_id,category_id,post_date,private_public,heading,content) values (" + user_id + "," + category_id + ", current_timestamp ,'" + private_public + "','" + heading + "','" + content + "') RETURNING post_id;"
	fmt.Print(user_id, category_id, private_public, heading, content)
	var id string
	var msg string
	err = db.QueryRow(stmt).Scan(&id)
	if id == "" {
		id = "0"
		msg = "Failed"
	} else {
		fmt.Print(err)
		msg = "Added"
	}
	obj := []byte(`{ "Id": ` + id + `, "message": ` + msg + `}`)
	defer db.Close()
	return obj
}

func other_preference(preference_name string) []byte {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	stmt, err := db.Prepare("INSERT INTO preference(preference_name) SELECT $1 WHERE NOT EXISTS (SELECT preference_id FROM preference WHERE preference_name = $1) RETURNING preference_id;")
	defer stmt.Close()
	var id string
	var msg string
	err = stmt.QueryRow(preference_name).Scan(&id)
	if id == "" {
		id = "0"
		msg = "Failed"
	} else {
		msg = "Added"
	}
	obj := []byte(`{ "Id": ` + id + `, "message": ` + msg + `}`)

	fmt.Println(id)
	defer db.Close()
	return obj
}

func login(user_mail_id string, passsword string) []byte {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	stmt := "select user_name from users where user_mail_id = $1 AND password = $2;"
	var user_name string
	var bit string
	row := db.QueryRow(stmt, user_mail_id, passsword)
	switch err := row.Scan(&user_name); err {
	case sql.ErrNoRows:
		bit = "0"
		user_name = "No rows were returned!"
	case nil:
		bit = "1"
		fmt.Println(user_name)
	default:
		panic(err)
	}
	obj := []byte(`{ "bit": ` + bit + `, "Status": "` + user_name + ` Exists" }`)
	fmt.Println(user_name, bit)
	defer db.Close()
	return obj

}

type Occupation struct {
	Occupation string `json:"occ"`
	Occ_Id     int    `json:"occ_id"`
}

var occ struct {
	Occ_Data []Occupation `json:"Occ_data"`
}

type SubOccupation struct {
	Sub_Occupation string `json:"sub_occ"`
	Sub_Id         int    `json:"sub_id"`
	Occ_id         int    `json:"occ_id"`
}

var sub struct {
	Sub_Data []SubOccupation `json:"Sub_data"`
}

func occupation() []byte {
	occ.Occ_Data = nil
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("select occ.occupation_name,occ.occupation_id from occupation as occ ;")
	defer rows.Close()
	for rows.Next() {
		// Scan one customer record
		var c Occupation
		if err := rows.Scan(&c.Occupation, &c.Occ_Id); err != nil {
			// handle error
		}
		// fmt.Print(c.sub_occupation)
		occ.Occ_Data = append(occ.Occ_Data, c)
	}
	if rows.Err() != nil {
		// handle error
		fmt.Println("Failed")
	}

	obj, err := json.Marshal(occ)
	defer db.Close()
	//fmt.Print(string(obj))
	return obj
}

type Preference struct {
	Preference string `json:"pre"`
	Pre_Id     int    `json:"pre_id"`
}

var pre struct {
	Pre_Data []Preference `json:"Pre_data"`
}

func preference() []byte {
	pre.Pre_Data = nil
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("select preference_name,preference_id from preference ;")
	defer rows.Close()
	for rows.Next() {
		// Scan one customer record
		var c Preference
		if err := rows.Scan(&c.Preference, &c.Pre_Id); err != nil {
			// handle error
		}
		// fmt.Print(pre.Pre_Data)
		pre.Pre_Data = append(pre.Pre_Data, c)
	}
	if rows.Err() != nil {
		// handle error
		fmt.Println("Failed")
	}

	obj, err := json.Marshal(pre)
	defer db.Close()
	//fmt.Print(string(obj))
	return obj
}

type all_post struct {
	Post_id        int    `json:"p_id"`
	User_id        int    `json:"u_id"`
	Category_id    int    `json:"c_id"`
	Post_date      string `json:"p_date"`
	Private_public string `json:"prt_public"`
	Heading        string `json:"head"`
	Content        string `json:"cont"`
	// status         string `json:"status"`
}

var all_post_b struct {
	All_post_Data []all_post `json:"All_post_data"`
}

func get_all_post() []byte {
	all_post_b.All_post_Data = nil
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("select post_id,user_id,category_id,post_date,private_public,heading,content from posts where private_public = '0' ;")
	defer rows.Close()
	for rows.Next() {
		// Scan one customer record
		var c all_post
		if err := rows.Scan(&c.Post_id, &c.User_id, &c.Category_id, &c.Post_date, &c.Private_public, &c.Heading, &c.Content); err != nil {
			// handle error
		}
		//fmt.Print(all_post_b.All_post_Data)
		all_post_b.All_post_Data = append(all_post_b.All_post_Data, c)
	}

	obj, err := json.Marshal(all_post_b)
	//fmt.Println(obj)
	// obj := []byte(`{ "status": ` + status + `, "post_id": ` + post_id + `, "user_id": ` + user_id + `, "post_date": "` + post_date + `" , "private_public":` + private_public + `, "heading" : "` + heading + `" , "content" : "` + content + `" }`)
	defer db.Close()
	//fmt.Print(string(obj))
	return obj
}
func get_posts(P_Id string, U_Id string, C_Id string) []byte {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	fmt.Print(P_Id, U_Id, C_Id)
	stmt := "select post_id,user_id,category_id,post_date,private_public,heading,content from posts where post_id=$1 AND user_id=$2 AND category_id=$3;"
	var post_id string
	var user_id string
	var category_id string
	var post_date string
	var private_public string
	var heading string
	var content string
	// var q_as string
	// var rating_avg string &q_as, &rating_avg //, "q_as" : ` + q_as + `, "rating_avg": ` + rating_avg + `
	var status string
	var obj []byte
	row := db.QueryRow(stmt, P_Id, U_Id, C_Id)
	switch err := row.Scan(&post_id, &user_id, &category_id, &post_date, &private_public, &heading, &content); err {
	case sql.ErrNoRows:
		status = "0"
	case nil:
		status = "1"
	default:
		panic(err)
	}
	//fmt.Print(content)
	obj = []byte(`{ "status": ` + status + `, "post_id": ` + post_id + `, "user_id": ` + user_id + `, "post_date": "` + post_date + `" , "private_public":` + private_public + `, "heading" : "` + heading + `" , "content" : "` + content + `" }`)
	defer db.Close()
	//fmt.Print(string(obj))
	return obj
}

func sub_occupation() []byte {
	sub.Sub_Data = nil
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("select sub.sub_occupation_name,sub.sub_occupation_id,sub.occupation_id from sub_occupation as sub ;")
	defer rows.Close()
	for rows.Next() {
		// Scan one customer record
		var c SubOccupation
		if err := rows.Scan(&c.Sub_Occupation, &c.Sub_Id, &c.Occ_id); err != nil {
			// handle error
		}
		// fmt.Print(c.sub_occupation)
		sub.Sub_Data = append(sub.Sub_Data, c)
	}
	if rows.Err() != nil {
		// handle error
		fmt.Println("Failed")
	}

	obj, err := json.Marshal(sub)
	defer db.Close()
	//fmt.Print(string(obj))
	return obj
}
func init() {

	// var (
	// 	mail     string
	// 	password string
	// )
	// rows, err := db.Query("select admin_email,admin_password from admin ")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	err := rows.Scan(&mail, &password)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(mail, password)
	// }
	// err = rows.Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	fmt.Println("Successfully connected!")
}
