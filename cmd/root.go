package cmd

import (
	"fmt"
	"os"

	"github.com/antsankov/go-live/lib"
	homedir "github.com/mitchellh/go-homedir"
	b "github.com/pkg/browser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-live",
	Short: "A simple server to host files in a directory.",
	Long:  `A simple server to host a directory. Can be used either for local development or for production`,
	Run: func(cmd *cobra.Command, args []string) {
		// Use a non-cobra check so we get more ways of catching version.
		if len(os.Args) > 1 {
			if os.Args[1] == "version" || os.Args[1] == "-v" || os.Args[1] == "--version" {
				printVersion()
				return
			}
		}

		dir, _ := cmd.Flags().GetString("dir")
		if dir == "" {
			dir = "./"
		}
		if dir[len(dir)-1] != '/' {
			dir += "/"
		}

		port, _ := cmd.Flags().GetString("port")
		if port == "" {
			port = "9000"
		}
		port = ":" + port
		url := fmt.Sprintf("http://localhost%s/", port)

		go lib.Printer(dir, port, url)

		quiet, _ := cmd.Flags().GetBool("quiet")
		if quiet == false {
			b.OpenURL(url)
		}

		cache, _ := cmd.Flags().GetBool("browser-cache")
		lib.StartServer(dir, port, cache)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-live.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("quiet", "q", false, "Open up the browser when started. Useful for development. Default: False")
	rootCmd.Flags().BoolP("use-browser-cache", "c", false, "Allow browser to cache page assets. Turn on for prod if not using nginx. Default: False")
	rootCmd.Flags().StringP("port", "p", "", "Set port to run on. Default: 9000")
	rootCmd.Flags().BoolP("version", "v", true, "Print the version of go-live.")
	rootCmd.Flags().StringP("dir", "d", "", "Set the directory to serve. Default: ./")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".go-live" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".go-live")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
