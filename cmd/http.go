package cmd

import (
	"strings"

	"github.com/ckeyer/api/types"
	"github.com/ckeyer/attack/http"
	"github.com/ckeyer/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(httpCmd())
}

func httpCmd() *cobra.Command {
	var (
		headers []string
		opt     = types.HTTPOption{
			Headers: map[string]string{},
		}
	)

	cmd := &cobra.Command{
		Use:   "http",
		Short: "start a http attack.",
		PreRun: func(cmd *cobra.Command, args []string) {
			opt.Method = strings.ToUpper(opt.Method)
			if len(args) != 1 {
				logrus.Fatalln("url required.")
			}
			opt.Url = args[0]
			for _, hdr := range headers {
				sli := strings.SplitN(hdr, "=", 2)
				if len(sli) != 2 {
					logrus.Fatalf("invalid header %s", hdr)
				}
				opt.Headers[sli[0]] = sli[1]
			}
			logrus.Debugf("http option: %+v", opt)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := http.Execute(opt); err != nil {
				logrus.Fatalln(err)
			}
		},
	}

	cmd.Flags().StringVarP(&opt.Method, "method", "m", "GET", "http method.")
	cmd.Flags().Int64VarP(&opt.Goroutine, "client", "c", 2, "http clients use.")
	cmd.Flags().Int64VarP(&opt.Count, "count", "n", 10, "every client sent request times.")
	cmd.Flags().StringArrayVarP(&headers, "header", "H", []string{}, "custom headers.")
	cmd.Flags().BoolVarP(&opt.RandUA, "rand-useragent", "u", false, "use random user-agent for every http client.")

	return cmd
}
