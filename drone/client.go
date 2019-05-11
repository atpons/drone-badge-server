package drone

import (
	"github.com/atpons/drone-badge-server/database"
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/go-memdb"
	"golang.org/x/oauth2"
)

type Drone struct {
	Db     *memdb.MemDB
	client drone.Client
}

func NewDrone(host, token string) (*Drone, error) {
	config := new(oauth2.Config)
	auther := config.Client(
		oauth2.NoContext,
		&oauth2.Token{
			AccessToken: token,
		},
	)
	c := drone.NewClient(host, auther)
	db, err := db.NewDb()
	if err != nil {
		return nil, err
	}
	client := &Drone{db,c}
	return client, err
}

func (c *Drone) SetStage(namespace, repo string, build int) error {
	repoId, err := c.getRepoId(namespace, repo)
	if err != nil {
		return err
	}
	stages, err := c.getStagesFromBuild(namespace, repo, build)
	if err != nil {
		return err
	}
	for _, stage := range stages {
		s := &db.Stage{
			Id: string(stage.ID),
			RepoId: repoId,
			Num: stage.Number,
			Status: statusExpression(stage.Status),
		}
		err := db.SetValue(c.Db, s)
		if err != nil{
			return err
		}
	}
	return nil
}

func (c *Drone) getStagesFromBuild (namespace, repo string, build int) ([]*drone.Stage, error) {
	b, err := c.client.Build(namespace, repo, build)
	stages := b.Stages
	return stages, err
}

func (c *Drone) getRepoId (namespace, repo string) (int, error) {
	r, err := c.client.Repo(namespace, repo)
	if err != nil {
		return 0, err
	}
	return int(r.ID), nil
}

func statusExpression (status string) bool {
	if status == "success" {
		return true
	}
	return false
}