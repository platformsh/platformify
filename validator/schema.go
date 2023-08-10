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

//go:embed schema/upsun.json
var upsunSchemaContents string

var (
	routesSchema      *gojsonschema.Schema
	servicesSchema    *gojsonschema.Schema
	applicationSchema *gojsonschema.Schema

	upsunSchema *gojsonschema.Schema
)

func init() {
	routesSchema, _ = gojsonschema.NewSchema(gojsonschema.NewStringLoader(routesSchemaContents))
	servicesSchema, _ = gojsonschema.NewSchema(gojsonschema.NewStringLoader(servicesSchemaContents))
	applicationSchema, _ = gojsonschema.NewSchema(gojsonschema.NewStringLoader(applicationSchemaContents))

	upsunSchema, _ = gojsonschema.NewSchema(gojsonschema.NewStringLoader(upsunSchemaContents))
}
