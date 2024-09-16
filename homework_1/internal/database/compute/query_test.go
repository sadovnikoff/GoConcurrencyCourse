package compute

import (
	"testing"
)

func TestQuery_Command(t *testing.T) {
	tests := []struct {
		name string
		q    Query
		want string
	}{
		{
			name: "Empty",
			q:    NewQuery("", "", ""),
			want: "",
		},
		{
			name: "NotEmpty",
			q:    NewQuery("SomeCmd", "arg1", "agr2"),
			want: "SomeCmd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.q.Command()
			if cmd != tt.want {
				t.Errorf("want %q; got %q", tt.want, cmd)
			}
		})
	}
}

func TestQuery_KeyArgument(t *testing.T) {
	tests := []struct {
		name string
		q    Query
		want string
	}{
		{
			name: "Empty",
			q:    NewQuery("", "", ""),
			want: "",
		},
		{
			name: "NotEmpty",
			q:    NewQuery("CMD", "KeyArgument", "ValArgument"),
			want: "KeyArgument",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.q.KeyArgument()
			if cmd != tt.want {
				t.Errorf("want %q; got %q", tt.want, cmd)
			}
		})
	}
}

func TestQuery_ValueArgument(t *testing.T) {
	tests := []struct {
		name string
		q    Query
		want string
	}{
		{
			name: "Empty",
			q:    NewQuery("", "", ""),
			want: "",
		},
		{
			name: "NotEmpty",
			q:    NewQuery("CMD", "KeyArgument", "ValArgument"),
			want: "ValArgument",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.q.ValueArgument()
			if cmd != tt.want {
				t.Errorf("want %q; got %q", tt.want, cmd)
			}
		})
	}
}
