package genx

import "strings"

type FuncInfo struct {
	Name       string
	Parameters []string
	Returns    []string
	Packages   []string
}
type Logic struct {
	Folder string
	Name   string      //逻辑名称
	Main   bool        //是否是主逻辑
	Base   string      //逻辑小写名称
	Sub    bool        //是否是子逻辑
	Funcs  []*FuncInfo //此逻辑下的方法
}
type LogicData struct {
	Name   string
	Main   *Logic
	Source *Logic
	Data   map[string]*Logic
}

func (r *LogicData) Packages() []string {
	var pack = make(map[string]string)
	for _, v := range r.Data {
		for _, v1 := range v.Funcs {
			for _, v2 := range v1.Packages {
				pack[v2] = v2
			}
		}
	}
	var data = make([]string, 0)
	for k := range pack {
		data = append(data, k)
	}
	return data
}
func (r *LogicData) getMain() Logic {
	for _, v := range r.Data {
		if v.Main {
			return *v
		}
	}
	return Logic{
		Folder: r.Name,
		Name:   strings.Title(r.Name),
		Main:   false,
		Base:   r.Name,
		Sub:    false,
		Funcs:  make([]*FuncInfo, 0),
	}
}
func (r *LogicData) GetMain() *Logic {
	var data = r.getMain()
	for _, v := range r.Data {
		if !v.Sub {
			continue
		}
		data.Funcs = append(data.Funcs, &FuncInfo{
			Name:       strings.Title(v.Base),
			Parameters: []string{},
			Returns:    []string{"Face I" + strings.Title(v.Name)},
			Packages:   []string{},
		})
	}
	return &data
}
func (r *LogicData) GetSource() *Logic {
	var data = r.getMain()
	if len(data.Funcs) == 0 {
		return nil
	}
	data.Name = strings.Title(r.Name) + "_"
	data.Main = true
	data.Sub = false
	return &data
}
