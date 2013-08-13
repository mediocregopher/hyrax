package config

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
)

// dumpExampleConfig returns the string representation of a configuration file
// with all parameters filled in with their default values
func dumpExampleConfig() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n")
	buffer.WriteString("#####################################\n")
	buffer.WriteString("### Hyrax configuration\n")
	buffer.WriteString("#####################################\n")
	buffer.WriteString("\n")

	for name, param := range params {
		buffer.WriteString("# " + param.Description + "\n")

		var def string
		if param.Type == STRING {
			def = param.Default.(string)
		} else {
			def = strconv.Itoa(param.Default.(int))
		}

		buffer.WriteString(name + ": " + def + "\n")
		buffer.WriteString("\n")
	}
	return buffer.String()
}

// readConfig returns a map of the key/values found in a given configuration file.
// Note: these key/values don't have to actually correspond to expected parameters,
// that parsing is done elsewhere
func readConfig(file string) (map[string]string, error) {
	fi, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(fi)
	ret := map[string]string{}
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		line = strings.TrimRight(line, "\n")
		if len(line) > 0 && line[0] != '#' {
			spl := strings.Split(line, ":")
			name := strings.Trim(spl[0], " \t")
			val := strings.Trim(spl[1], " \t")
			ret[name] = val
		}
	}

	return ret, nil

}
