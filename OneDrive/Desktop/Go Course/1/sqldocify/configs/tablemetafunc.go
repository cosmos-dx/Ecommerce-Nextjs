package configs

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func GetMetaTableInstance() *MetaTableList {
	once.Do(func() {
		instance = &MetaTableList{
			ExistingTables: make(map[string]MetaTableDetails),
		}
		if err := loadActiveMetaTables(instance); err != nil {
			panic(err)
		}
	})
	return instance
}
func (t *MetaTableList) FindMetaTable(name string) *MetaTableDetails {
	t.mu.Lock()
	defer t.mu.Unlock()

	if details, exists := t.ExistingTables[name]; exists {
		return &details
	}
	return nil
}

func (tl *MetaTableList) UpdateMetaTable(tableName string, details MetaTableDetails) {
	tl.mu.Lock()
	defer tl.mu.Unlock()
	tl.ExistingTables[tableName] = details
	saveActiveMetaTables(tl, "activetables.json")
}

func (tl *MetaTableList) RemoveMetaTable(tableName string) {
	tl.mu.Lock()
	defer tl.mu.Unlock()
	delete(tl.ExistingTables, tableName)
	saveActiveMetaTables(tl, "activetables.json")
}

func loadActiveMetaTables(tl *MetaTableList) error {
	fileName := "activetables.json"

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer file.Close()

		tl.ExistingTables = make(map[string]MetaTableDetails)
		return saveActiveMetaTables(tl, fileName)
	}

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(byteValue, &tl.ExistingTables)
}

func saveActiveMetaTables(tl *MetaTableList, fileName string) error {
	data, err := json.MarshalIndent(tl.ExistingTables, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, data, 0644)
}
