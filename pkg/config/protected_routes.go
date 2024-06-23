package config

var PROTECTED_GET = []string{
	"/users",
	"/institutions/sync/:country",
}
var PROTECTED_POST = []string{
	"/users/:dni/role",
	"/users",
	"/institutions",
	"/autentia/institutions",
	"/autentia/persons",
}
var PROTECTED_PUT = []string{
	"/users/:dni",
	"/autentia/institutions/:name",
	"/autentia/passwords/:dni",
	"/autentia/persons/:dni",
}

var PROTECTED_DELETE = []string{
	"/roles/:userID",
}
