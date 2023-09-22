package gfen

import (
	"github.com/hinego/gfen/genx"
	"github.com/hinego/gfen/horm"
	_ "github.com/hinego/gfen/logic"
	"github.com/hinego/gfen/ssr"
)

func Skip(skip bool) {
	ssr.Gen().Skip(skip)
}
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
func Horm(data *horm.Input) (err error) {
	return ssr.Horm().Generate(data)
}
