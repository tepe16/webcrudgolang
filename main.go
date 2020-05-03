package main 

import (
    "fmt"
    "log"
    "net/http"
	"text/template"
	"io"
	"github.com/gorilla/sessions"
	"os"
	"database/sql"
	_"github.com/go-sql-driver/mysql"


)
type Karyawan struct {
	Id int
	Id_Kar int
	NamaKaryawan string
	Alamat string
	Gambar string
	Email string
	No_Hp string
}

func kon() (db *sql.DB) {
	NamaDriver := "mysql"
	DbUser :="root"
	DbPass := ""
	NamaDb := "crudgolang"

	db, err := sql.Open(NamaDriver, DbUser+":"+DbPass+"@/"+NamaDb)

	if err != nil {
		panic(err.Error())

	}

	return db
}
var tmpl = template.Must(template.ParseGlob("view/*"))

func Index(w http.ResponseWriter, r *http.Request) {

tmpl.ExecuteTemplate(w, "Index", nil)
}



//Karyawan
func InputKaryawan(w http.ResponseWriter, r* http.Request){
	tmpl.ExecuteTemplate(w, "InputKaryawan", nil)
}

func SaveKaryawan(w http.ResponseWriter, r *http.Request){
  db := kon()
    tambah_karyawan :=r.FormValue("tambah_karyawan")
	nama_karyawan := r.FormValue("nama_karyawan")
	alamat := r.FormValue("alamat")
	no_hp := r.FormValue("no_hp")
	password := r.FormValue("password")
	email := r.FormValue("email")
	gambar, header, _ := r.FormFile("gambar")

	if tambah_karyawan =="Tambah Karyawan"{
		if nama_karyawan ==""{
			fmt.Fprintln(w, "isi nama Karyawan")
		} else if alamat ==""{
			fmt.Fprintln(w, "Isi alamat")
		} else if no_hp ==""{
			fmt.Fprintln(w, "Isi no hp")
		} else if password =="" {
			fmt.Fprintln(w, "Isi Password")
		} else if email ==""{
			fmt.Fprintln(w, "Isi Email")
		} else {
		 insert, err := db.Prepare("INSERT INTO karyawan (nama_karyawan,alamat,no_hp,gambar,password,email) VALUES (?,?,?,?,?,?)")
		 if err != nil {
			 panic(err.Error())

		 }
          insert.Exec(nama_karyawan,alamat,no_hp,header.Filename,password,email)
		  out, _  := os.Create("./images/" + header.Filename)
		  _, _ =io.Copy(out, gambar)
		  gambar.Close()
		  out.Close()
		  defer db.Close()

		  http.Redirect(w,r, "/", 301)
		}
	}

}
 func LihatKaryawan(w http.ResponseWriter, r *http.Request){
	 db := kon()
   tampil, err := db.Query("SELECT id_karyawan,nama_karyawan,alamat,gambar,email,no_hp FROM Karyawan")
     if err != nil {
		 panic(err.Error())
	 }
	 kar := Karyawan{}
	 res := []Karyawan{}
	 i := 1
	 for tampil.Next() {
		var nama_karyawan, alamat, email, no_hp, gambar string
		var id_karyawan int
		err = tampil.Scan(&id_karyawan, &nama_karyawan, &alamat, &gambar, &email, &no_hp)
		if err != nil {
			panic(err.Error())
		}
		kar.Id = i
		kar.Id_Kar = id_karyawan
		kar.NamaKaryawan = nama_karyawan
		kar.Alamat = alamat
		kar.Gambar = gambar
		kar.Email = email
		kar.No_Hp = no_hp
		
		
 		res = append(res, kar)
		 i++

	 }
	 
	 tmpl.ExecuteTemplate(w, "LihatKaryawan", res)
	 defer db.Close()
 }

 func HapusKaryawan(w http.ResponseWriter, r *http.Request){
	 db := kon()
	 id := r.URL.Query().Get("id_karyawan")
	hapus, err := db.Prepare("delete from karyawan where id_karyawan=?")
	
	if err != nil {
		panic(err.Error())
	}

	hapus.Exec(id)
	defer db.Close()
	http.Redirect(w, r, "/LihatKaryawan", 301)

 }

 func DetailKaryawan(w http.ResponseWriter, r *http.Request){
	 db := kon()
	 
	 id := r.URL.Query().Get("id_karyawan")
	 detail, err := db.Query("select nama_karyawan,email,alamat,no_hp,gambar from karyawan where id_karyawan=?" , id)
      if err != nil {
		  panic(err.Error())
	  }
	  det := Karyawan{}
	  for detail.Next(){
	  var nama_karyawan, email, alamat, no_hp, gambar string
	 err =  detail.Scan(&nama_karyawan,&email,&alamat,&no_hp,&gambar)
	   if err != nil {
		   panic(err.Error())
	   }

	   det.NamaKaryawan = nama_karyawan
	   det.Email = email
	   det.Alamat = alamat
	   det.No_Hp = no_hp
	   det.Gambar = gambar

	  }
	 tmpl.ExecuteTemplate(w, "DetailKaryawan", det)
	 defer db.Close()
 }


 func EditKaryawan(w http.ResponseWriter, r *http.Request){
	 db := kon()
	
	 id := r.URL.Query().Get("id_karyawan")

	 edit, err:= db.Query("select id_karyawan,nama_karyawan, alamat, no_hp, gambar,email from karyawan where id_karyawan = ?", id)
      if err != nil {
		  panic(err.Error())
	  }
	  
	  edt := Karyawan{}

	  for edit.Next(){
		  var  nama_karyawan, email, no_hp, gambar, alamat string
          var id_karyawan int
		  err = edit.Scan(&id_karyawan,&nama_karyawan, &alamat, &no_hp,&gambar, &email)

		  if err != nil {
			  panic(err.Error())
		  }
          edt.Id_Kar = id_karyawan
		  edt.NamaKaryawan  = nama_karyawan
		  edt.Email = email
		  edt.No_Hp = no_hp
		  edt.Gambar = gambar
		  edt.Alamat = alamat

	  }

	  tmpl.ExecuteTemplate(w, "EditKaryawan", edt)
	  defer db.Close()
 }

func UpdateKaryawan(w http.ResponseWriter, r *http.Request){
	db := kon()
    if r.Method == "POST" {
        id_karyawan := r.FormValue("id_karyawan")
		nama_karyawan := r.FormValue("nama_karyawan")
		alamat := r.FormValue("alamat")
		no_hp := r.FormValue("no_hp")
		email := r.FormValue("email")
		files, _, _  := r.FormFile("gambar")
		
		if files!= nil {
			data := Karyawan{}
   var gambar string
			selectUp := "select gambar from karyawan where id_karyawan=?"
			err := db.QueryRow(selectUp, id_karyawan).Scan(&gambar)

			if err != nil {
				panic(err.Error())
			}
              data.Gambar = gambar
			

			if gambar!="" && gambar!="default.jpg" {
                _ = os.Remove("./images/" + gambar)
			}
			gbr, header, _ := r.FormFile("gambar")

			UpdGambar, err := db.Prepare("update karyawan set nama_karyawan=?,alamat=?,no_hp=?,email=?, gambar=? where id_karyawan=?")
             if  err != nil {
				 panic(err.Error())
			 }

			 UpdGambar.Exec(nama_karyawan,alamat,no_hp,email,header.Filename,id_karyawan)
		     out, _ := os.Create("./images/" + header.Filename)
				_, _ = io.Copy(out, gbr)
				gbr.Close()
				out.Close()
			} else {
			  
				upd, err := db.Prepare("update karyawan set nama_karyawan=?,alamat=?,no_hp=?,email=? where id_karyawan=?")
                if err != nil {
					panic(err.Error())
				}

				upd.Exec(nama_karyawan,alamat,no_hp,email,id_karyawan)
			}
			http.Redirect(w,r,"/",301)
	}

}

// end karyawan

// login
func ProsesLogin(w http.ResponseWriter, r *http.Request){
   db := kon()
   defer db.Close()
   if r.Method == "POST"{
	 email := r.FormValue("email")
	 password := r.FormValue("password")

	 sql := "SELECT * from karyawan where email=? and password=?"
      data := Karyawan{}
	 err := db.QueryRow(sql,email,password).Scan(&data.Email)
     if err != nil {
		http.Redirect(w,r,"user?err=Email dan Password anda Salah",301)
	 } else {
		 
		session := sessions.Start(w, r)
		session.Set("suserid", data.Email)
		http.Redirect(w,r,"Header", 301)
		
	 }

   }
 

}

//end login

func Home(w http.ResponseWriter, r *http.Request){
  db := kon()
	session := sessions.Start(w, r)
	var suserid = session.GetString("suserid")
	var data = make(map[string]string)
	data["suserid"] = suserid
	data["err"] = r.URL.Query().Get("err")
	if suserid != ""{
		http.Redirect(w,r,"Home", 301)

	} else {
		http.Redirect(w,r,"login?err=Harap login terlebih dahulu",301)
	}

}

func main() {
	log.Println("Server started on: http://localhost:8080")

	 http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	 http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	 http.HandleFunc("/", Index)
	 http.HandleFunc("/Home", Home)
	 http.HandleFunc("/input_karyawan", InputKaryawan)
	 http.HandleFunc("/save_karyawan", SaveKaryawan) 
	 http.HandleFunc("/lihat_karyawan", LihatKaryawan)
	 http.HandleFunc("/hapus_karyawan", HapusKaryawan)
	 http.HandleFunc("/detail_karyawan", DetailKaryawan)
	 http.HandleFunc("/edit_karyawan", EditKaryawan)
	 http.HandleFunc("/update_karyawan", UpdateKaryawan)
	 http.ListenAndServe(":8080", nil)
}