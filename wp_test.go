package gowp

import (
	"context"
	"errors"
	"testing"
	"time"
)

const (
	testDefaultNumTasks   = 10
	testDefaultNumWorkers = 2
)

var testNoOpFunc = func() error { return nil }

var testErr = errors.New("test error")

var testFuncWithErr = func() error { return testErr }

func TestNew(t *testing.T) {
	type args struct {
		numTasks int
		opts     []Option
	}
	tests := []struct {
		name    string
		args    args
		errVal  error
		wantErr bool
	}{
		{
			name:    "invalid number of tasks",
			args:    args{numTasks: -1},
			errVal:  ErrInvalidBuffer,
			wantErr: true,
		},
		{
			name:    "invalid context",
			args:    args{numTasks: testDefaultNumTasks, opts: []Option{WithContext(nil)}},
			errVal:  ErrNilContext,
			wantErr: true,
		},
		{
			name:    "invalid number of workers",
			args:    args{numTasks: testDefaultNumTasks, opts: []Option{WithNumWorkers(-1)}},
			errVal:  ErrInvalidWorkerCnt,
			wantErr: true,
		},
		{
			name:    "valid configuration",
			args:    args{numTasks: testDefaultNumTasks, opts: []Option{WithExitOnError(true)}},
			errVal:  nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.numTasks, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !errors.Is(err, tt.errVal) {
				t.Errorf("New() = %v, want %v", err, tt.errVal)
			}
		})
	}
}

func TestPool_Submit(t *testing.T) {
	type args struct {
		t Task
	}
	tests := []struct {
		name    string
		p       *Pool
		args    args
		wantErr bool
		errVal  error
		setup   func(p *Pool)
	}{
		{
			name:    "nil task",
			p:       newPool(context.Background(), testDefaultNumWorkers, testDefaultNumTasks, true),
			args:    args{t: nil},
			wantErr: true,
			errVal:  ErrNilTask,
		},
		{
			name:    "submit task on closed pool",
			p:       newPool(context.Background(), testDefaultNumWorkers, testDefaultNumTasks, true),
			args:    args{t: testNoOpFunc},
			wantErr: true,
			errVal:  ErrPoolClosed,
			setup: func(p *Pool) {
				_ = p.Wait()
			},
		},
		{
			name:    "submit task while pool is closing",
			p:       newPool(context.Background(), testDefaultNumWorkers, testDefaultNumTasks, true),
			args:    args{t: testNoOpFunc},
			wantErr: true,
			errVal:  ErrInvalidSend,
			setup: func(p *Pool) {
				close(p.in)
			},
		},
		{
			name:    "submit task after exhausting buffer",
			p:       newPool(context.Background(), 1, 1, true),
			args:    args{t: testNoOpFunc},
			wantErr: true,
			errVal:  ErrNoBuffer,
			setup: func(p *Pool) {
				_ = p.Submit(func() error {
					time.Sleep(time.Second)
					return nil
				})

				_ = p.Submit(func() error {
					time.Sleep(time.Second)
					return nil
				})
			},
		},
		{
			name:    "submit task with no error",
			p:       newPool(context.Background(), testDefaultNumWorkers, testDefaultNumTasks, true),
			args:    args{t: testNoOpFunc},
			wantErr: false,
			errVal:  nil,
			setup:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt.p)
			}

			err := tt.p.Submit(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pool.Submit() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !errors.Is(err, tt.errVal) {
				t.Errorf("Pool.Submit() = %v, want %v", err, tt.errVal)
			}
		})
	}
}

func TestPool_Wait(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	tests := []struct {
		name    string
		p       *Pool
		wantErr bool
		errVal  error
		setup   func(p *Pool)
	}{
		{
			name:    "error reported by one of the task",
			p:       newPool(context.Background(), testDefaultNumWorkers, testDefaultNumTasks, true),
			wantErr: true,
			errVal:  testErr,
			setup: func(p *Pool) {
				err := p.Submit(testFuncWithErr)
				if err != nil {
					t.Errorf("Pool.Submit() error = %v", err)
				}
			},
		},
		{
			name:    "context cancelled",
			p:       newPool(ctx, testDefaultNumWorkers, testDefaultNumTasks, true),
			wantErr: true,
			errVal:  context.Canceled,
			setup: func(p *Pool) {
				cancel()
				err := p.Submit(testNoOpFunc)
				if err != nil {
					t.Errorf("Pool.Submit() error = %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt.p)
			}

			err := tt.p.Wait()
			if (err != nil) != tt.wantErr {
				t.Errorf("Pool.Wait() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !errors.Is(err, tt.errVal) {
				t.Errorf("Pool.Wait() = %v, want %v", err, tt.errVal)
			}
		})
	}
}
