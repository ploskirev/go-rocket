package model

import "database/sql"

type Order struct {
	OrderUUID       sql.NullString
	UserUUID        sql.NullString
	PartUUIDs       sql.NullString // массив строк в моделе
	TotalPrice      sql.NullFloat64
	TransactionUUID sql.NullString
	PaymentMethod   sql.NullInt32 // enum в моделе
	Status          sql.NullString
	CreatedAt       sql.NullTime
	UpdatedAt       sql.NullTime
}
