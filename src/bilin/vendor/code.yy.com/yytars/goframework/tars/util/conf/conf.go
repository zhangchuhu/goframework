/*
有bug,当配置项不存在,返回的结果有问题
*/
package conf

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"os"
	"strconv"
	"strings"
)

type Conf struct {
	f string
}

func NewConf(f string) *Conf {
	c := new(Conf)
	c.Init(f)
	return c
}
func (c *Conf) Init(file string) {
	c.f = file
}

func (c *Conf) GetDomain(d string) []string {
	f, _ := os.Open(c.f)
	defer f.Close()
	decoder := xml.NewDecoder(f)
	l := strings.Split(d, "/")
	length := len(l)
	pos := 1
	var domains []string
EXIT:
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch element := t.(type) {
		case xml.StartElement:
			name := element.Name.Local
			if pos >= length {
				domains = append(domains, name)
				break
			}
			if name == l[pos] {
				continue
			}

		case xml.EndElement:
			name := element.Name.Local
			if name == l[length-1] {
				break EXIT
			}
		default:
		}
		pos++
	}
	return domains
}

func (c *Conf) GetInt(path string) int {
	s := c.GetString(path)
	r, err := strconv.Atoi(s)
	if err != nil {
		r = 0
	}
	return r
}

func (c *Conf) GetString(path string) string {
	// /taf/application/server<node>
	e := strings.Split(path, "<")
	kv := c.GetMap(e[0])
	key := strings.TrimRight(e[1], ">")
	return kv[key]
}

func (c *Conf) GetMap(path string) map[string]string {
	f, _ := os.Open(c.f)
	defer f.Close()
	decoder := xml.NewDecoder(f)
	l := strings.Split(path, "/")
	length := len(l)
	pos := 1
	var find bool
	kv := make(map[string]string)
EXIT:
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch element := t.(type) {
		case xml.StartElement:
			name := element.Name.Local
			if pos >= length || name == l[pos] {
				continue
			}
			find = true
		case xml.EndElement:
			name := element.Name.Local
			if name == l[length-1] {
				break EXIT
			}
		case xml.CharData:
			if !find {
				break
			}
			scanner := bufio.NewScanner(bytes.NewReader(element))
			for scanner.Scan() {
				line := scanner.Text()
				if strings.TrimSpace(line) == "" {
					continue
				}
				newline := strings.TrimSpace(line)
				linelist := strings.Split(newline, "#")
				k := strings.Split(linelist[0], "=")
				if len(k) < 2 {
					continue
				}
				kv[k[0]] = k[1]
			}
		}
	}
	pos++
	return kv
}
