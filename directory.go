package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tealeg/xlsx/v3"
)

const (
	Sheet_Least_Row  = 4
	Sheet_Max_Column = 1000
)

func fetchDirectoryContents(directory string) map[string]string {
	fileList, _ := getAllFiles(directory)
	return loadAllFiles(fileList)
}

func loadAllFiles(filesMap map[string]string) map[string]string {
	result := make(map[string]string)

	for _, fp := range filesMap {
		r, err := loadExcel(fp)
		if err == nil {
			for key, value := range r {
				result[key] = value
			}
		}
	}
	return result
}

func loadExcel(fp string) (map[string]string, error) {
	wb, err := xlsx.OpenFile(fp)
	if err != nil {
		return nil, err
	}

	filename := getFileBaseName(fp)
	if len(wb.Sheets) == 0 {
		return nil, errors.New("excel has no sheet.")
	}

	sheet := wb.Sheets[0]
	maxRow := sheet.MaxRow
	maxCol := sheet.MaxCol

	if maxRow == 0 {
		return nil, errors.New("excel sheet has not row.")
	}

	if maxCol >= Sheet_Max_Column {
		return nil, errors.New("excel sheet has max column.")
	}

	if maxRow == Sheet_Least_Row {
		return nil, nil
	}

	keys := make([]string, maxCol)

	keyRow, _ := sheet.Row(0)
	for i := 0; i < maxCol; i++ {
		cell := keyRow.GetCell(i)
		if cell != nil {
			keys[i] = strings.TrimSpace(cell.String())
		}
	}

	excelData := make(map[string]string)
	for i := 4; i < maxRow; i++ {
		row, err := sheet.Row(i)
		if err != nil {
			return nil, err
		}
		id := strings.TrimSpace(row.GetCell(0).String())
		if id == "" {
			continue
		}

		for j := 0; j < maxCol; j++ {
			if len(keys[j]) == 0 {
				continue
			}

			key := keys[j]
			value := row.GetCell(j).String()
			value = strings.Replace(value, "\n", "", -1)
			value = strings.Replace(value, "\t", "", -1)
			value = strings.Replace(value, " ", "", -1)
			saveKey := fmt.Sprintf("%s.%s.%s.R%d.L%d", filename, id, key, i, j)
			excelData[saveKey] = value
		}
	}
	return excelData, nil
}

func getFileBaseName(fp string) string {
	filenameWithSuffix := filepath.Base(fp)
	fileSuffix := filepath.Ext(filenameWithSuffix)
	return strings.TrimSuffix(filenameWithSuffix, fileSuffix)
}

func getAllFiles(dirPath string) (map[string]string, error) {
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	PathSep := string(os.PathSeparator)

	fm := make(map[string]string)

	for _, fi := range dir {
		realPath := fmt.Sprintf("%s%s%s", dirPath, PathSep, fi.Name())
		if fi.IsDir() {
			t, _ := getAllFiles(realPath)
			for k, v := range t {
				fm[k] = v
			}
		} else {
			ok := (strings.HasSuffix(fi.Name(), ".xlsx") || strings.HasSuffix(fi.Name(), ".xls")) && !strings.HasPrefix(fi.Name(), "~$")
			if ok {
				fm[getFileBaseName(fi.Name())] = realPath
			}
		}
	}

	return fm, nil
}
