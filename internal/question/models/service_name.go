package models

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

type ServiceName struct {
	Name        string
	Type        string
	Description string
	Disk        bool
	Docs        struct {
		Relationship string
		URL          string
	}
	Endpoint    string
	MinDiskSize *int
	Versions    struct {
		Supported []string
	}
	Runtime bool
}

func (s *ServiceName) String() string {
	return s.Type
}

func (s *ServiceName) Title() string {
	return s.Name
}

func (s *ServiceName) IsPersistent() bool {
	return s.Disk
}

func (s *ServiceName) DefaultVersion() string {
	if len(s.Versions.Supported) > 0 {
		return s.Versions.Supported[0]
	}

	return ""
}

type ServiceNameList []*ServiceName

func (s *ServiceNameList) WriteAnswer(_ string, value interface{}) error {
	switch answer := value.(type) {
	case []survey.OptionAnswer: // MultiSelect
		for _, item := range answer {
			service, err := ServiceNames.ServiceByTitle(item.Value)
			if err != nil {
				return err
			}
			*s = append(*s, service)
		}
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}

func (s *ServiceNameList) AllTitles() []string {
	titles := make([]string, 0, len(*s))
	for _, service := range *s {
		titles = append(titles, service.Title())
	}
	return titles
}

func (s ServiceNameList) ServiceByTitle(title string) (*ServiceName, error) {
	for _, service := range s {
		if service.Title() == title {
			return service, nil
		}
	}
	return nil, fmt.Errorf("service name by title is not found")
}
