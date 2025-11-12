package repository

import "time"

type TypeObject int
const(
	TypeDir TypeObject = iota 
	TypeFile
)
const accessRights = 0o755

type Object struct{
	Name string
	Type TypeObject
	Size int
	TimeModification time.Time 
}

type Options struct{
	ConflictingName bool
}

type BrowseFiles struct{
	Root string
	Opt Options
}