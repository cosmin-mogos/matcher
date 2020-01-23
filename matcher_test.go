package matcher

import "testing"

type test1 struct {
	a int
}

type test2 struct {
	a string
}

type leg struct {
	distance float64
	duration float64
	steps    *string
}

type route struct {
	distance float64
	duration float64
	geometry string
	legs     []leg
}

type directions struct {
	routes []route
}

func TestSuccess(t *testing.T) {
	tests := []struct {
		name     string
		template string
		value    interface{}
	}{
		{
			name:     "concrete value",
			template: `{"a": 10}`,
			value:    test1{a: 10},
		},
		{
			name:     "any value",
			template: `{"a": ?}`,
			value:    test1{a: 7},
		},
		{
			name:     "any or omit value, value present",
			template: `{"a": *}`,
			value:    test1{a: 7},
		},
		{
			name:     "any or omit value, value not present",
			template: `{"a": *}`,
			value:    struct{}{},
		},
		{
			name:     "array with any object, int element",
			template: `{"a": [?]}`,
			value: struct {
				a []int
			}{
				a: []int{10},
			},
		},
		{
			name:     "array with any object, string element",
			template: `{"a": [?]}`,
			value: struct {
				a []string
			}{
				a: []string{"string"},
			},
		},
		{
			name: "complex example",
			template: `{
				"routes": [
					{
						"distance": 2287.5,
						"duration": 351,
						"geometry": ?,
						"legs": [
							{
								"distance": ?,
								"duration": 351,
								"steps": *
							}
						]
					}
				]
			}`,
			value: directions{
				routes: []route{
					{
						distance: 2287.6,
						duration: 351,
						geometry: "234JaRbMeKuRxImGqEqIzYaSzFrK",
						legs: []leg{
							{
								distance: 132,
								duration: 351,
								steps:    nil,
							},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			matches := match(test.template, test.value)

			if !matches {
				t.Fail()
			}
		})
	}
}

func TestFail(t *testing.T) {
	tests := []struct {
		name     string
		template string
		value    interface{}
	}{
		{
			name:     "concrete value",
			template: `{"a": 10}`,
			value:    test1{a: 11},
		},
		{
			name:     "concrete value of wrong type",
			template: `{"a": 10}`,
			value:    test2{a: "string"},
		},
		{
			name:     "wrong array size, any element",
			template: `{"a": [?,?]}`,
			value: struct {
				a []int
			}{
				a: []int{10},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			matches := match(test.template, test.value)

			if matches {
				t.Fail()
			}
		})
	}
}
