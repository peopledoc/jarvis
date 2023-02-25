module jarvis

go 1.13

require (
	github.com/golang/mock v1.4.4
	github.com/gorilla/mux v1.7.4
	github.com/lunixbochs/vtclean v1.0.0 // indirect
	github.com/manifoldco/promptui v0.8.0
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/pelletier/go-toml v1.6.0 // indirect
	github.com/relex/aini v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.6.1
	github.com/stretchr/testify v1.2.2
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/text v0.3.2 // indirect
	gopkg.in/ini.v1 v1.51.1 // indirect
	gopkg.in/yaml.v2 v2.2.7
)

replace github.com/relex/aini => ./internal/pkg/aini
