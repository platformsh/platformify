package validator

import (
	_ "embed"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed schema/platformsh.routes.json
var routesSchemaContents string

//go:embed schema/platformsh.services.json
var servicesSchemaContents string

//go:embed schema/platformsh.application.json
var applicationSchemaContents string

var (
	routesSchema      *gojsonschema.Schema
	servicesSchema    *gojsonschema.Schema
	applicationSchema *gojsonschema.Schema
)

func init() {
	routesSchema, _ = gojsonschema.NewSchema(gojsonschema.NewStringLoader(routesSchemaContents))
	servicesSchema, _ = gojsonschema.NewSchema(gojsonschema.NewStringLoader(servicesSchemaContents))
	applicationSchema, _ = gojsonschema.NewSchema(gojsonschema.NewStringLoader(applicationSchemaContents))
}
