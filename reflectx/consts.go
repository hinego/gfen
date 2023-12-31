package reflectx

const EnumTypeTemplate = `import { request } from '@umijs/max';
import Decimal from 'decimal.js';  {{range $k, $v := .Data}} 




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
	{{end}} {{if $v.Data}}  
	export const DataColums: Func.Column[] = [ {{range $v.Data.Data}} 
			{
				name:"{{.Json}}",
				desc:"{{.Desc}}",
				type:"{{.Typescript}}",
				filter:"{{.Filter}}",
				table:"{{.Table}}",
				config:"{{.Config}}",
				ts:"{{.Ts}}",
				{{if .HideTable}}hideTable:true,{{else}}hideTable:false,{{end}}
				{{if .Ellipsis}}ellipsis:true,{{else}}ellipsis:false,{{end}}
				tooltip:"{{.Tooltip}}",
			}, {{end}}
	] {{end}} {{if $v.Create}}  
	export const CreateColums: Func.Column[] = [ {{range $v.Create.Data}} 
			{
				name:"{{.Json}}",
				desc:"{{.Desc}}",
				type:"{{.Typescript}}",
				filter:"{{.Filter}}",
				table:"{{.Table}}",
				config:"{{.Config}}",
				ts:"{{.Ts}}",
				{{if .HideTable}}hideTable:true,{{else}}hideTable:false,{{end}}
				{{if .Ellipsis}}ellipsis:true,{{else}}ellipsis:false,{{end}}
				tooltip:"{{.Tooltip}}",
			}, {{end}}
	] {{end}} {{if $v.Update}}  
	export const UpdateColums: Func.Column[] = [ {{range $v.Update.Data}} 
			{
				name:"{{.Json}}",
				desc:"{{.Desc}}",
				type:"{{.Typescript}}",
				table:"{{.Table}}",
				config:"{{.Config}}",
				filter:"{{.Filter}}",
				ts:"{{.Ts}}",
				{{if .HideTable}}hideTable:true,{{else}}hideTable:false,{{end}}
				{{if .Ellipsis}}ellipsis:true,{{else}}ellipsis:false,{{end}}
				tooltip:"{{.Tooltip}}",
			}, {{end}}
	] {{end}} {{range .Fields}}
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
}  {{end}}
export namespace Func {
	export interface Column
    {
        name: string;
        desc: string;
		tooltip: string;
        type: string;
        sorter?: boolean;
        required?: boolean;
		table?: string;
		ts: string;
		hideTable: boolean;
		ellipsis: boolean;
		filter: string;
		config: string;
    }
	export interface Sorter 
	{
		name: string;
		type: string;
	}
	export const GetLocation = (): string =>
	{
		const pathname = location.pathname.replaceAll("//", "/")
		const split = pathname.split("/")
		if (split.length <= 4) {
			return pathname
		}
		return "/"+split.slice(-3).join("/")
	}
	export const MapFind = (name: string,table?: string) => {
		const path = Func.GetLocation().toLowerCase();
		const arr = path.split('/');
		const key = table ? "/" + table + "/" + name : "/" + arr[arr.length - 1] + "/" + name;
		switch (key) { {{range $k, $v := .Data}}   {{range .Enum}}
			case '/{{$v.Name | lower}}/{{.Name | lower}}': return {{$v.Name | ToName}}.{{.Name | ToName}}Map; {{end}} {{end}}
			default: return undefined;
		}
	}
	export const ArrayFind = (name: string,table?: string) => {
		const path = Func.GetLocation().toLowerCase();
		const arr = path.split('/');
		const key = table ? "/" + table + "/" + name : "/" + arr[arr.length - 1] + "/" + name;
		switch (path) { {{range $k, $v := .Data}}  {{range .Enum}}
			case '/{{$v.Name | lower}}/{{.Name | lower}}': return {{$v.Name | ToName}}.{{.Name | ToName}}Array; {{end}} {{end}}
			default: return undefined;
		}
	}
	export const ColumsFind = (name: string): Column[] | undefined  => {
		const path = Func.GetLocation().toLowerCase()+"/"+name.toLowerCase();
		switch (path) { {{range $k, $v := .Data}}   {{if $v.Data}}
			case '/{{$v.Version}}/{{$v.API}}/{{$v.File | lower}}/data': return {{$v.Name | ToName}}.DataColums; {{end}} {{if .Create}}
			case '/{{$v.Version}}/{{$v.API}}/{{$v.File | lower}}/create': return {{$v.Name | ToName}}.CreateColums; {{end}} {{if .Update}}
			case '/{{$v.Version}}/{{$v.API}}/{{$v.File | lower}}/update': return {{$v.Name | ToName}}.UpdateColums; {{end}} {{end}}
			default: return undefined;
		}
	}
	export const FunctionFind = (path: string)  => {
		path = path.toLowerCase();
		switch (path) { {{range $k, $v := .Data}}  {{range .Func}} 
			case '{{.Path | lower}}': return {{$v.Name | ToName}}.{{.Fun | ToName}}; {{end}} {{end}}
			default: return undefined;
		}
	}
	export const FetchFind = (name?: string)  => {
		const path = name ? name + "/fetch" : Func.GetLocation().toLowerCase()+"/fetch";
		switch (path) { {{range $k, $v := .Data}}  {{range .Func}} {{if eq .Fun "fetch"}}
			case '{{.Path | lower}}': return {{$v.Name | ToName}}.{{.Fun | ToName}}; {{end}} {{end}} {{end}}
			default: return undefined;
		}
	}
	export const GetFind = ()  => {
		const path = Func.GetLocation().toLowerCase()+"/get";
		switch (path) { {{range $k, $v := .Data}}  {{range .Func}} {{if eq .Fun "get"}}
			case '{{.Path | lower}}': return {{$v.Name | ToName}}.{{.Fun | ToName}}; {{end}} {{end}} {{end}}
			default: return undefined;
		}
	}
	export const UpdateFind = ()  => {
		const path = Func.GetLocation().toLowerCase()+"/update";
		switch (path) { {{range $k, $v := .Data}}  {{range .Func}} {{if eq .Fun "update"}}
			case '{{.Path | lower}}': return {{$v.Name | ToName}}.{{.Fun | ToName}}; {{end}} {{end}} {{end}}
			default: return undefined;
		}
	}
	export const CreateFind = ()  => {
		const path = Func.GetLocation().toLowerCase()+"/create";
		switch (path) { {{range $k, $v := .Data}}  {{range .Func}} {{if eq .Fun "create"}}
			case '{{.Path | lower}}': return {{$v.Name | ToName}}.{{.Fun | ToName}}; {{end}} {{end}} {{end}}
			default: return undefined;
		}
	}
	export const DeleteFind = ()  => {
		const path = Func.GetLocation().toLowerCase()+"/delete";
		switch (path) { {{range $k, $v := .Data}}  {{range .Func}} {{if eq .Fun "delete"}}
			case '{{.Path | lower}}': return {{$v.Name | ToName}}.{{.Fun | ToName}}; {{end}} {{end}} {{end}}
			default: return undefined;
		}
	}
	
}
export default { {{range $k, $v := .Data}} 
	{{$v.Name | ToName}}, {{end}}
	Func,
 };

`
