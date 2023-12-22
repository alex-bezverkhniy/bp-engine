package validators

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/alex-bezverkhniy/bp-engine/internal/config"
	"github.com/alex-bezverkhniy/bp-engine/internal/model"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

type (
	Validator interface {
		Validate(process model.ProcessDTO, newStatus model.ProcessStatusDTO) error
		CompileJsonSchema() error
	}

	BasicValidator struct {
		conf        config.ProcessConfigList
		jsonSchemas map[string]*jsonschema.Schema
	}
)

var ErrUnknownStatus = errors.New("unknown status")
var ErrNotAllowedStatus = errors.New("not allowed status")
var ErrPayloadValidation = errors.New("payload validation error: ")

func NewBasicValidator(conf []config.ProcessConfig) Validator {
	return &BasicValidator{
		conf: conf,
	}
}

func (bv *BasicValidator) Validate(process model.ProcessDTO, newStatus model.ProcessStatusDTO) error {
	// Check if status defined
	_, err := bv.conf.GetStatusConfig(process.Code, newStatus.Name)
	if err != nil {
		return ErrUnknownStatus
	}

	// Check current status config
	currentStatusCfg, err := bv.conf.GetStatusConfig(process.Code, process.CurrentStatus.Name)
	if err != nil {
		return ErrUnknownStatus
	}
	found := false
	for _, s := range currentStatusCfg.Next {
		if s == newStatus.Name {
			found = true
		}
	}
	if !found {
		return ErrNotAllowedStatus
	}

	schemaKey := bv.schemaKey(process.Code, newStatus.Name)
	schema := bv.jsonSchemas[schemaKey]
	if schema != nil {
		var m interface{}
		m, err = newStatus.Payload.ToStringKeys(newStatus.Payload["data"])
		if err != nil {
			return err
		}
		err = schema.Validate(m)
		if err != nil {
			return errors.Join(ErrPayloadValidation, bv.formatErrMsg(err))
		}
	}

	return nil
}

// CompileJsonSchema - Compiles JSON Schemas and adds it into map
func (bv *BasicValidator) CompileJsonSchema() error {
	compiler := jsonschema.NewCompiler()

	jsonSchemas := map[string]*jsonschema.Schema{}
	for _, pc := range bv.conf {
		for _, s := range pc.Statuses {
			if len(s.Schema) > 0 {
				schemaKey := bv.schemaKey(pc.Name, s.Name)
				err := compiler.AddResource(schemaKey, strings.NewReader(s.Schema))
				if err != nil {
					return err
				}
				schema, err := compiler.Compile(schemaKey)
				if err != nil {
					return err
				}
				jsonSchemas[schemaKey] = schema
			}
		}
	}
	bv.jsonSchemas = jsonSchemas
	return nil
}

func (bv *BasicValidator) schemaKey(processName, statusName string) string {
	return fmt.Sprintf("%s-%s", processName, statusName)
}

func (bv *BasicValidator) formatErrMsg(srcErr error) error {
	var re = regexp.MustCompile(`^jsonschema:(.*)(file:\/\/(.*):)(.*)$`)

	msg := re.ReplaceAllString(srcErr.Error(), `$1 JSON Schema:$4`)
	return errors.New(strings.Trim(msg, " "))
}
