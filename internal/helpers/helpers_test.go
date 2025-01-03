// AnhCao 2024
package helpers

import (
	"reflect"
	"testing"
)

func TestMapInterfaceToStruct(t *testing.T) {
	type SampleStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	tests := []struct {
		name    string
		input   interface{}
		want    *SampleStruct
		wantErr bool
	}{
		{
			name: "valid input",
			input: map[string]interface{}{
				"name":  "test",
				"value": 123,
			},
			want: &SampleStruct{
				Name:  "test",
				Value: 123,
			},
			wantErr: false,
		},
		{
			name: "invalid input",
			input: map[string]interface{}{
				"name":  "test",
				"value": "invalid",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := MapInterfaceToStruct[SampleStruct](test.input)
			if (err != nil) != test.wantErr {
				t.Errorf("MapInterfaceToStruct() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !test.wantErr && !reflect.DeepEqual(got, test.want) {
				t.Errorf("MapInterfaceToStruct() = %v, want %v", got, test.want)
			}
		})
	}
}
