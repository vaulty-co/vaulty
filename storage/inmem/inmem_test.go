package inmem

import "testing"

func TestLoadFromFile(t *testing.T) {
	st := NewStorage()
	st.LoadFromFile("./test-fixtures/routes.json")

}

// config, err := LoadConfig(filepath.Join(FixturePath, "config.hcl"))
// if err != nil {
// 	t.Fatalf("err: %s", err)
// }
