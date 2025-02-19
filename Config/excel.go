package Config

import (
	"github.com/extrame/xls"
	"github.com/tealeg/xlsx"
	"github.com/xuri/excelize/v2"
	"main/Logger"
)

func LoadExcel(excelPath string) ([][]string, error) {
	var rows [][]string
	f, err := excelize.OpenFile(excelPath)
	if err != nil {
		Logger.LogE("can not open file '%s': %v", excelPath, err)
		return rows, err
	}
	sheetMap := f.GetSheetMap()
	sheetName := ""
	for _, _name := range sheetMap {
		sheetName = _name
		break
	}
	rows, err = f.GetRows(sheetName)
	if err != nil {
		Logger.LogE("load sheet '%s' failed: %v", sheetName, err)
		return rows, err
	}
	return rows, nil
}

func TransformXls2Xlsx(from string, to string) error {
	xlsFile, err := xls.Open(from, "utf-8")
	if err != nil {
		Logger.LogE("can not open .xls file '%s': %v", from, err)
		return err
	}

	xlsxFile := xlsx.NewFile()

	for i := 0; i < xlsFile.NumSheets(); i++ {
		sheet := xlsFile.GetSheet(i)

		newSheet, err := xlsxFile.AddSheet(sheet.Name)
		if err != nil {
			Logger.LogE("can not add sheet '%s': %v", sheet.Name, err)
			return err
		}

		for row := 0; row < int(sheet.MaxRow); row++ {
			currentRow := sheet.Row(row)
			newRow := newSheet.AddRow()

			for j := 0; j < currentRow.LastCol(); j++ {
				cell := currentRow.Col(j)
				newCell := newRow.AddCell()
				newCell.Value = cell
			}
		}
	}
	err = xlsxFile.Save(to)
	if err != nil {
		Logger.LogE("save .xlsx file '%s' failed: %v", to, err)
		return err
	}
	return nil
}
