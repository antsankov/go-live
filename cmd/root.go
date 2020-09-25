package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-live",
	Short: "A simple server to host files in a directory.",
	Long:  `A simple server to host a directory. Can be used either for local development or for production`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		if dir == "" {
			dir = "."
		}
		port, _ := cmd.Flags().GetString("port")
		if port == "" {
			port = "8000"
		}
		port = ":" + port
		go printer(dir, port)
		startServer(dir, port)
	},
}

func printer(dir string, port string) {
	start := time.Now()
	for {
		fmt.Println("\033[2J")
		fmt.Println("go-live\n--")
		fmt.Println("Serving: " + dir)
		fmt.Println("Port: " + port)
		fmt.Println(time.Since(start).Round(time.Second))
		time.Sleep(100 * time.Millisecond)
	}
}

func startServer(dir string, port string) {
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
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
	rootCmd.Flags().StringP("port", "p", "", "Set port to run on.")
	rootCmd.Flags().StringP("dir", "d", "", "Set the directory to serve.")
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
