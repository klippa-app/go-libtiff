package imports

import (
	"context"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/emscripten"
)

// Instantiate instantiates the "env" module used by Emscripten into the
// runtime default namespace.
//
// # Notes
//
//   - Closing the wazero.Runtime has the same effect as closing the result.
//   - To add more functions to the "env" module, use FunctionExporter.
//   - To instantiate into another wazero.Namespace, use FunctionExporter.
func Instantiate(ctx context.Context, r wazero.Runtime, mod wazero.CompiledModule) (api.Closer, error) {
	builder := r.NewHostModuleBuilder("env")
	exporter, err := emscripten.NewFunctionExporterForModule(mod)
	if err != nil {
		return nil, err
	}
	exporter.ExportFunctions(builder)
	NewFunctionExporter().ExportFunctions(builder)
	return builder.Instantiate(ctx)
}

// FunctionExporter configures the functions in the "env" module used by
// Emscripten.
type FunctionExporter interface {
	// ExportFunctions builds functions to export with a wazero.HostModuleBuilder
	// named "env".
	ExportFunctions(builder wazero.HostModuleBuilder)
}

// NewFunctionExporter returns a FunctionExporter object with trace disabled.
func NewFunctionExporter() FunctionExporter {
	return &functionExporter{}
}

type functionExporter struct{}

// ExportFunctions implements FunctionExporter.ExportFunctions
func (e *functionExporter) ExportFunctions(b wazero.HostModuleBuilder) {
	b.NewFunctionBuilder().WithGoModuleFunction(TIFFReadProcGoCB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).Export("TIFFReadProcGoCB")
	b.NewFunctionBuilder().WithGoModuleFunction(TIFFWriteProcGoCB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).Export("TIFFWriteProcGoCB")
	b.NewFunctionBuilder().WithGoModuleFunction(TIFFSeekProcGoCB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI64, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI64}).Export("TIFFSeekProcGoCB")
	b.NewFunctionBuilder().WithGoModuleFunction(TIFFCloseProcGoCB{}, []api.ValueType{api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).Export("TIFFCloseProcGoCB")
	b.NewFunctionBuilder().WithGoModuleFunction(TIFFSizeProcGoCB{}, []api.ValueType{api.ValueTypeI32}, []api.ValueType{api.ValueTypeI64}).Export("TIFFSizeProcGoCB")
	b.NewFunctionBuilder().WithGoModuleFunction(TIFFMapFileProcGoCB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).Export("TIFFMapFileProcGoCB")
	b.NewFunctionBuilder().WithGoModuleFunction(TIFFUnmapFileProcGoCB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI64}, []api.ValueType{}).Export("TIFFUnmapFileProcGoCB")
	b.NewFunctionBuilder().WithGoModuleFunction(TIFFOpenOptionsSetErrorHandlerExtRGoCB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).Export("TIFFOpenOptionsSetErrorHandlerExtRGoCB")
	b.NewFunctionBuilder().WithGoModuleFunction(TIFFOpenOptionsSetWarningHandlerExtRGoCB{}, []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).Export("TIFFOpenOptionsSetWarningHandlerExtRGoCB")
}
