package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net"
	"os/exec"
	"strconv"
	"time"
)

type applicationInfo struct {
	Id              primitive.ObjectID `bson:"_id"`
	ApplicationName string             `bson:"applicationName"`
}

type bindingInfo struct {
	Port int
	Ip   string
}

const (
	IP1 = "10.200.67.5"
	IP2 = "10.200.67.6"
	IP3 = "10.200.67.10"
)

var IPS = [3]string{IP1, IP3, IP2}

func ArmeriaDebugCommand() *cobra.Command {
	var projectName string
	cmd := &cobra.Command{
		Use:   "armeria",
		Short: "在Chrome 中打开 Armeria 项目的 debug 页面",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(projectName) <= 0 {
				return errors.New("项目名称不能为空")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			openWithProjectName(projectName)
		},
	}
	cmd.Flags().StringVarP(&projectName, "projectName", "p", "", "项目名称")
	if err := cmd.MarkFlagRequired("projectName"); err != nil {
		panic(err)
	}
	return cmd
}

func openWithProjectName(name string) {
	bindingInfo := queryBindingInfoFromFile(name)
	if bindingInfo == nil {
		bindingInfo = queryBindingInfoFromMongo(name)
	}
	openChromeNewTab(bindingInfo.Ip, bindingInfo.Port)
}

func queryBindingInfoFromFile(name string) *bindingInfo {
	path := "/Users/reasonknow/Downloads/应用信息规划表.xlsx"
	excelFile, err := excelize.OpenFile(path)
	if err != nil {
		return nil
	}

	sheetName := excelFile.GetSheetName(1)
	rows := excelFile.GetRows(sheetName)
	for _, row := range rows {
		if row[1] == name {
			port, err := strconv.Atoi(row[4])
			if err != nil {
				fmt.Println("转换失败")
				return nil
			}
			for _, ip := range IPS {
				if checkIpAndPort(ip, port) {
					return &bindingInfo{
						Port: port,
						Ip:   ip,
					}
				}
			}
		}
	}
	return nil
}

func checkIpAndPort(ip string, port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 3*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func queryBindingInfoFromMongo(name string) *bindingInfo {
	// 连接MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://10.200.48.16:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	machine := client.Database("machine")
	application := machine.Collection("application")
	binding := machine.Collection("binding")
	var applicationResult applicationInfo
	err = application.FindOne(context.Background(), bson.D{{"applicationName", name}}).Decode(&applicationResult)
	if err != nil {
		log.Fatal("查询 mongo 项目信息失败")
	}
	var bindingInfo bindingInfo
	err = binding.FindOne(context.Background(), bson.D{{"applicationId", applicationResult.Id.Hex()}, {"role", "prepare"}}).Decode(&bindingInfo)
	if err != nil {
		log.Fatal("查询 mongo 绑定信息失败")
	}
	return &bindingInfo
}

func openChromeNewTab(ip string, port int) {
	url := fmt.Sprintf("http://%s:%s/internal/docs/#/", ip, strconv.Itoa(port+7))
	cmd := exec.Command("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--new-tab", url)
	err := cmd.Run()
	if err != nil {
		println(err.Error())
	}
}
