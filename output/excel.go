package output

import (
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

type OutputExcel struct {
	Wb       *excelize.File
	filename string
}

func NewExcel(filename string) *OutputExcel {
	e := new(OutputExcel)
	e.Wb = excelize.NewFile()
	e.filename = filename
	return e
}

func (oe *OutputExcel) FillSheet(sheetname string, data [][]string) error {
	_, err := oe.Wb.NewSheet(sheetname)
	if err != nil {
		return err
	}
	sheet, err := oe.Wb.NewStreamWriter(sheetname)
	if err != nil {
		return err
	}

	for rIdx, listRow := range data {
		// SetRowが[]interface{}型のみ受け付けるので、スライスをそのまま使うことはできないしコピーして移すこともできない
		// ので、愚直にループでひとつずついれる
		irow := make([]interface{}, len(listRow))
		for idx, v := range listRow {
			irow[idx] = v
		}
		cell, err := excelize.CoordinatesToCellName(1, rIdx+1)
		if err != nil {
			return err
		}
		if err = sheet.SetRow(cell, irow); err != nil {
			return err
		}
	}
	return nil
}

func (oe *OutputExcel) Save(path string) error {
	if err := oe.Wb.SaveAs(filepath.FromSlash(filepath.Join(path, oe.filename))); err != nil {
		return err
	}
	return nil
}
