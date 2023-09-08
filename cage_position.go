package main

type Position struct {
	HouseID string `json:"house_id" gorm:"primaryKey;comment:房间号"`
	ShelfID string `json:"shelf_id" gorm:"primaryKey;comment:笼架号"`
	X       string `json:"x" gorm:"primaryKey;comment:横坐标"`
	Y       string `json:"y" gorm:"primaryKey;comment:纵坐标"`
}

type AssociateCagePosition struct {
	Position
	CageID uint `json:"cage_id" gorm:"comment:笼位号"`
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
	FieldIgnored         string    `gorm:"-"`
	FieldIgnoreMigration string    `gorm:"<-:false;-:migration;"`
}
