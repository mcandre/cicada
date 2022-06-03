package cicada

// Version is semver.
const Version = "0.0.6"

// DefaultLeadMonths provides additional time for engineers
// to implement version migrations prior to final end of life.
//
// Not too short that developers fail to migrate,
// not too long that developers forget to migrate.
const DefaultLeadMonths = 1
