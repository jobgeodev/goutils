package MainApp

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
	// "utils/util/Parser"
)

type MainApp struct {
	LogDirName       string //
	LogFileName      string //
	LogFilePath      string //
	ServerConfigPath string // 配置信息
	// ServerConfig     Parser.ServerConfig // 配置信息
	Paras   []Para // 参数列表
	AppName string // 应用名称
}

type Para struct {
	Value interface{} `json:"value"` //
	Name  string      `json:"name"`  //
	Usage string      `json:"usage"` //
}

func NewMainApp(app_name string, log_dir string) *MainApp {
	return &MainApp{
		AppName:    app_name,
		LogDirName: log_dir,
	}
}

func (me *MainApp) SetLog(log_dir string, log_file string) error {
	me.LogDirName = log_dir
	me.LogFileName = log_file

	if err := os.MkdirAll(me.LogDirName, os.ModePerm); err != nil {
		return err
	}

	me.LogFilePath = filepath.Join(me.LogDirName, me.LogFileName)

	if logFile, logErr := os.OpenFile(me.LogFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666); logErr != nil {
		log.Printf("Fail to OpenFile Error:[%v]", logErr)
		return logErr
	} else {
		multi_writer := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multi_writer)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}

	return nil
}

// func (me *MainApp) LoadServerConfig(run_env string) error {
// 	// server_config_files := map[string]string{
// 	// 	"dev":     "server_config_dev.json",
// 	// 	"test":    "server_config_test.json",
// 	// 	"release": "server_config_release.json",
// 	// }

// 	switch run_env {
// 	case "dev":
// 		me.ServerConfigPath = "server_config_dev.json"
// 	case "test":
// 		me.ServerConfigPath = "server_config_test.json"
// 	case "release":
// 		me.ServerConfigPath = "server_config_release.json"
// 	default:
// 		tip := fmt.Sprintf("fail to parse run_env [%v] not in [dev,test,release] ", run_env)
// 		return errors.New(tip)
// 	}

// 	cfg_client := Parser.NewParserClient(me.ServerConfigPath)

// 	if cfg, err := cfg_client.ParseServerConfig(); err != nil {
// 		return err
// 	} else {
// 		me.ServerConfig = *cfg
// 	}

// 	return nil
// }

// func (me *MainApp) LoadParaConfig(content []byte) error {
// 	if err := json.Unmarshal(content, &me.Paras); err != nil {
// 		log.Printf("json.Unmarshal [%s] Error [%s]", content, err)
// 		return err
// 	}

// 	// paras := make([]Para, 0)
// 	// paras = append(paras, Para{
// 	// 	Type:    "bool",
// 	// 	Name:    "h",
// 	// 	Default: false,
// 	// 	Value:   &run_help,
// 	// 	Usage:   "Usage Help     ",
// 	// })
// 	// paras = append(paras, MainApp.Para{
// 	// 	Type:    "int",
// 	// 	Name:    "p",
// 	// 	Default: 2310,
// 	// 	Value:   &run_port,
// 	// 	Usage:   "Service Port   ",
// 	// })
// 	// paras = append(paras, MainApp.Para{
// 	// 	Type:    "string",
// 	// 	Name:    "e",
// 	// 	Default: "test",
// 	// 	Value:   &run_env,
// 	// 	Usage:   "Run Environment",
// 	// })
// 	return nil
// }

// func (me *MainApp) AppendRunPara(t string, n string, d interface{}, v *interface{}, u string, clear bool) {
// 	if clear {
// 		me.Paras = make([]Para, 0)
// 	}

// 	me.Paras = append(me.Paras, Para{
// 		Type:    t,
// 		Name:    n,
// 		Default: d,
// 		Value:   v,
// 		Usage:   u,
// 	})

// 	switch t {
// 	case "bool":
// 		flag.BoolVar(&run_help, "h", run_help, Paras["h"])

// 	}
// }

func (me *MainApp) AppendBoolVar(p *bool, name string, value bool, usage string) {
	flag.BoolVar(p, name, value, usage)

	para := Para{
		Value: interface{}(p),
		Name:  name,
		Usage: usage,
	}
	me.Paras = append(me.Paras, para)
}

func (me *MainApp) AppendIntVar(p *int, name string, value int, usage string) {
	flag.IntVar(p, name, value, usage)

	para := Para{
		Value: interface{}(p),
		Name:  name,
		Usage: usage,
	}
	me.Paras = append(me.Paras, para)
}

func (me *MainApp) AppendStringVar(p *string, name string, value string, usage string) {
	flag.StringVar(p, name, value, usage)

	para := Para{
		Value: interface{}(p),
		Name:  name,
		Usage: usage,
	}
	me.Paras = append(me.Paras, para)
}

// func (me *MainApp) FlagParse() {
// 	flag.Parse()
// }

func (me *MainApp) PrintUsage() {
	log.Println("\nUsage :")
	flag.PrintDefaults()
}

func (me *MainApp) PrintParameters() {
	log.Println("Starting application...")
	log.Println("PrintParameters Start...")
	for _, para := range me.Paras {

		value_txt := ""
		switch para.Value.(type) {
		case *bool:
			value_txt = fmt.Sprintf("%v", *(para.Value.(*bool)))
		case *int:
			value_txt = fmt.Sprintf("%v", *(para.Value.(*int)))
		case *string:
			value_txt = fmt.Sprintf("%v", *(para.Value.(*string)))
		}
		log.Printf("%s     : [%s]", para.Usage, value_txt)
	}

	log.Println("PrintParameters Finish.")
}

func (me *MainApp) SimpleCheckParameters() {
	flag.Parse()

	log_file := fmt.Sprintf("%s-%s.log", me.AppName, time.Now().Format("20060102_150405"))
	if err := me.SetLog(me.LogDirName, log_file); err != nil {
		log.Printf("Init Client Error:[%v]", err)
		os.Exit(1)
	}

	help := false
	env := "test"
	for _, para := range me.Paras {
		if para.Name == "h" {
			help = *(para.Value.(*bool))
		} else if para.Name == "e" {
			env = *(para.Value.(*string))
		}
	}

	if help {
		me.PrintUsage()
		os.Exit(1)
	} else {
		me.PrintParameters()
	}

	// // 解析配置文件
	// if err := me.LoadServerConfig(env); err != nil {
	// 	log.Printf("Parse Config Error:[%v]", err)
	// 	os.Exit(1)
	// }

	// me.ServerConfig.Print()
}
