package Log

type Category string
type SubCategory string
type Extrakey string

const (
	General         Category = "General"
	Internal        Category = "Internal"
	Postgres        Category = "Postgres"
	Redis           Category = "Redis"
	Validation      Category = "Validation"
	RequestResponse Category = "RequestResponse"
)

const (

	// General
	Startup         SubCategory = "Startup"
	ExternalService SubCategory = "ExternalService"

	// Postgres
	Select    SubCategory = "Select"
	Rolback   SubCategory = "Rollback"
	Update    SubCategory = "Update"
	Delete    SubCategory = "Delete"
	Insert    SubCategory = "Insert"
	Migration SubCategory = "Migration"
	Commit    SubCategory = "Commit"

	// Internal
	Api                 SubCategory = "Api"
	HashPassword        SubCategory = "HashPassword"
	DefaultRoleNotFound SubCategory = "DefaultRoleNotFound"

	// Validation
	MobileValidation   SubCategory = "MobileValidation"
	PasswordValidation SubCategory = "PasswordValidation"
)

const (
	AppName      Extrakey = "AppName"
	LoggerName   Extrakey = "LoggerName"
	ClientIp     Extrakey = "ClientIp"
	HostIp       Extrakey = "HostIp"
	Method       Extrakey = "Method"
	StatusCode   Extrakey = "StatusCode"
	BodySize     Extrakey = "BodySize"
	Path         Extrakey = "Path"
	Latency      Extrakey = "Latency"
	Body         Extrakey = "Body"
	ErrorMessage Extrakey = "ErrorMessage"
	RequestBody  Extrakey = "RequestBody"
	ResponseBody Extrakey = "ResponseBody"
)
