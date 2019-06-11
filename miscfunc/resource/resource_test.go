package resource

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestAccess(t *testing.T) {
	type args struct {
		ctx      context.Context
		proxyReq events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Access(tt.args.ctx, tt.args.proxyReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Access() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Access() = %v, want %v", got, tt.want)
			}
		})
	}
}
