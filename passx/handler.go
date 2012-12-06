package passx

import (
	"html/template"
	"math"
	"net/http"
	"strconv"
	texttemplate "text/template"
)

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/schedules", schedules)
}

func root(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("root.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func schedules(w http.ResponseWriter, r *http.Request) {
	minUnits, err := strconv.Atoi(r.FormValue("minUnits"))
	if err != nil {
		minUnits = 0
	}

	maxUnits, err := strconv.Atoi(r.FormValue("maxUnits"))
	if err != nil {
		maxUnits = math.MaxInt32
	}

	schedules, err := generateSchedulesJSON([]byte(r.FormValue("classes")),
		minUnits, maxUnits)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := texttemplate.ParseFiles("schedules.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, schedules)
}
