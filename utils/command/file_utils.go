// Package utils /*
/*
 文件工具，用户 在各种文件之间进行格式转换
*/
package utils

import (
	"bufio"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

var cols = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

// utils csv2excel -c /Users/reasonknow/Desktop/test.csv -e /Users/reasonknow/Desktop/test.xlsx -s ',&'

func FileCommand() *cobra.Command {
	var separator string
	var csvReadPath string
	var excelWritePath string
	cmd := &cobra.Command{
		Use:   "csv2excel",
		Short: "c2e",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(separator) <= 0 {
				separator = ","
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			csvToExcel(csvReadPath, excelWritePath, separator)
		},
	}
	cmd.Flags().StringVarP(&separator, "separator", "s", "", "分隔符")
	cmd.Flags().StringVarP(&csvReadPath, "csvReadPath", "c", "", "csv文件位置")
	cmd.Flags().StringVarP(&excelWritePath, "excelWritePath", "e", "", "Excel 文件的写入位置")
	if err := cmd.MarkFlagRequired("csvReadPath"); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired("excelWritePath"); err != nil {
		panic(err)
	}
	return cmd
}

func csvToExcel(csvReadPath string, excelWritePath string, separator string) string {

	readFile, err := os.Open(csvReadPath)
	if err != nil {
		return "打开文件失败"
	}
	defer func(readFile *os.File) {
		err := readFile.Close()
		if err != nil {
			panic("关闭文件失败")
		}
	}(readFile)

	scanner := bufio.NewScanner(readFile)
	excelFile := excelize.NewFile()

	index := 0
	for scanner.Scan() {
		text := scanner.Text()
		split := strings.Split(text, separator)
		index++
		for i, s := range split {
			idx := cols[i] + strconv.Itoa(index)
			excelFile.SetCellValue("Sheet1", idx, s)
		}
	}
	if err = excelFile.SaveAs(excelWritePath); err != nil {
		return "保存文件失败"
	}
	return "转换成功"
}
