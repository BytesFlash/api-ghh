package config

const (
	CHILE    = "CHILE"
	COLOMBIA = "COLOMBIA"
	ADD      = "A"
	LIST     = "L"
	DELETE   = "E"
	READ     = "R"
	CHANGE   = "C"
	PASSWORD = "P"

	DEFAULT_UBICATION = "Sin definir"

	ADMIN = "Admin"
	JP    = "Jp"

	DEFAULT_OWNER = "NO ASIGNADO"
	ALL_OWNERS    = "TODOS"

	STATUS_ACTIVE   = "Activo"
	STATUS_INACTIVE = "Inactivo"
	STATUS_LOCKED   = "Bloqueado"
)

var (
	OWNERS            = []string{"ACEPTA", "I-MED", "AUTENTIA", ALL_OWNERS, DEFAULT_OWNER}
	COUNTRIES         = []string{"CHILE", "ECUADOR", "COLOMBIA", "RDOMINICANA", "MEXICO"}
	STUATUS_ACTIVOS   = []string{STATUS_ACTIVE}
	STUATUS_INACTIVOS = []string{STATUS_INACTIVE, STATUS_LOCKED}
	ALL_STATUS        = []string{STATUS_INACTIVE, STATUS_LOCKED, STATUS_ACTIVE}
	ALL_GROUPS        = []string{ADMIN, JP}
)
