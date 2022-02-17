package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type myCfg struct {
	Mysql mysqlCfg // 变量需要pulic，即首字母大写
	Kafka kafkaCfg
}

type mysqlCfg struct {
	Addr string
	Port string
}

type kafkaCfg struct {
	Addr string
	Port string
}

func main() {

	//设置默认值
	viper.SetDefault("ContentDir", "content")
	viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
	fmt.Println(viper.GetString("ContentDir"))
	fmt.Println(viper.GetStringMapString("Taxonomies"))

	//从配置文件读取
	viper.SetConfigName("config")         // 文件名
	viper.SetConfigType("yaml")           // 文件后缀
	viper.AddConfigPath("/etc/appname/")  // 文件路径1
	viper.AddConfigPath("$HOME/.appname") // 文件路径2
	viper.AddConfigPath(".")              // 文件路径3-当前工作空间
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("file not found")
		} else {
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}
	}
	fmt.Println(viper.GetString("testkey"))

	//生成配置文件
	//无参数表示写入到已存在的文件，该文件需要上文中提到的方法事先声明
	err := viper.WriteConfig()
	if err != nil {
		fmt.Println(err.Error())
	}
	//_ = viper.SafeWriteConfig()
	err = viper.WriteConfigAs("./config2")
	if err != nil {
		fmt.Println(err.Error())
	}
	//_ =viper.SafeWriteConfigAs("/path/to/my/.other_config")

	//设置值
	viper.Set("Verbose", true)
	viper.Set("testkey", "testvalue_new")
	fmt.Println(viper.GetBool("Verbose"))
	fmt.Println(viper.GetString("testkey"))

	//设置别名
	viper.RegisterAlias("nameA", "nameB") // nameB是nameA的别名
	viper.Set("nameA", "erer")
	fmt.Println(viper.GetString("nameB")) //erer

	//环境变量
	viper.SetEnvPrefix("WANGJUAN")
	_ = viper.BindEnv("id")       // 只有一个参数，会去查找 WANGJUAN_ID, 默认全转成大写字母
	os.Setenv("WANGJUAN_ID", "1") // 一般是外部设置环境变量
	id := viper.GetInt("id")
	fmt.Println(id) // 1

	_ = viper.BindEnv("id2", "ID") // 两个参数，不会自动加前缀，只会去查找 ID
	os.Setenv("ID", "2")
	id2 := viper.GetInt("id2")
	fmt.Println(id2) // 2

	os.Setenv("ID", "22")            // 动态变化的
	fmt.Println(viper.GetInt("id2")) //22

	os.Setenv("WANGJUAN_ID_AUTO", "3")
	viper.AutomaticEnv()                 //有这步就不需要上述的显示Bind了，会自动去环境变量中查找，默认转成大写并加前缀
	fmt.Println(viper.GetInt("id_auto")) //3

	os.Setenv("WANGJUAN_ID_REPLACER", "4")
	my_replacer := strings.NewReplacer(".", "_") //把key中的 . 替换成 _ 再去匹配
	viper.SetEnvKeyReplacer(my_replacer)
	fmt.Println(viper.GetInt("id.replacer")) //4

	viper.AllowEmptyEnv(false) // 默认为false，表示对于空值的环境变量，会尝试往下一个可能的地方查询

	//从flag读取
	flag.Int("flagname", 1234, "help message for flagname")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine) // 转换到pflag
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
	flagname := viper.GetInt("flagname")
	fmt.Println(flagname) //1234

	//从viper获取值
	viper.Set("my_key", "my_value") // GET & SET 大小写不敏感
	fmt.Println(viper.GetString("MY_KEY"))
	if viper.GetBool("my_bool") { // 判断是否存在
		fmt.Println("my_bool enabled")
	} else {
		fmt.Println("my_bool unabled")
	}
	fmt.Println(viper.AllSettings()) //所有配置项

	//Unmarshaling
	viper.Reset() //重置所有配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig()

	fmt.Println(viper.GetString("mysql.addr") + ":" + viper.GetString("mysql.port"))
	fmt.Println(viper.GetString("mysql.addr") + ":" + viper.GetString("mysql.port"))
	var C myCfg
	err = viper.Unmarshal(&C)
	if err != nil {
		panic(err)
	}

	fmt.Println(C.Mysql.Addr + ":" + C.Mysql.Port)
	fmt.Println(C.Kafka.Addr + ":" + C.Kafka.Port)

	//多实例viper
	viper_A := viper.New()
	viper_B := viper.New()
	viper_A.SetDefault("ContentDir", "content")
	viper_B.SetDefault("ContentDir", "foobar")

	fmt.Println(viper_A.GetString("ContentDir")) //content
	fmt.Println(viper_B.GetString("ContentDir")) //foobar

	//测试优先级
	viper.Reset() //重置所有配置
	//配置文件中为 testkey:testvalue
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig()

	//与先后顺序无关
	viper.SetEnvPrefix("WANGJUAN") //设置了外部环境变量 export WANGJUAN_TESTKEY=aaa
	viper.AutomaticEnv()
	fmt.Println(viper.GetString("testkey")) // aaa
}
