package badge

import "net/http"

type Stage map[int][]string

type Repo map[int]Stage

func ReturnImage(w http.ResponseWriter, r *http.Request, s Repo, repoId int, stageId int, status bool) {
	w.Header().Set("Content-type", "image/svg+xml")
	var fileName string
	if status {
		fileName = s[repoId][stageId][0]
	} else {
		fileName = s[repoId][stageId][1]
	}
	http.ServeFile(w, r, fileName)
}
