package model

import (
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegistryType string

const (
	GITHUB RegistryType = "github"
	GITLAB RegistryType = "gitlab"
)

type RegistryContent struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ModuleID      primitive.ObjectID `json:"moduleId" bson:"moduleId"`
	RegistryType  RegistryType       `json:"registryType" bson:"registryType"`
	Content       map[string]string  `json:"content" bson:"content"`
	ParsedContent *tfconfig.Module   `json:"parsedContent" bson:"parsedContent"`
}

type ExecuteCommand struct {
	Command string `json:"command" bson:"-"`
}
