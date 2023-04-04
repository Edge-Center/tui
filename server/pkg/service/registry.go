package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/anaxaim/tui/server/pkg/model"
	"github.com/anaxaim/tui/server/pkg/repository"
	"github.com/anaxaim/tui/server/pkg/utils"
)

type registryService struct {
	registryRepository repository.RegistryRepository
	moduleRepository   repository.ModuleRepository
}

func NewRegistryService(registryRepository repository.RegistryRepository, moduleRepository repository.ModuleRepository) RegistryService {
	return &registryService{
		registryRepository: registryRepository,
		moduleRepository:   moduleRepository,
	}
}

func (r *registryService) ImportModuleContentByID(id string) (*model.RegistryContent, error) {
	module, err := r.moduleRepository.GetModuleByID(id)
	if err != nil {
		return nil, err
	}

	registry := new(model.RegistryContent)
	registry.RegistryType = module.RegistryDetails.RegistryType
	registry.ModuleID = module.ID

	repo, err := utils.CloneGitRepo(module)
	if err != nil {
		return nil, err
	}

	tree, err := utils.GetCommitTree(repo)
	if err != nil {
		return nil, err
	}

	content, err := utils.GetModuleContent(tree)
	if err != nil {
		return nil, err
	}

	registry.Content = content

	tmpModuleDir := filepath.Join(os.TempDir(), fmt.Sprintf("module_%s", module.ID))
	defer os.RemoveAll(tmpModuleDir)

	err = utils.ExtractRegistryFilesToDir(content, tmpModuleDir)
	if err != nil {
		return nil, err
	}

	tfModule, err := utils.LoadTFModule(tmpModuleDir)
	if err != nil {
		return nil, err
	}

	if len(module.Variables) > 0 {
		tfModule.Variables = utils.MergeTfVariables(module.Variables, tfModule.Variables)
	}

	registry.ParsedContent = tfModule

	registryInfo, err := r.registryRepository.Save(registry)
	if err != nil {
		module.Status = model.ERROR
	} else {
		module.Status = model.RUNNING
		module.RegistryDetails.RegistryID = registryInfo.ID
	}

	_, err = r.moduleRepository.Update(module)
	if err != nil {
		return nil, err
	}

	return registryInfo, nil
}

func (r *registryService) GetModuleContentByID(id string) (*model.RegistryContent, error) {
	return r.registryRepository.Get(id)
}
