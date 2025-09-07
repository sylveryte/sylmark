package lsp

import (
	"path/filepath"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type DocumentURI string

func (d DocumentURI) GetFileName() string {
	return filepath.Base(string(d))
}

func (d DocumentURI) LocationOfFile() Location {
	return Location{
		URI: d,
	}
}

func (d DocumentURI) LocationOf(node *tree_sitter.Node) Location {
	return Location{
		URI:   d,
		Range: GetRange(node),
	}
}

type InitializeParams struct {
	ProcessID             int                `json:"processId,omitempty"`
	RootURI               DocumentURI        `json:"rootUri,omitempty"`
	InitializationOptions *InitializeOptions `json:"initializationOptions,omitempty"`
	Capabilities          ClientCapabilities `json:"capabilities,omitempty"`
	Trace                 string             `json:"trace,omitempty"`
}

type SemanticTokensLegend struct {
	TokenTypes     []SemanticTokenType     `json:"tokenTypes"`
	TokenModifiers []SemanticTokenModifier `json:"tokenModifiers"`
}
type ExecuteCommandOptions struct {
	Commands []string `json:"commands"`
}
type SemanticTokensOptions struct {
	Legend SemanticTokensLegend `json:"legend"`
	Range  bool                 `json:"range"`
	Full   bool                 `json:"full"`
}

type SemantiTokens struct {
	ResultId string `json:"resultId"` // optional
	Data     []uint `json:"data"`
}

type InitializeOptions struct {
	DocumentFormatting bool `json:"documentFormatting"`
	RangeFormatting    bool `json:"documentRangeFormatting"`
	Hover              bool `json:"hover"`
	DocumentSymbol     bool `json:"documentSymbol"`
	CodeAction         bool `json:"codeAction"`
	Completion         bool `json:"completion"`
}

type ClientCapabilities struct{}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities,omitempty"`
}

type TextDocumentSyncKind int

const (
	TDSKNone TextDocumentSyncKind = iota
	TDSKFull
	TDSKIncremental
)

type CompletionProvider struct {
	ResolveProvider   bool     `json:"resolveProvider,omitempty"`
	TriggerCharacters []string `json:"triggerCharacters"`
}

type WorkspaceFoldersServerCapabilities struct {
	Supported           bool `json:"supported"`
	ChangeNotifications bool `json:"changeNotifications"`
}

type ServerCapabilitiesWorkspace struct {
	WorkspaceFolders WorkspaceFoldersServerCapabilities `json:"workspaceFolders"`
}
type DiagnosticOptions struct {
	InterFileDependencies bool `json:"interFileDependencies"`
	WorkspaceDiagnostics  bool `json:"workspaceDiagnostics"`
}

type ServerCapabilities struct {
	TextDocumentSync           TextDocumentSyncKind         `json:"textDocumentSync,omitempty"`
	DocumentSymbolProvider     bool                         `json:"documentSymbolProvider,omitempty"`
	CompletionProvider         *CompletionProvider          `json:"completionProvider,omitempty"`
	DefinitionProvider         bool                         `json:"definitionProvider,omitempty"`
	ReferencesProvider         bool                         `json:"referencesProvider,omitempty"`
	SemanticTokensProvider     SemanticTokensOptions        `json:"semanticTokensProvider"`
	DocumentFormattingProvider bool                         `json:"documentFormattingProvider,omitempty"`
	DiagnosticProvider         DiagnosticOptions            `json:"diagnosticProvider"`
	RangeFormattingProvider    bool                         `json:"documentRangeFormattingProvider,omitempty"`
	HoverProvider              bool                         `json:"hoverProvider,omitempty"`
	CodeActionProvider         bool                         `json:"codeActionProvider,omitempty"`
	ExecuteCommandProvider     ExecuteCommandOptions        `json:"executeCommandProvider"`
	Workspace                  *ServerCapabilitiesWorkspace `json:"workspace,omitempty"`
}

type TextDocumentIdentifier struct {
	URI DocumentURI `json:"uri"`
}

type TextDocumentItem struct {
	URI        DocumentURI `json:"uri"`
	LanguageID string      `json:"languageId"`
	Version    int         `json:"version"`
	Text       string      `json:"text"`
}

type ExecuteCommandParams struct {
	Command   string   `json:"command"`
	Arguments []string `json:"arguments"`
}
type ShowDocumentParams struct {
	URI       DocumentURI `json:"uri"`
	External  bool        `json:"external"`
	Selection Range       `json:"selection"`
	TakeFocus bool        `json:"takeFocus"`
}

type ShowDocumentResult struct {
	Success bool `json:"success"`
}

type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type TextDocumentContentChangeEvent struct {
	Range       Range  `json:"range"`
	RangeLength int    `json:"rangeLength"`
	Text        string `json:"text"`
}

type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int `json:"version"`
}

type SemanticTokensParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type SemanticTokensRangeParams struct {
	TextDocumentIdentifier
	Range Range `json:"range"`
}

type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type DidSaveTextDocumentParams struct {
	Text         *string                `json:"text"`
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type TextDocumentPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

type DocumentDiagnosticParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}
type CodeActionContext struct {
	Diagnostics []Diagnostic `json:"diagnostics"`
}
type CodeAction struct {
	Title       string       `json:"title"`
	Diagnostics []Diagnostic `json:"diagnostics"` // that is resolved by
	Command     Command      `json:"command"`
}
type CodeActionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
	Context      CodeActionContext      `json:"context"`
}

type CompletionParams struct {
	TextDocumentPositionParams
	CompletionContext CompletionContext `json:"context"`
}

type CompletionContext struct {
	TriggerKind      int     `json:"triggerKind"`
	TriggerCharacter *string `json:"triggerCharacter"`
}

type HoverParams struct {
	TextDocumentPositionParams
}

type DefinitionParams struct {
	TextDocumentPositionParams
}

type ReferencesParams struct {
	TextDocumentPositionParams
}

type Location struct {
	URI   DocumentURI `json:"uri"`
	Range Range       `json:"range"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// Hover is
type Hover struct {
	Contents any    `json:"contents"`
	Range    *Range `json:"range"`
}

type DiagnosticResult struct {
	Kind  DiagnosticReportKind `json:"kind"`
	Items []Diagnostic         `json:"items"`
}

type Diagnostic struct {
	Range    *Range             `json:"range"`
	Severity DiagnosticSeverity `json:"severity"`
	Tags     []DiagnosticTag    `json:"tags"`
	Message  string             `json:"message"`
}

type DiagnosticReportKind string

type DiagnosticSeverity int

type DiagnosticTag int

type CompletionItemKind int

type CompletionItemTag int

type InsertTextFormat int

type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}

type Command struct {
	Title     string `json:"title"`
	Command   string `json:"command"`
	Arguments []any  `json:"arguments,omitempty"`
	OS        string `json:"-"`
}

type CompletionItem struct {
	Label               string              `json:"label"`
	Kind                CompletionItemKind  `json:"kind,omitempty"`
	Tags                []CompletionItemTag `json:"tags,omitempty"`
	Detail              string              `json:"detail,omitempty"`
	Documentation       string              `json:"documentation,omitempty"` // string | MarkupContent
	Deprecated          bool                `json:"deprecated,omitempty"`
	Preselect           bool                `json:"preselect,omitempty"`
	SortText            string              `json:"sortText,omitempty"`
	FilterText          string              `json:"filterText,omitempty"`
	InsertText          string              `json:"insertText,omitempty"`
	InsertTextFormat    InsertTextFormat    `json:"insertTextFormat,omitempty"`
	TextEdit            *TextEdit           `json:"textEdit,omitempty"`
	AdditionalTextEdits []TextEdit          `json:"additionalTextEdits,omitempty"`
	CommitCharacters    []string            `json:"commitCharacters,omitempty"`
	Command             *Command            `json:"command,omitempty"`
	Data                any                 `json:"data,omitempty"`
}

const (
	DiagnosticSeverityError       DiagnosticSeverity = 1
	DiagnosticSeverityWarning     DiagnosticSeverity = 2
	DiagnosticSeverityInformation DiagnosticSeverity = 3
	DiagnosticSeverityHint        DiagnosticSeverity = 4
)

const (
	DiagnosticTagUnnecessary DiagnosticTag = 1
	DiagnosticTagDeprecated  DiagnosticTag = 2
)

const (
	DiagnosticReportFull      DiagnosticReportKind = "full"
	DiagnosticReportUnchanged DiagnosticReportKind = "unchanged"
)

const (
	TextCompletion          CompletionItemKind = 1
	MethodCompletion        CompletionItemKind = 2
	FunctionCompletion      CompletionItemKind = 3
	ConstructorCompletion   CompletionItemKind = 4
	FieldCompletion         CompletionItemKind = 5
	VariableCompletion      CompletionItemKind = 6
	ClassCompletion         CompletionItemKind = 7
	InterfaceCompletion     CompletionItemKind = 8
	ModuleCompletion        CompletionItemKind = 9
	PropertyCompletion      CompletionItemKind = 10
	UnitCompletion          CompletionItemKind = 11
	ValueCompletion         CompletionItemKind = 12
	EnumCompletion          CompletionItemKind = 13
	KeywordCompletion       CompletionItemKind = 14
	SnippetCompletion       CompletionItemKind = 15
	ColorCompletion         CompletionItemKind = 16
	FileCompletion          CompletionItemKind = 17
	ReferenceCompletion     CompletionItemKind = 18
	FolderCompletion        CompletionItemKind = 19
	EnumMemberCompletion    CompletionItemKind = 20
	ConstantCompletion      CompletionItemKind = 21
	StructCompletion        CompletionItemKind = 22
	EventCompletion         CompletionItemKind = 23
	OperatorCompletion      CompletionItemKind = 24
	TypeParameterCompletion CompletionItemKind = 25
)

type SemanticTokenType string

const (
	NamespaceSematicTokenType     SemanticTokenType = "namespace"
	TypeSematicTokenType          SemanticTokenType = "type"
	ClassSematicTokenType         SemanticTokenType = "class"
	EnumSematicTokenType          SemanticTokenType = "enum"
	InterfaceSematicTokenType     SemanticTokenType = "interface"
	StructSematicTokenType        SemanticTokenType = "struct"
	TypeParameterSematicTokenType SemanticTokenType = "typeParameter"
	ParameterSematicTokenType     SemanticTokenType = "parameter"
	VariableSematicTokenType      SemanticTokenType = "variable"
	PropertySematicTokenType      SemanticTokenType = "property"
	EnumMemberSematicTokenType    SemanticTokenType = "enumMember"
	EventSematicTokenType         SemanticTokenType = "event"
	FunctionSematicTokenType      SemanticTokenType = "function"
	MethodSematicTokenType        SemanticTokenType = "method"
	MacroSematicTokenType         SemanticTokenType = "macro"
	KeywordSematicTokenType       SemanticTokenType = "keyword"
	ModifierSematicTokenType      SemanticTokenType = "modifier"
	CommentSematicTokenType       SemanticTokenType = "comment"
	StringSematicTokenType        SemanticTokenType = "string"
	NumberSematicTokenType        SemanticTokenType = "number"
	RegexpSematicTokenType        SemanticTokenType = "regexp"
	OperatorSematicTokenType      SemanticTokenType = "operator"
	DecoratorSematicTokenType     SemanticTokenType = "decorator"
)

type SemanticTokenModifier string

const (
	DeclarationSemanticTokenModifier    SemanticTokenModifier = "declaration"
	DefinitionSemanticTokenModifier     SemanticTokenModifier = "definition"
	ReadonlySemanticTokenModifier       SemanticTokenModifier = "readonly"
	StaticSemanticTokenModifier         SemanticTokenModifier = "static"
	DeprecatedSemanticTokenModifier     SemanticTokenModifier = "deprecated"
	AbstractSemanticTokenModifier       SemanticTokenModifier = "abstract"
	AsyncSemanticTokenModifier          SemanticTokenModifier = "async"
	ModificationSemanticTokenModifier   SemanticTokenModifier = "modification"
	DocumentationSemanticTokenModifier  SemanticTokenModifier = "documentation"
	DefaultLibrarySemanticTokenModifier SemanticTokenModifier = "defaultLibrary"
)
