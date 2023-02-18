package main

type Result struct {
	Status   uint16
	Warnings uint16

	InsertId     uint64
	AffectedRows uint64

	*Resultset
}
