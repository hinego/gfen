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
	])  {{end}}   {{range .Fields}}
	export interface {{.TypeName | ToName}} 
	{ {{range .Data}}
		 {{.Json}}{{if .Optional}}?{{end}}: {{.TypeNameArray}};  {{end}}
	}   {{end}}  {{range .Func}}  
	export async function {{.Fun | ToName}} ({{if .In.Have}}body: {{.In.TypeNameArray | ToName}} , {{end}}options?: { [key: string]: any })
	{
		 return request<{{.Out.TypeNameArray | ToName}}>('{{.Path}}', {
			  method: 'POST',
			  headers: { 'Content-Type': 'application/json' },
			  {{if .In.Have}}data: body,{{end}}
			  ...(options || {}),
		 });
	} {{end}}
} 
{{end}}

 export default { {{range $k, $v := .Data}} 
	{{$v.Name | ToName}}, {{end}}
 };

`
