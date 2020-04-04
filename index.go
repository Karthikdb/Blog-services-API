package main

import (
	//"encoding/json"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

type App struct {
	user_id        string `json:"user_id"`
	category_id    string `json:"category_id"`
	content        string `json:"content"`
	heading        string `json:"heading"`
	private_public string `json:"private_public"`
}

func main() {
	http1 := http.NewServeMux()
	http1.HandleFunc("/register_user", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println(json.Marshal(register("1")));
		r.ParseForm()
		fmt.Println(r.Form)
		w.Write(register(r.FormValue("user_name"), r.FormValue("email"), r.FormValue("password"), r.FormValue("mobile"), r.FormValue("dob"), r.FormValue("occ"), r.FormValue("sub"), r.FormValue("hint")))
	})
	http1.HandleFunc("/login_user", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println(json.Marshal(register("1")));
		r.ParseForm()
		fmt.Println(r.Form)
		w.Write(login(r.FormValue("email"), r.FormValue("password")))

	})
	http1.HandleFunc("/get_occupations", func(w http.ResponseWriter, r *http.Request) {
		w.Write(occupation())
	})
	http1.HandleFunc("/get_sub_occupations", func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Content-Type", "application/json")
		w.Write(sub_occupation())
	})
	http1.HandleFunc("/view_posts", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fmt.Println(r.Form)
		w.Write(get_posts(r.FormValue("post_id"), r.FormValue("user_id"), r.FormValue("category_id")))
	})
	http1.HandleFunc("/preference_list", func(w http.ResponseWriter, r *http.Request) {
		w.Write(preference())
	})
	http1.HandleFunc("/add_preference", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fmt.Println(r.Form)
		w.Write(other_preference(r.FormValue("other_pref")))
	})
	http1.HandleFunc("/insert_post", func(w http.ResponseWriter, r *http.Request) {
		v := make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			panic(err)
		}
		fmt.Sprint(v["content"])
		userid := fmt.Sprint(v["user_id"])
		catid := fmt.Sprint(v["category_id"])
		ptpub := fmt.Sprint(v["private_public"])
		head := fmt.Sprint(v["heading"])
		cont := fmt.Sprint(v["content"])
		//w.Write(add_post(v.user_id, v.category_id, v.private_public, v.heading, v.content))
		// w.Write(add_post(v["user_id"], v["category_id"], v["private_public"], v["heading"], v["content"]))
		w.Write(add_post(userid, catid, ptpub, head, cont))
	})
	http1.HandleFunc("/all_posts", func(w http.ResponseWriter, r *http.Request) {
		w.Write(get_all_post())
	})

	handler := cors.Default().Handler(http1)
	http.ListenAndServe(":4200", handler)

}
