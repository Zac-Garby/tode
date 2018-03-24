function example() {
	var input = document.getElementById("search-field")

	var values = [
		"F = ma",
		"chain rule",
		"#503",
		"sin(x)^2",
		"pythagoras",
		"euler",
		"physics",
		"calculus",
		"E = 1/2 mv^2",
	]
	
	if (form("regex").checked) {
		values = [
			`^F = .*`,
			`^.* = hf$`,
			`a .{1,2} b`,
			`(sin|cos|tan)\\(x\\)\\^2`,
			`.*`,
			`[a-f] + [g-z]`,
			`euler|py(thagoras)?`,
		]
	}
	
	var value = values[Math.floor(Math.random() * values.length)]

	input.value = value
}

function form(name) {
	var form = document.querySelector("form#search")
	return form[name]
}

form("search-type").value = "normal"