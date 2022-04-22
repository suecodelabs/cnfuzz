package config

type ResultProcessConfig struct {
	ContainerName string
	Namespace     string
	EndpointUrl   string
	ReportBucket  string
	Image         string
	AccessKey     string
	SecretKey     string
}

func CreateResultProcessConfig(jobName string, namespace string) *ResultProcessConfig {
	return &ResultProcessConfig{
		ContainerName: jobName,
		Namespace:     namespace,
		// FIXME get these values from config
		EndpointUrl:  "http://minio-fs:9000",
		ReportBucket: "s3://restler-reports",
		Image:        "amazon/aws-cli",

		AccessKey: "minio",    //viper.GetString(cmd.AccessKey),
		SecretKey: "minio123", //viper.GetString(cmd.SecretKey),
	}
}
