package enricher

import (
	"bytes"
	"log"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type FieldMapperEnricherFieldConfig struct {
	SourceField string      `yaml:"source_field"`
	TargetField string      `yaml:"target_field"`
	Mapping     map[any]any `yaml:"mapping"`
	Template    string      `yaml:"template"`
}

type FieldMapperEnricherConfig struct {
	Fields []*FieldMapperEnricherFieldConfig `yaml:"fields"`
}

type FieldMapperEnricher struct {
	Config    *FieldMapperEnricherConfig
	templates map[*FieldMapperEnricherFieldConfig]*template.Template
}

func NewFieldMapperEnricher(config *FieldMapperEnricherConfig) FieldMapperEnricher {
	if config == nil {
		config = &FieldMapperEnricherConfig{}
	}
	templates := map[*FieldMapperEnricherFieldConfig]*template.Template{}
	for _, field := range config.Fields {
		templates[field] = nil
		if field.Template != "" {
			tmpl, err := template.New("field_mapper_enricher").Funcs(sprig.FuncMap()).Parse(field.Template)
			if err != nil {
				panic(err)
			}
			templates[field] = tmpl
		}
	}
	return FieldMapperEnricher{
		Config:    config,
		templates: templates,
	}
}

func (e *FieldMapperEnricher) Process(msg map[string]interface{}) map[string]interface{} {
	for i := range e.Config.Fields {
		msg = e.applyMapper(msg, e.Config.Fields[i])
	}
	return msg
}

func (e *FieldMapperEnricher) applyMapper(msg map[string]interface{}, fieldConfig *FieldMapperEnricherFieldConfig) map[string]interface{} {
	if e.templates[fieldConfig] != nil {
		return e.applyTemplateMapper(msg, fieldConfig)
	} else {
		return e.applySimpleMapper(msg, fieldConfig)
	}
}

func (e *FieldMapperEnricher) applyTemplateMapper(msg map[string]interface{}, fieldConfig *FieldMapperEnricherFieldConfig) map[string]interface{} {
	tmpl := e.templates[fieldConfig]
	data := struct {
		Config      *FieldMapperEnricherFieldConfig
		SourceField any
		Mapping     map[any]any
		Msg         map[string]interface{}
	}{
		Config:      fieldConfig,
		SourceField: msg[fieldConfig.SourceField],
		Mapping:     fieldConfig.Mapping,
		Msg:         msg,
	}
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		log.Printf("error executing field mapper template for source_field=%s, target_field=%s: %v", fieldConfig.SourceField, fieldConfig.TargetField, err)
		return msg
	}
	value := strings.TrimSpace(buf.String())
	if value != "" {
		msg[fieldConfig.TargetField] = value
	}
	return msg
}

func (e *FieldMapperEnricher) applySimpleMapper(msg map[string]interface{}, fieldConfig *FieldMapperEnricherFieldConfig) map[string]interface{} {
	sourceValue, ok := msg[fieldConfig.SourceField]
	if !ok {
		return msg
	}
	value, ok := fieldConfig.Mapping[sourceValue]
	if !ok {
		return msg
	}
	msg[fieldConfig.TargetField] = value
	return msg
}
