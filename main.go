package gfen

import (
	"github.com/hinego/gfen/genx"
	_ "github.com/hinego/gfen/logic"
	"github.com/hinego/gfen/ssr"
)

func Parse(in *genx.LogicInput) (err error) {
	return ssr.Logic().Parse(in)
}
func API(in *genx.ApiInput) (err error) {
	return ssr.Ctrl().Generate(in)
}
func Dao(data *genx.DaoInput) (err error) {
	return ssr.Dao().Generate(data)
}
func GenModel(data *genx.DaoInput) (err error) {
	return ssr.Dao().GenModel(data)
}
