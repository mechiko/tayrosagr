package main

import (
	"os"
	"tayrosagr/utility"

	"github.com/mechiko/dbscan"
)

const THEMES_DIR = "./themes/"

func main() {
	listDbs := make(dbscan.ListDbInfoForScan)
	listDbs[dbscan.Config] = &dbscan.DbInfo{}
	listDbs[dbscan.TrueZnak] = &dbscan.DbInfo{}
	dbs, err := dbscan.New(listDbs, ".")
	if err != nil {
		utility.MessageBox("ошибка", err.Error())
		os.Exit(1)
	}

	_, _ = startDialog(dbs)
}
