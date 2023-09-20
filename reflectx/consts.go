package reflectx

const EnumTypeTemplate = `import { request } from '@umijs/max';
import Decimal from 'decimal.js'; 
{{range $k, $v := .Data}} 
export namespace {{$v.Name | ToName}} { {{range .Enum}} 
	export enum {{.Name | ToName}} { {{range .Enums}}   	
		{{.Name | ToName}} = {{.Typescript}}, {{end}}
	} 
	export const {{.Name | ToName}}Map = new Map([ {{range .Enums}}
		[{{.Typescript}}, { text: '{{.Desc}}' }],{{end}}
	])
	export const {{.Name | ToName}}Array = [ {{range .Enums}}
		{
			value: {{.Typescript}},
			label: '{{.Desc}}',
	  	},{{end}}
	] {{end}}   {{range .Enum}} 
	{{end}}
	
	{{range .Fields}}
	export interface {{.TypeName | ToName}} 
	{ {{range .Data}}
		 {{.Json}}{{if .Optional}}?{{end}}: {{.TypeNameArray}};  {{end}}
	}   {{end}}  {{range .Func}}  
	export async function {{.Fun | ToName}} ({{if .In.Have}}body{{if .In.IsOptional}}?{{end}}: {{.In.TypeNameArray | ToName}} , {{end}}options?: { [key: string]: any })
	{
		 return request<{{.Out.TypeNameArray | ToName}}>('{{.Path}}', {
			  method: 'POST',
			  headers: { 'Content-Type': 'application/json' },
			  {{if .In.Have}}data: body{{if .In.IsOptional}} || {}{{end}},{{end}}
			  ...(options || {}),
		 });
	} {{end}}
} 
{{end}}

export namespace Func {
	export const MapFind = (table: string,name: string) => {
		name = name.toLowerCase();
		table = table.toLowerCase();
		const text = table + '.' + name;
		switch (name) { {{range $k, $v := .Data}}   {{range .Enum}}
			case '{{$v.Name | lower}}.{{.Name | lower}}': return {{$v.Name | ToName}}.{{.Name | ToName}}Map; {{end}} {{end}}
			default: return undefined;
		}
	}
	export const ArrayFind = (table: string,name: string) => {
		name = name.toLowerCase();
		table = table.toLowerCase();
		const text = table + '.' + name;
		switch (name) { {{range $k, $v := .Data}}  {{range .Enum}}
			case '{{$v.Name | lower}}.{{.Name | lower}}': return {{$v.Name | ToName}}.{{.Name | ToName}}Array; {{end}} {{end}}
			default: return undefined;
		}
	}
}



 export default { {{range $k, $v := .Data}} 
	{{$v.Name | ToName}}, {{end}}
	Func,
 };

`
