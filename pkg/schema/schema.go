package schema

type Location string

const (
	LocationUserGlobal  Location = "user_global"
	LocationProjectRoot Location = "project_root"
	LocationSubdir      Location = "subdir"
	LocationImport      Location = "import"
)

type Type string

const (
	TypeRules       Type = "rules"
	TypeIgnore      Type = "ignore"
	TypeMCP         Type = "mcp"
	TypeCommands    Type = "commands"
	TypeSubagents   Type = "subagents"
	TypeSkills      Type = "skills"
	TypeHooks       Type = "hooks"
	TypePermissions Type = "permissions"
	TypePolicies    Type = "policies"
)

type Envelope struct {
	SchemaVersion int    `json:"schema_version"`
	Command       string `json:"command"`
	CWD           string `json:"cwd"`
}

type Source struct {
	ID       string   `json:"id"`
	Path     string   `json:"path"`
	Location Location `json:"location"`
	Type     Type     `json:"type"`
	Bytes    int      `json:"bytes"`
}

type Chunk struct {
	Content            string `json:"content"`
	SourceID           string `json:"source_id"`
	AssembledLineStart int    `json:"assembled_line_start"`
	AssembledLineEnd   int    `json:"assembled_line_end"`
	SourceLineStart    int    `json:"source_line_start"`
	SourceLineEnd      int    `json:"source_line_end"`
}

type Conflict struct {
	Type      string   `json:"type"` // duplicate_heading | near_duplicate_paragraph | contradicting_imperative
	Heading   string   `json:"heading,omitempty"`
	SourceIDs []string `json:"source_ids"`
}

type Finding struct {
	Type       string `json:"type"` // secret | pii | internal_host | prohibited_content | policy_violation
	Severity   string `json:"severity"`
	Message    string `json:"message"`
	SourceID   string `json:"source_id"`
	SourceFile string `json:"source_file"`
	LineStart  int    `json:"line_start"`
	LineEnd    int    `json:"line_end"`
}

type ScanResult struct {
	Envelope
	Sources   []Source   `json:"sources"`
	Assembled []Chunk    `json:"assembled"`
	Conflicts []Conflict `json:"conflicts"`
	Findings  []Finding  `json:"findings"`
}

type SectionIncluded struct {
	Heading string `json:"heading"`
	Source  string `json:"source"`
}

type SectionExcluded struct {
	Heading string `json:"heading"`
	Reason  string `json:"reason"`
}

type GenResult struct {
	Envelope
	OutputPath       string            `json:"output_path"`
	Role             string            `json:"role"`
	Tags             []string          `json:"tags"`
	SectionsIncluded []SectionIncluded `json:"sections_included"`
	SectionsExcluded []SectionExcluded `json:"sections_excluded"`
	ManifestPath     string            `json:"manifest_path"`
	Content          string            `json:"content"`
}

type Change struct {
	Type    string   `json:"type"` // added | removed | modified
	Section string   `json:"section"`
	Role    string   `json:"role"`
	Tags    []string `json:"tags"`
	Before  string   `json:"before"`
	After   string   `json:"after"`
}

type DiffResult struct {
	Envelope
	InSync  bool     `json:"in_sync"`
	Changes []Change `json:"changes"`
}

type ConvertResult struct {
	Envelope
	From       string `json:"from"`
	To         string `json:"to"`
	Type       Type   `json:"type"`
	OutputPath string `json:"output_path"`
	Content    string `json:"content"`
}

type UpdateResult struct {
	Envelope
	TargetsUpdated []string `json:"targets_updated"`
	Changes        []Change `json:"changes"`
}

type CleanResult struct {
	Envelope
	RemovedFiles []string `json:"removed_files"`
}

type Section struct {
	Heading    string            `json:"heading"`
	Content    string            `json:"content"`
	Dimensions map[string]string `json:"dimensions,omitempty"`
	SourceFile string            `json:"source_file,omitempty"`
	LineStart  int               `json:"line_start,omitempty"`
	LineEnd    int               `json:"line_end,omitempty"`
}
