package models

type Pagination struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}

type DocType string

// Doc Type
const (
	Nectec DocType = "nectec"
	Farmer DocType = "farmer"
	Gap DocType = "gap"
	Gmp DocType = "gmp"
	Exporter DocType = "exporter"
	Packer DocType = "packer"
	Packing DocType = "packing"
	Packaging DocType = "packaging"
)