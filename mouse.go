package main

import (
	"gorm.io/gorm"
	"time"
)

type Gender uint8

const (
	GenderUnknown Gender = iota
	GenderMale
	GenderFemale
)

type Status uint8

const (
	StatusUnknown Status = iota
	StatusFeeding
	StatusGenotyping
	StatusGenotypeChecking
	StatusSeparatedGenotypeUnidentified
	StatusBreedingWaiting
	StatusExperimentWaiting
	StatusTransferWaiting
	StatusNewBorn
	StatusBreeding
	StatusExperimentOnGoing
	StatusSacrificeWaiting
	StatusBreedingFinished
	StatusDead
	StatusExperimentTerminal
	StatusSacrificed
	StatusPerformingExternalExperiment
	StatusTissueCollection
	StatusGenotypesUnknown
)

type StrainType struct {
	gorm.Model
	Name      string `gorm:"uniqueIndex;comment:基因类型" json:"name"`
	CreatedBy uint   `gorm:"index" json:"created_by,omitempty"`
}

type Genotype struct {
	gorm.Model
	StrainType string `gorm:"uniqueIndex:ui_sg,priority:1;comment:基因类型" json:"strain_type,omitempty"`
	Genotype   string `gorm:"uniqueIndex:ui_sg,priority:2;comment:基因型" json:"genotype,omitempty"`
	CreatedBy  uint   `json:"created_by,omitempty"`
}

type Strain struct {
	gorm.Model
	ProjectID   string        `json:"project_id" gorm:"comment:项目号"`
	Name        string        `json:"name" gorm:"uniqueIndex;comment:品系名"`
	StrainTypes []*StrainType `json:"strain_types,omitempty" gorm:"many2many:ass_strain_type;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedBy   uint          `json:"created_by" gorm:"comment:创建人"`
	Leader      uint          `json:"leader" gorm:"index;comment:品系负责人"`
	IsWT        bool          `json:"is_wt,omitempty" gorm:"index;comment:是否WT"`
}

type IdentifiedGenotypes struct {
	ID                 string      `gorm:"primaryKey;comment:主键" json:"-"`
	TaskID             uint        `gorm:"index;comment:任务id" json:"task_id,omitempty"`
	Identifier         uint        `gorm:"index;comment:鉴定员" json:"identifier,omitempty"`
	MouseID            string      `gorm:"index;comment:小鼠编号" json:"mouse_id,omitempty"`
	IsFinal            bool        `gorm:"index;default:false;comment:是否最终结果" json:"is_final,omitempty"`
	Genotypes          []*Genotype `gorm:"many2many:ass_identified_genotypes_genotype" json:"genotypes,omitempty"`
	StrainID           uint        `gorm:"index;comment:品系id" json:"strain_id"`
	Strain             Strain      `json:"strain" gorm:"foreignKey:StrainID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	IdentificationDate time.Time   `gorm:"autoUpdateTime;comment:鉴定日期" json:"identification_date"`
}

type Mouse struct {
	ID           string         `gorm:"primaryKey;comment:小鼠编号;" json:"id"`
	Year         int            `json:"year,omitempty" gorm:"index;comment:小鼠年份"`
	Number       int            `json:"number,omitempty" gorm:"index;comment:小鼠编号数字部分"`
	CreatedAt    time.Time      `json:"created_at,omitempty"`
	UpdatedAt    time.Time      `json:"updated_at,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Gender       Gender         `json:"gender" swaggertype:"primitive,string" enums:"male,female" gorm:"comment:性别"`
	ExperimentID string         `gorm:"index;default:NULL;comment:实验项目号" json:"experiment_id,omitempty"`
	// BirthdayOrArrivalDate is not necessarily the date of birth,
	// it can also be the arrival date depending on if the `WeekAgeOnArrival` is greater than 0
	BirthdayOrArrivalDate time.Time `json:"boad" gorm:"column:boad;comment:出生日期/到货日期"`
	// WeekAgeOnArrival is the week age of the mouse when it arrives at the facility,
	// if it's greater than 0 then the `BirthdayOrArrivalDate` means the arrival date
	WeekAgeOnArrival int                  `json:"waoa,omitempty" gorm:"column:waoa;comment:到货周龄"`
	ProjectID        string               `gorm:"index;comment:项目号" json:"project_id,omitempty"`
	StrainID         uint                 `json:"strain_id" gorm:"index;comment:品系号"`
	Strain           Strain               `json:"strain,omitempty" gorm:"foreignKey:StrainID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Genotype         *IdentifiedGenotypes `json:"genotype,omitempty" gorm:"foreignKey:MouseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Generation       int                  `json:"generation" gorm:"comment:代数"`
	Father           string               `json:"father,omitempty" gorm:"comment:父亲"`
	Mother1          string               `json:"mother1,omitempty" gorm:"comment:母亲1"`
	Mother2          string               `json:"mother2,omitempty" gorm:"comment:母亲2"`
	Supplier         string               `gorm:"index;default:'Sironax';comment:供应商" json:"supplier,omitempty"`
	Status           Status               `gorm:"index;comment:小鼠状态" json:"status" swaggertype:"primitive,string" enums:"feeding,genotyping,genotypechecking,unidentified,breedingwaiting,experimentwaiting,transferwaiting,newborn,breeding,experimentongoing,sacrificewaiting,dead,breedingfinished,endpoint,sacrificed,performingexternalexperiment,tissuecollection,genotypesunknown"`
	CreatorID        uint                 `json:"creator_id,omitempty" gorm:"comment:创建人"`
	External         string               `json:"external,omitempty" gorm:"index;default:NULL;comment:外部小鼠记录"`
	Remarks          string               `json:"remarks,omitempty"`
}
