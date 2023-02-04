package files

import "testing"

const testFilePath = "templates/logs/logs.go.tpl"

func TestBoxProviderImplementsProvider(t *testing.T) {
	var _ Provider = &EmbedProvider{}
}

func TestBoxProviderConstructs(t *testing.T) {
	p := NewBoxProvider()
	if p == nil {
		t.Errorf("failed to build EmbedProvider")
	}
}

func TestBoxProvider_Get(t *testing.T) {
	bytes, err := NewBoxProvider().GetFile(testFilePath)
	if err != nil || len(bytes) == 0 {
		t.Errorf("failed to GetFile %s", err)
	}
}

func TestBoxProvider_MustGet(t *testing.T) {
	bytes := NewBoxProvider().MustGetFile(testFilePath)
	if len(bytes) == 0 {
		t.Fail()
	}
}

func TestBoxProvider_GetTemplate(t *testing.T) {
	tpl, err := NewBoxProvider().GetTemplate(testFilePath)
	if err != nil || tpl == nil {
		t.Errorf("failed to GetTemplate %s", err)
	}
}

func TestBoxProvider_GetTemplateCache(t *testing.T) {
	p := NewBoxProvider()
	_, _ = p.GetTemplate(testFilePath)
	if len(p.cache) != 1 {
		t.Fail()
	}
}

func TestBoxProvider_MustGetTemplate(t *testing.T) {
	tpl := NewBoxProvider().MustGetTemplate(testFilePath)
	if tpl == nil {
		t.Fail()
	}
}
