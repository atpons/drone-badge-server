package db

import (
	"fmt"

	"github.com/hashicorp/go-memdb"
)

type Stage struct {
	Id     string
	RepoId int
	Num    int
	Status bool
}

func NewDb() (*memdb.MemDB, error) {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"stage": {
				Name: "stage",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
					},
					"repoid": {
						Name:    "repoid",
						Indexer: &memdb.IntFieldIndex{Field: "RepoId"},
					},
				},
			},
		},
	}
	db, err := memdb.NewMemDB(schema)
	return db, err
}

func SetValue(db *memdb.MemDB, stage *Stage) error {
	txn := db.Txn(true)
	defer txn.Abort()
	err := txn.Insert("stage", stage)
	txn.Commit()
	return err
}

func GetValue(db *memdb.MemDB, repoId int, stageNum int) (bool, error) {
	txn := db.Txn(false)
	defer txn.Abort()
	raw, err := txn.Get("stage", "repoid", repoId)
	if err != nil {
		return false, nil
	}
	for item := raw.Next(); item != nil; item = raw.Next() {
		if item.(*Stage).Num == stageNum {
			return item.(*Stage).Status, nil
		}
	}
	return false, fmt.Errorf("get error: record not found")
}
