# gormaid

This is a tool for go projects to automatically generate boilerplate code for CRUD methods using gorm.

## The problem
Suppose in the following code file we are declaring gorm structs
```go
type Position struct {
    HouseID string `json:"house_id" gorm:"primaryKey;comment:房间号"`
    ShelfID string `json:"shelf_id" gorm:"primaryKey;comment:笼架号"`
    X       string `json:"x" gorm:"primaryKey;comment:横坐标"`
    Y       string `json:"y" gorm:"primaryKey;comment:纵坐标"`
}

type Transfer struct {
	TaskID               uint      `json:"task_id" gorm:"primaryKey;comment:任务id"`
	Mice                 []string  `json:"mice,omitempty" binding:"required" gorm:"serializer:json;comment:小鼠编号"`
	SourceCageID         uint      `json:"src_cage_id,omitempty" gorm:"column:src_cage_id;primaryKey;comment:来源笼位"`
	SourcePosition       *Position `json:"src,omitempty" gorm:"embedded;embeddedPrefix:src_"`
	TargetCageID         uint      `json:"dest_cage_id,omitempty" gorm:"column:dest_cage_id;primaryKey;comment:目的笼位id"`
	TargetPosition       *Position `json:"dest,omitempty" gorm:"embedded;embeddedPrefix:dest_"`
	SourceFeeder         uint      `json:"source_feeder,omitempty" gorm:"comment:来源房间饲养员"`
	TargetFeeder         uint      `json:"target_feeder,omitempty" gorm:"comment:目的房间饲养员"`
}
```
whose schemas will be automatically created thanks to [GORM](gorm.io).
However, professional go developers tend to find it tedious and error-prone work to exhaustively write code for CRUD methods.
For instance, in order to query for a `Transfer` record using any field, one may exhaust all 14 fields including the one in the embedded `Position` struct.
The development work amount gets exploded when more structs are coming and they have to be thoroughly tested prior to going production.

## The goal
We've talked about the pain, and we need a solution. What about a `go:generate` command that takes away all our pains?
Things like
```go
//go:generate gormaid -struct Transfer -package tasktransfer -o postgres/transfer_crud.go
type Transfer struct {
	TaskID               uint      `json:"task_id" gorm:"primaryKey;comment:任务id"`
	Mice                 []string  `json:"mice,omitempty" binding:"required" gorm:"serializer:json;comment:小鼠编号"`
	SourceCageID         uint      `json:"src_cage_id,omitempty" gorm:"column:src_cage_id;primaryKey;comment:来源笼位"`
	SourcePosition       *Position `json:"src,omitempty" gorm:"embedded;embeddedPrefix:src_"`
	TargetCageID         uint      `json:"dest_cage_id,omitempty" gorm:"column:dest_cage_id;primaryKey;comment:目的笼位id"`
	TargetPosition       *Position `json:"dest,omitempty" gorm:"embedded;embeddedPrefix:dest_"`
	SourceFeeder         uint      `json:"source_feeder,omitempty" gorm:"comment:来源房间饲养员"`
	TargetFeeder         uint      `json:"target_feeder,omitempty" gorm:"comment:目的房间饲养员"`
}
```
where `-struct Transfer` means the struct to be parsed;
`-package tasktransfer` means the output package name;
`-o postgres/transfer_crud.go` means the relative output path of the generated file.
