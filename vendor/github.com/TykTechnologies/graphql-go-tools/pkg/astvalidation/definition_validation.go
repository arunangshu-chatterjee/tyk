package astvalidation

import (
	"github.com/TykTechnologies/graphql-go-tools/pkg/ast"
	"github.com/TykTechnologies/graphql-go-tools/pkg/astvisitor"
	"github.com/TykTechnologies/graphql-go-tools/pkg/operationreport"
)

func DefaultDefinitionValidator() *DefinitionValidator {
	return NewDefinitionValidator(
		PopulatedTypeBodies(),
		UniqueOperationTypes(),
		UniqueTypeNames(),
		UniqueFieldDefinitionNames(),
		UniqueEnumValueNames(),
		UniqueUnionMemberTypes(),
		KnownTypeNames(),
		RequireDefinedTypesForExtensions(),
		ImplementTransitiveInterfaces(),
		ImplementingTypesAreSupersets(),
	)
}

func NewDefinitionValidator(rules ...Rule) *DefinitionValidator {
	validator := &DefinitionValidator{
		walker: astvisitor.NewWalker(48),
	}

	for _, rule := range rules {
		validator.RegisterRule(rule)
	}

	return validator
}

type DefinitionValidator struct {
	walker astvisitor.Walker
}

func (d *DefinitionValidator) RegisterRule(rule Rule) {
	rule(&d.walker)
}

func (d *DefinitionValidator) Validate(definition *ast.Document, report *operationreport.Report) ValidationState {
	if report == nil {
		report = &operationreport.Report{}
	}

	d.walker.Walk(definition, nil, report)

	if report.HasErrors() {
		return Invalid
	}
	return Valid
}