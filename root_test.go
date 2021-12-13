package main

//import "testing"
//
//func TestBash(t *testing.T) {
//	vals, err := InvokeBash("invoke-completion -s ")
//	if err != nil {
//		t.Error(err.Error())
//	}
//	assertContains(t, vals, rawValue{Value: "bash"})
//}
//
//func TestElvish(t *testing.T) {
//	vals, err := InvokeElvish("invoke-completion -s ")
//	if err != nil {
//		t.Error(err.Error())
//	}
//	assertContains(t, vals, rawValue{Value: "elvish "})
//}
//
//func TestFish(t *testing.T) {
//	vals, err := InvokeFish("invoke-completion -s ")
//	if err != nil {
//		t.Error(err.Error())
//	}
//	assertContains(t, vals, rawValue{Value: "fish"})
//}
//
//func TestOil(t *testing.T) {
//	vals, err := InvokeOil("invoke-completion -s ")
//	if err != nil {
//		t.Error(err.Error())
//	}
//	assertContains(t, vals, rawValue{Value: "oil"})
//}
//
//func TestPowershell(t *testing.T) {
//	vals, err := InvokePowershell("invoke-completion -s ")
//	if err != nil {
//		t.Error(err.Error())
//	}
//	assertContains(t, vals, rawValue{Value: "powershell "})
//}
//
//func TestXonsh(t *testing.T) {
//	vals, err := InvokeXonsh("invoke-completion -s ")
//	if err != nil {
//		t.Error(err.Error())
//	}
//	assertContains(t, vals, rawValue{Value: "xonsh "})
//}
//
//func TestZsh(t *testing.T) {
//	vals, err := InvokeZsh("invoke-completion -s ")
//	if err != nil {
//		t.Error(err.Error())
//	}
//	assertContains(t, vals, rawValue{Value: "zsh "})
//}
//
//func assertContains(t *testing.T, vals []*rawValue, expected rawValue) {
//	for _, v := range vals {
//		if (expected.Value == "" || v.Value == expected.Value) &&
//			(expected.Display == "" || v.Display == expected.Display) &&
//			(expected.Description == "" || v.Description == expected.Description) {
//			return
//		}
//	}
//	t.Errorf("expected %#v", expected)
//}
