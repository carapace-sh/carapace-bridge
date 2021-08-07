package main

import "testing"

func TestBash(t *testing.T) {
	vals, err := invokeBash("invoke-completion -s ")
	if err != nil {
		t.Error(err.Error())
	}
	assertContains(t, vals, rawValue{Value: "bash"})
}

func TestElvish(t *testing.T) {
	vals, err := invokeElvish("invoke-completion -s ")
	if err != nil {
		t.Error(err.Error())
	}
	assertContains(t, vals, rawValue{Value: "elvish "})
}

func TestFish(t *testing.T) {
	vals, err := invokeFish("invoke-completion -s ")
	if err != nil {
		t.Error(err.Error())
	}
	assertContains(t, vals, rawValue{Value: "fish"})
}

func TestOil(t *testing.T) {
	vals, err := invokeOil("invoke-completion -s ")
	if err != nil {
		t.Error(err.Error())
	}
	assertContains(t, vals, rawValue{Value: "oil"})
}

func TestPowershell(t *testing.T) {
	vals, err := invokePowershell("invoke-completion -s ")
	if err != nil {
		t.Error(err.Error())
	}
	assertContains(t, vals, rawValue{Value: "powershell "})
}

func TestXonsh(t *testing.T) {
	vals, err := invokeXonsh("invoke-completion -s ")
	if err != nil {
		t.Error(err.Error())
	}
	assertContains(t, vals, rawValue{Value: "xonsh "})
}

func TestZsh(t *testing.T) {
	vals, err := invokeZsh("invoke-completion -s ")
	if err != nil {
		t.Error(err.Error())
	}
	assertContains(t, vals, rawValue{Value: "zsh "})
}

func assertContains(t *testing.T, vals []*rawValue, expected rawValue) {
	for _, v := range vals {
		if (expected.Value == "" || v.Value == expected.Value) &&
			(expected.Display == "" || v.Display == expected.Display) &&
			(expected.Description == "" || v.Description == expected.Description) {
			return
		}
	}
	t.Errorf("expected %#v", expected)
}
