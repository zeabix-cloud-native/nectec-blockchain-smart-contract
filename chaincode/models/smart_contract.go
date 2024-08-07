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
	Regulator DocType = "regulator"
	Gap DocType = "gap"
	Gmp DocType = "gmp"
	Exporter DocType = "exporter"
	Packer DocType = "packer"
	Packing DocType = "packing"
	Packaging DocType = "packaging"
	Hscode DocType = "hscode"
	FormE DocType = "formE"
	PlantType DocType = "plantType"
)