	{{if $v.Data}}  
	export const DataColums: Func.Column[] = [ {{range $v.Data.Data}} 
			{
				name:"{{.Json}}",
				desc:"{{.Desc}}",
				type:"{{.Typescript}}",
			}, {{end}}
	] {{end}} {{if .Create}}  
	export const CreateColums: Func.Column[] = [ {{range $v.Data.Data}} 
			{
				name:"{{.Json}}",
				desc:"{{.Desc}}",
				type:"{{.Typescript}}",
			}, {{end}}
	] {{end}} {{if .Update}}  
	export const UpdateColums: Func.Column[] = [ {{range $v.Data.Data}} 
			{
				name:"{{.Json}}",
				desc:"{{.Desc}}",
				type:"{{.Typescript}}",
			}, {{end}}
	] {{end}}




    export const ColumsFind = (table: string,action string) => {
		table = table.toLowerCase();
		action = action.toLowerCase();
		const text = table + '.' + action;
		switch (action) { {{range $k, $v := .Data}}   {{if $v.Data}}
			case '{{$v.Name | lower}}.data': return {{$v.Name | ToName}}.DataColums; {{end}} {{if .Create}}
			case '{{$v.Name | lower}}.create': return {{$v.Name | ToName}}.CreateColums; {{end}} {{if .Update}}
			case '{{$v.Name | lower}}.update': return {{$v.Name | ToName}}.UpdateColums; {{end}}
			default: return undefined;
			{{end}}
	{{end}} }