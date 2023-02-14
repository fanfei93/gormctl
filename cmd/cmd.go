package cmd

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"gitlab.2345.cn/gomod/gormctl/config"
	"gitlab.2345.cn/gomod/gormctl/view/gtools"
	"os"
)

var rootCmd = &cobra.Command{
	// 命令名称
	Use:   "main",
	Short: "gorm mysql reflect tools",
	// 命令执行过程
	Run: func(cmd *cobra.Command, args []string) {
		h, _ := cmd.Flags().GetString("host")
		u, _ := cmd.Flags().GetString("user")
		p, _ := cmd.Flags().GetString("password")
		d, _ := cmd.Flags().GetString("database")
		t, _ := cmd.Flags().GetString("table")
		o, _ := cmd.Flags().GetString("outdir")
		port, _ := cmd.Flags().GetInt("port")
		packageName, _ := cmd.Flags().GetString("p")
		conf := config.Config{
			Host:        h,
			User:        u,
			Password:    p,
			Database:    d,
			Port:        port,
			Table:       t,
			OutDir:      o,
			PackageName: packageName,
		}
		err := validator.New().Struct(conf)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		gtools.Execute(conf)
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("host", "H", "tst", "数据库地址.")
	rootCmd.MarkFlagRequired("host")
	rootCmd.PersistentFlags().StringP("user", "U", "", "用户名.")
	rootCmd.MarkFlagRequired("user")
	rootCmd.PersistentFlags().StringP("password", "P", "", "密码.")
	rootCmd.MarkFlagRequired("password")
	rootCmd.PersistentFlags().StringP("database", "D", "", "数据库名.")
	rootCmd.MarkFlagRequired("database")
	rootCmd.PersistentFlags().StringP("table", "T", "", "表名.")
	rootCmd.MarkFlagRequired("table")
	rootCmd.PersistentFlags().StringP("outdir", "O", "", "输出目录.")
	rootCmd.MarkFlagRequired("outdir")

	rootCmd.Flags().Int("port", 3306, "端口号")
	rootCmd.Flags().String("p", "model", "模型包名")
}
func Execute() {
	// 执行命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Errorf("error:%s", err)
		os.Exit(1)
	}
}
