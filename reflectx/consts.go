package reflectx

const EnumTypeTemplate = `import { request } from '@umijs/max';
import Decimal from 'decimal.js';

{{range $k, $v := .Data}} 
export namespace {{$v.Name | ToName}} {
	{{range .Enum}} 
	export enum {{.Name | ToName}} { {{range .Enums}}   	
		{{.Name | ToName}} = {{.Typescript}}, {{end}}
	} 
	export const {{.Name | ToName}}Map = new Map([ {{range .Enums}}
		[{{.Typescript}}, { text: '{{.Desc}}' }],{{end}}
	])  {{end}}
	{{range .Func}}  {{range .Fields}}
	export interface {{.TypeName | ToName}} 
	{ {{range .Data}}
		 {{.Json}}?: {{.TypeName}};  {{end}}
	}
	{{end}}
	export async function {{.Fun | ToName}} (body: {{.In.TypeName | ToName}} , options?: { [key: string]: any })
	{
		 return request<{{.Out.TypeName | ToName}}>('{{.Path}}', {
			  method: 'POST',
			  headers: { 'Content-Type': 'application/json' },
			  data: body,
			  ...(options || {}),
		 });
	} {{end}}
} 
{{end}}


 export default { {{range $k, $v := .Data}} 
	{{$v.Name | ToName}}, {{end}}
 };

`
