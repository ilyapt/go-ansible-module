package ansible_module

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type Module struct {
	argumentData []byte
	responseLock *sync.Mutex
	response map[string]interface{}
	logLock *sync.Mutex
	logLines []string
}

func New(input interface{}) *Module {
	m := &Module{
		responseLock: &sync.Mutex{},
		response: make(map[string]interface{}),
		logLock: &sync.Mutex{},
	}

	if len(os.Args) != 2 {
		m.FailWithErrorf("no argument file provided")
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		m.FailWithErrorf("couldn't read argument file %s: %s", os.Args[1], err.Error())
	}
	m.argumentData = data
	m.ParseArgs(input)
	return m
}

func (m *Module) ParseArgs(input interface{}) {
	err := json.Unmarshal(m.argumentData, input)
	if err != nil {
		m.FailWithErrorf("argument file %s not valid json: %s", os.Args[1], err.Error())
	}
}

func (m *Module) FailWithErrorf(format string, a ...interface{}) {
	m.Set("err", fmt.Sprintf(format, a...))
	m.printEndExit(false, true)
}

func (m *Module) FailWithError(err error) {
	m.FailWithErrorf(err.Error())
}

func (m *Module) FailIfError(err error) {
	if err != nil {
		m.FailWithErrorf(err.Error())
	}
}

func (m *Module) Done(changed bool) {
	m.printEndExit(changed, false)
}

func (m *Module) Set(key string, value interface{}) {
	m.responseLock.Lock()
	defer m.responseLock.Unlock()
	m.response[key] = value
}

func (m *Module) LogPrint(a ...interface{}) {
	m.logLock.Lock()
	defer m.logLock.Unlock()
	for _, l := range a {
		m.logLines = append(m.logLines, fmt.Sprint(l))
	}
}

func (m *Module) LogPrintf(format string, a ...interface{}) {
	m.logLock.Lock()
	defer m.logLock.Unlock()
	m.logLines = append(m.logLines, fmt.Sprintf(format, a...))
}

func (m *Module) printEndExit(changed, failed bool) {
	m.responseLock.Lock()
	defer m.responseLock.Unlock()
	m.response["changed"] = changed
	m.response["failed"] = failed
	m.response["logs"] = m.logLines
	data, err := json.Marshal(m.response)
	if err != nil {
		fmt.Println("{\"changed\":false, \"failed\":true, \"err\":\""+err.Error()+"\"}")
		os.Exit(1)
	}
	fmt.Println(string(data))
	if failed {
		os.Exit(1)
	}
	os.Exit(0)
}

