package templates

import (
	"testing"
)

func TestExecute(t *testing.T) {
	type args struct {
		text string
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "sestatus",
			args: args{
				text: `{{repl regexSplit ":\\s+" (regexFind "(?m)^Current mode:.+$" .CmdSestatus) 2 | last}}`,
				data: map[string]interface{}{
					"CmdSestatus": `SELinux status:                 enabled
SELinuxfs mount:                /sys/fs/selinux
SELinux root directory:         /etc/selinux
Loaded policy name:             targeted
Current mode:                   enforcing
Mode from config file:          enforcing
Policy MLS status:              enabled
Policy deny_unknown status:     allowed
Max kernel policy version:      31
`,
				},
			},
			want: "enforcing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Execute(tt.args.text, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
