package main

import "github.com/atpons/drone-badge-server/badge"

const (
	BuildStageSuccessImage  = "bs.png"
	BuildStageFailureImage  = "bf.png"
	DeployStageSuccessImage = "ds.png"
	DeployStageFailureImage = "df.png"
)

var BadgeRoute = badge.Repo{
	264: {
		1: {BuildStageSuccessImage, BuildStageFailureImage},
		2: {DeployStageSuccessImage, DeployStageFailureImage},
	},
}
